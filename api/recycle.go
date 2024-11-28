package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RecycleResponse struct {
	Data  RecycleData  `json:"data"`
	Error ErrorMessage `json:"error"`
}

type RecycleData struct {
	Cooldown  Cooldown       `json:"cooldown"`
	Details   RecycleDetails `json:"details"`
	Character Character      `json:"character"`
}

type RecycleDetails struct {
	Items []SimpleItem `json:"items"`
}

func (c *Svc) RecycleItems(characterName string) error {
	inventory := c.Characters[characterName].Inventory
	for _, item := range inventory {
		if item.Code == "" {
			continue
		}

		i := c.GetItem(item.Code)
		contentCode := i.Craft.Skill
		coords := c.GetCoordinatesByCode(contentCode)
		if _, err := c.MoveCharacter(characterName, coords[0].X, coords[0].Y); err != nil {
			return fmt.Errorf("moving to bank: %w", err)
		}

		path := fmt.Sprintf("/my/%s/action/recycling", characterName)
		body := SimpleItem{
			Code:     item.Code,
			Quantity: item.Quantity,
		}
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshalling body: %w", err)
		}
		resp, err := c.Client.Do(http.MethodPost, path, nil, bodyBytes)
		if err != nil {
			return fmt.Errorf("executing recycle %s request: %w", item.Code, err)
		}

		recycleResp := RecycleResponse{}
		if err := json.Unmarshal(resp, &recycleResp); err != nil {
			return fmt.Errorf("unmarshalling response: %w", err)
		}

		c.Characters[characterName] = &recycleResp.Data.Character
		c.Characters[characterName].WaitForCooldown()
	}

	return nil
}
