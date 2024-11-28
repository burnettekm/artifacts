package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type ResourceResponse struct {
	Data  []ResourceData `json:"data"`
	Error ErrorMessage   `json:"error"`
}

type ResourceData struct {
	Name  string `json:"name"`
	Code  string `json:"code"`
	Skill string `json:"skill"`
	Level int    `json:"level"`
	Drops []Drop `json:"drops"`

	Total int `json:"total"`
	Page  int `json:"page"`
	Size  int `json:"size"`
	Pages int `json:"pages"`
}

type Drop struct {
	Code        string `json:"code"`
	Rate        int    `json:"rate"`
	MinQuantity int    `json:"min_quantity"`
	MaxQuantity int    `json:"max_quantity"`
}

func (c *ArtifactsClient) GetResources(pageNumber int) ([]ResourceData, error) {
	path := "/resources"
	params := map[string]string{
		"size": strconv.Itoa(100),
		"page": strconv.Itoa(pageNumber),
	}

	resp, err := c.Do(http.MethodGet, path, params, nil)
	if err != nil {
		return nil, fmt.Errorf("executing gather request: %w", err)
	}

	resourceResponse := ResourceResponse{}
	if err := json.Unmarshal(resp, &resourceResponse); err != nil {
		return nil, fmt.Errorf("unmarshalling payload: %w", err)
	}

	return resourceResponse.Data, nil
}

func (c *ArtifactsClient) Gather(characterName string) (*SkillData, error) {
	fmt.Println("Gathering!")
	path := fmt.Sprintf("/my/%s/action/gathering", characterName)
	resp, err := c.Do(http.MethodPost, path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("executing gather request: %w", err)
	}

	gatherResp := SkillResponse{}
	if err := json.Unmarshal(resp, &gatherResp); err != nil {
		return nil, fmt.Errorf("unmarshalling payload: %w", err)
	}

	return &gatherResp.Data, nil
}
