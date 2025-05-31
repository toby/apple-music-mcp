package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "apple-music-mcp",
	Short: "Apple Music MCP CLI tool",
	Long:  `A command line tool for Apple Music authentication and account management.`,
}

func init() {
	// Add global flags here if needed
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}