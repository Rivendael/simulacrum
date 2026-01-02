#!/bin/bash

# Define certificate/key files
PUBLIC_KEYS_FILE="public_keys.pem"
CERT_FILE="server.crt"
KEY_FILE="server.key"
CA_CERT_FILE="ca.crt"

# Build the image
echo "Building Docker image..."
docker build \
  --build-arg PUBLIC_KEYS_FILE="$PUBLIC_KEYS_FILE" \
  --build-arg CERT_FILE="$CERT_FILE" \
  --build-arg KEY_FILE="$KEY_FILE" \
  --build-arg CA_CERT_FILE="$CA_CERT_FILE" \
  --tag simulacrum \
  .

# Run with environment variables (no config file mount)
# Files (keys, certs) are baked into the image via Dockerfile
echo "Running Simulacrum with environment variables (Production + TLS)..."
docker run --rm \
  --publish 8080:8080 \
  --env SERVER_PORT=8080 \
  --env SERVER_ENVIRONMENT=production \
  --env GIN_MODE=release \
  --env TLS_ENABLED=true \
  --env TLS_REQUIRE_CLIENT_CERT=true \
  simulacrum
