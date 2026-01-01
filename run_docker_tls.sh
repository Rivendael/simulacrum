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
  -t simulacrum .

# Run with environment variables (no config file mount)
# Files (keys, certs) are baked into the image via Dockerfile
echo "Running Simulacrum with environment variables (Production + TLS)..."
docker run --rm -p 8080:8080 \
  -e SERVER_PORT=8080 \
  -e SERVER_ENVIRONMENT=production \
  -e GIN_MODE=release \
  -e TLS_ENABLED=true \
  -e TLS_REQUIRE_CLIENT_CERT=true \
  simulacrum
