#!/bin/bash

# Build script for cloudru-containerapps-mcp
set -e  # Exit on any error

echo "Building cloudru-containerapps-mcp..."

# Remove existing binary if it exists
if [ -f "cloudru-containerapps-mcp" ]; then
    echo "Removing existing binary..."
    rm cloudru-containerapps-mcp
fi

# Build the MCP server from the cmd directory
echo "Compiling cmd/cloudru-containerapps-mcp/main.go..."
go build -o cloudru-containerapps-mcp cmd/cloudru-containerapps-mcp/main.go

if [ $? -eq 0 ]; then
    echo "Build successful!"
    echo "Binary created: cloudru-containerapps-mcp"
    echo "Binary size: $(du -h cloudru-containerapps-mcp | cut -f1)"
else
    echo "Build failed!"
    exit 1
fi