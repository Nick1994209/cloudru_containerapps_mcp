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

// GetContainerApp gets a specific ContainerApp from Cloud.ru API
func (c *ContainerAppsApplication) GetContainerApp(projectID string, containerAppName string, credentials domain.Credentials) (*domain.ContainerApp, error) {
	// Get access token using KEY_ID and KEY_SECRET
	token, err := c.getAccessToken(credentials.KeyID, credentials.KeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Make request to ContainerApps API
	url := fmt.Sprintf("https://containers.api.cloud.ru/v1/containers/%s?projectId=%s", containerAppName, projectID)
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

	// Check if body is empty
	if len(body) == 0 {
		return nil, fmt.Errorf("API returned empty response body with status %d", resp.StatusCode)
	}

	// Parse response
	var containerApp domain.ContainerApp
	if err := json.Unmarshal(body, &containerApp); err != nil {
		return nil, fmt.Errorf("failed to parse containerapp response: %w body: %s", err, body)
	}

	return &containerApp, nil
}

// CreateContainerApp creates a new ContainerApp in Cloud.ru
func (c *ContainerAppsApplication) CreateContainerApp(projectID string, containerAppName string, containerAppPort int, containerAppImage string, credentials domain.Credentials) (*domain.ContainerApp, error) {
	// Get access token using KEY_ID and KEY_SECRET
	token, err := c.getAccessToken(credentials.KeyID, credentials.KeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Prepare the request payload
	payload := map[string]interface{}{
		"name":        containerAppName,
		"projectId":   projectID,
		"description": fmt.Sprintf("Container App %s created via MCP", containerAppName),
		"template": map[string]interface{}{
			"containers": []map[string]interface{}{
				{
					"name":          containerAppName,
					"image":         containerAppImage,
					"containerPort": containerAppPort,
					"env": []map[string]string{
						{
							"name":  "CONTAINERAPP_NAME",
							"value": containerAppName,
						},
					},
				},
			},
		},
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Make request to ContainerApps API
	url := "https://containers.api.cloud.ru/v2/containers/"
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonPayload)))
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

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Check if body is empty
	if len(body) == 0 {
		return nil, fmt.Errorf("API returned empty response body with status %d", resp.StatusCode)
	}

	// Parse response
	var containerApp domain.ContainerApp
	if err := json.Unmarshal(body, &containerApp); err != nil {
		return nil, fmt.Errorf("failed to parse containerapp response: %w body: %s", err, body)
	}

	return &containerApp, nil
}

// getAccessToken gets an access token using KEY_ID and KEY_SECRET
func (c *ContainerAppsApplication) getAccessToken(keyID, keySecret string) (string, error) {
	url := "https://iam.api.cloud.ru/api/v1/auth/token"

	payload := strings.NewReader(fmt.Sprintf(`{"keyId": "%s","secret": "%s"}`, keyID, keySecret))

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
		return "", fmt.Errorf("failed to read containerapps response body: %w", err)
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
