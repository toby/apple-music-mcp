package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/toby/apple-music-mcp/pkg/auth"
	"github.com/toby/apple-music-mcp/pkg/config"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with Apple Music API",
	Long:  `Authenticate with Apple Music API using your developer credentials.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		// Validate configuration
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("configuration validation failed: %w", err)
		}

		// Create auth manager and authenticate
		authManager := auth.NewAuthManager(cfg)
		if err := authManager.Authenticate(); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}