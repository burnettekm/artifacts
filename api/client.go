package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	basePath   string
	AuthToken  string
	httpClient *http.Client
}

func NewClient(basePath, authToken string) *Client {
	return &Client{
		basePath:   basePath,
		AuthToken:  authToken,
		httpClient: http.DefaultClient,
	}
}

func (c *Client) GetCharacter(name string) (*CharacterResponse, error) {
	u, err := url.Parse(c.basePath)
	if err != nil {
		return nil, fmt.Errorf("parsing base path: %w", err)
	}
	u.Path = fmt.Sprintf("/characters/%s", name)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("preparing request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.AuthToken))

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %w", err)
	}

	charResp := CharacterResponse{}
	if err := json.Unmarshal(bytes, &charResp); err != nil {
		return nil, fmt.Errorf("unmarshalling resp payload: %w", err)
	}
	if charResp.Error.Code != 0 {
		return nil, fmt.Errorf("error response received: status code: %d, error message: %s", charResp.Error.Code, charResp.Error.Message)
	}
	return &charResp, nil
}

func (c *Client) MoveCharacter(characterName string, x, y int) (*MoveResponse, error) {
	u, err := url.Parse(c.basePath)
	if err != nil {
		return nil, fmt.Errorf("parsing base path: %w", err)
	}
	u.Path = fmt.Sprintf("/my/%s/action/move", characterName)

	reqBody := MoveRequestBody{
		X: x,
		Y: y,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshalling body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("preparing request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.AuthToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %w", err)
	}

	moveResp := MoveResponse{}
	if err := json.Unmarshal(respBytes, &moveResp); err != nil {
		return nil, fmt.Errorf("unmarshalling resp payload: %w", err)
	}

	if moveResp.Error.Code != 0 {
		return nil, fmt.Errorf("error response received: status code: %d, error message: %s", moveResp.Error.Code, moveResp.Error.Message)
	}

	return &moveResp, nil
}

func (c *Client) Fight(name string) (*FightResponse, error) {
	u, err := url.Parse(c.basePath)
	if err != nil {
		return nil, fmt.Errorf("parsing base path: %w", err)
	}
	u.Path = fmt.Sprintf("/my/%s/action/fight", name)

	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("preparing request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.AuthToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %w", err)
	}

	fightResp := FightResponse{}
	if err := json.Unmarshal(respBytes, &fightResp); err != nil {
		return nil, fmt.Errorf("unmarshalling resp payload: %w", err)
	}

	if fightResp.Error.Code != 0 {
		return nil, fmt.Errorf("error response received: status code: %d, error message: %s", fightResp.Error.Code, fightResp.Error.Message)
	}

	return &fightResp, nil
}
