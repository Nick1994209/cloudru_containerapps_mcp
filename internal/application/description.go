package application

import (
	"cloudru-containerapps-mcp/internal/config"
	"cloudru-containerapps-mcp/internal/domain"
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
2. cloudru_containerapps_docker_login(registry_name, key_id, key_secret) - Login to Docker registry
3. cloudru_containerapps_docker_push(registry_name, repository_name, image_version, key_id, key_secret) - Build and push Docker image

Environment variables can be used as fallbacks for parameters:
- REGISTRY_NAME: Registry name (e.g., "registry.cloud.ru")
- KEY_ID: Service account key ID for authentication
- KEY_SECRET: Service account key secret for authentication
- REPOSITORY_NAME: Repository name (defaults to current directory name if not set)
- DOCKERFILE: Path to Dockerfile (defaults to "Dockerfile" if not set)

Current configuration values:
- REGISTRY_NAME: (` + cfg.RegistryName + `) (Registry for storing Docker images)
- REPOSITORY_NAME: (` + cfg.RepositoryName + `) (Name of the repository in the registry)
- DOCKERFILE: (` + cfg.Dockerfile + `) (Path to the Dockerfile to build the image, by default Dockerfile)
- KEY_ID: (` + maskSensitiveInfo(cfg.KeyID) + `) (Authentication key identifier)
- KEY_SECRET: (` + maskSensitiveInfo(cfg.KeySecret) + `) (Authentication key secret)
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
