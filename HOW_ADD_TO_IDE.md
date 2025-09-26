# How to Add Cloud.ru Container Apps MCP to IDEs and AI Coding Assistants

This guide explains how to integrate the Cloud.ru Container Apps MCP with popular IDEs and AI coding assistants.

## Prerequisites

Before adding the MCP to your IDE, ensure you have:

1. Either build the MCP server: `go build -o cloudru-containerapps-mcp`
   OR install it using: `go install github.com/Nick1994209/cloudru_containerapps_mcp/cmd/cloudru-containerapps-mcp@latest`
   (if using the installed version, make sure GOPATH/bin is in your PATH)
   Note: This method works for released versions that include the proper directory structure.
   For local development or if the remote repository doesn't have the cmd directory yet,
   use the build from source method.
2. Set up your Cloud.ru credentials as environment variables:
   Alternatively, you can create a `.env` file in the project root directory with these variables.
   
   **Required environment variables:**
   - `CLOUDRU_KEY_ID` - Your service account key ID (required)
   - `CLOUDRU_KEY_SECRET` - Your service account key secret (required)
   
   To obtain access keys for authentication, please follow the instructions at:
   https://cloud.ru/docs/console_api/ug/topics/quickstart
   
   You will need a Key ID and Key Secret to use this service.
   
   **Optional environment variables:**
   - `CLOUDRU_REGISTRY_NAME` - Your Cloud.ru registry name
   - `CLOUDRU_PROJECT_ID` - Your Cloud.ru project ID
   - `CLOUDRU_DOCKERFILE` - Path to Dockerfile (optional)
   - `CLOUDRU_DOCKERFILE_TARGET` - Target stage in a multi-stage Dockerfile (optional, defaults to '-' which means no target)
   - `CLOUDRU_DOCKERFILE_FOLDER` - Dockerfile folder (build context, optional, defaults to '.' which means current directory)

## Kilo Code

### Adding the MCP Server

1. Open Kilo Code settings
2. Navigate to "Tools" → "MCP Servers"
3. Click "Add New Server"
4. Configure with:
   - Name: `Cloud.ru Container Apps`
   - Command: `./cloudru-containerapps-mcp` (or just `cloudru-containerapps-mcp` if you've installed it via `go install` and added GOPATH/bin to your PATH)
   - Working Directory: Path to this project directory (only needed if using the local binary)

### Setting Environment Variables in Kilo Code

In Kilo Code, you can set environment variables directly in the MCP server configuration:

1. In the MCP server configuration dialog, look for the "Environment Variables" section
2. Add the required variables in JSON format:

```json
{
  "CLOUDRU_REGISTRY_NAME": "your-registry-name",
  "CLOUDRU_KEY_ID": "your-service-account-key-id",
  "CLOUDRU_KEY_SECRET": "your-service-account-key-secret",
  "CLOUDRU_PROJECT_ID": "your-project-id",
  "CLOUDRU_REPOSITORY_NAME": "your-repository-name",
  "CLOUDRU_DOCKERFILE": "Dockerfile",
  "CLOUDRU_DOCKERFILE_TARGET": "-",
  "CLOUDRU_DOCKERFILE_FOLDER": "."
}
```

3. Alternatively, you can set them in the "Environment" field if it accepts key-value pairs:

```
CLOUDRU_REGISTRY_NAME=your-registry-name
CLOUDRU_KEY_ID=your-service-account-key-id
CLOUDRU_KEY_SECRET=your-service-account-key-secret
CLOUDRU_PROJECT_ID=your-project-id
CLOUDRU_REPOSITORY_NAME=your-repository-name
CLOUDRU_DOCKERFILE=Dockerfile
CLOUDRU_DOCKERFILE_TARGET=-
CLOUDRU_DOCKERFILE_FOLDER=.
```

### Example Configuration JSON

Here's a complete example of how to configure this MCP server in Kilo Code:

```json
{
  "mcpServers": {
    "cloudru-containerapps-mcp": {
      "command": "./cloudru-containerapps-mcp",  // or just "cloudru-containerapps-mcp" if installed via go install
      "args": [],
      "env": {
        "CLOUDRU_REGISTRY_NAME": "your-registry-name",
        "CLOUDRU_KEY_ID": "your-service-account-key-id",
        "CLOUDRU_KEY_SECRET": "your-service-account-key-secret",
        "CLOUDRU_PROJECT_ID": "your-project-id",
        "CLOUDRU_REPOSITORY_NAME": "your-repository-name",
        "CLOUDRU_DOCKERFILE": "Dockerfile",
        "CLOUDRU_DOCKERFILE_TARGET": "-",
        "CLOUDRU_DOCKERFILE_FOLDER": "."
      }
    }
  }
}
```

4. Click "Save and Connect"
5. The tools should now be available in Kilo Code with access to your environment variables

## Roo Code

1. Open Roo Code preferences
2. Go to "Extensions" → "MCP Integration"
3. Click "Add MCP Server"
4. Fill in the details:
   - Server Name: `Cloud.ru Container Apps`
   - Executable Path: Full path to `cloudru-containerapps-mcp` (or just `cloudru-containerapps-mcp` if installed via `go install` and GOPATH/bin is in your PATH)
   - Arguments: Leave empty
5. Configure environment variables in the server settings if supported
6. Click "Test Connection" to verify
7. Enable the server and restart Roo Code if needed

## Claude (Anthropic Console)

1. Access the Anthropic Console
2. Navigate to "Tools" → "Custom MCPs"
3. Click "Register New MCP"
4. Provide the following configuration:
   - Tool Name: `cloudru-containerapps`
   - Execution Method: `subprocess`
   - Command: `["./cloudru-containerapps-mcp"]` (or just `["cloudru-containerapps-mcp"]` if installed via `go install`)
   - Working Directory: Project path (only needed if using the local binary)
   - Environment Variables: Add as key-value pairs in the configuration interface
5. Save the configuration
6. The tools will be available in Claude prompts

## Cursor

1. Open Cursor settings
2. Go to "Extensions" → "MCP Tools"
3. Click "Add External Tool"
4. Enter configuration:
   - Tool Identifier: `cloudru-containerapps-mcp`
   - Executable: Path to the `cloudru-containerapps-mcp` binary
   - Environment Variables: Set in the tool configuration if supported
   - Auto-start: Enabled
5. Restart Cursor to load the tool
6. Verify integration by opening the tools panel

## General Usage Notes

Once integrated, you can use the following tools:

1. `cloudru_containerapps_description()` - Get usage instructions
2. `cloudru_docker_login()` - Authenticate with Cloud.ru registry
3. `cloudru_docker_push()` - Build and push Docker images (supports dockerfile_path, dockerfile_target, and dockerfile_folder parameters)

### Example Prompts

- "Use cloudru_containerapps_description to tell me about this tool"
- "Run cloudru_docker_login with my registry credentials"
- "Execute cloudru_docker_push to deploy my application with version v1.2.3"

## Troubleshooting

If the MCP doesn't appear in your IDE:

1. Verify the binary is executable: `chmod +x cloudru-containerapps-mcp`
2. Test the binary directly: `echo '{"method":"initialize","id":1}' | ./cloudru-containerapps-mcp`
3. Check that all required environment variables are set
4. Ensure your IDE has permission to execute subprocesses
5. Consult your IDE's documentation for MCP/tool integration specifics

## Security Considerations

- The MCP requires access to Docker daemon
- Authentication credentials are handled through environment variables
- Ensure your IDE/assistant has appropriate permissions for container operations
- Never commit sensitive credentials to version control