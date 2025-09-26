# Cloud.ru Container Apps MCP

A Model Context Protocol (MCP) server for interacting with Cloud.ru Container Apps and Artifact Registry.

## What is MCP?

Model Context Protocol (MCP) is an open standard that enables seamless integration between AI assistants and external tools or data sources. It allows AI models to interact with your applications, services, and data in a secure and controlled manner.

With MCP, you can:
- Extend the capabilities of AI assistants beyond their training data
- Enable real-time interaction with live systems and APIs
- Provide contextually relevant information from your own data sources
- Execute complex workflows without leaving your AI assistant interface

### Example Usage

Instead of manually performing tasks, you can simply ask your AI assistant:

```
"Deploy my latest application changes to Cloud.ru Container Apps"
```

Your AI assistant, using this MCP server, can then:
1. Build a Docker image of your application
2. Push it to Cloud.ru Artifact Registry
3. Update your Container App with the new image
4. Report back the status of the deployment

All of this happens automatically through natural language commands, making complex DevOps tasks accessible to everyone.

## Features

This MCP provides the following functions:

1. `cloudru_containerapps_description()` - Returns usage instructions for this MCP
2. `cloudru_docker_login(registry_name)` - Login to Cloud.ru Docker registry
3. `cloudru_docker_push(registry_name, repository_name, image_version, dockerfile_path, dockerfile_target, dockerfile_folder)` - Build and push Docker image to Cloud.ru Artifact Registry
4. `cloudru_get_list_containerapps(project_id)` - Get list of Container Apps from Cloud.ru. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru
5. `cloudru_get_containerapp(project_id, containerapp_name)` - Get a specific Container App from Cloud.ru by name. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru
6. `cloudru_create_containerapp(project_id, containerapp_name, containerapp_port, containerapp_image)` - Create a new Container App in Cloud.ru
7. `cloudru_delete_containerapp(project_id, containerapp_name)` - Delete a Container App from Cloud.ru. WARNING: This action cannot be undone!
8. `cloudru_start_containerapp(project_id, containerapp_name)` - Start a Container App in Cloud.ru
9. `cloudru_stop_containerapp(project_id, containerapp_name)` - Stop a Container App in Cloud.ru
10. `cloudru_get_list_docker_registries(project_id)` - Get list of Docker Registries from Cloud.ru. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru
11. `cloudru_create_docker_registry(project_id, registry_name, is_public)` - Create a new Docker Registry in Cloud.ru

## Installation cloudru-containerapps-mcp to your system
[docs/INSTALLATION.md](docs/INSTALLATION.md)

## Add cloudru-containerapps-mcp to your IDE. For example VisualStudioCode or Cursor
[docs/HOW_ADD_TO_IDE.md](docs/HOW_ADD_TO_IDE.md)

## MCP Environment variables
[docs/ENVIRONMENT_VARIABLES.md](docs/ENVIRONMENT_VARIABLES.md)

### Functions

#### cloudru_containerapps_description()

Returns usage instructions for this MCP.

#### cloudru_docker_login(registry_name)

Logs into the Cloud.ru Docker registry using the provided credentials.

Parameters:
- `registry_name`: Name of the registry (falls back to CLOUDRU_REGISTRY_NAME env var)

If login fails, you'll need to:
1. Go to Cloud.ru Evolution Artifact Registry
2. Create a registry
3. Obtain access keys
4. See documentation: https://cloud.ru/docs/container-apps-evolution/ug/topics/tutorials__before-work

#### cloudru_docker_push(registry_name, repository_name, image_version, dockerfile_path, dockerfile_target, dockerfile_folder)

Builds a Docker image and pushes it to Cloud.ru Artifact Registry.

Parameters:
- `registry_name`: Name of the registry (falls back to CLOUDRU_REGISTRY_NAME env var)
- `repository_name`: Name of the repository (falls back to CLOUDRU_REPOSITORY_NAME env var, then to current directory name)
- `image_version`: Version/tag for the image
- `dockerfile_path`: Path to Dockerfile (optional, defaults to 'Dockerfile')
- `dockerfile_target`: Target stage in a multi-stage Dockerfile (optional, defaults to '-' which means no target)
- `dockerfile_folder`: Dockerfile folder (build context, defaults to '.' which means current directory)

If Docker push fails due to authentication issues and CLOUDRU_KEY_ID/CLOUDRU_KEY_SECRET environment variables are set, the function will attempt to re-login and retry the push operation.

## Running the MCP Server

To start the MCP server, simply run:

#### cloudru_get_list_containerapps(project_id)

Gets a list of Container Apps from Cloud.ru. Project ID can be set via CLOUDRU_PROJECT_ID environment variable and obtained from console.cloud.ru.

Parameters:
- `project_id`: Project ID in Cloud.ru (falls back to CLOUDRU_PROJECT_ID env var)

#### cloudru_get_containerapp(project_id, containerapp_name)

Gets a specific Container App from Cloud.ru by name. Project ID can be set via CLOUDRU_PROJECT_ID environment variable and obtained from console.cloud.ru.

Parameters:
- `project_id`: Project ID in Cloud.ru (falls back to CLOUDRU_PROJECT_ID env var)
- `containerapp_name`: Name of the Container App to retrieve

#### cloudru_create_containerapp(project_id, containerapp_name, containerapp_port, containerapp_image)

Creates a new Container App in Cloud.ru.

Parameters:
- `project_id`: Project ID in Cloud.ru (falls back to CLOUDRU_PROJECT_ID env var)
- `containerapp_name`: Name of the Container App to create
- `containerapp_port`: Port number for the Container App
- `containerapp_image`: Image for the Container App

#### cloudru_delete_containerapp(project_id, containerapp_name)

Deletes a Container App from Cloud.ru. WARNING: This action cannot be undone!

Parameters:
- `project_id`: Project ID in Cloud.ru (falls back to CLOUDRU_PROJECT_ID env var)
- `containerapp_name`: Name of the Container App to delete

#### cloudru_start_containerapp(project_id, containerapp_name)

Starts a Container App in Cloud.ru.

Parameters:
- `project_id`: Project ID in Cloud.ru (falls back to CLOUDRU_PROJECT_ID env var)
- `containerapp_name`: Name of the Container App to start

#### cloudru_stop_containerapp(project_id, containerapp_name)

Stops a Container App in Cloud.ru.

Parameters:
- `project_id`: Project ID in Cloud.ru (falls back to CLOUDRU_PROJECT_ID env var)
- `containerapp_name`: Name of the Container App to stop

#### cloudru_get_list_docker_registries(project_id)

Gets a list of Docker Registries from Cloud.ru. Project ID can be set via CLOUDRU_PROJECT_ID environment variable and obtained from console.cloud.ru.

Parameters:
- `project_id`: Project ID in Cloud.ru (falls back to CLOUDRU_PROJECT_ID env var)

#### cloudru_create_docker_registry(project_id, registry_name, is_public)

Creates a new Docker Registry in Cloud.ru.

Parameters:
- `project_id`: Project ID in Cloud.ru (falls back to CLOUDRU_PROJECT_ID env var)
- `registry_name`: Name of the Docker Registry to create
- `is_public`: Boolean flag indicating if the registry should be public (true) or private (false)

## Running the MCP Server

To start the MCP server, you can use either the locally built binary or the Go-installed binary:

### Using the locally built binary:

```bash
git clone <this repo>
cd cloudru-containerapps-mcp
go build -o cloudru-containerapps-mcp cmd/cloudru-containerapps-mcp/main.go
./cloudru-containerapps-mcp
```
or [docs/INSTALLATION.md](docs/INSTALLATION.md)

The server will listen for JSON-RPC messages on stdin/stdout.

## Documentation

For more information about Cloud.ru Container Apps, see:
https://cloud.ru/docs/container-apps-evolution/ug/topics/tutorials__before-work
