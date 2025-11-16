package config

import (
	"fmt"
	"os"
)

// Config holds all application configuration.
type Config struct {
	Server  ServerConfig
	AWS     AWSConfig
	Cognito CognitoConfig
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Host string
	Port string
}

// AWSConfig holds AWS-specific configuration.
type AWSConfig struct {
	Region  string
	Profile string
}

// CognitoConfig holds AWS Cognito configuration.
type CognitoConfig struct {
	Region       string
	UserPoolID   string
	ClientID     string
	ClientSecret string
}

// Load loads configuration from environment variables with defaults.
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host: getEnvOrDefault("SERVER_HOST", "localhost"),
			Port: getEnvOrDefault("SERVER_PORT", "8080"),
		},
		AWS: AWSConfig{
			Region:  getEnvOrDefault("AWS_REGION", "us-east-1"),
			Profile: getEnvOrDefault("AWS_PROFILE", ""),
		},
		Cognito: CognitoConfig{
			Region:       getEnvOrDefault("AWS_COGNITO_REGION", getEnvOrDefault("AWS_REGION", "us-east-1")),
			UserPoolID:   os.Getenv("AWS_COGNITO_USER_POOL_ID"),
			ClientID:     os.Getenv("AWS_COGNITO_CLIENT_ID"),
			ClientSecret: os.Getenv("AWS_COGNITO_CLIENT_SECRET"),
		},
	}

	// Validate configuration
	if cfg.Server.Port == "" {
		return nil, fmt.Errorf("SERVER_PORT is required")
	}

	// Validate Cognito configuration
	if cfg.Cognito.UserPoolID == "" {
		return nil, fmt.Errorf("AWS_COGNITO_USER_POOL_ID is required")
	}
	if cfg.Cognito.ClientID == "" {
		return nil, fmt.Errorf("AWS_COGNITO_CLIENT_ID is required")
	}
	if cfg.Cognito.ClientSecret == "" {
		return nil, fmt.Errorf("AWS_COGNITO_CLIENT_SECRET is required")
	}

	return cfg, nil
}

// getEnvOrDefault returns the value of an environment variable or a default value.
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
