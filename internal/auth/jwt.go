package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// PublicKeyManager holds multiple public keys for JWT validation
type PublicKeyManager struct {
	Keys []*rsa.PublicKey
}

// LoadPublicKeysFromFile reads a PEM file containing one or more public keys
// Each public key should be in PEM format (-----BEGIN PUBLIC KEY-----)
func LoadPublicKeysFromFile(filepath string) (*PublicKeyManager, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public keys file: %w", err)
	}

	manager := &PublicKeyManager{
		Keys: []*rsa.PublicKey{},
	}

	// Parse multiple PEM blocks from the file
	for len(data) > 0 {
		block, rest := pem.Decode(data)
		if block == nil {
			break
		}

		if block.Type != "PUBLIC KEY" {
			data = rest
			continue
		}

		pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse public key: %w", err)
		}

		rsaKey, ok := pubKey.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("public key is not RSA")
		}

		manager.Keys = append(manager.Keys, rsaKey)
		data = rest
	}

	if len(manager.Keys) == 0 {
		return nil, fmt.Errorf("no valid public keys found in file")
	}

	return manager, nil
}

// ValidateToken validates a JWT token against any of the stored public keys
// Returns the token claims if valid, error otherwise
func (pm *PublicKeyManager) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	tokenString = strings.TrimSpace(tokenString)

	if tokenString == "" {
		return nil, fmt.Errorf("empty token")
	}

	var lastErr error

	// Try each public key
	for _, pubKey := range pm.Keys {
		token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (any, error) {
			// Verify signing method is RSA
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return pubKey, nil
		})

		if err != nil {
			lastErr = err
			continue
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return claims, nil
		}

		lastErr = fmt.Errorf("token claims invalid")
	}

	if lastErr != nil {
		return nil, fmt.Errorf("failed to validate token: %w", lastErr)
	}

	return nil, fmt.Errorf("failed to validate token with any public key")
}

// VerifyTokenSignature verifies the token signature directly
func (pm *PublicKeyManager) VerifyTokenSignature(tokenString string) error {
	_, err := pm.ValidateToken(tokenString)
	return err
}
