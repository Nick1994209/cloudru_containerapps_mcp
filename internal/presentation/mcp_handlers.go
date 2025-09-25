package presentation

import (
	"cloudru-containerapps-mcp/internal/config"
	"cloudru-containerapps-mcp/internal/domain"
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// MCPServer holds the application services
type MCPServer struct {
	descriptionService   domain.DescriptionService
	dockerService        domain.DockerService
	containerAppsService domain.ContainerAppsService

	mappedFields map[string]struct {
		envValue     string
		description  string
		defaultValue string
		required     bool
	}
	cfg *config.Config
}

// NewMCPServer creates a new MCP server with the required services
func NewMCPServer(descriptionService domain.DescriptionService, dockerService domain.DockerService, containerAppsService domain.ContainerAppsService) *MCPServer {
	cfg := config.LoadConfig()
	return &MCPServer{
		descriptionService:   descriptionService,
		dockerService:        dockerService,
		containerAppsService: containerAppsService,
		mappedFields: map[string]struct {
			envValue     string
			description  string
			defaultValue string
			required     bool
		}{
			"registry_name": {
				envValue:    cfg.RegistryName,
				description: "Registry name",
				required:    true,
			},
			"key_id": {
				envValue:    cfg.KeyID,
				description: "Service account key ID",
				required:    true,
			},
			"key_secret": {
				envValue:    cfg.KeySecret,
				description: "Service account key secret",
				required:    true,
			},
			"repository_name": {
				envValue:     cfg.RepositoryName,
				description:  "Repository name",
				defaultValue: cfg.CurrentDir,
				required:     true,
			},
			"image_version": {
				description:  "Image version",
				defaultValue: "latest",
				required:     true,
			},
			"dockerfile_path": {
				envValue:     cfg.Dockerfile,
				description:  "Repository name",
				defaultValue: "Dockerfile",
				required:     false,
			},
			"dockerfile_target": {
				envValue:     cfg.DockerfileTarget,
				description:  "Dockerfile target stage",
				defaultValue: "-",
				required:     false,
			},
			"dockerfile_folder": {
				envValue:     cfg.DockerfileFolder,
				description:  "Dockerfile folder (build context)",
				defaultValue: ".",
				required:     false,
			},
			"project_id": {
				envValue:    cfg.ProjectID,
				description: "Project ID for Container Apps (can be set via PROJECT_ID environment variable)",
				required:    true,
			},
		},
	}
}

func (s *MCPServer) getMCPFieldsOptions(description string, fields ...string) []mcp.ToolOption {
	result := []mcp.ToolOption{
		mcp.WithDescription(description),
	}
	for _, field := range fields {
		fieldData := s.mappedFields[field]
		if fieldData.envValue == "" {
			opts := []mcp.PropertyOption{
				mcp.Description(fieldData.description),
				mcp.DefaultString(fieldData.defaultValue),
			}
			if fieldData.required {
				opts = append(opts, mcp.Required())
			}
			result = append(result, mcp.WithString(field, opts...))
		}
	}
	return result
}

func (s *MCPServer) getMCPFieldValue(field string, request mcp.CallToolRequest) (string, error) {
	fieldData := s.mappedFields[field]
	// If we have an environment variable value, use it
	if fieldData.envValue != "" {
		return fieldData.envValue, nil
	}

	// Try to get the value from the request
	result, err := request.RequireString(field)
	if err != nil && fieldData.defaultValue == "" {
		// If there's an error and no default or env value, return the error
		return "", err
	}

	// If we got a value from the request, use it
	if result != "" {
		return result, nil
	}

	// Otherwise, use the default value if available
	if fieldData.defaultValue != "" {
		return fieldData.defaultValue, nil
	}

	// Otherwise, use the default value if available
	if !fieldData.required {
		return "", nil
	}

	// If we get here, return whatever we have (likely empty)
	return "", fmt.Errorf("field %s is empty: %s", field, fieldData.description)
}

// RegisterDescriptionTool registers the description tool with the MCP server
func (s *MCPServer) RegisterDescriptionTool(server *server.MCPServer) {
	descriptionTool := mcp.NewTool("cloudru_containerapps_description",
		mcp.WithDescription("Returns usage instructions for Cloud.ru Container Apps MCP"),
	)

	server.AddTool(descriptionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText(s.descriptionService.GetDescription()), nil
	})
}

// RegisterDockerLoginTool registers the docker login tool with the MCP server
func (s *MCPServer) RegisterDockerLoginTool(server *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions("Login to Cloud.ru Artifact registry (Docker registry)", "registry_name", "key_id", "key_secret")
	dockerLoginTool := mcp.NewTool("cloudru_docker_login", toolOptions...)

	server.AddTool(dockerLoginTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Using helper functions for type-safe argument access
		registryName, err := s.getMCPFieldValue("registry_name", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		keyID, err := s.getMCPFieldValue("key_id", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		keySecret, err := s.getMCPFieldValue("key_secret", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		credentials := domain.Credentials{
			RegistryName: registryName,
			KeyID:        keyID,
			KeySecret:    keySecret,
		}

		err = s.dockerService.Login(credentials)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText("Successfully logged into Cloud.ru Docker registry"), nil
	})
}

// RegisterDockerPushTool registers the docker push tool with the MCP server
func (s *MCPServer) RegisterDockerPushTool(server *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Build and push Docker image to Cloud.ru Artifact Registry (Docker registry)",
		"registry_name",
		"repository_name",
		"image_version",
		"dockerfile_path",
		"dockerfile_target",
		"dockerfile_folder",
		"key_id",
		"key_secret",
	)
	dockerPushTool := mcp.NewTool("cloudru_docker_push", toolOptions...)

	server.AddTool(dockerPushTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		registryName, err := s.getMCPFieldValue("registry_name", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		keyID, err := s.getMCPFieldValue("key_id", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		keySecret, err := s.getMCPFieldValue("key_secret", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		repositoryName, err := s.getMCPFieldValue("repository_name", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		imageVersion, _ := s.getMCPFieldValue("image_version", request)
		dockerfilePath, _ := request.RequireString("dockerfile_path")
		dockerfileTarget, _ := request.RequireString("dockerfile_target")
		dockerfileFolder, _ := request.RequireString("dockerfile_folder")

		image := domain.DockerImage{
			RegistryName:     registryName,
			RepositoryName:   repositoryName,
			ImageVersion:     imageVersion,
			DockerfilePath:   dockerfilePath,
			DockerfileTarget: dockerfileTarget,
			DockerfileFolder: dockerfileFolder,
		}

		credentials := domain.Credentials{
			RegistryName: registryName,
			KeyID:        keyID,
			KeySecret:    keySecret,
		}

		imageTag := fmt.Sprintf("%s.cr.cloud.ru/%s:%s", registryName, repositoryName, imageVersion)
		fmt.Printf("Starting Docker build and push process for image: %s\n", imageTag)
		err = s.dockerService.BuildAndPush(image, credentials)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully built and pushed Docker image to Cloud.ru Artifact Registry: %s", imageTag)), nil
	})
}

// RegisterGetListContainerAppsTool registers the get list container apps tool with the MCP server
func (s *MCPServer) RegisterGetListContainerAppsTool(server *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Get list of Container Apps from Cloud.ru. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru",
		"project_id",
		"key_id",
		"key_secret",
	)
	getListContainerAppsTool := mcp.NewTool("cloudru_get_list_containerapps", toolOptions...)

	server.AddTool(getListContainerAppsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get project ID
		projectID, err := s.getMCPFieldValue("project_id", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get credentials
		keyID, err := s.getMCPFieldValue("key_id", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		keySecret, err := s.getMCPFieldValue("key_secret", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		credentials := domain.Credentials{
			KeyID:     keyID,
			KeySecret: keySecret,
		}

		// Call the service
		containerApps, err := s.containerAppsService.GetListContainerApps(projectID, credentials)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert to JSON for output
		result, err := json.MarshalIndent(containerApps, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	})
}
