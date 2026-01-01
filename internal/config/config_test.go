package config

import (
	"os"
	"testing"
)

func TestLoadConfigValid(t *testing.T) {
	os.Setenv("SERVER_PORT", "9000")
	os.Setenv("SERVER_ENVIRONMENT", "production")
	os.Setenv("GIN_MODE", "release")
	os.Setenv("AUTH_PUBLIC_KEYS_FILE", "/path/to/keys.pem")
	os.Setenv("TLS_ENABLED", "true")
	os.Setenv("TLS_CERT_FILE", "server.crt")
	os.Setenv("TLS_KEY_FILE", "server.key")
	os.Setenv("TLS_CA_CERT_FILE", "ca.crt")
	os.Setenv("TLS_REQUIRE_CLIENT_CERT", "true")
	defer os.Clearenv()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if cfg.Server.Port != "9000" {
		t.Errorf("Expected port 9000, got %s", cfg.Server.Port)
	}
	if cfg.Server.Environment != "production" {
		t.Errorf("Expected production, got %s", cfg.Server.Environment)
	}
	if !cfg.TLS.Enabled {
		t.Errorf("Expected TLS enabled")
	}
	if cfg.TLS.RequireClientCert != true {
		t.Errorf("Expected RequireClientCert true")
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	os.Clearenv()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if cfg.Server.Port != "8080" {
		t.Errorf("Expected default port 8080, got %s", cfg.Server.Port)
	}
	if cfg.Server.Environment != "development" {
		t.Errorf("Expected default development, got %s", cfg.Server.Environment)
	}
	if cfg.Server.GinMode != "debug" {
		t.Errorf("Expected default debug, got %s", cfg.Server.GinMode)
	}
	if cfg.Auth.PublicKeysFile != "public_keys.pem" {
		t.Errorf("Expected default public_keys.pem, got %s", cfg.Auth.PublicKeysFile)
	}
}

func TestLoadConfigProductionMode(t *testing.T) {
	os.Clearenv()
	os.Setenv("SERVER_ENVIRONMENT", "production")
	defer os.Clearenv()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if cfg.Server.GinMode != "release" {
		t.Errorf("Expected release mode for production, got %s", cfg.Server.GinMode)
	}
}
