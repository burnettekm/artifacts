package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Service interface {
	MoveCharacter(x, y int) (*MoveResponse, error)
	Fight() (*FightResponse, error)
	FightLoop() error
	ContinuousFightLoop() error
	Rest() error

	AcceptTask() (*AcceptTaskResponse, error)
	CompleteTask() (*CompleteTaskResponse, error)

	WaitForCooldown()
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

func (c *CharacterSvc) WaitForCooldown() {
	if c.Character.Cooldown == 0 {
		return
	}

	fmt.Printf("On cooldown for %d seconds\n", c.Character.Cooldown)

	time.Sleep(time.Duration(c.Character.Cooldown) * time.Second)
	fmt.Println("cooldown ended...")
	c.Character.Cooldown = 0
	return
}

func (c *CharacterSvc) MoveCharacter(x, y int) (*MoveResponse, error) {
	if c.Character.X == x && c.Character.Y == y {
		fmt.Printf("character already at %d, %d\n", x, y)
		return nil, nil
	}

	fmt.Printf("Moving to %d, %d\n", x, y)
	path := fmt.Sprintf("/my/%s/action/move", c.Character.Name)
	reqBody := MoveRequestBody{
		X: x,
		Y: y,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshalling body: %w", err)
	}

	respBytes, err := c.Client.Do(http.MethodPost, path, nil, bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}

	moveResp := MoveResponse{}
	if err := json.Unmarshal(respBytes, &moveResp); err != nil {
		return nil, fmt.Errorf("unmarshalling resp payload: %w", err)
	}

	if moveResp.Error.Code != 0 {
		return nil, fmt.Errorf("error response received: status code: %d, error message: %s", moveResp.Error.Code, moveResp.Error.Message)
	}

	c.Character = &moveResp.Data.Character
	c.WaitForCooldown()

	return &moveResp, nil
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
		return nil, fmt.Errorf("executing request: %w", err)
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

	c.WaitForCooldown()

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

	c.WaitForCooldown()
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
		return fmt.Errorf("executing request: %w", err)
	}

	restResp := RestResponse{}
	if err := json.Unmarshal(respBytes, &restResp); err != nil {
		return fmt.Errorf("unmarshalling resp payload: %w", err)
	}

	if restResp.Error.Code != 0 {
		return fmt.Errorf("error response received: status code: %d, error message: %s", restResp.Error.Code, restResp.Error.Message)
	}

	c.Character = &restResp.Rest.Character
	c.WaitForCooldown()

	return nil
}

func (c *CharacterSvc) AcceptTask() (*AcceptTaskResponse, error) {
	fmt.Printf("Accepting task\n")
	path := fmt.Sprintf("/my/%s/action/task/new", c.Character.Name)
	respBytes, err := c.Client.Do(http.MethodPost, path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
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
	c.WaitForCooldown()

	return &acceptTaskResp, nil
}

func (c *CharacterSvc) CompleteTask() (*CompleteTaskResponse, error) {
	fmt.Printf("Accepting task\n")
	path := fmt.Sprintf("/my/%s/action/task/complete", c.Character.Name)
	respBytes, err := c.Client.Do(http.MethodPost, path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
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
	c.WaitForCooldown()

	return &completeTaskResponse, nil
}
