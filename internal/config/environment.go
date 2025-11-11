package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Environment represents the application environment
type Environment string

const (
	Development Environment = "development"
	Staging     Environment = "staging"
	Production  Environment = "production"
)

// EnvironmentConfig manages environment-specific configurations
type EnvironmentConfig struct {
	Current Environment
	EnvFile string
}

// LoadEnvironment loads environment-specific configuration
func LoadEnvironment() (*EnvironmentConfig, error) {
	// Get environment from CLI flag, env var, or default
	env := getEnvironment()

	envConfig := &EnvironmentConfig{
		Current: env,
		EnvFile: fmt.Sprintf(".env.%s", env),
	}

	// Load environment file if it exists
	if err := envConfig.loadEnvFile(); err != nil {
		return nil, fmt.Errorf("failed to load environment file: %w", err)
	}

	// Set environment-specific defaults
	if err := envConfig.setEnvironmentDefaults(); err != nil {
		return nil, fmt.Errorf("failed to set environment defaults: %w", err)
	}

	return envConfig, nil
}

// getEnvironment determines the current environment
func getEnvironment() Environment {
	// Check CLI environment flag (set by main.go)
	if env := viper.GetString("environment"); env != "" {
		return Environment(env)
	}

	// Check environment variable
	if env := os.Getenv("FLEEKS_ENVIRONMENT"); env != "" {
		return Environment(env)
	}

	if env := os.Getenv("ENVIRONMENT"); env != "" {
		return Environment(env)
	}

	// Default to development
	return Development
}

// loadEnvFile loads the environment-specific .env file
func (e *EnvironmentConfig) loadEnvFile() error {
	// Get the directory containing the CLI executable or current working directory
	execDir, err := os.Executable()
	if err != nil {
		execDir, _ = os.Getwd()
	} else {
		execDir = filepath.Dir(execDir)
	}

	// Look for .env file in executable directory first
	envPath := filepath.Join(execDir, e.EnvFile)
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		// Fallback to current working directory
		cwd, _ := os.Getwd()
		envPath = filepath.Join(cwd, e.EnvFile)
	}

	// Check if file exists
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		// Environment file doesn't exist, that's OK for production
		if e.Current == Production {
			return nil
		}
		return fmt.Errorf("environment file not found: %s", envPath)
	}

	// Read environment file
	file, err := os.Open(envPath)
	if err != nil {
		return fmt.Errorf("failed to open environment file: %w", err)
	}
	defer file.Close()

	// Parse environment variables
	return e.parseEnvFile(file)
}

// parseEnvFile parses environment variables from file
func (e *EnvironmentConfig) parseEnvFile(file *os.File) error {
	// Read file content
	content, err := os.ReadFile(file.Name())
	if err != nil {
		return err
	}

	// Parse line by line
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE format
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if len(value) >= 2 {
			if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
				(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
				value = value[1 : len(value)-1]
			}
		}

		// Set environment variable
		os.Setenv(key, value)

		// Also set in viper for configuration override
		viperKey := strings.ToLower(strings.ReplaceAll(key, "FLEEKS_", ""))
		viperKey = strings.ReplaceAll(viperKey, "_", ".")
		viper.Set(viperKey, value)
	}

	return nil
}

// setEnvironmentDefaults sets environment-specific default values
func (e *EnvironmentConfig) setEnvironmentDefaults() error {
	switch e.Current {
	case Development:
		return e.setDevelopmentDefaults()
	case Staging:
		return e.setStagingDefaults()
	case Production:
		return e.setProductionDefaults()
	default:
		return fmt.Errorf("unknown environment: %s", e.Current)
	}
}

// setDevelopmentDefaults sets development environment defaults
func (e *EnvironmentConfig) setDevelopmentDefaults() error {
	// API defaults for development
	viper.SetDefault("api.base_url", "http://localhost:8000")
	viper.SetDefault("api.timeout", "30s")
	viper.SetDefault("api.debug", true)
	viper.SetDefault("api.tls_verify", false)

	// WebSocket defaults
	viper.SetDefault("websocket.base_url", "ws://localhost:8000")
	viper.SetDefault("websocket.timeout", "10s")

	// Service endpoints
	viper.SetDefault("services.lsp_url", "http://localhost:8001")
	viper.SetDefault("services.mcp_url", "http://localhost:8002")

	// Development features
	viper.SetDefault("dev.mode", true)
	viper.SetDefault("dev.verbose", true)
	viper.SetDefault("dev.mock_apis", false)
	viper.SetDefault("dev.log_level", "debug")

	return nil
}

// setStagingDefaults sets staging environment defaults
func (e *EnvironmentConfig) setStagingDefaults() error {
	// API defaults for staging
	viper.SetDefault("api.base_url", "https://staging-api.fleeks.dev")
	viper.SetDefault("api.timeout", "45s")
	viper.SetDefault("api.debug", false)
	viper.SetDefault("api.tls_verify", true)

	// WebSocket defaults
	viper.SetDefault("websocket.base_url", "wss://staging-api.fleeks.dev")
	viper.SetDefault("websocket.timeout", "15s")

	// Service endpoints
	viper.SetDefault("services.lsp_url", "https://staging-lsp.fleeks.dev")
	viper.SetDefault("services.mcp_url", "https://staging-mcp.fleeks.dev")

	// Staging features
	viper.SetDefault("dev.mode", false)
	viper.SetDefault("dev.verbose", false)
	viper.SetDefault("dev.log_level", "info")

	return nil
}

// setProductionDefaults sets production environment defaults
func (e *EnvironmentConfig) setProductionDefaults() error {
	// API defaults for production
	viper.SetDefault("api.base_url", "https://api.fleeks.dev")
	viper.SetDefault("api.timeout", "60s")
	viper.SetDefault("api.debug", false)
	viper.SetDefault("api.tls_verify", true)

	// WebSocket defaults
	viper.SetDefault("websocket.base_url", "wss://api.fleeks.dev")
	viper.SetDefault("websocket.timeout", "20s")

	// Service endpoints
	viper.SetDefault("services.lsp_url", "https://lsp.fleeks.dev")
	viper.SetDefault("services.mcp_url", "https://mcp.fleeks.dev")

	// Production features
	viper.SetDefault("dev.mode", false)
	viper.SetDefault("dev.verbose", false)
	viper.SetDefault("dev.log_level", "warn")

	return nil
}

// GetEnvironmentInfo returns information about the current environment
func (e *EnvironmentConfig) GetEnvironmentInfo() map[string]interface{} {
	return map[string]interface{}{
		"environment":   string(e.Current),
		"env_file":      e.EnvFile,
		"api_base_url":  viper.GetString("api.base_url"),
		"ws_base_url":   viper.GetString("websocket.base_url"),
		"lsp_service":   viper.GetString("services.lsp_url"),
		"mcp_service":   viper.GetString("services.mcp_url"),
		"dev_mode":      viper.GetBool("dev.mode"),
		"debug_enabled": viper.GetBool("api.debug"),
		"tls_verify":    viper.GetBool("api.tls_verify"),
	}
}

// String returns the string representation of the environment
func (e Environment) String() string {
	return string(e)
}

// IsValid checks if the environment value is valid
func (e Environment) IsValid() bool {
	switch e {
	case Development, Staging, Production:
		return true
	default:
		return false
	}
}
