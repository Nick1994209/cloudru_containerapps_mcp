package main

import (
	"log"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/application"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/config"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/domain"
)

func main() {
	cfg := config.LoadConfig()

	getListContainerApps(cfg)
	// createContainerApp(
	// 	cfg,
	// 	"testme",
	// 	"quickstart.cr.cloud.ru/react-helloworld@sha256:a1a1e0a11299668c5f05a299f74b3943236ca3390a6fda64e98cc2498064c266",
	// 	8080,
	// )
	// getContainerApp(cfg, "testme")
}

func getListContainerApps(cfg *config.Config) {
	ca := application.NewContainerAppsApplication()

	log.Println("Testing GetListContainerApps...")
	cas, err := ca.GetListContainerApps(
		cfg.ProjectID,
		domain.Credentials{
			KeyID:     cfg.KeyID,
			KeySecret: cfg.KeySecret,
		})
	if err != nil {
		log.Printf("GetListContainerApps error: %v", err)
	} else {
		log.Printf("GetListContainerApps success: found %d container apps", len(cas))
		log.Printf("Container apps: %+v", cas)
	}
}

func getContainerApp(cfg *config.Config, name string) {
	ca := application.NewContainerAppsApplication()

	log.Println("Testing GetContainerApp...")
	cas_, err := ca.GetContainerApp(
		cfg.ProjectID,
		name,
		domain.Credentials{
			KeyID:     cfg.KeyID,
			KeySecret: cfg.KeySecret,
		})
	if err != nil {
		log.Printf("GetContainerApp error: %v", err)
	} else {
		log.Printf("GetContainerApp success: %+v", cas_)
	}
}

func createContainerApp(cfg *config.Config, name, image string, port int) {
	ca := application.NewContainerAppsApplication()

	// Test CreateContainerApp
	containerApp, err := ca.CreateContainerApp(
		cfg.ProjectID,
		name,
		port,
		image,
		domain.Credentials{
			KeyID:     cfg.KeyID,
			KeySecret: cfg.KeySecret,
		})
	log.Print(containerApp)
	if err != nil {
		log.Fatalf("CreateContainerApp error: %v", err.Error())
	}
}
