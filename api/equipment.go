package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type UnequipBody struct {
	Slot     string `json:"slot"`
	Quantity int    `json:"quantity"`
}

type EquipBody struct {
	Slot     string `json:"slot"`
	Quantity int    `json:"quantity"`
	Code     string `json:"code"`
}

type UnequipResponse struct {
	Data  UnequipData  `json:"data"`
	Error ErrorMessage `json:"error"`
}

type UnequipData struct {
	Cooldown  Cooldown      `json:"cooldown"`
	Slot      string        `json:"slot"`
	Item      CraftableItem `json:"item"`
	Character Character     `json:"character"`
}

type EquipResponse struct {
	Data  EquipData    `json:"data"`
	Error ErrorMessage `json:"error"`
}

type EquipData struct {
	Cooldown  Cooldown      `json:"cooldown"`
	Slot      string        `json:"slot"`
	Item      CraftableItem `json:"item"`
	Character Character     `json:"character"`
}

func (c *Svc) Unequip(item CraftableItem) error {
	fmt.Printf("Unquipping item: %s\n", item.Name)
	unequipResp, err := c.Client.Unequip(c.Character.Name, item)
	if err != nil {
		return fmt.Errorf("unequipping item: %w", err)
	}
	c.Character = &unequipResp.Character
	c.Character.WaitForCooldown()
	return nil
}

func (c *Svc) Equip(item CraftableItem) error {
	fmt.Printf("Equipping item: %s\n", item.Name)
	equipResp, err := c.Client.Equip(c.Character.Name, item)
	if err != nil {
		return fmt.Errorf("equipping item: %w", err)
	}
	c.Character = &equipResp.Character
	c.Character.WaitForCooldown()
	return nil
}

func (c *ArtifactsClient) Unequip(characterName string, item CraftableItem) (*UnequipData, error) {
	path := fmt.Sprintf("/my/%s/action/unequip", characterName)
	bodyStruct := UnequipBody{
		Slot:     item.Type,
		Quantity: 1,
	}

	bodyBytes, err := json.Marshal(bodyStruct)
	if err != nil {
		return nil, fmt.Errorf("marshalling body: %w", err)
	}

	resp, err := c.Do(http.MethodPost, path, nil, bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("executing gather request: %w", err)
	}

	unequipResp := UnequipResponse{}
	if err := json.Unmarshal(resp, &unequipResp); err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}

	return &unequipResp.Data, nil
}

func (c *ArtifactsClient) Equip(characterName string, item CraftableItem) (*EquipData, error) {
	fmt.Printf("Equipping item: %s\n", item.Name)
	path := fmt.Sprintf("/my/%s/action/equip", characterName)
	bodyStruct := EquipBody{
		Slot:     item.Type,
		Quantity: 1,
		Code:     item.Code,
	}

	bodyBytes, err := json.Marshal(bodyStruct)
	if err != nil {
		return nil, fmt.Errorf("marshalling body: %w", err)
	}

	resp, err := c.Do(http.MethodPost, path, nil, bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("executing gather request: %w", err)
	}

	equipResp := EquipResponse{}
	if err := json.Unmarshal(resp, &equipResp); err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}

	return &equipResp.Data, nil
}
