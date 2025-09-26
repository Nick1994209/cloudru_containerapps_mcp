# Cloud.ru Container Apps MCP

A Model Context Protocol (MCP) server for interacting with Cloud.ru Container Apps and Artifact Registry.

## Features

This MCP provides the following functions:

1. `cloudru_containerapps_description()` - Returns usage instructions for this MCP
2. `cloudru_docker_login(registry_name, key_id, key_secret)` - Login to Cloud.ru Docker registry
3. `cloudru_docker_push(registry_name, repository_name, image_version, dockerfile_path, dockerfile_target, dockerfile_folder, key_id, key_secret)` - Build and push Docker image to Cloud.ru Artifact Registry
4. `cloudru_get_list_containerapps(project_id, key_id, key_secret)` - Get list of Container Apps from Cloud.ru. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru
5. `cloudru_get_containerapp(project_id, containerapp_name, key_id, key_secret)` - Get a specific Container App from Cloud.ru by name. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru
6. `cloudru_create_containerapp(project_id, containerapp_name, containerapp_port, containerapp_image, key_id, key_secret)` - Create a new Container App in Cloud.ru
7. `cloudru_delete_containerapp(project_id, containerapp_name, key_id, key_secret)` - Delete a Container App from Cloud.ru. WARNING: This action cannot be undone!
8. `cloudru_start_containerapp(project_id, containerapp_name, key_id, key_secret)` - Start a Container App in Cloud.ru
9. `cloudru_stop_containerapp(project_id, containerapp_name, key_id, key_secret)` - Stop a Container App in Cloud.ru

## Prerequisites

- Go 1.22 or later
- Docker installed and configured
- Access to Cloud.ru Container Apps service
- Service account credentials with appropriate permissions

## Installation

1. Clone this repository
2. Run `go build -o cloudru-containerapps-mcp` to build the binary
3. Make sure Docker is installed and running on your system

## Usage

### Environment Variables

The following environment variables can be used as fallbacks for function parameters:

- `REGISTRY_NAME`: Registry name
- `KEY_ID`: Service account key ID
- `KEY_SECRET`: Service account key secret
- `REPOSITORY_NAME`: Repository name (defaults to current directory name if not set)
- `PROJECT_ID`: Project ID for Container Apps (can be obtained from console.cloud.ru)
- `CONTAINERAPP_NAME`: Container App name (optional)
- `DOCKERFILE`: Path to Dockerfile (defaults to 'Dockerfile' if not set)
- `DOCKERFILE_TARGET`: Target stage in a multi-stage Dockerfile (optional, defaults to '-' which means no target)
- `DOCKERFILE_FOLDER`: Dockerfile folder (build context, defaults to '.' which means current directory)

### Functions

#### cloudru_containerapps_description()

Returns usage instructions for this MCP.

#### cloudru_docker_login(registry_name, key_id, key_secret)

Logs into the Cloud.ru Docker registry using the provided credentials.

Parameters:
- `registry_name`: Name of the registry (falls back to REGISTRY_NAME env var)
- `key_id`: Service account key ID (falls back to KEY_ID env var)
- `key_secret`: Service account key secret (falls back to KEY_SECRET env var)

If login fails, you'll need to:
1. Go to Cloud.ru Evolution Artifact Registry
2. Create a registry
3. Obtain access keys
4. See documentation: https://cloud.ru/docs/container-apps-evolution/ug/topics/tutorials__before-work

#### cloudru_docker_push(registry_name, repository_name, image_version, dockerfile_path, dockerfile_target, dockerfile_folder, key_id, key_secret)

Builds a Docker image and pushes it to Cloud.ru Artifact Registry.

Parameters:
- `registry_name`: Name of the registry (falls back to REGISTRY_NAME env var)
- `repository_name`: Name of the repository (falls back to REPOSITORY_NAME env var, then to current directory name)
- `image_version`: Version/tag for the image
- `dockerfile_path`: Path to Dockerfile (optional, defaults to 'Dockerfile')
- `dockerfile_target`: Target stage in a multi-stage Dockerfile (optional, defaults to '-' which means no target)
- `dockerfile_folder`: Dockerfile folder (build context, defaults to '.' which means current directory)
- `key_id`: Service account key ID (falls back to KEY_ID env var)
- `key_secret`: Service account key secret (falls back to KEY_SECRET env var)

If Docker push fails due to authentication issues and KEY_ID/KEY_SECRET environment variables are set, the function will attempt to re-login and retry the push operation.

## Running the MCP Server

To start the MCP server, simply run:

#### cloudru_get_list_containerapps(project_id, key_id, key_secret)

Gets a list of Container Apps from Cloud.ru. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru.

Parameters:
- `project_id`: Project ID in Cloud.ru (falls back to PROJECT_ID env var)
- `key_id`: Service account key ID (falls back to KEY_ID env var)
- `key_secret`: Service account key secret (falls back to KEY_SECRET env var)

#### cloudru_get_containerapp(project_id, containerapp_name, key_id, key_secret)

Gets a specific Container App from Cloud.ru by name. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru.

Parameters:
- `project_id`: Project ID in Cloud.ru (falls back to PROJECT_ID env var)
- `containerapp_name`: Name of the Container App to retrieve
- `key_id`: Service account key ID (falls back to KEY_ID env var)
- `key_secret`: Service account key secret (falls back to KEY_SECRET env var)

#### cloudru_create_containerapp(project_id, containerapp_name, containerapp_port, containerapp_image, key_id, key_secret)

Creates a new Container App in Cloud.ru.

Parameters:
- `project_id`: Project ID in Cloud.ru (falls back to PROJECT_ID env var)
- `containerapp_name`: Name of the Container App to create
- `containerapp_port`: Port number for the Container App
- `containerapp_image`: Image for the Container App
- `key_id`: Service account key ID (falls back to KEY_ID env var)
- `key_secret`: Service account key secret (falls back to KEY_SECRET env var)

#### cloudru_delete_containerapp(project_id, containerapp_name, key_id, key_secret)

Deletes a Container App from Cloud.ru. WARNING: This action cannot be undone!

Parameters:
- `project_id`: Project ID in Cloud.ru (falls back to PROJECT_ID env var)
- `containerapp_name`: Name of the Container App to delete
- `key_id`: Service account key ID (falls back to KEY_ID env var)
- `key_secret`: Service account key secret (falls back to KEY_SECRET env var)

#### cloudru_start_containerapp(project_id, containerapp_name, key_id, key_secret)

Starts a Container App in Cloud.ru.

Parameters:
- `project_id`: Project ID in Cloud.ru (falls back to PROJECT_ID env var)
- `containerapp_name`: Name of the Container App to start
- `key_id`: Service account key ID (falls back to KEY_ID env var)
- `key_secret`: Service account key secret (falls back to KEY_SECRET env var)

#### cloudru_stop_containerapp(project_id, containerapp_name, key_id, key_secret)

Stops a Container App in Cloud.ru.

Parameters:
- `project_id`: Project ID in Cloud.ru (falls back to PROJECT_ID env var)
- `containerapp_name`: Name of the Container App to stop
- `key_id`: Service account key ID (falls back to KEY_ID env var)
- `key_secret`: Service account key secret (falls back to KEY_SECRET env var)

## Running the MCP Server

To start the MCP server, simply run:

```bash
./cloudru-containerapps-mcp
```

The server will listen for JSON-RPC messages on stdin/stdout.

## Documentation

For more information about Cloud.ru Container Apps, see:
https://cloud.ru/docs/container-apps-evolution/ug/topics/tutorials__before-work