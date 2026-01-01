# Simulacrum

A deterministic data obscuration API that transforms personal information into false identities using cryptographic hashing.

## Features

- **Deterministic Obscuration**: Same input always produces the same fake identity
- **Irreversible**: Uses SHA-256 hashing so original data cannot be recovered
- **ID-based Seeding**: Uses customer ID as the randomization seed
- **Comprehensive Coverage**: Obscures names, emails, addresses, and phone numbers

## Project Structure

```
simulacrum/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── data/
│   │   ├── data.go          # Core obscuration logic and PersonalData struct
│   │   ├── data_test.go     # Data obscuration tests
│   │   └── names.go         # Name lists for generation
│   └── handlers/
│       ├── handlers.go      # HTTP request handlers
│       └── handlers_test.go # Handler tests
├── go.mod
├── go.sum
└── README.md
```

## Building

```bash
go mod tidy
go build -o simulacrum ./cmd/server
```

## Running

```bash
./simulacrum
```

The server will start on `:8080`.

## API

### POST /obscure

Takes personal data and returns obscured data with a deterministic false identity.

**Request:**
```json
{
  "id": "customer123",
  "name": "John Doe",
  "email": "john@example.com",
  "address": "123 Main St",
  "phone_number": "555-1234"
}
```

**Response:**
```json
{
  "id": "customer123",
  "name": "Sarah Anderson",
  "email": "robert.williams@test.org",
  "address": "4782 Oak Ave",
  "phone_number": "555-284-6371"
}
```

## Testing

```bash
go test ./...
```

All tests verify:
- Deterministic behavior (same input → same output)
- Proper obscuration (output ≠ input)
- ID-based uniqueness (different IDs → different outputs)
- Empty value preservation
- API validation and error handling
