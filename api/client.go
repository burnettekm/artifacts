package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client interface {
	Do(method, path string, params map[string]string, body []byte) ([]byte, error)

	GetCharacter(name string) (*CharacterResponse, error)
	GetCharacters() ([]*Character, error)
	MoveCharacter(name string, x, y int) (*MoveResponse, error)

	Unequip(characterName string, item CraftableItem) (*UnequipData, error)
	Equip(characterName string, item CraftableItem) (*EquipData, error)

	GetItem(code string) (*CraftableItem, error)
	GetItems(pageNum int) ([]CraftableItem, error)
	CraftItem(characterName, code string, quantity int) (*SkillData, error)

	GetResources(pageNumber int) ([]ResourceData, error)
	Gather(characterName string) (*SkillData, error)

	GetMaps(contentCode, contentType *string) ([]Map, error)
	GetMonsters(pageNumber int) ([]MonsterData, error)
}

type ArtifactsClient struct {
	basePath   string
	AuthToken  string
	httpClient *http.Client
}

func NewClient(authToken string) Client {
	return &ArtifactsClient{
		basePath:   "https://api.artifactsmmo.com",
		AuthToken:  authToken,
		httpClient: http.DefaultClient,
	}
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

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("got error response: %d: %s", resp.StatusCode, string(respBytes))
	}

	return respBytes, nil
}
