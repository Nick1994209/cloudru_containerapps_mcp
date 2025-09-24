package application

import (
	"cloudru-containerapps-mcp/internal/domain"
	"fmt"
	"os/exec"
)

// DockerApplication implements the DockerService interface using actual Docker commands
type DockerApplication struct{}

// NewDockerApplication creates a new DockerApplication
func NewDockerApplication() domain.DockerService {
	return &DockerApplication{}
}

// Login logs into the Cloud.ru Docker registry using Docker CLI
func (d *DockerApplication) Login(credentials domain.Credentials) error {
	loginTarget := fmt.Sprintf("%s.cr.cloud.ru", credentials.RegistryName)
	cmd := exec.Command("docker", "login", loginTarget, "-u", credentials.KeyID, "-p", credentials.KeySecret)

	// Execute the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker login failed: %w\nOutput: %s\n\nPlease ensure:\n1. The registry exists in Cloud.ru Evolution Artifact Registry\n2. You have created a registry and obtained access keys\n3. See documentation: https://cloud.ru/docs/container-apps-evolution/ug/topics/tutorials__before-work", err, string(output))
	}

	return nil
}

// BuildAndPush builds and pushes a Docker image to Cloud.ru Artifact Registry
func (d *DockerApplication) BuildAndPush(image domain.DockerImage, credentials domain.Credentials) error {
	imageTag := fmt.Sprintf("%s.cr.cloud.ru/%s:%s", credentials.RegistryName, image.RepositoryName, image.ImageVersion)

	// Build the Docker image
	buildCmd := exec.Command("docker", "build", "--platform", "linux/amd64", "-t", imageTag, "-f", image.DockerfilePath, ".")
	buildOutput, buildErr := buildCmd.CombinedOutput()
	if buildErr != nil {
		return fmt.Errorf("failed to build Docker image: %w\nOutput: %s", buildErr, string(buildOutput))
	}

	// Push the Docker image
	pushCmd := exec.Command("docker", "push", imageTag)
	pushOutput, pushErr := pushCmd.CombinedOutput()

	if pushErr != nil {
		// If push fails, try to re-login if credentials are available
		if (credentials.KeyID != "") || (credentials.KeySecret != "") {
			loginErr := d.Login(credentials)
			if loginErr != nil {
				return fmt.Errorf("docker push failed and re-login unsuccessful: %w\nOutput: %s\n\nTo resolve this issue:\n1. Set KEY_ID and KEY_SECRET environment variables\n2. Or run the cloudru_containerapps_docker_login function\n3. See documentation: https://cloud.ru/docs/container-apps-evolution/ug/topics/tutorials__before-work", loginErr, string(pushOutput))
			}

			// Retry push after re-login
			pushCmd = exec.Command("docker", "push", imageTag)
			pushOutput, pushErr = pushCmd.CombinedOutput()
			if pushErr != nil {
				return fmt.Errorf("docker push still failed after re-login: %w\nOutput: %s", pushErr, string(pushOutput))
			}
		} else {
			return fmt.Errorf("docker push failed: %w\nOutput: %s\n\nTo resolve this issue:\n1. Set KEY_ID and KEY_SECRET environment variables\n2. Or run the cloudru_containerapps_docker_login function\n3. See documentation: https://cloud.ru/docs/container-apps-evolution/ug/topics/tutorials__before-work", pushErr, string(pushOutput))
		}
	}

	return nil
}
