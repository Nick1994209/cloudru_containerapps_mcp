package application

import (
	"cloudru-containerapps-mcp/internal/domain"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ContainerAppsApplication implements the ContainerAppsService interface
type ContainerAppsApplication struct{}

// NewContainerAppsApplication creates a new ContainerAppsApplication
func NewContainerAppsApplication() domain.ContainerAppsService {
	return &ContainerAppsApplication{}
}

// GetListContainerApps gets a list of ContainerApps from Cloud.ru API
func (c *ContainerAppsApplication) GetListContainerApps(projectID string, credentials domain.Credentials) ([]domain.ContainerApp, error) {
	// Get access token using KEY_ID and KEY_SECRET
	token, err := c.getAccessToken(credentials.KeyID, credentials.KeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Make request to ContainerApps API
	url := fmt.Sprintf("https://containers.api.cloud.ru/v1/containers?projectId=%s", projectID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response struct {
		Data []domain.ContainerApp `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse containerapps response: %w body: %s", err, body)
	}

	return response.Data, nil
}

// getAccessToken gets an access token using KEY_ID and KEY_SECRET
func (c *ContainerAppsApplication) getAccessToken(keyID, keySecret string) (string, error) {
	url := "https://iam.api.cloud.ru/auth/token"

	payload := strings.NewReader(fmt.Sprintf(`{
		"auth": {
			"identity": {
				"methods": ["api_key"],
				"api_key": {
					"id": "%s",
					"secret": "%s"
				}
			}
		}
	}`, keyID, keySecret))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("authentication failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response to get token
	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w %s", err, body)
	}

	return result.AccessToken, nil
}
