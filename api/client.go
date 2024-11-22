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
