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
	path := fmt.Sprintf("/my/%s/action/fight", c.Characters[characterName].Name)
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
	fmt.Printf("Character level: %d\n", c.Characters[characterName].Level)
	fmt.Printf("XP to level: %d\n", c.Characters[characterName].MaxXP-c.Characters[characterName].XP)
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
