package main

import (
	"fmt"
	"log"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/application"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/domain"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/presentation"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create infrastructure layer
	dockerInfrastructure := application.NewDockerApplication()
	containerAppsService := application.NewContainerAppsApplication()

	// Create application layer
	descriptionService := application.NewDescriptionApplication()

	// Log the application description
	log.Println("Application Description:")
	log.Println(descriptionService.GetDescription())

	// Create presentation layer
	mcpServer := presentation.NewMCPServer(descriptionService, dockerInfrastructure, containerAppsService, containerAppsService.(domain.DockerRegistryService))

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
	mcpServer.RegisterGetListDockerRegistriesTool(s)
	mcpServer.RegisterCreateDockerRegistryTool(s)

	// Start the server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
