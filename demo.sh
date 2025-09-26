#!/bin/bash

# Demo script for Cloud.ru Container Apps MCP

echo "Cloud.ru Container Apps MCP Demo"
echo "================================="

# Check if the MCP binary exists
if [ ! -f "./cloudru-containerapps-mcp" ]; then
    echo "Building MCP server..."
    go build -o cloudru-containerapps-mcp
fi

echo "Starting MCP server..."
echo "You can now communicate with the MCP server using JSON-RPC over stdin/stdout."
echo ""
echo "Example commands:"
echo "1. Initialize: {\"method\":\"initialize\",\"id\":1}"
echo "2. Get description: {\"method\":\"cloudru_containerapps_description\",\"id\":2}"
echo "3. Docker login: {\"method\":\"cloudru_docker_login\",\"params\":{\"registry_name\":\"myregistry\",\"key_id\":\"mykeyid\",\"key_secret\":\"mykeysecret\"},\"id\":3}"
echo "4. Docker push: {\"method\":\"cloudru_docker_push\",\"params\":{\"registry_name\":\"myregistry\",\"repository_name\":\"myrepo\",\"image_version\":\"v1.0.0\"},\"id\":4}"
echo ""
echo "To exit, press Ctrl+C"

# Start the MCP server
./cloudru-containerapps-mcp