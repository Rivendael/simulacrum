package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server ServerConfig
	Auth   AuthConfig
	TLS    TLSConfig
}

type ServerConfig struct {
	Port        string
	Environment string
	GinMode     string
}

type AuthConfig struct {
	PublicKeysFile string
	Enabled        bool
}

type TLSConfig struct {
	Enabled           bool
	CertFile          string
	KeyFile           string
	CACertFile        string
	RequireClientCert bool
	MinVersion        string
}

func LoadConfig() (*Config, error) {
	var cfg Config

	// Load from environment variables
	if v := os.Getenv("SERVER_PORT"); v != "" {
		cfg.Server.Port = v
	}
	if v := os.Getenv("SERVER_ENVIRONMENT"); v != "" {
		cfg.Server.Environment = v
	}
	if v := os.Getenv("GIN_MODE"); v != "" {
		cfg.Server.GinMode = v
	}
	if v := os.Getenv("AUTH_PUBLIC_KEYS_FILE"); v != "" {
		cfg.Auth.PublicKeysFile = v
	}
	if v := os.Getenv("TLS_ENABLED"); v != "" {
		if boolVal, err := strconv.ParseBool(v); err == nil {
			cfg.TLS.Enabled = boolVal
		}
	}
	if v := os.Getenv("TLS_CERT_FILE"); v != "" {
		cfg.TLS.CertFile = v
	}
	if v := os.Getenv("TLS_KEY_FILE"); v != "" {
		cfg.TLS.KeyFile = v
	}
	if v := os.Getenv("TLS_CA_CERT_FILE"); v != "" {
		cfg.TLS.CACertFile = v
	}
	if v := os.Getenv("TLS_REQUIRE_CLIENT_CERT"); v != "" {
		if boolVal, err := strconv.ParseBool(v); err == nil {
			cfg.TLS.RequireClientCert = boolVal
		}
	}

	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}
	if cfg.Server.Environment == "" {
		cfg.Server.Environment = "development"
	}
	if cfg.Server.GinMode == "" {
		if cfg.Server.Environment == "production" {
			cfg.Server.GinMode = "release"
		} else {
			cfg.Server.GinMode = "debug"
		}
	}
	if cfg.Auth.PublicKeysFile == "" {
		cfg.Auth.PublicKeysFile = "public_keys.pem"

	}
	if cfg.TLS.MinVersion == "" {
		cfg.TLS.MinVersion = "1.2"
	}

	return &cfg, nil
}
