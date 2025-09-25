package main

import (
	"cloudru-containerapps-mcp/internal/application"
	"cloudru-containerapps-mcp/internal/domain"
	"log"
)

func main() {
	ca := application.NewContainerAppsApplication()
	cas, err := ca.GetListContainerApps("a9e46dcd-b00a-4a87-8ec2-******", domain.Credentials{
		KeyID:     "******",
		KeySecret: "******",
	})
	log.Print(cas)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
