package auth

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/toby/apple-music-mcp/pkg/config"
)

func TestTokenStorage(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "apple-music-mcp-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tokenPath := filepath.Join(tmpDir, "tokens.json")
	storage := NewTokenStorage(tokenPath)

	// Test saving tokens
	tokenData := &TokenData{
		DeveloperToken: "test-developer-token",
		UserToken:      "test-user-token",
		ExpiresAt:      time.Now().Add(time.Hour),
		CreatedAt:      time.Now(),
	}

	err = storage.SaveTokens(tokenData)
	if err != nil {
		t.Fatalf("Failed to save tokens: %v", err)
	}

	// Test loading tokens
	loadedData, err := storage.LoadTokens()
	if err != nil {
		t.Fatalf("Failed to load tokens: %v", err)
	}

	if loadedData.DeveloperToken != tokenData.DeveloperToken {
		t.Errorf("Expected developer token %s, got %s", tokenData.DeveloperToken, loadedData.DeveloperToken)
	}

	if loadedData.UserToken != tokenData.UserToken {
		t.Errorf("Expected user token %s, got %s", tokenData.UserToken, loadedData.UserToken)
	}
}

func TestAuthManager(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "apple-music-mcp-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &config.Config{
		AppleMusic: config.AppleMusicConfig{
			TeamID:         "TEST123",
			KeyID:          "TEST_KEY",
			PrivateKeyPath: "/nonexistent/path.p8",
		},
		Storage: config.StorageConfig{
			DatabasePath: filepath.Join(tmpDir, "test.db"),
		},
	}

	authManager := NewAuthManager(cfg)
	if authManager == nil {
		t.Fatal("NewAuthManager returned nil")
	}

	// Test that GetClient returns error when no tokens exist
	_, err = authManager.GetClient()
	if err == nil {
		t.Error("Expected error when no tokens exist, but got none")
	}
}