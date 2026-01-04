#!/bin/bash

# Build the image
echo "Building Docker image..."
docker build \
  --tag simulacrum \
  .

# Run with environment variables
# We don't mount a config file, so the app will use defaults + these overrides
echo "Running Simulacrum with environment variables..."
docker run --rm \
  --publish 8080:8080 \
  --volume "$(pwd)/public_keys.pem:/app/public_keys.pem" \
  --env SERVER_ENVIRONMENT=production \
  --env GIN_MODE=release \
  --env TLS_ENABLED=false \
  simulacrum
