package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
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
	ContainerAppName string
	CurrentDir       string
}

// EnvVarNames contains the names of environment variables
const (
	EnvRegistryName     = "CLOUDRU_REGISTRY_NAME"
	EnvKeyID            = "CLOUDRU_KEY_ID"
	EnvKeySecret        = "CLOUDRU_KEY_SECRET"
	EnvRepositoryName   = "CLOUDRU_REPOSITORY_NAME"
	EnvProjectID        = "CLOUDRU_PROJECT_ID"
	EnvContainerAppName = "CLOUDRU_CONTAINERAPP_NAME"
	Dockerfile          = "CLOUDRU_DOCKERFILE"
	DockerfileTarget    = "CLOUDRU_DOCKERFILE_TARGET"
	DockerfileFolder    = "CLOUDRU_DOCKERFILE_FOLDER"
)

// LoadConfig loads configuration from environment variables and .env file
func LoadConfig() *Config {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables only")
	}

	// Check for required environment variables
	keyID := os.Getenv(EnvKeyID)
	keySecret := os.Getenv(EnvKeySecret)

	if keyID == "" || keySecret == "" {
		log.Fatal(`CLOUDRU_KEY_ID and CLOUDRU_KEY_SECRET environment variables must be set.
		
To obtain access keys for authentication, please follow the instructions at:
https://cloud.ru/docs/console_api/ug/topics/quickstart

You will need a Key ID and Key Secret to use this service.`)
	}

	dir, err := os.Getwd()
	if err != nil {
		dir = "default"
	}
	projectDirName := filepath.Base(dir)

	return &Config{
		RegistryName:     os.Getenv(EnvRegistryName),
		KeyID:            keyID,
		KeySecret:        keySecret,
		RepositoryName:   os.Getenv(EnvRepositoryName),
		ProjectID:        os.Getenv(EnvProjectID),
		ContainerAppName: os.Getenv(EnvContainerAppName),
		Dockerfile:       os.Getenv(Dockerfile),
		DockerfileTarget: os.Getenv(DockerfileTarget),
		DockerfileFolder: os.Getenv(DockerfileFolder),
		CurrentDir:       projectDirName,
	}
}
