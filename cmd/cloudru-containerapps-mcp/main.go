package main

import (
	"fmt"

	"github.com/Nick1994209/cloudru_containerapps_mcp/internal/application"
	"github.com/Nick1994209/cloudru_containerapps_mcp/internal/presentation"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create infrastructure layer
	dockerInfrastructure := application.NewDockerApplication()
	containerAppsService := application.NewContainerAppsApplication()

	// Create application layer
	descriptionService := application.NewDescriptionApplication()

	// Create presentation layer
	mcpServer := presentation.NewMCPServer(descriptionService, dockerInfrastructure, containerAppsService)

	// Create a new MCP server
	s := server.NewMCPServer(
		"Cloud.ru Container Apps MCP",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithRecovery(),
	)

	// Register tools with the MCP server
	mcpServer.RegisterDescriptionTool(s)
	mcpServer.RegisterDockerLoginTool(s)
	mcpServer.RegisterDockerPushTool(s)
	mcpServer.RegisterGetListContainerAppsTool(s)
	mcpServer.RegisterGetContainerAppTool(s)
	mcpServer.RegisterCreateContainerAppTool(s)
	mcpServer.RegisterDeleteContainerAppTool(s)
	mcpServer.RegisterStartContainerAppTool(s)
	mcpServer.RegisterStopContainerAppTool(s)

	// Start the server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
