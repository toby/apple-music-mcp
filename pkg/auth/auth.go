package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/minchao/go-apple-music"
	"github.com/toby/apple-music-mcp/pkg/config"
)

// TokenStorage handles storing and retrieving tokens
type TokenStorage struct {
	databasePath string
}

// NewTokenStorage creates a new token storage instance
func NewTokenStorage(databasePath string) *TokenStorage {
	return &TokenStorage{databasePath: databasePath}
}

// TokenData represents stored token information
type TokenData struct {
	DeveloperToken string    `json:"developer_token"`
	UserToken      string    `json:"user_token,omitempty"`
	ExpiresAt      time.Time `json:"expires_at"`
	CreatedAt      time.Time `json:"created_at"`
}

// SaveTokens saves the tokens to storage
func (ts *TokenStorage) SaveTokens(data *TokenData) error {
	dir := filepath.Dir(ts.databasePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create storage directory: %w", err)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token data: %w", err)
	}

	if err := os.WriteFile(ts.databasePath, jsonData, 0600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}

// LoadTokens loads tokens from storage
func (ts *TokenStorage) LoadTokens() (*TokenData, error) {
	data, err := os.ReadFile(ts.databasePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no tokens found, please authenticate first")
		}
		return nil, fmt.Errorf("failed to read token file: %w", err)
	}

	var tokenData TokenData
	if err := json.Unmarshal(data, &tokenData); err != nil {
		return nil, fmt.Errorf("failed to parse token file: %w", err)
	}

	return &tokenData, nil
}

// AuthManager handles Apple Music authentication
type AuthManager struct {
	config  *config.Config
	storage *TokenStorage
}

// NewAuthManager creates a new authentication manager
func NewAuthManager(cfg *config.Config) *AuthManager {
	storage := NewTokenStorage(cfg.Storage.DatabasePath)
	return &AuthManager{
		config:  cfg,
		storage: storage,
	}
}

// generateDeveloperToken creates a JWT developer token for Apple Music API
func (am *AuthManager) generateDeveloperToken() (string, error) {
	// Read the private key
	keyData, err := os.ReadFile(am.config.AppleMusic.PrivateKeyPath)
	if err != nil {
		return "", fmt.Errorf("failed to read private key: %w", err)
	}

	// Parse the PEM-encoded private key
	block, _ := pem.Decode(keyData)
	if block == nil {
		return "", fmt.Errorf("failed to decode PEM block")
	}

	// Parse the ECDSA private key
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	ecdsaKey, ok := privateKey.(*ecdsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("private key is not ECDSA")
	}

	// Create JWT token
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": am.config.AppleMusic.TeamID,
		"iat": now.Unix(),
		"exp": now.Add(180 * 24 * time.Hour).Unix(), // 6 months
	})

	// Set the key ID in the header
	token.Header["kid"] = am.config.AppleMusic.KeyID

	// Sign the token
	tokenString, err := token.SignedString(ecdsaKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// Authenticate performs the authentication flow
func (am *AuthManager) Authenticate() error {
	// Generate developer token
	developerToken, err := am.generateDeveloperToken()
	if err != nil {
		return fmt.Errorf("failed to generate developer token: %w", err)
	}

	// For this basic implementation, we'll just store the developer token
	// In a full implementation, this would involve OAuth flow for user token
	tokenData := &TokenData{
		DeveloperToken: developerToken,
		ExpiresAt:      time.Now().Add(180 * 24 * time.Hour), // 6 months
		CreatedAt:      time.Now(),
	}

	if err := am.storage.SaveTokens(tokenData); err != nil {
		return fmt.Errorf("failed to save tokens: %w", err)
	}

	fmt.Println("Authentication successful! Developer token generated and saved.")
	fmt.Println("Note: This is a basic implementation with developer token only.")
	fmt.Println("For full functionality, user authentication would be required.")

	return nil
}

// GetClient returns an authenticated Apple Music API client
func (am *AuthManager) GetClient() (*applemusic.Client, error) {
	tokenData, err := am.storage.LoadTokens()
	if err != nil {
		return nil, err
	}

	// Check if token is expired
	if time.Now().After(tokenData.ExpiresAt) {
		return nil, fmt.Errorf("token expired, please authenticate again")
	}

	// Create transport with the developer token
	transport := &applemusic.Transport{
		Token: tokenData.DeveloperToken,
	}

	// If we have a user token, add it
	if tokenData.UserToken != "" {
		transport.MusicUserToken = tokenData.UserToken
	}

	// Create the Apple Music client
	client := applemusic.NewClient(transport.Client())

	return client, nil
}

// GetUserInfo fetches information about the authenticated user
func (am *AuthManager) GetUserInfo(ctx context.Context) (map[string]interface{}, error) {
	client, err := am.GetClient()
	if err != nil {
		return nil, err
	}

	// Try to get user's storefront (this is available with just developer token)
	storefronts, resp, err := client.Me.GetStorefront(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get user storefront: %w", err)
	}

	info := map[string]interface{}{
		"authentication_status": "authenticated",
		"token_type":           "developer_token",
	}

	if storefronts != nil && len(storefronts.Data) > 0 {
		storefront := storefronts.Data[0]
		info["storefront"] = map[string]interface{}{
			"id":                  storefront.Id,
			"type":               storefront.Type,
			"name":               storefront.Attributes.Name,
			"default_language":   storefront.Attributes.DefaultLanguageTag,
			"supported_languages": storefront.Attributes.SupportedLanguageTags,
		}
	}

	if resp != nil {
		info["api_response_status"] = resp.Status
	}

	// Get token info
	tokenData, err := am.storage.LoadTokens()
	if err == nil {
		info["token_created"] = tokenData.CreatedAt
		info["token_expires"] = tokenData.ExpiresAt
		info["has_user_token"] = tokenData.UserToken != ""
	}

	return info, nil
}