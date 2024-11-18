package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type MapResponse struct {
	Data  []Map        `json:"data"`
	Error ErrorMessage `json:"error"`
}

type Map struct {
	Name    string  `json:"name"`
	Skin    string  `json:"skin"`
	X       int     `json:"x"`
	Y       int     `json:"y"`
	Content Content `json:"content"`
	Total   int     `json:"total"`
	Page    int     `json:"page"`
	Size    int     `json:"size"`
	Pages   int     `json:"pages"`
}

func (c *ArtifactsClient) GetMaps(contentCode, contentType *string) (*MapResponse, error) {
	p := map[string]string{}
	if contentCode != nil {
		p["content_code"] = *contentCode
	}
	if contentType != nil {
		p["content_type"] = *contentType
	}
	respBytes, err := c.Do(http.MethodGet, "/maps", p, nil)
	if err != nil {
		return nil, fmt.Errorf("executing GetMaps request: %w", err)
	}
	mapResp := MapResponse{}
	if err := json.Unmarshal(respBytes, &mapResp); err != nil {
		return nil, fmt.Errorf("unmarhsalling body: %w", err)
	}
	return &mapResp, nil
}

// query params content_code (name of thing), content_type (monster, resource, etc)
