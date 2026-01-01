package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Helper function to generate RSA key pair for testing
func generateTestKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

// Helper function to create a test JWT token
func createTestToken(privateKey *rsa.PrivateKey, claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

// Helper function to write public key to PEM file
func writePublicKeyToFile(pubKey *rsa.PublicKey, filepath string) error {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return err
	}

	pubKeyPEM := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return pem.Encode(file, pubKeyPEM)
}

// Helper function to write multiple public keys to a file
func writeMultiplePublicKeysToFile(pubKeys []*rsa.PublicKey, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, pubKey := range pubKeys {
		pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
		if err != nil {
			return err
		}

		pubKeyPEM := &pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: pubKeyBytes,
		}

		if err := pem.Encode(file, pubKeyPEM); err != nil {
			return err
		}
	}

	return nil
}

// TestLoadPublicKeysFromFile tests loading a single public key from file
func TestLoadPublicKeysFromFile(t *testing.T) {
	// Generate a test key pair
	_, pubKey, err := generateTestKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate test key pair: %v", err)
	}

	// Write to temporary file
	tmpFile := "test_pubkey.pem"
	defer os.Remove(tmpFile)

	if err := writePublicKeyToFile(pubKey, tmpFile); err != nil {
		t.Fatalf("Failed to write public key to file: %v", err)
	}

	// Load the key
	pkm, err := LoadPublicKeysFromFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load public keys: %v", err)
	}

	if len(pkm.Keys) != 1 {
		t.Errorf("Expected 1 key, got %d", len(pkm.Keys))
	}
}

// TestLoadMultiplePublicKeysFromFile tests loading multiple public keys from file
func TestLoadMultiplePublicKeysFromFile(t *testing.T) {
	// Generate multiple test key pairs
	var pubKeys []*rsa.PublicKey
	for i := 0; i < 3; i++ {
		_, pubKey, err := generateTestKeyPair()
		if err != nil {
			t.Fatalf("Failed to generate test key pair: %v", err)
		}
		pubKeys = append(pubKeys, pubKey)
	}

	// Write to temporary file
	tmpFile := "test_pubkeys_multi.pem"
	defer os.Remove(tmpFile)

	if err := writeMultiplePublicKeysToFile(pubKeys, tmpFile); err != nil {
		t.Fatalf("Failed to write public keys to file: %v", err)
	}

	// Load the keys
	pkm, err := LoadPublicKeysFromFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load public keys: %v", err)
	}

	if len(pkm.Keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(pkm.Keys))
	}
}

// TestLoadPublicKeysFromFileNotFound tests error handling for missing file
func TestLoadPublicKeysFromFileNotFound(t *testing.T) {
	_, err := LoadPublicKeysFromFile("nonexistent_file.pem")
	if err == nil {
		t.Errorf("Expected error for nonexistent file")
	}
}

// TestValidateTokenValid tests validating a valid token
func TestValidateTokenValid(t *testing.T) {
	// Generate test key pair
	privKey, pubKey, err := generateTestKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate test key pair: %v", err)
	}

	// Create a token with valid claims
	claims := jwt.MapClaims{
		"sub": "user123",
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	tokenString, err := createTestToken(privKey, claims)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	// Write public key to file
	tmpFile := "test_pubkey_validate.pem"
	defer os.Remove(tmpFile)
	if err := writePublicKeyToFile(pubKey, tmpFile); err != nil {
		t.Fatalf("Failed to write public key to file: %v", err)
	}

	// Load and validate
	pkm, err := LoadPublicKeysFromFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load public keys: %v", err)
	}

	resultClaims, err := pkm.ValidateToken(tokenString)
	if err != nil {
		t.Errorf("Failed to validate valid token: %v", err)
	}

	if resultClaims["sub"] != "user123" {
		t.Errorf("Expected subject 'user123', got %v", resultClaims["sub"])
	}
}

// TestValidateTokenWithBearerPrefix tests validating token with Bearer prefix
func TestValidateTokenWithBearerPrefix(t *testing.T) {
	// Generate test key pair
	privKey, pubKey, err := generateTestKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate test key pair: %v", err)
	}

	// Create a token
	claims := jwt.MapClaims{
		"sub": "user456",
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	tokenString, err := createTestToken(privKey, claims)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	// Write public key to file
	tmpFile := "test_pubkey_bearer.pem"
	defer os.Remove(tmpFile)
	if err := writePublicKeyToFile(pubKey, tmpFile); err != nil {
		t.Fatalf("Failed to write public key to file: %v", err)
	}

	// Load and validate with Bearer prefix
	pkm, err := LoadPublicKeysFromFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load public keys: %v", err)
	}

	tokenWithBearer := fmt.Sprintf("Bearer %s", tokenString)
	resultClaims, err := pkm.ValidateToken(tokenWithBearer)
	if err != nil {
		t.Errorf("Failed to validate token with Bearer prefix: %v", err)
	}

	if resultClaims["sub"] != "user456" {
		t.Errorf("Expected subject 'user456', got %v", resultClaims["sub"])
	}
}

// TestValidateTokenExpired tests validating an expired token
func TestValidateTokenExpired(t *testing.T) {
	// Generate test key pair
	privKey, pubKey, err := generateTestKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate test key pair: %v", err)
	}

	// Create an expired token
	claims := jwt.MapClaims{
		"sub": "user789",
		"exp": time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
		"iat": time.Now().Add(-2 * time.Hour).Unix(),
	}

	tokenString, err := createTestToken(privKey, claims)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	// Write public key to file
	tmpFile := "test_pubkey_expired.pem"
	defer os.Remove(tmpFile)
	if err := writePublicKeyToFile(pubKey, tmpFile); err != nil {
		t.Fatalf("Failed to write public key to file: %v", err)
	}

	// Load and validate
	pkm, err := LoadPublicKeysFromFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load public keys: %v", err)
	}

	_, err = pkm.ValidateToken(tokenString)
	if err == nil {
		t.Errorf("Expected error for expired token")
	}
}

// TestValidateTokenInvalid tests validating an invalid token
func TestValidateTokenInvalid(t *testing.T) {
	// Generate test key pair
	_, pubKey, err := generateTestKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate test key pair: %v", err)
	}

	// Write public key to file
	tmpFile := "test_pubkey_invalid.pem"
	defer os.Remove(tmpFile)
	if err := writePublicKeyToFile(pubKey, tmpFile); err != nil {
		t.Fatalf("Failed to write public key to file: %v", err)
	}

	// Load and try to validate invalid token
	pkm, err := LoadPublicKeysFromFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load public keys: %v", err)
	}

	_, err = pkm.ValidateToken("invalid.token.here")
	if err == nil {
		t.Errorf("Expected error for invalid token")
	}
}

// TestValidateTokenEmptyToken tests validating an empty token
func TestValidateTokenEmptyToken(t *testing.T) {
	// Generate test key pair
	_, pubKey, err := generateTestKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate test key pair: %v", err)
	}

	// Write public key to file
	tmpFile := "test_pubkey_empty.pem"
	defer os.Remove(tmpFile)
	if err := writePublicKeyToFile(pubKey, tmpFile); err != nil {
		t.Fatalf("Failed to write public key to file: %v", err)
	}

	// Load and try to validate empty token
	pkm, err := LoadPublicKeysFromFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load public keys: %v", err)
	}

	_, err = pkm.ValidateToken("")
	if err == nil {
		t.Errorf("Expected error for empty token")
	}
}

// TestValidateTokenMultipleKeys tests validating with multiple public keys
func TestValidateTokenMultipleKeys(t *testing.T) {
	// Generate two test key pairs
	privKey1, pubKey1, err := generateTestKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate first key pair: %v", err)
	}

	_, pubKey2, err := generateTestKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate second key pair: %v", err)
	}

	// Create a token with the first key
	claims := jwt.MapClaims{
		"sub": "user_multi",
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	tokenString, err := createTestToken(privKey1, claims)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	// Write both public keys to file
	tmpFile := "test_pubkeys_multi_validate.pem"
	defer os.Remove(tmpFile)
	if err := writeMultiplePublicKeysToFile([]*rsa.PublicKey{pubKey1, pubKey2}, tmpFile); err != nil {
		t.Fatalf("Failed to write public keys to file: %v", err)
	}

	// Load and validate - should succeed even though token uses pubKey1
	pkm, err := LoadPublicKeysFromFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load public keys: %v", err)
	}

	resultClaims, err := pkm.ValidateToken(tokenString)
	if err != nil {
		t.Errorf("Failed to validate token with multiple keys: %v", err)
	}

	if resultClaims["sub"] != "user_multi" {
		t.Errorf("Expected subject 'user_multi', got %v", resultClaims["sub"])
	}
}
