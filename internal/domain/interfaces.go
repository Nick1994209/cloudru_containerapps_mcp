package domain

// DescriptionService provides usage instructions for the MCP
type DescriptionService interface {
	GetDescription() string
}

// DockerService handles Docker operations
type DockerService interface {
	Login(registryName string, credentials Credentials) (string, error)
	BuildAndPush(image DockerImage, credentials Credentials) (string, error)
}

// ContainerAppsService handles Cloud.ru Container Apps API operations
type ContainerAppsService interface {
	GetListContainerApps(projectID string, credentials Credentials) ([]ContainerApp, error)
	GetContainerApp(projectID string, containerAppName string, credentials Credentials) (*ContainerApp, error)
	CreateContainerApp(projectID string, containerAppName string, containerAppPort int, containerAppImage string, credentials Credentials) (*ContainerApp, error)
	DeleteContainerApp(projectID string, containerAppName string, credentials Credentials) error
	StartContainerApp(projectID string, containerAppName string, credentials Credentials) error
	StopContainerApp(projectID string, containerAppName string, credentials Credentials) error
}
