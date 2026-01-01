# JWT Authentication Setup

Your API now requires JWT tokens for the `/obscure` endpoint. Here's how to use it:

## 1. Generate Test Keys & Tokens

```bash
go run generate_test_keys.go
```

This creates:
- `test_private.key` - Private key (keep secret, used to sign tokens)
- `public_keys.pem` - Public key (given to server, used to verify tokens)
- `test_token.txt` - A valid test token (24-hour expiration)

## 2. Start the Server

```bash
go run ./cmd/server -keys public_keys.pem -port 8080
```

Or with custom options:
```bash
go run ./cmd/server -keys my_keys.pem -port 3000
```

## 3. Make API Requests

### Health Check (No Auth Required)
```bash
curl http://localhost:8080/health
```

### Obscure Data (Requires JWT)
```bash
# Get the token from test_token.txt
TOKEN=$(cat test_token.txt)

# Make request with Bearer token
curl -X POST http://localhost:8080/obscure \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "user123",
    "name": "John Doe",
    "email": "john@example.com",
    "phone_number": "+1-555-0123"
  }'
```

## Architecture

- **JWT Validation**: `internal/auth/jwt.go`
  - Loads public keys from PEM file
  - Supports multiple keys for key rotation
  - Validates token signatures and expiration
  
- **Middleware**: `internal/auth/middleware.go`
  - `JWTMiddleware()` - Gin middleware for automatic validation
  - Returns 401 Unauthorized on invalid tokens
  - Stores claims in context for handler access

## Token Claims

Your test tokens include:
```json
{
  "sub": "test-user",
  "user": "testuser",
  "exp": 1767342143,
  "iat": 1767255743
}
```

Access claims in handlers:
```go
claims := auth.GetClaims(c)
user := auth.GetClaimString(c, "user")
```

## Multiple Public Keys (Key Rotation)

To support key rotation, add multiple public keys to the PEM file:

```bash
cat key1.pem key2.pem key3.pem > public_keys.pem
```

The server will try each key until one validates the token.

## Creating Custom Tokens

Use the private key to create custom tokens:

```bash
# In Go
privateKey, _ := ioutil.ReadFile("test_private.key")
// ... parse as RSA private key
claims := jwt.MapClaims{
  "sub": "your-user-id",
  "exp": time.Now().Add(24 * time.Hour).Unix(),
}
token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
tokenString, _ := token.SignedString(privateKey)
```

## Error Responses

**Missing Token:**
```json
{"error": "missing authorization header"}
```

**Invalid Token:**
```json
{"error": "failed to validate token: ..."}
```

**Expired Token:**
```json
{"error": "failed to validate token: token is expired"}
```
