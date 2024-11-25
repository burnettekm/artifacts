package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client interface {
	Do(method, path string, params map[string]string, body []byte) ([]byte, error)

	MoveCharacter(name string) (*MoveResponse, error)

	Unequip(characterName string, item CraftableItem) (*UnequipData, error)
	Equip(characterName string, item CraftableItem) (*EquipData, error)

	GetItem(code string) (*CraftableItem, error)
	CraftItem(characterName, code string, quantity int) (*SkillData, error)

	GetResource(code string) ([]ResourceData, error)
	Gather(characterName string) (*SkillData, error)

	GetMaps(contentCode, contentType *string) (*MapResponse, error)
}

type ArtifactsClient struct {
	basePath   string
	AuthToken  string
	httpClient *http.Client
}

func NewClient(authToken string) *ArtifactsClient {
	return &ArtifactsClient{
		basePath:   "https://api.artifactsmmo.com",
		AuthToken:  authToken,
		httpClient: http.DefaultClient,
	}
}

// todo: move to character file?
func (c *ArtifactsClient) MoveCharacter(name string, x, y int) (*MoveResponse, error) {
	fmt.Printf("Moving to %d, %d\n", x, y)
	path := fmt.Sprintf("/my/%s/action/move", name)
	reqBody := MoveRequestBody{
		X: x,
		Y: y,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshalling body: %w", err)
	}

	respBytes, err := c.Do(http.MethodPost, path, nil, bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("executing move request: %w", err)
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

func (c *ArtifactsClient) Do(method, path string, params map[string]string, body []byte) ([]byte, error) {
	u, err := url.Parse(c.basePath)
	if err != nil {
		return nil, fmt.Errorf("parsing base path: %w", err)
	}
	u.Path = path
	v := url.Values{}
	for key, value := range params {
		v.Add(key, value)
	}
	u.RawQuery = v.Encode()

	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(body))
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

	if resp.StatusCode > 200 {
		return nil, fmt.Errorf("got error response: %d", resp.StatusCode)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %w", err)
	}

	return respBytes, nil
}
