# Example Dockerfile for testing the Cloud.ru Container Apps MCP
FROM alpine:latest

# Install a simple web server
RUN apk add --no-cache python3

# Create a simple HTML file
RUN echo "<html><body><h1>Hello from Cloud.ru Container Apps!</h1></body></html>" > /index.html

# Expose port 8000
EXPOSE 8000

# Run a simple HTTP server
CMD ["python3", "-m", "http.server", "8000"]