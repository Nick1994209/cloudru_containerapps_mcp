package main

import (
	"fmt"

	"cloudru-containerapps-mcp/internal/application"
)

func main() {
	// Create the services
	descriptionService := application.NewDescriptionApplication()
	_ = application.NewDockerApplication()

	// Test that we can create the services without errors
	fmt.Println("SUCCESS: Created DescriptionApplication and DockerApplication services")

	// Test the description service
	description := descriptionService.GetDescription()
	if len(description) > 0 {
		fmt.Println("SUCCESS: DescriptionService returned a non-empty description")
	} else {
		fmt.Println("FAILURE: DescriptionService returned an empty description")
	}
}
