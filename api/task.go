package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AcceptTaskResponse struct {
	Data  AcceptTaskData `json:"data"`
	Error ErrorMessage   `json:"error"`
}

type CompleteTaskResponse struct {
	Data  CompleteTaskData `json:"data"`
	Error ErrorMessage     `json:"error"`
}

type AcceptTaskData struct {
	Cooldown  Cooldown  `json:"cooldown"`
	Character Character `json:"character"`
	Task      Task      `json:"task"`
}

type CompleteTaskData struct {
	Cooldown  Cooldown    `json:"cooldown"`
	Character Character   `json:"character"`
	Rewards   TaskRewards `json:"rewards"`
}

type Task struct {
	Code    string      `json:"code"`
	Type    string      `json:"type"`
	Total   int         `json:"total"`
	Rewards TaskRewards `json:"rewards"`
}

type TaskRewards struct {
	Gold  int          `json:"gold"`
	Items []SimpleItem `json:"items"`
}

func (c *Svc) AcceptTask(characterName string) (*AcceptTaskResponse, error) {
	fmt.Printf("Accepting task\n")
	path := fmt.Sprintf("/my/%s/action/task/new", characterName)
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

	c.Characters[characterName] = &acceptTaskResp.Data.Character
	fmt.Printf("Task code: %s\n", acceptTaskResp.Data.Task.Code)
	fmt.Printf("Task type: %s\n", acceptTaskResp.Data.Task.Type)
	fmt.Printf("Task total: %d\n", acceptTaskResp.Data.Task.Total)
	fmt.Printf("Task rewards: %s\n", acceptTaskResp.Data.Task.Rewards)
	c.Characters[characterName].WaitForCooldown()

	return &acceptTaskResp, nil
}

func (c *Svc) CompleteTask(characterName string) (*CompleteTaskResponse, error) {
	fmt.Printf("Accepting task\n")
	path := fmt.Sprintf("/my/%s/action/task/complete", c.Characters[characterName].Name)
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

	c.Characters[characterName] = &completeTaskResponse.Data.Character
	fmt.Printf("Task rewards: %v\n", completeTaskResponse.Data.Rewards)
	c.Characters[characterName].WaitForCooldown()

	return &completeTaskResponse, nil
}
