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

func (c *Svc) Fight() (*FightResponse, error) {
	percentHealth := float64(c.Character.Hp) / float64(c.Character.MaxHP) * 100.0
	if percentHealth < 25 {
		fmt.Printf("Character HP below 25 percent: %.2f\n", percentHealth)
		return nil, nil
	}

	fmt.Println("Fighting!")
	path := fmt.Sprintf("/my/%s/action/fight", c.Character.Name)
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

	c.Character = &fightResp.Data.Character
	fmt.Printf("Result: %s\n", fightResp.Data.Fight.Result)
	fmt.Printf("XP Gained: %d\n", fightResp.Data.Fight.Xp)
	fmt.Printf("Character level: %d\n", c.Character.Level)
	fmt.Printf("XP to level: %d\n", c.Character.MaxXP-c.Character.XP)
	fmt.Printf("Drops received: %v\n", fightResp.Data.Fight.Drops)
	fmt.Printf("Gold received: %v\n", fightResp.Data.Fight.Gold)
	fmt.Printf("Character HP: %d\n", fightResp.Data.Character.Hp)
	fmt.Printf("Cooldown: %d seconds\n", fightResp.Data.Cooldown.TotalSeconds)

	c.Character.WaitForCooldown()

	return &fightResp, nil
}

func (c *Svc) ContinuousFightLoop() error {
	percentHealth := float64(c.Character.Hp) / float64(c.Character.MaxHP) * 100.0
	if percentHealth < 25 {
		fmt.Printf("Character HP below 25 percent: %.2f, HP: %d MaxHP: %d\n", percentHealth, c.Character.Hp, c.Character.MaxHP)
		if err := c.Rest(); err != nil {
			return fmt.Errorf("executing rest request: %w", err)
		}
		if err := c.ContinuousFightLoop(); err != nil {
			return fmt.Errorf("recursive rest fightloop: %w", err)
		}
	}

	fightResp, err := c.Fight()
	if err != nil {
		return fmt.Errorf("executing fight request: %w", err)
	}

	c.Character = &fightResp.Data.Character
	if err := c.ContinuousFightLoop(); err != nil {
		return fmt.Errorf("recursive fightloop: %w", err)
	}

	return nil
}
