package domain

// DescriptionService provides usage instructions for the MCP
type DescriptionService interface {
	GetDescription() string
}

// DockerService handles Docker operations
type DockerService interface {
	Login(credentials Credentials) error
	BuildAndPush(image DockerImage, credentials Credentials) error
}
