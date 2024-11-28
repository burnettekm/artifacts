package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type FightResponse struct {
	Data  FightData    `json:"data"`
	Error ErrorMessage `json:"error"`
}

type FightData struct {
	Cooldown  Cooldown  `json:"cooldown"`
	Fight     Fight     `json:"fight"`
	Character Character `json:"character"`
}

type Fight struct {
	Xp                 int          `json:"xp"`
	Gold               int          `json:"gold"`
	Drops              []SimpleItem `json:"drops"`
	Turns              int          `json:"turns"`
	MonsterBlockedHits BlockedHits  `json:"monster_blocked_hits"`
	PlayerBlockedHits  BlockedHits  `json:"player_blocked_hits"`
	Logs               []string     `json:"logs"`
	Result             string       `json:"result"`
}

type BlockedHits struct {
	Fire  int `json:"fire"`
	Earth int `json:"earth"`
	Water int `json:"water"`
	Air   int `json:"air"`
	Total int `json:"total"`
}

func (c *Svc) Fight(characterName string) (*FightResponse, error) {
	percentHealth := float64(c.Characters[characterName].Hp) / float64(c.Characters[characterName].MaxHP) * 100.0
	if percentHealth < 25 {
		fmt.Printf("Character HP below 25 percent: %.2f\n", percentHealth)
		return nil, nil
	}

	fmt.Println("Fighting!")
	path := fmt.Sprintf("/my/%s/action/fight", characterName)
	respBytes, err := c.Client.Do(http.MethodPost, path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("executing fight request: %w", err)
	}

	fightResp := FightResponse{}
	if err := json.Unmarshal(respBytes, &fightResp); err != nil {
		return nil, fmt.Errorf("unmarshalling resp payload: %w", err)
	}

	if fightResp.Error.Code != 0 {
		return nil, fmt.Errorf("error response received: status code: %d, error message: %s", fightResp.Error.Code, fightResp.Error.Message)
	}

	c.Characters[characterName] = &fightResp.Data.Character
	fmt.Printf("Result: %s\n", fightResp.Data.Fight.Result)
	fmt.Printf("XP Gained: %d\n", fightResp.Data.Fight.Xp)
	fmt.Printf("Character level: %d\n", fightResp.Data.Character.Level)
	fmt.Printf("XP to level: %d\n", fightResp.Data.Character.MaxXP-fightResp.Data.Character.XP)
	fmt.Printf("Drops received: %v\n", fightResp.Data.Fight.Drops)
	fmt.Printf("Gold received: %v\n", fightResp.Data.Fight.Gold)
	fmt.Printf("Character HP: %d\n", fightResp.Data.Character.Hp)
	fmt.Printf("Cooldown: %d seconds\n", fightResp.Data.Cooldown.TotalSeconds)

	c.Characters[characterName].WaitForCooldown()

	return &fightResp, nil
}

func (c *Svc) ContinuousFightLoop(characterName string) error {
	percentHealth := float64(c.Characters[characterName].Hp) / float64(c.Characters[characterName].MaxHP) * 100.0
	if percentHealth < 25 {
		fmt.Printf("Character HP below 25 percent: %.2f, HP: %d MaxHP: %d\n", percentHealth, c.Characters[characterName].Hp, c.Characters[characterName].MaxHP)
		if err := c.Rest(characterName); err != nil {
			return fmt.Errorf("executing rest request: %w", err)
		}
		if err := c.ContinuousFightLoop(characterName); err != nil {
			return fmt.Errorf("recursive rest fightloop: %w", err)
		}
	}

	fightResp, err := c.Fight(characterName)
	if err != nil {
		return fmt.Errorf("executing fight request: %w", err)
	}

	c.Characters[characterName] = &fightResp.Data.Character
	if err := c.ContinuousFightLoop(characterName); err != nil {
		return fmt.Errorf("recursive fightloop: %w", err)
	}

	return nil
}

func (c *Svc) FightForCrafting(characterName, dropCode string, quantity *int) error {
	maxLevel := c.Characters[characterName].Level
	monsters := c.GetMonsterByDrop(dropCode)

	// find monster with highest drop rate
	bestMonsterCode := ""
	maxDropRate := 0
	for _, monster := range monsters {
		if monster.Level > maxLevel {
			continue
		}
		for _, drop := range monster.Drops {
			if drop.Code != dropCode {
				continue
			}
			if drop.Rate > maxDropRate {
				maxDropRate = drop.Rate
				bestMonsterCode = monster.Code
			}
		}
	}

	// find selected monster
	coords := c.GetCoordinatesByCode(bestMonsterCode)
	if _, err := c.MoveCharacter(characterName, coords[0].X, coords[0].Y); err != nil {
		return fmt.Errorf("moving to bank: %w", err)
	}

	wantQuantity := 1000
	if quantity != nil {
		wantQuantity = *quantity
	}
	if err := c.ContinuousFightLoopForCrafting(characterName, dropCode, wantQuantity); err != nil {
		return fmt.Errorf("ContinuousFightLoopForCrafting: %w", err)
	}

	return nil
}

func (c *Svc) ContinuousFightLoopForCrafting(characterName, dropCode string, wantQuantity int) error {
	// if we have the quantity we want, stop
	runningTotal := 0

	// check bank
	bankQuantity, err := c.CheckBankForItem(dropCode)
	if err != nil {
		return fmt.Errorf("checking bank for item: %w", err)
	}
	runningTotal += bankQuantity

	// check inventory
	_, invQuantity := c.Characters[characterName].FindItemInInventory(dropCode)
	runningTotal += invQuantity

	if runningTotal >= wantQuantity {
		return nil
	}

	percentHealth := float64(c.Characters[characterName].Hp) / float64(c.Characters[characterName].MaxHP) * 100.0
	if percentHealth < 25 {
		fmt.Printf("Character HP below 25 percent: %.2f, HP: %d MaxHP: %d\n", percentHealth, c.Characters[characterName].Hp, c.Characters[characterName].MaxHP)
		if err := c.Rest(characterName); err != nil {
			return fmt.Errorf("executing rest request: %w", err)
		}
		//if err := c.ContinuousFightLoopForCrafting(characterName, dropCode, wantQuantity); err != nil {
		//	return fmt.Errorf("recursive rest fightloop: %w", err)
		//}
	}

	fightResp, err := c.Fight(characterName)
	if err != nil {
		return fmt.Errorf("executing fight request: %w", err)
	}
	c.Characters[characterName] = &fightResp.Data.Character

	if c.Characters[characterName].IsInventoryFull() {
		if err := c.DepositAllItems(characterName); err != nil {
			return fmt.Errorf("depositing all items: %w", err)
		}
	}

	if err := c.ContinuousFightLoopForCrafting(characterName, dropCode, wantQuantity); err != nil {
		return fmt.Errorf("recursive fightloop: %w", err)
	}

	return nil
}
