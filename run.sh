#!/bin/bash

set -e

echo "Building simulacrum..."
go build -o simulacrum ./cmd/server

echo "Starting server on :8080..."
./simulacrum
