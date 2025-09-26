package presentation

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Nick1994209/cloudru_containerapps_mcp/internal/config"
	"github.com/Nick1994209/cloudru_containerapps_mcp/internal/domain"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// MCPServer holds the application services
type MCPServer struct {
	descriptionService    domain.DescriptionService
	dockerService         domain.DockerService
	containerAppsService  domain.ContainerAppsService
	dockerRegistryService domain.DockerRegistryService

	mappedFields map[string]struct {
		envValue     string
		description  string
		defaultValue string
		title        string
		required     bool
	}
	cfg *config.Config
}

// NewMCPServer creates a new MCP server with the required services
func NewMCPServer(descriptionService domain.DescriptionService, dockerService domain.DockerService, containerAppsService domain.ContainerAppsService, dockerRegistryService domain.DockerRegistryService) *MCPServer {
	cfg := config.LoadConfig()

	defaultRepoName := cfg.CurrentDir
	if cfg.DockerfileTarget != "" && cfg.DockerfileTarget != "-" {
		defaultRepoName = defaultRepoName + "-" + cfg.DockerfileTarget
	}

	containerappImage := fmt.Sprintf("%s.cr.cloud.ru/%s:%s", cfg.RegistryName, cfg.RepositoryName, "latest")
	return &MCPServer{
		descriptionService:    descriptionService,
		dockerService:         dockerService,
		containerAppsService:  containerAppsService,
		dockerRegistryService: dockerRegistryService,

		mappedFields: map[string]struct {
			envValue     string
			description  string
			defaultValue string
			title        string
			required     bool
		}{
			"registry_name": {
				envValue:    cfg.RegistryName,
				description: "Registry name",
				required:    true,
			},
			"repository_name": {
				envValue:     cfg.RepositoryName,
				description:  "Repository name",
				defaultValue: defaultRepoName,
				required:     true,
			},
			"image_version": {
				description: "Image version",
				title:       "For example: latest or v0.0.1",
				required:    true,
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
			"containerapp_name": {
				envValue:     cfg.ContainerAppName,
				description:  "Container App name (can be set via CONTAINERAPP_NAME environment variable)",
				required:     false,
				defaultValue: cfg.CurrentDir,
				title:        "You can use example: " + cfg.CurrentDir,
			},
			"containerapp_port": {
				description: "Container App port number",
				required:    true,
				title:       "You can use example: 8000",
			},
			"containerapp_image": {
				description: "Container App image",
				required:    true,
				title:       "Example image: " + containerappImage,
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
			}
			if fieldData.required {
				opts = append(opts, mcp.Required())
			}
			if fieldData.title != "" {
				opts = append(opts, mcp.Title(fieldData.title))
			}
			if fieldData.defaultValue != "" {
				opts = append(opts, mcp.DefaultString(fieldData.defaultValue))
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
	toolOptions := s.getMCPFieldsOptions("Login to Cloud.ru Artifact registry (Docker registry)", "registry_name")
	dockerLoginTool := mcp.NewTool("cloudru_docker_login", toolOptions...)

	server.AddTool(dockerLoginTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Using helper functions for type-safe argument access
		registryName, err := s.getMCPFieldValue("registry_name", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		credentials := domain.Credentials{
			KeyID:     s.cfg.KeyID,
			KeySecret: s.cfg.KeySecret,
		}

		result, err := s.dockerService.Login(registryName, credentials)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Successfully login to Cloud.ru Artifact Registry: %s", result)), nil
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
	)
	dockerPushTool := mcp.NewTool("cloudru_docker_push", toolOptions...)

	server.AddTool(dockerPushTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		registryName, err := s.getMCPFieldValue("registry_name", request)
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
			KeyID:     s.cfg.KeyID,
			KeySecret: s.cfg.KeySecret,
		}

		imageTag := fmt.Sprintf("%s.cr.cloud.ru/%s:%s", registryName, repositoryName, imageVersion)
		fmt.Printf("Starting Docker build and push process for image: %s\n", imageTag)
		result, err := s.dockerService.BuildAndPush(image, credentials)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully built and pushed Docker image to Cloud.ru Artifact Registry: %s", result)), nil
	})
}

// RegisterGetListContainerAppsTool registers the get list container apps tool with the MCP server
func (s *MCPServer) RegisterGetListContainerAppsTool(server *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Get list of Container Apps from Cloud.ru. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru",
		"project_id",
	)
	getListContainerAppsTool := mcp.NewTool("cloudru_get_list_containerapps", toolOptions...)

	server.AddTool(getListContainerAppsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get project ID
		projectID, err := s.getMCPFieldValue("project_id", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		credentials := domain.Credentials{
			KeyID:     s.cfg.KeyID,
			KeySecret: s.cfg.KeySecret,
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

// RegisterGetContainerAppTool registers the get container app tool with the MCP server
func (s *MCPServer) RegisterGetContainerAppTool(server *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Get a specific Container App from Cloud.ru by name. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru",
		"project_id",
		"containerapp_name",
	)
	getContainerAppTool := mcp.NewTool("cloudru_get_containerapp", toolOptions...)

	server.AddTool(getContainerAppTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get project ID
		projectID, err := s.getMCPFieldValue("project_id", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get container app name
		containerAppName, err := s.getMCPFieldValue("containerapp_name", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		credentials := domain.Credentials{
			KeyID:     s.cfg.KeyID,
			KeySecret: s.cfg.KeySecret,
		}

		// Call the service
		containerApp, err := s.containerAppsService.GetContainerApp(projectID, containerAppName, credentials)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert to JSON for output
		result, err := json.MarshalIndent(containerApp, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	})
}

// RegisterCreateContainerAppTool registers the create container app tool with the MCP server
func (s *MCPServer) RegisterCreateContainerAppTool(server *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Create a new Container App in Cloud.ru",
		"project_id",
		"containerapp_name",
		"containerapp_port",
		"containerapp_image",
	)
	createContainerAppTool := mcp.NewTool("cloudru_create_containerapp", toolOptions...)

	server.AddTool(createContainerAppTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get project ID
		projectID, err := s.getMCPFieldValue("project_id", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get container app name
		containerAppName, err := s.getMCPFieldValue("containerapp_name", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get container app port
		containerAppPortStr, err := s.getMCPFieldValue("containerapp_port", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert port to integer
		var containerAppPort int
		fmt.Sscanf(containerAppPortStr, "%d", &containerAppPort)

		// Get container app image
		containerAppImage, err := s.getMCPFieldValue("containerapp_image", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		credentials := domain.Credentials{
			KeyID:     s.cfg.KeyID,
			KeySecret: s.cfg.KeySecret,
		}

		// Call the service
		containerApp, err := s.containerAppsService.CreateContainerApp(projectID, containerAppName, containerAppPort, containerAppImage, credentials)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert to JSON for output
		result, err := json.MarshalIndent(containerApp, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully created Container App: %s\n%s", containerAppName, string(result))), nil
	})
}

// RegisterDeleteContainerAppTool registers the delete container app tool with the MCP server
func (s *MCPServer) RegisterDeleteContainerAppTool(server *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Delete a Container App from Cloud.ru. WARNING: This action cannot be undone!",
		"project_id",
		"containerapp_name",
	)
	deleteContainerAppTool := mcp.NewTool("cloudru_delete_containerapp", toolOptions...)

	server.AddTool(deleteContainerAppTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get project ID
		projectID, err := s.getMCPFieldValue("project_id", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get container app name
		containerAppName, err := s.getMCPFieldValue("containerapp_name", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		credentials := domain.Credentials{
			KeyID:     s.cfg.KeyID,
			KeySecret: s.cfg.KeySecret,
		}

		// Confirmation prompt - in MCP context, we'll add a warning in the description
		// but the actual confirmation would typically happen in the client UI

		// Call the service
		err = s.containerAppsService.DeleteContainerApp(projectID, containerAppName, credentials)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully deleted Container App: %s", containerAppName)), nil
	})
}

// RegisterStartContainerAppTool registers the start container app tool with the MCP server
func (s *MCPServer) RegisterStartContainerAppTool(server *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Start a Container App in Cloud.ru",
		"project_id",
		"containerapp_name",
	)
	startContainerAppTool := mcp.NewTool("cloudru_start_containerapp", toolOptions...)

	server.AddTool(startContainerAppTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get project ID
		projectID, err := s.getMCPFieldValue("project_id", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get container app name
		containerAppName, err := s.getMCPFieldValue("containerapp_name", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		credentials := domain.Credentials{
			KeyID:     s.cfg.KeyID,
			KeySecret: s.cfg.KeySecret,
		}

		// Call the service
		err = s.containerAppsService.StartContainerApp(projectID, containerAppName, credentials)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully started Container App: %s", containerAppName)), nil
	})
}

// RegisterStopContainerAppTool registers the stop container app tool with the MCP server
func (s *MCPServer) RegisterStopContainerAppTool(server *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Stop a Container App in Cloud.ru",
		"project_id",
		"containerapp_name",
	)
	stopContainerAppTool := mcp.NewTool("cloudru_stop_containerapp", toolOptions...)

	server.AddTool(stopContainerAppTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get project ID
		projectID, err := s.getMCPFieldValue("project_id", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get container app name
		containerAppName, err := s.getMCPFieldValue("containerapp_name", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		credentials := domain.Credentials{
			KeyID:     s.cfg.KeyID,
			KeySecret: s.cfg.KeySecret,
		}

		// Call the service
		err = s.containerAppsService.StopContainerApp(projectID, containerAppName, credentials)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully stopped Container App: %s", containerAppName)), nil
	})
}

// RegisterGetListDockerRegistriesTool registers the get list docker registries tool with the MCP server
func (s *MCPServer) RegisterGetListDockerRegistriesTool(server *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Get list of Docker Registries from Cloud.ru. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru",
		"project_id",
	)
	getListDockerRegistriesTool := mcp.NewTool("cloudru_get_list_docker_registries", toolOptions...)

	server.AddTool(getListDockerRegistriesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get project ID
		projectID, err := s.getMCPFieldValue("project_id", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		credentials := domain.Credentials{
			KeyID:     s.cfg.KeyID,
			KeySecret: s.cfg.KeySecret,
		}

		// Call the service
		dockerRegistries, err := s.dockerRegistryService.GetListDockerRegistries(projectID, credentials)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert to JSON for output
		result, err := json.MarshalIndent(dockerRegistries, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	})
}

// RegisterCreateDockerRegistryTool registers the create docker registry tool with the MCP server
func (s *MCPServer) RegisterCreateDockerRegistryTool(server *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Create a new Docker Registry in Cloud.ru",
		"project_id",
		"registry_name",
		"is_public",
	)
	createDockerRegistryTool := mcp.NewTool("cloudru_create_docker_registry", toolOptions...)

	server.AddTool(createDockerRegistryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get project ID
		projectID, err := s.getMCPFieldValue("project_id", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get registry name
		registryName, err := s.getMCPFieldValue("registry_name", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get is_public flag
		isPublicStr, err := s.getMCPFieldValue("is_public", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert is_public to boolean
		var isPublic bool
		if isPublicStr == "true" || isPublicStr == "1" {
			isPublic = true
		} else if isPublicStr == "false" || isPublicStr == "0" {
			isPublic = false
		} else {
			return mcp.NewToolResultError("is_public must be 'true' or 'false'"), nil
		}

		credentials := domain.Credentials{
			KeyID:     s.cfg.KeyID,
			KeySecret: s.cfg.KeySecret,
		}

		// Call the service
		dockerRegistry, err := s.dockerRegistryService.CreateDockerRegistry(projectID, registryName, isPublic, credentials)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert to JSON for output
		result, err := json.MarshalIndent(dockerRegistry, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully created Docker Registry: %s\n%s", registryName, string(result))), nil
	})
}
