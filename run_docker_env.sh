#!/bin/bash

# Build the image
echo "Building Docker image..."
docker build -t simulacrum .

# Run with environment variables
# We don't mount a config file, so the app will use defaults + these overrides
echo "Running Simulacrum with environment variables..."
docker run --rm -p 8080:8080 \
  -v "$(pwd)/public_keys.pem:/app/public_keys.pem" \
  -e SERVER_ENVIRONMENT=production \
  -e GIN_MODE=release \
  -e TLS_ENABLED=false \
  simulacrum
