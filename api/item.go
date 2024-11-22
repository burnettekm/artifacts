package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ItemResponse struct {
	Item  CraftableItem `json:"data"`
	Error ErrorMessage  `json:"error"`
}

type CraftableItem struct {
	Name        string   `json:"name"`
	Code        string   `json:"code"`
	Level       int      `json:"level"`
	Type        string   `json:"type"`
	Subtype     string   `json:"subtype"`
	Description string   `json:"description"`
	Effects     []Effect `json:"effects"`
	Craft       Craft    `json:"craft"`
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
