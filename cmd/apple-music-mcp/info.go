package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/toby/apple-music-mcp/pkg/auth"
	"github.com/toby/apple-music-mcp/pkg/config"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display authenticated user information",
	Long:  `Display information about the authenticated Apple Music user account.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		// Create auth manager
		authManager := auth.NewAuthManager(cfg)

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Get user info
		userInfo, err := authManager.GetUserInfo(ctx)
		if err != nil {
			return fmt.Errorf("failed to get user info: %w", err)
		}

		// Format and display the information
		jsonData, err := json.MarshalIndent(userInfo, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format user info: %w", err)
		}

		fmt.Println("Apple Music Account Information:")
		fmt.Println(string(jsonData))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}