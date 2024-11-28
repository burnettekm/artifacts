package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type ItemResponse struct {
	Item  CraftableItem `json:"data"`
	Error ErrorMessage  `json:"error"`
}

type GetItemsResponse struct {
	Data  []CraftableItem `json:"data"`
	Total int             `json:"total"`
	Page  int             `json:"page"`
	Size  int             `json:"size"`
	Pages int             `json:"pages"`
}

type CraftableItem struct {
	Name        string   `json:"name"`
	Code        string   `json:"code"`
	Level       int      `json:"level"`
	Type        string   `json:"type"`
	Subtype     string   `json:"subtype"`
	Description string   `json:"description"`
	Effects     []Effect `json:"effects"`
	Craft       *Craft   `json:"craft"`
	Tradeable   bool     `json:"tradeable"`
}

type Effect struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type Craft struct {
	Skill    string       `json:"skill"`
	Level    int          `json:"level"`
	Items    []SimpleItem `json:"items"`
	Quantity int          `json:"quantity"`
}

func (c *ArtifactsClient) GetItem(code string) (*CraftableItem, error) {
	params := map[string]string{
		"code": code,
		"size": strconv.Itoa(100),
	}
	resp, err := c.Do(http.MethodGet, fmt.Sprintf("/items/%s", code), params, nil)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	itemResp := ItemResponse{}
	if err := json.Unmarshal(resp, &itemResp); err != nil {
		return nil, fmt.Errorf("unmarshalling resp: %w", err)
	}

	return &itemResp.Item, nil
}

func (c *ArtifactsClient) GetItems(pageNum int) ([]CraftableItem, error) {
	params := map[string]string{
		"page": strconv.Itoa(pageNum),
	}
	resp, err := c.Do(http.MethodGet, "/items", params, nil)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	itemResp := GetItemsResponse{}
	if err := json.Unmarshal(resp, &itemResp); err != nil {
		return nil, fmt.Errorf("unmarshalling resp: %w", err)
	}

	return itemResp.Data, nil
}

func (c *ArtifactsClient) CraftItem(characterName, code string, quantity int) (*SkillData, error) {
	path := fmt.Sprintf("/my/%s/action/crafting", characterName)
	bodyStruct := SimpleItem{
		Code:     code,
		Quantity: quantity,
	}
	bodyBytes, err := json.Marshal(bodyStruct)
	resp, err := c.Do(http.MethodPost, path, nil, bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("executing crafting request: %w", err)
	}

	craftingResp := SkillResponse{}
	if err := json.Unmarshal(resp, &craftingResp); err != nil {
		return nil, fmt.Errorf("unmarshalling payload: %w", err)
	}

	return &craftingResp.Data, nil
}
