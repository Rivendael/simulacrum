# Simulacrum

Simulacrum is a high-performance Go service designed to anonymize and "obscure"
sensitive data within JSON payloads. It replaces Personally Identifiable
Information (PII) with realistic, deterministically generated fake data while
preserving the original JSON structure and data types.

This tool is ideal for developers and QA engineers who need to generate
realistic test data from production datasets without exposing sensitive user
information.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Features](#features)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [Docker](#docker)
- [Usage](#usage)
- [Development](#development)
- [License](#license)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Features

- **Smart Obfuscation**: Automatically detects and replaces sensitive fields
such as:
  - Names (First, Last, Middle)
  - Emails & Phone Numbers
  - Addresses (Street, City, State, Zip, Country)
  - IDs (SSN, Passport, Driver's License, Tax ID)
  - Financial Data (Bank Accounts)

- **Deterministic Generation**: Uses consistent hashing to ensure that the same
    input value always produces the same obscured output. This is crucial for
    maintaining referential integrity across datasets.
- **High Performance**: Built on the [Gin](https://github.com/gin-gonic/gin)
    framework and uses [fastjson](https://github.com/valyala/fastjson) for
    efficient JSON processing.
- **Secure**:
  - **JWT Authentication**: Secures the API using JSON Web Tokens (RSA signed).
  - **TLS/mTLS Support**: Supports HTTPS and Mutual TLS for secure communication.
- **Configurable**: Flexible configuration via Environment Variables.

## Getting Started

### Prerequisites

- Go 1.22 or higher
- Docker (For easy deploy and testing)

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/yourusername/simulacrum.git
    cd simulacrum
    ```

2. Download dependencies:

    ```bash
    go mod download
    ```

3. Build the executables:

    ```bash
    # Build the server
    go build \
      -o server \
      cmd/server/main.go

    # Build the utility tools
    go build \
      -o generate-keys \
      cmd/generate-keys/generate_keys.go
    go build \
      -o generate-certs \
      cmd/generate-certs/generate_certs.go
    ```

## Configuration

Simulacrum is configured using Environment Variables.

### Environment Variables

| Variable                  | Description                             | Default                               |
|---------------------------|-----------------------------------------|---------------------------------------|
| `SERVER_PORT`             | Port to listen on                       | `8080`                                |
| `SERVER_ENVIRONMENT`      | Environment name                        | `development`                         |
| `GIN_MODE`                | Gin framework mode (`debug`, `release`) | `debug` (or `release` if env is prod) |
| `AUTH_PUBLIC_KEYS_FILE`   | Path to public keys file                | `public_keys.pem`                     |
| `TLS_ENABLED`             | Enable HTTPS (`true`/`false`)           | `false`                               |
| `TLS_CERT_FILE`           | Path to server certificate              |                                       |
| `TLS_KEY_FILE`            | Path to server private key              |                                       |
| `TLS_CA_CERT_FILE`        | Path to CA certificate (for mTLS)       |                                       |
| `TLS_REQUIRE_CLIENT_CERT` | Require mTLS (`true`/`false`)           | `false`                               |

## Docker

Simulacrum includes a Dockerfile for easy deployment.

### 1. Build the Image

```bash
docker build \
  --tag simulacrum \
  .
```

### 2. Run the Container

You must mount the public keys file into the container.

**Basic Usage:**

```bash
docker run \
  --publish 8080:8080 \
  --volumne $(pwd)/public_keys.pem:/app/public_keys.pem \
  simulacrum
```

**With Environment Variables:**

```bash
docker run \
  --publish 8080:8080 \
  --volume $(pwd)/public_keys.pem:/app/public_keys.pem \
  --env SERVER_ENVIRONMENT=production \
  --env GIN_MODE=release \
  --env TLS_ENABLED=false \
  simulacrum
```

**With Custom Config File:**

```bash
docker run \
  --publish 8080:8080 \
  --volume $(pwd)/public_keys.pem:/app/public_keys.pem \
  --volume $(pwd)/config.prod.yaml:/app/config.yaml \
  simulacrum
```

## Usage

### 1. Generate Keys & Certificates

Before running the server, you need to generate the necessary keys for JWT
authentication and optionally TLS certificates.

```bash
# Generate JWT RSA keys (private key for signing tokens, public key for the 
# server)
./generate-keys

# (Optional) Generate TLS certificates for HTTPS/mTLS
./generate-certs
```

### 2. Run the Server

```bash
./server -config config.dev.yaml
```

The server will start on the configured port (default `:8080`).

### 3. API Endpoints

#### `GET /health`

Health check endpoint. No authentication required.

**Response:**

```json
{
  "status": "ok"
}
```

#### `POST /obscure`

The main endpoint to obscure data.

- **Headers**: `Authorization: Bearer <JWT_TOKEN>`
- **Body**: Arbitrary JSON object.

**Example Request:**

```bash
curl --request POST \
  "http://localhost:8080/obscure" \
  --header "Authorization: Bearer <YOUR_JWT_TOKEN>" \
  --header "Content-Type: application/json" \
  --data '{
    "id": "12345",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "address": {
        "street": "123 Main St",
        "city": "New York"
    }
  }'
```

**Example Response:**

```json
{
  "id": "84921",
  "name": "Alice Smith",
  "email": "alice.smith@test.com",
  "address": {
      "street": "456 Elm St",
      "city": "Los Angeles"
  }
}
```

## Development

### Running Tests

```bash
go test ./...
```

### Project Structure

- `cmd/`: Entry points for the server and utility tools.
- `internal/auth/`: JWT handling and middleware.
- `internal/config/`: Configuration loading logic.
- `internal/data/`: Data generation logic (names, addresses, etc.).
- `internal/handlers/`: HTTP request handlers.
- `bruno/`: API collection for [Bruno](https://www.usebruno.com/) (useful for
    testing).

## License

[MIT](LICENSE)
