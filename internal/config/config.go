package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Config holds the application configuration
type Config struct {
	Provider  string `json:"provider"`
	APIKey    string `json:"api_key"`
	Model     string `json:"model"`
	Onboarded bool   `json:"onboarded"`
}

// configDir returns the platform-specific config directory
func configDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".octrafic")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return configDir, nil
}

// configPath returns the full path to the config file
func configPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// Load loads the configuration from disk
func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil // No config yet
		}
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Save saves the configuration to disk
func (c *Config) Save() error {
	path, err := configPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// IsFirstLaunch checks if this is the first launch (no config or not onboarded)
func IsFirstLaunch() (bool, error) {
	config, err := Load()
	if err != nil {
		return true, nil // Treat errors as first launch
	}
	return !config.Onboarded, nil
}

// GetEnvVarName returns the environment variable name for a config key
func GetEnvVarName(key string) string {
	return "OCTRAFIC_" + strings.ToUpper(key)
}

// GetEnv retrieves an environment variable with Octrafic prefix
func GetEnv(key string) string {
	return os.Getenv(GetEnvVarName(key))
}
