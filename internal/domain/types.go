package domain

// Credentials represents the authentication credentials for Cloud.ru
type Credentials struct {
	RegistryName string
	KeyID        string
	KeySecret    string
}

// DockerImage represents a Docker image to be built and pushed
type DockerImage struct {
	RegistryName   string
	RepositoryName string
	ImageVersion   string
	DockerfilePath string
}
