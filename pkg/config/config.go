package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	AppleMusic AppleMusicConfig `yaml:"apple_music"`
	Storage    StorageConfig    `yaml:"storage"`
}

// AppleMusicConfig holds Apple Music API credentials
type AppleMusicConfig struct {
	TeamID         string `yaml:"team_id"`
	KeyID          string `yaml:"key_id"`
	PrivateKeyPath string `yaml:"private_key_path"`
}

// StorageConfig holds storage configuration
type StorageConfig struct {
	DatabasePath string `yaml:"database_path"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	return &Config{
		Storage: StorageConfig{
			DatabasePath: filepath.Join(homeDir, ".local", "share", "apple-music-mcp", "data.db"),
		},
	}
}

// Load reads configuration from the default location
func Load() (*Config, error) {
	return LoadFromPath(DefaultConfigPath())
}

// DefaultConfigPath returns the default configuration file path
func DefaultConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "apple-music-mcp", "config.yaml")
}

// LoadFromPath reads configuration from the specified path
func LoadFromPath(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Fill in defaults for missing values
	defaultConfig := DefaultConfig()
	if config.Storage.DatabasePath == "" {
		config.Storage.DatabasePath = defaultConfig.Storage.DatabasePath
	}

	return &config, nil
}

// Save writes the configuration to the default location
func (c *Config) Save() error {
	return c.SaveToPath(DefaultConfigPath())
}

// SaveToPath writes the configuration to the specified path
func (c *Config) SaveToPath(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.AppleMusic.TeamID == "" {
		return fmt.Errorf("apple_music.team_id is required")
	}
	if c.AppleMusic.KeyID == "" {
		return fmt.Errorf("apple_music.key_id is required")
	}
	if c.AppleMusic.PrivateKeyPath == "" {
		return fmt.Errorf("apple_music.private_key_path is required")
	}

	// Check if private key file exists
	if _, err := os.Stat(c.AppleMusic.PrivateKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("private key file not found: %s", c.AppleMusic.PrivateKeyPath)
	}

	return nil
}