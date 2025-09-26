package main

import (
	"log"

	"github.com/Nick1994209/cloudru_containerapps_mcp/internal/application"
	"github.com/Nick1994209/cloudru_containerapps_mcp/internal/config"
	"github.com/Nick1994209/cloudru_containerapps_mcp/internal/domain"
)

func main() {
	cfg := config.LoadConfig()

	getListDockerRegistries(cfg)
	createDockerRegistry(cfg, "test-registry", false)
}

func getListDockerRegistries(cfg *config.Config) {
	ca := application.NewContainerAppsApplication()

	log.Println("Testing GetListDockerRegistries...")
	registries, err := ca.(domain.DockerRegistryService).GetListDockerRegistries(
		cfg.ProjectID,
		domain.Credentials{
			KeyID:     cfg.KeyID,
			KeySecret: cfg.KeySecret,
		})
	if err != nil {
		log.Printf("GetListDockerRegistries error: %v", err)
	} else {
		log.Printf("GetListDockerRegistries success: found %d registries", len(registries))
		log.Printf("Registries: %+v", registries)
	}
}

func createDockerRegistry(cfg *config.Config, name string, isPublic bool) {
	ca := application.NewContainerAppsApplication()

	log.Printf("Testing CreateDockerRegistry with name: %s, isPublic: %v...", name, isPublic)
	registry, err := ca.(domain.DockerRegistryService).CreateDockerRegistry(
		cfg.ProjectID,
		name,
		isPublic,
		domain.Credentials{
			KeyID:     cfg.KeyID,
			KeySecret: cfg.KeySecret,
		})
	if err != nil {
		log.Printf("CreateDockerRegistry error: %v", err)
	} else {
		log.Printf("CreateDockerRegistry success: %+v", registry)
	}
}
