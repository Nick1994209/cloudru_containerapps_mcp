package config

import (
	"os"
	"path/filepath"
)

// Config holds the configuration for the MCP
type Config struct {
	RegistryName     string
	KeyID            string
	KeySecret        string
	RepositoryName   string
	Dockerfile       string
	DockerfileTarget string
	DockerfileFolder string
	ProjectID        string
	CurrentDir       string
}

// EnvVarNames contains the names of environment variables
const (
	EnvRegistryName   = "REGISTRY_NAME"
	EnvKeyID          = "KEY_ID"
	EnvKeySecret      = "KEY_SECRET"
	EnvRepositoryName = "REPOSITORY_NAME"
	EnvProjectID      = "PROJECT_ID"
	Dockerfile        = "DOCKERFILE"
	DockerfileTarget  = "DOCKERFILE_TARGET"
	DockerfileFolder  = "DOCKERFILE_FOLDER"
)

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	dir, err := os.Getwd()
	if err != nil {
		dir = "default"
	}
	projectDirName := filepath.Base(dir)

	return &Config{
		RegistryName:     os.Getenv(EnvRegistryName),
		KeyID:            os.Getenv(EnvKeyID),
		KeySecret:        os.Getenv(EnvKeySecret),
		RepositoryName:   os.Getenv(EnvRepositoryName),
		ProjectID:        os.Getenv(EnvProjectID),
		Dockerfile:       os.Getenv(Dockerfile),
		DockerfileTarget: os.Getenv(DockerfileTarget),
		DockerfileFolder: os.Getenv(DockerfileFolder),
		CurrentDir:       projectDirName,
	}
}
