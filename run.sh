#!/bin/bash

set -e

echo "Building simulacrum..."
go build \
  -o simulacrum \
  ./cmd/server

echo "Starting server with environment variables..."
export SERVER_PORT=8080
export SERVER_ENVIRONMENT=development
export GIN_MODE=debug
export AUTH_PUBLIC_KEYS_FILE=public_keys.pem
./simulacrum
