package application

import (
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/config"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/domain"
)

// DescriptionApplication implements the DescriptionService interface
type DescriptionApplication struct{}

// NewDescriptionApplication creates a new DescriptionApplication
func NewDescriptionApplication() domain.DescriptionService {
	return &DescriptionApplication{}
}

// GetDescription returns usage instructions for this MCP
func (d *DescriptionApplication) GetDescription() string {
	cfg := config.LoadConfig()

	return `Cloud.ru Container Apps MCP provides functions to interact with Cloud.ru Artifact Registry:

1. cloudru_containerapps_description() - Returns usage instructions for this MCP
2. cloudru_get_list_docker_registries(project_id, key_id, key_secret) - Get list of Docker Registries
3. cloudru_create_docker_registry(project_id, registry_name, is_public, key_id, key_secret) - Create a new Docker Registry
4. cloudru_docker_login(registry_name, key_id, key_secret) - Login to Docker registry
5. cloudru_docker_push(registry_name, repository_name, image_version, key_id, key_secret) - Build and push Docker image
6. cloudru_get_list_containerapps(project_id, key_id, key_secret) - Get list of Container Apps
7. cloudru_get_containerapp(project_id, containerapp_name, key_id, key_secret) - Get a specific Container App by name
8. cloudru_create_containerapp(project_id, containerapp_name, containerapp_port, containerapp_image, key_id, key_secret) - Create a new Container App
9. cloudru_delete_containerapp(project_id, containerapp_name, key_id, key_secret) - Delete a Container App (WARNING: This action cannot be undone!)
10. cloudru_start_containerapp(project_id, containerapp_name, key_id, key_secret) - Start a Container App
11. cloudru_stop_containerapp(project_id, containerapp_name, key_id, key_secret) - Stop a Container App

Environment variables can be used as fallbacks for parameters:

**Required environment variables:**
- CLOUDRU_KEY_ID: Service account key ID for authentication (required)
- CLOUDRU_KEY_SECRET: Service account key secret for authentication (required)

To obtain access keys for authentication, please follow the instructions at:
https://cloud.ru/docs/console_api/ug/topics/quickstart

You will need a Key ID and Key Secret to use this service.

**Optional environment variables:**
- CLOUDRU_REGISTRY_NAME: Registry name (e.g., "registry.cloud.ru")
- CLOUDRU_PROJECT_ID: Project ID for Container Apps (can be obtained from console.cloud.ru)
- CLOUDRU_CONTAINERAPP_NAME: Container App name (optional)
- CLOUDRU_REPOSITORY_NAME: Repository name (defaults to current directory name if not set)
- CLOUDRU_PROJECT_ID: Project ID for Container Apps (can be obtained from console.cloud.ru)
- CLOUDRU_CONTAINERAPP_NAME: Container App name (optional)
- CLOUDRU_DOCKERFILE: Path to Dockerfile (defaults to "Dockerfile" if not set)

Current configuration values:
- CLOUDRU_REGISTRY_NAME: (` + cfg.RegistryName + `) (Registry for storing Docker images)
- CLOUDRU_REPOSITORY_NAME: (` + cfg.RepositoryName + `) (Name of the repository in the registry)
- CLOUDRU_PROJECT_ID: (` + cfg.ProjectID + `) (Project ID for Container Apps)
- CLOUDRU_DOCKERFILE: (` + cfg.Dockerfile + `) (Path to the Dockerfile to build the image, by default Dockerfile)
- CLOUDRU_KEY_ID: (` + maskSensitiveInfo(cfg.KeyID) + `) (Authentication key identifier)
- CLOUDRU_KEY_SECRET: (` + maskSensitiveInfo(cfg.KeySecret) + `) (Authentication key secret)
- Current directory: ` + cfg.CurrentDir + ` (Name of the current working directory)

For more details see: https://cloud.ru/docs/container-apps-evolution/ug/topics/tutorials__before-work`
}

// maskSensitiveInfo replaces the middle of a string with asterisks for sensitive data
func maskSensitiveInfo(value string) string {
	if len(value) == 0 {
		return ""
	}

	if len(value) <= 4 {
		return "***"
	}

	// Show first 2 and last 2 characters, replace the rest with asterisks
	start := value[:3]
	end := value[len(value)-3:]
	middle := "***"

	return start + middle + end
}
