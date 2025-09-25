package presentation

import (
	"cloudru-containerapps-mcp/internal/config"
	"cloudru-containerapps-mcp/internal/domain"
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// MCPServer holds the application services
type MCPServer struct {
	descriptionService domain.DescriptionService
	dockerService      domain.DockerService

	mappedFields map[string]struct {
		envValue     string
		description  string
		defaultValue string
	}
	cfg *config.Config
}

// NewMCPServer creates a new MCP server with the required services
func NewMCPServer(descriptionService domain.DescriptionService, dockerService domain.DockerService) *MCPServer {
	cfg := config.LoadConfig()
	return &MCPServer{
		descriptionService: descriptionService,
		dockerService:      dockerService,
		mappedFields: map[string]struct {
			envValue     string
			description  string
			defaultValue string
		}{
			"registry_name": {
				envValue:    cfg.RegistryName,
				description: "Registry name",
			},
			"key_id": {
				envValue:    cfg.KeyID,
				description: "Service account key ID",
			},
			"key_secret": {
				envValue:    cfg.KeySecret,
				description: "Service account key secret",
			},
			"repository_name": {
				envValue:     cfg.RepositoryName,
				description:  "Repository name",
				defaultValue: cfg.CurrentDir,
			},
			"image_version": {
				description:  "Image version",
				defaultValue: "latest",
			},
			"dockerfile_path": {
				envValue:     cfg.Dockerfile,
				description:  "Repository name",
				defaultValue: "Dockerfile",
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
			result = append(result, mcp.WithString(
				field,
				mcp.Description(fieldData.description),
				mcp.Required(),
				mcp.DefaultString(fieldData.defaultValue),
			))
		}
	}
	return result
}

func (s *MCPServer) getMCPFieldValue(field string, request mcp.CallToolRequest) (string, error) {
	fieldData := s.mappedFields[field]

	// Try to get the value from the request
	result, err := request.RequireString(field)
	if err != nil && fieldData.defaultValue == "" && fieldData.envValue == "" {
		// If there's an error and no default or env value, return the error
		return "", err
	}

	// If we got a value from the request, use it
	if result != "" {
		return result, nil
	}

	// If we have an environment variable value, use it
	if fieldData.envValue != "" {
		return fieldData.envValue, nil
	}

	// Otherwise, use the default value if available
	if fieldData.defaultValue != "" {
		return fieldData.defaultValue, nil
	}

	// If we get here, return whatever we have (likely empty)
	return result, nil
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
	toolOptions := s.getMCPFieldsOptions("Login to Cloud.ru Docker registry", "registry_name", "key_id", "key_secret")
	dockerLoginTool := mcp.NewTool("cloudru_containerapps_docker_login", toolOptions...)

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
		"Build and push Docker image to Cloud.ru Artifact Registry",
		"registry_name",
		"repository_name",
		"image_version",
		"dockerfile_path",
		"key_id",
		"key_secret",
	)
	dockerPushTool := mcp.NewTool("cloudru_containerapps_docker_push", toolOptions...)

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

		image := domain.DockerImage{
			RegistryName:   registryName,
			RepositoryName: repositoryName,
			ImageVersion:   imageVersion,
			DockerfilePath: dockerfilePath,
		}

		credentials := domain.Credentials{
			RegistryName: registryName,
			KeyID:        keyID,
			KeySecret:    keySecret,
		}

		err = s.dockerService.BuildAndPush(image, credentials)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText("Successfully built and pushed Docker image to Cloud.ru Artifact Registry"), nil
	})
}
