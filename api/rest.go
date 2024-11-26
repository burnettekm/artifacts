package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RestResponse struct {
	Rest  Rest         `json:"data"`
	Error ErrorMessage `json:"error"`
}

type Rest struct {
	Cooldown   Cooldown  `json:"cooldown"`
	HpRestored int       `json:"hp_restored"`
	Character  Character `json:"character"`
}

func (c *Svc) Rest() error {
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
