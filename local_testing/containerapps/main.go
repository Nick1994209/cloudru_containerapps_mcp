package main

import (
	"log"

	"github.com/Nick1994209/cloudru_containerapps_mcp/internal/application"
	"github.com/Nick1994209/cloudru_containerapps_mcp/internal/domain"
)

const project string = "****"
const key_id string = "****"
const key_secret string = "****"

func main() {
	getListContainerApps()
	createContainerApp("testme", "image.cr.cloud.ru/image", 8080)
	getContainerApp("testme")
}

func getListContainerApps() {
	ca := application.NewContainerAppsApplication()

	log.Println("Testing GetListContainerApps...")
	cas, err := ca.GetListContainerApps(
		project,
		domain.Credentials{
			KeyID:     key_id,
			KeySecret: key_secret,
		})
	if err != nil {
		log.Printf("GetListContainerApps error: %v", err)
	} else {
		log.Printf("GetListContainerApps success: found %d container apps", len(cas))
		log.Printf("Container apps: %+v", cas)
	}
}

func getContainerApp(name string) {
	ca := application.NewContainerAppsApplication()

	log.Println("Testing GetContainerApp...")
	cas_, err := ca.GetContainerApp(
		project,
		name,
		domain.Credentials{
			KeyID:     key_id,
			KeySecret: key_secret,
		})
	if err != nil {
		log.Printf("GetContainerApp error: %v", err)
	} else {
		log.Printf("GetContainerApp success: %+v", cas_)
	}
}

func createContainerApp(name, image string, port int) {
	ca := application.NewContainerAppsApplication()

	// Test CreateContainerApp
	containerApp, err := ca.CreateContainerApp(
		project,
		name,
		port,
		image,
		domain.Credentials{
			KeyID:     key_id,
			KeySecret: key_secret,
		})
	log.Print(containerApp)
	if err != nil {
		log.Fatalf("CreateContainerApp error: %v", err.Error())
	}
}
