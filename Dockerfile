# check=skip=SecretsUsedInArgOrEnv
# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the server binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server cmd/server/main.go

# Final stage
FROM scratch

WORKDIR /app

# Copy CA certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary from builder
COPY --from=builder /app/server .

# Define build arguments for config files
ARG PUBLIC_KEYS_FILE=public_keys.pem
ARG CERT_FILE=server.crt
ARG KEY_FILE=server.key
ARG CA_CERT_FILE=ca.crt

# Set environment variables
ENV AUTH_PUBLIC_KEYS_FILE=${PUBLIC_KEYS_FILE}
ENV TLS_CERT_FILE=${CERT_FILE}
ENV TLS_KEY_FILE=${KEY_FILE}
ENV TLS_CA_CERT_FILE=${CA_CERT_FILE}

# Copy configuration files
COPY ${PUBLIC_KEYS_FILE} .
COPY ${CERT_FILE} .
COPY ${KEY_FILE} .
COPY ${CA_CERT_FILE} .

# Expose the application port
EXPOSE 8080

# Create a volume mount point for keys
VOLUME /app/keys

# Run the server
CMD ["./server"]
