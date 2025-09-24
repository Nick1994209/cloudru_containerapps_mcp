# Cloud.ru Container Apps MCP

A Model Context Protocol (MCP) server for interacting with Cloud.ru Container Apps and Artifact Registry.

## Features

This MCP provides the following functions:

1. `cloudru_containerapps_description()` - Returns usage instructions for this MCP
2. `cloudru_containerapps_docker_login(registry_name, key_id, key_secret)` - Login to Cloud.ru Docker registry
3. `cloudru_containerapps_docker_push(registry_name, repository_name, image_version, key_id, key_secret)` - Build and push Docker image to Cloud.ru Artifact Registry

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

### Functions

#### cloudru_containerapps_description()

Returns usage instructions for this MCP.

#### cloudru_containerapps_docker_login(registry_name, key_id, key_secret)

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

#### cloudru_containerapps_docker_push(registry_name, repository_name, image_version, key_id, key_secret)

Builds a Docker image and pushes it to Cloud.ru Artifact Registry.

Parameters:
- `registry_name`: Name of the registry (falls back to REGISTRY_NAME env var)
- `repository_name`: Name of the repository (falls back to REPOSITORY_NAME env var, then to current directory name)
- `image_version`: Version/tag for the image
- `key_id`: Service account key ID (falls back to KEY_ID env var)
- `key_secret`: Service account key secret (falls back to KEY_SECRET env var)
- `dockerfile_path`: Path to Dockerfile (optional, defaults to 'Dockerfile')

If Docker push fails due to authentication issues and KEY_ID/KEY_SECRET environment variables are set, the function will attempt to re-login and retry the push operation.

## Running the MCP Server

To start the MCP server, simply run:

```bash
./cloudru-containerapps-mcp
```

The server will listen for JSON-RPC messages on stdin/stdout.

## Documentation

For more information about Cloud.ru Container Apps, see:
https://cloud.ru/docs/container-apps-evolution/ug/topics/tutorials__before-work