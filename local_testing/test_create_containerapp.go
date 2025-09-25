package main

import (
	"cloudru-containerapps-mcp/internal/application"
	"cloudru-containerapps-mcp/internal/domain"
	"log"
)

func main() {
	ca := application.NewContainerAppsApplication()

	// Test CreateContainerApp
	containerApp, err := ca.CreateContainerApp(
		"a9e46dcd-b00a-4a87****028b31931c7b",
		"test-container-app",
		8080,
		"nginx:latest",
		domain.Credentials{
			KeyID:     "****",
			KeySecret: "****",
		})
	log.Print(containerApp)
	if err != nil {
		log.Fatalf("CreateContainerApp error: %v", err.Error())
	}
}
