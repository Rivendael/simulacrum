package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	"simulacrum/internal/auth"
	"simulacrum/internal/config"
	"simulacrum/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration from environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Load public keys for JWT validation
	pkm, err := auth.LoadPublicKeysFromFile(cfg.Auth.PublicKeysFile)
	if err != nil {
		log.Fatalf("Failed to load public keys: %v", err)
	}

	r := gin.Default()

	// Apply JWT middleware to /obscure endpoint
	r.POST("/obscure", auth.JWTMiddleware(pkm), handlers.HandleObscure)

	// Health check endpoint (no auth required)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	fmt.Printf("Server starting on :%s (%s)...\n", cfg.Server.Port, cfg.Server.Environment)
	fmt.Printf("Using public keys from: %s\n", cfg.Auth.PublicKeysFile)
	fmt.Println("Endpoints:")
	fmt.Println("  GET  /health        - Health check (no auth)")
	fmt.Println("  POST /obscure       - Obscure data (requires JWT)")

	// Setup TLS if enabled
	if cfg.TLS.Enabled && cfg.TLS.CertFile != "" && cfg.TLS.KeyFile != "" {
		tlsConfig := &tls.Config{
			MinVersion: parseTLSVersion(cfg.TLS.MinVersion),
		}

		// Add client certificate verification if required
		if cfg.TLS.RequireClientCert && cfg.TLS.CACertFile != "" {
			caCertPEM, err := os.ReadFile(cfg.TLS.CACertFile)
			if err != nil {
				log.Fatalf("Failed to read CA certificate file: %v", err)
			}

			caCertPool := x509.NewCertPool()
			if !caCertPool.AppendCertsFromPEM(caCertPEM) {
				log.Fatal("Failed to parse CA certificate")
			}

			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
			tlsConfig.ClientCAs = caCertPool
			fmt.Printf("mTLS enabled - client certificates required\n")
		} else {
			fmt.Printf("HTTPS enabled (without client verification)\n")
		}

		if err := r.RunTLS(":"+cfg.Server.Port, cfg.TLS.CertFile, cfg.TLS.KeyFile); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := r.Run(":" + cfg.Server.Port); err != nil {
			log.Fatal(err)
		}
	}
}

// parseTLSVersion converts string version to tls.Version constant
func parseTLSVersion(version string) uint16 {
	switch version {
	case "1.0":
		return tls.VersionTLS10
	case "1.1":
		return tls.VersionTLS11
	case "1.2":
		return tls.VersionTLS12
	case "1.3":
		return tls.VersionTLS13
	default:
		return tls.VersionTLS12
	}
}
