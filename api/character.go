package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Service interface {
	Fight() (*FightResponse, error)
	FightLoop() error
	ContinuousFightLoop() error
	Rest() error

	AcceptTask() (*AcceptTaskResponse, error)
	CompleteTask() (*CompleteTaskResponse, error)

	CraftItem(code string) error
	Craft(code string, quantity int) error
	Gather() error

	WaitForCooldown()
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

func (c *CharacterSvc) CraftItem(code string) error {
	fmt.Printf("Attempting to craft item %s", code)

	item, err := c.Client.GetItem(code)
	if err != nil {
		return fmt.Errorf("getting item: %w", err)
	}

	// verify character can craft item
	if !c.Character.AbleToCraft(item.Craft.Skill, item.Craft.Level) {
		return fmt.Errorf("unable to craft item: required level: %d", item.Craft.Level)
	}

	// get dependent items
	requiredItems := []*CraftableItem{}
	for _, subItem := range item.Craft.Items {
		craftable, err := c.Client.GetItem(subItem.Code)
		if err != nil {
			return fmt.Errorf("getting item: %s: %w", subItem.Code, err)
		}

		// check if item in inventory
		found, quantity := c.FindItemInInventory(subItem.Code)
		if found && quantity >= subItem.Quantity {
			continue
		}

		if !c.Character.AbleToCraft(craftable.Craft.Skill, craftable.Craft.Level) {
			return fmt.Errorf("unable to craft subitem: %s: needs %s level: %d", craftable.Name, craftable.Craft.Skill, craftable.Craft.Level)
		}
		requiredItems = append(requiredItems, craftable)
	}

	// let's assume we're gathering for now
	fmt.Printf("Gathering required items to craft %s\n", code)
	for _, reqItem := range requiredItems {
		fmt.Printf("Gathering %v\n", reqItem)
		// find location of item
		contentType := "resource"
		mapResp, err := c.Client.GetMaps(&reqItem.Code, &contentType)
		if err != nil {
			return fmt.Errorf("finding item: %s: %w", reqItem.Code, err)
		}

		// move to item
		if _, err := c.MoveCharacter(mapResp.Data[0].X, mapResp.Data[0].Y); err != nil {
			return fmt.Errorf("moving to item: %w", err)
		}

		for i := 0; i <= reqItem.Craft.Quantity; i++ {
			// gather item
			if err := c.Gather(); err != nil {
				return fmt.Errorf("attempting to gather %s #%d: %w", reqItem.Name, i, err)
			}
		}
	}

	fmt.Println("Ready to craft item...")

	if err := c.Craft(code, 1); err != nil {
		return fmt.Errorf("crafting final item: %w", err)
	}

	fmt.Println("Successfully crafted item!")
	return nil
}

func (c *CharacterSvc) FindItemInInventory(code string) (bool, int) {
	for _, slot := range c.Character.Inventory {
		if slot.Code == code {
			return true, slot.Quantity
		}
	}

	return false, 0
}

func (c *CharacterSvc) Craft(code string, quantity int) error {
	fmt.Printf("Crafting %s!\n", code)
	path := fmt.Sprintf("/my/%s/action/crafting", c.Character.Name)
	bodyStruct := SimpleItem{
		Code:     code,
		Quantity: quantity,
	}
	bodyBytes, err := json.Marshal(bodyStruct)
	resp, err := c.Client.Do(http.MethodPost, path, nil, bodyBytes)
	if err != nil {
		return fmt.Errorf("executing crafting request: %w", err)
	}

	craftingResp := SkillResponse{}
	if err := json.Unmarshal(resp, &craftingResp); err != nil {
		return fmt.Errorf("unmarshalling payload: %w", err)
	}
	fmt.Printf("received %v", craftingResp.Data.Details.Items)
	c.Character = &craftingResp.Data.Character
	c.WaitForCooldown()
	return nil
}

func (c *CharacterSvc) Gather() error {
	fmt.Println("Gathering!")
	path := fmt.Sprintf("/my/%s/action/gathering", c.Character.Name)
	resp, err := c.Client.Do(http.MethodPost, path, nil, nil)
	if err != nil {
		return fmt.Errorf("executing gather request: %w", err)
	}

	gatherResp := SkillResponse{}
	if err := json.Unmarshal(resp, &gatherResp); err != nil {
		return fmt.Errorf("unmarshalling payload: %w", err)
	}
	fmt.Printf("received %v", gatherResp.Data.Details.Items)
	c.Character = &gatherResp.Data.Character
	c.WaitForCooldown()
	return nil
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
		return nil, fmt.Errorf("executing move request: %w", err)
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
	c.WaitForCooldown()

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
	c.WaitForCooldown()

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
	c.WaitForCooldown()

	return &completeTaskResponse, nil
}
