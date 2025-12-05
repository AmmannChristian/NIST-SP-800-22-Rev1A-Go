package config

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all service configuration
type Config struct {
	// gRPC server configuration
	GRPCPort int

	// TLS configuration for gRPC
	TLSEnabled    bool
	TLSCertFile   string
	TLSKeyFile    string
	TLSCAFile     string
	TLSClientAuth string
	TLSMinVersion string

	// Metrics server configuration
	MetricsPort int

	// Logging configuration
	LogLevel string

	// Authentication configuration
	AuthEnabled  bool
	AuthIssuer   string
	AuthAudience string
	AuthJWKSURL  string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		GRPCPort:      getEnvInt("GRPC_PORT", 9090),
		TLSEnabled:    getEnvBool("TLS_ENABLED", false),
		TLSCertFile:   getEnvString("TLS_CERT_FILE", ""),
		TLSKeyFile:    getEnvString("TLS_KEY_FILE", ""),
		TLSCAFile:     getEnvString("TLS_CA_FILE", ""),
		TLSClientAuth: getEnvString("TLS_CLIENT_AUTH", "none"),
		TLSMinVersion: getEnvString("TLS_MIN_VERSION", "1.2"),
		MetricsPort:   getEnvInt("METRICS_PORT", 9091),
		LogLevel:      getEnvString("LOG_LEVEL", "info"),
		AuthEnabled:   getEnvBool("AUTH_ENABLED", false),
		AuthIssuer:    getEnvString("AUTH_ISSUER", ""),
		AuthAudience:  getEnvString("AUTH_AUDIENCE", ""),
		AuthJWKSURL:   getEnvString("AUTH_JWKS_URL", ""),
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.GRPCPort < 1 || c.GRPCPort > 65535 {
		return fmt.Errorf("invalid GRPC_PORT: %d (must be 1-65535)", c.GRPCPort)
	}

	if c.MetricsPort < 1 || c.MetricsPort > 65535 {
		return fmt.Errorf("invalid METRICS_PORT: %d (must be 1-65535)", c.MetricsPort)
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("invalid LOG_LEVEL: %s (must be debug/info/warn/error)", c.LogLevel)
	}

	if c.AuthEnabled {
		if c.AuthIssuer == "" {
			return fmt.Errorf("invalid AUTH_ISSUER: required when AUTH_ENABLED=true")
		}
		if c.AuthAudience == "" {
			return fmt.Errorf("invalid AUTH_AUDIENCE: required when AUTH_ENABLED=true")
		}
	}

	if c.TLSEnabled {
		if c.TLSCertFile == "" {
			return fmt.Errorf("invalid TLS_CERT_FILE: required when TLS_ENABLED=true")
		}
		if c.TLSKeyFile == "" {
			return fmt.Errorf("invalid TLS_KEY_FILE: required when TLS_ENABLED=true")
		}
		if _, err := parseTLSClientAuth(c.TLSClientAuth); err != nil {
			return err
		}
		if _, err := parseTLSMinVersion(c.TLSMinVersion); err != nil {
			return err
		}
	}

	return nil
}

// TLSClientAuthType returns the parsed tls.ClientAuthType from configuration.
func (c *Config) TLSClientAuthType() (tls.ClientAuthType, error) {
	return parseTLSClientAuth(c.TLSClientAuth)
}

// TLSMinVersionValue returns the configured minimum TLS version (defaults to TLS 1.2).
func (c *Config) TLSMinVersionValue() (uint16, error) {
	return parseTLSMinVersion(c.TLSMinVersion)
}

func parseTLSClientAuth(mode string) (tls.ClientAuthType, error) {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "", "none", "noclientcert":
		return tls.NoClientCert, nil
	case "request", "requestclientcert":
		return tls.RequestClientCert, nil
	case "requireany", "requireanyclientcert":
		return tls.RequireAnyClientCert, nil
	case "verifyifgiven", "verify_client_cert_if_given":
		return tls.VerifyClientCertIfGiven, nil
	case "requireandverify", "requireandverifyclientcert", "mtls":
		return tls.RequireAndVerifyClientCert, nil
	default:
		return tls.NoClientCert, fmt.Errorf("invalid TLS_CLIENT_AUTH: %s", mode)
	}
}

func parseTLSMinVersion(version string) (uint16, error) {
	switch strings.ToLower(strings.TrimSpace(version)) {
	case "", "default", "1.2", "tls1.2", "tls12":
		return tls.VersionTLS12, nil
	case "1.3", "tls1.3", "tls13":
		return tls.VersionTLS13, nil
	default:
		return 0, fmt.Errorf("invalid TLS_MIN_VERSION: %s (use 1.2 or 1.3)", version)
	}
}

// getEnvString reads a string from environment variable or returns default
func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt reads an integer from environment variable or returns default
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// getEnvBool reads a boolean from environment variable or returns default
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}
