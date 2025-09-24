package main

import (
	"fmt"
	"os"

	"cloudru-containerapps-mcp/internal/config"
	"cloudru-containerapps-mcp/internal/presentation"

	"github.com/mark3labs/mcp-go/application"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	// Create the MCP server
	descriptionService := application.NewDescriptionService()
	dockerService := application.NewDockerService()
	mcpServer := presentation.NewMCPServer(descriptionService, dockerService)

	// Create a mock request with no repository_name parameter
	request := mcp.NewCallToolRequest(map[string]interface{}{
		"registry_name": "test-registry",
		"key_id":        "test-key-id",
		"key_secret":    "test-key-secret",
		"image_version": "v1.0.0",
	})

	// Test getting the repository_name field value
	repoName, err := mcpServer.GetMCPFieldValue("repository_name", request)
	if err != nil {
		fmt.Printf("Error getting repository name: %v\n", err)
		os.Exit(1)
	}

	// Get the current directory name for comparison
	cfg := config.LoadConfig()

	fmt.Printf("Repository name from getMCPFieldValue: '%s'\n", repoName)
	fmt.Printf("Current directory name: '%s'\n", cfg.CurrentDir)
	fmt.Printf("Environment REPOSITORY_NAME: '%s'\n", cfg.RepositoryName)

	if repoName == cfg.CurrentDir {
		fmt.Println("SUCCESS: Repository name correctly defaulted to current directory name")
	} else {
		fmt.Println("FAILURE: Repository name did not default to current directory name")
	}
}
