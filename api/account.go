package api

import (
	"encoding/json"
	"fmt"
)

func (c *ArtifactsClient) GetCharacter(name string) (*CharacterResponse, error) {
	path := fmt.Sprintf("/characters/%s", name)
	respBytes, err := c.Do("GET", path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("doing request: %w", err)
	}

	charResp := CharacterResponse{}
	if err := json.Unmarshal(respBytes, &charResp); err != nil {
		return nil, fmt.Errorf("unmarshalling resp payload: %w", err)
	}
	if charResp.Error.Code != 0 {
		return nil, fmt.Errorf("error response received: status code: %d, error message: %s", charResp.Error.Code, charResp.Error.Message)
	}
	return &charResp, nil
}

func (c *ArtifactsClient) GetCharacters() (*ListCharactersResponse, error) {
	path := "/my/characters"
	respBytes, err := c.Do("GET", path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("doing request: %w", err)
	}

	charResp := ListCharactersResponse{}
	if err := json.Unmarshal(respBytes, &charResp); err != nil {
		return nil, fmt.Errorf("unmarshalling resp payload: %w", err)
	}
	if charResp.Error.Code != 0 {
		return nil, fmt.Errorf("error response received: status code: %d, error message: %s", charResp.Error.Code, charResp.Error.Message)
	}
	return &charResp, nil
}
