package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Service interface {
	Fight() (*FightResponse, error)
	FightLoop() error
	ContinuousFightLoop() error
	Rest() error

	AcceptTask() (*AcceptTaskResponse, error)
	CompleteTask() (*CompleteTaskResponse, error)

	MoveCharacter(x, y int) (*MoveResponse, error)
}

type CharacterSvc struct {
	Character *Character
	Client    *ArtifactsClient
}

type ListCharactersResponse struct {
	Characters []Character  `json:"data"`
	Error      ErrorMessage `json:"error"`
}

type CharacterResponse struct {
	Character Character    `json:"data"`
	Error     ErrorMessage `json:"error"`
}

func NewCharacterSvc(client *ArtifactsClient, char *Character) Service {
	return &CharacterSvc{
		Character: char,
		Client:    client,
	}
}

func (c *CharacterSvc) Unequip(item CraftableItem) error {
	fmt.Printf("Unquipping item: %s\n", item.Name)
	unequipResp, err := c.Client.Unequip(c.Character.Name, item)
	if err != nil {
		return fmt.Errorf("unequipping item: %w", err)
	}
	c.Character = &unequipResp.Character
	c.Character.WaitForCooldown()
	return nil
}

func (c *CharacterSvc) Equip(item CraftableItem) error {
	fmt.Printf("Equipping item: %s\n", item.Name)
	equipResp, err := c.Client.Equip(c.Character.Name, item)
	if err != nil {
		return fmt.Errorf("equipping item: %w", err)
	}
	c.Character = &equipResp.Character
	c.Character.WaitForCooldown()
	return nil
}

func (c *CharacterSvc) MoveCharacter(x, y int) (*MoveResponse, error) {
	if c.Character.X == x && c.Character.Y == y {
		fmt.Printf("character already at %d, %d\n", x, y)
		return nil, nil
	}

	moveResp, err := c.Client.MoveCharacter(c.Character.Name, x, y)
	if err != nil {
		return nil, fmt.Errorf("moving character: %w", err)
	}

	c.Character = &moveResp.Data.Character
	c.Character.WaitForCooldown()

	return moveResp, nil
}

func (c *CharacterSvc) Fight() (*FightResponse, error) {
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

func (c *CharacterSvc) FightLoop() error {
	percentHealth := float64(c.Character.Hp) / float64(c.Character.MaxHP) * 100.0
	if percentHealth < 25 {
		fmt.Printf("Character HP below 25 percent: %.2f, HP: %d MaxHP: %d\n", percentHealth, c.Character.Hp, c.Character.MaxHP)
		return nil
	}

	fightResp, err := c.Fight()
	if err != nil {
		return fmt.Errorf("executing fight request: %w", err)
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
	if err := c.FightLoop(); err != nil {
		return fmt.Errorf("recursive fightloop: %w", err)
	}

	return nil
}

func (c *CharacterSvc) ContinuousFightLoop() error {
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

func (c *CharacterSvc) Rest() error {
	fmt.Printf("Resting\n")
	path := fmt.Sprintf("/my/%s/action/rest", c.Character.Name)
	respBytes, err := c.Client.Do(http.MethodPost, path, nil, nil)
	if err != nil {
		return fmt.Errorf("executing rest request: %w", err)
	}

	restResp := RestResponse{}
	if err := json.Unmarshal(respBytes, &restResp); err != nil {
		return fmt.Errorf("unmarshalling resp payload: %w", err)
	}

	if restResp.Error.Code != 0 {
		return fmt.Errorf("error response received: status code: %d, error message: %s", restResp.Error.Code, restResp.Error.Message)
	}

	c.Character = &restResp.Rest.Character
	c.Character.WaitForCooldown()

	return nil
}

func (c *CharacterSvc) AcceptTask() (*AcceptTaskResponse, error) {
	fmt.Printf("Accepting task\n")
	path := fmt.Sprintf("/my/%s/action/task/new", c.Character.Name)
	respBytes, err := c.Client.Do(http.MethodPost, path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("executing accept task request: %w", err)
	}

	acceptTaskResp := AcceptTaskResponse{}
	if err := json.Unmarshal(respBytes, &acceptTaskResp); err != nil {
		return nil, fmt.Errorf("unmarshalling resp payload: %w", err)
	}

	if acceptTaskResp.Error.Code != 0 {
		return nil, fmt.Errorf("error response received: status code: %d, error message: %s", acceptTaskResp.Error.Code, acceptTaskResp.Error.Message)
	}

	c.Character = &acceptTaskResp.Data.Character
	fmt.Printf("Task code: %s\n", acceptTaskResp.Data.Task.Code)
	fmt.Printf("Task type: %s\n", acceptTaskResp.Data.Task.Type)
	fmt.Printf("Task total: %d\n", acceptTaskResp.Data.Task.Total)
	fmt.Printf("Task rewards: %s\n", acceptTaskResp.Data.Task.Rewards)
	c.Character.WaitForCooldown()

	return &acceptTaskResp, nil
}

func (c *CharacterSvc) CompleteTask() (*CompleteTaskResponse, error) {
	fmt.Printf("Accepting task\n")
	path := fmt.Sprintf("/my/%s/action/task/complete", c.Character.Name)
	respBytes, err := c.Client.Do(http.MethodPost, path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("executing complete task request: %w", err)
	}

	completeTaskResponse := CompleteTaskResponse{}
	if err := json.Unmarshal(respBytes, &completeTaskResponse); err != nil {
		return nil, fmt.Errorf("unmarshalling resp payload: %w", err)
	}

	if completeTaskResponse.Error.Code != 0 {
		return nil, fmt.Errorf("error response received: status code: %d, error message: %s", completeTaskResponse.Error.Code, completeTaskResponse.Error.Message)
	}

	c.Character = &completeTaskResponse.Data.Character
	fmt.Printf("Task rewards: %v\n", completeTaskResponse.Data.Rewards)
	c.Character.WaitForCooldown()

	return &completeTaskResponse, nil
}
