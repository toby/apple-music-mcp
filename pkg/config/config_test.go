package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	if config == nil {
		t.Fatal("DefaultConfig returned nil")
	}
	
	homeDir, _ := os.UserHomeDir()
	expectedPath := filepath.Join(homeDir, ".local", "share", "apple-music-mcp", "data.db")
	
	if config.Storage.DatabasePath != expectedPath {
		t.Errorf("Expected database path %s, got %s", expectedPath, config.Storage.DatabasePath)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		expectErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				AppleMusic: AppleMusicConfig{
					TeamID:         "TEST123",
					KeyID:          "KEY123",
					PrivateKeyPath: "/tmp/test.p8",
				},
			},
			expectErr: true, // Will fail because file doesn't exist
		},
		{
			name: "missing team_id",
			config: &Config{
				AppleMusic: AppleMusicConfig{
					KeyID:          "KEY123",
					PrivateKeyPath: "/tmp/test.p8",
				},
			},
			expectErr: true,
		},
		{
			name: "missing key_id",
			config: &Config{
				AppleMusic: AppleMusicConfig{
					TeamID:         "TEST123",
					PrivateKeyPath: "/tmp/test.p8",
				},
			},
			expectErr: true,
		},
		{
			name: "missing private_key_path",
			config: &Config{
				AppleMusic: AppleMusicConfig{
					TeamID: "TEST123",
					KeyID:  "KEY123",
				},
			},
			expectErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}