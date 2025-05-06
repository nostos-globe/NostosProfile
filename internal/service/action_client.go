package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

)

const (
	actionBaseURL = "https://actions.nostos-globe.me/api/actions"
)

type ActionClient struct {
	httpClient *http.Client
}

func NewActionClient() *ActionClient {
	return &ActionClient{
		httpClient: &http.Client{},
	}
}

func (c *ActionClient) CreateAction(action *models.Action) error {
	jsonData, err := json.Marshal(action)
	if err != nil {
		return fmt.Errorf("error marshaling action: %w", err)
	}

	req, err := http.NewRequest("POST", actionBaseURL+"/create", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}