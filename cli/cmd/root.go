package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gohttpx-cli",
	Short: "GoHTTPx CLI - Manage users and notes via REST API",
	Long: `GoHTTPx CLI is a command-line interface for the GoHTTPx REST API.
It allows you to manage users and notes, configure API endpoints, and interact
with your GoHTTPx server from the command line.

Examples:
  gohttpx-cli config set-url http://localhost:8080
  gohttpx-cli user create --name "John Doe"
  gohttpx-cli note create --title "Meeting" --content "Team standup notes"
  gohttpx-cli user list
  gohttpx-cli note list`,
	Version: "1.0.0",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().Bool("json", false, "Output in JSON format")
	rootCmd.PersistentFlags().Bool("verbose", false, "Enable verbose output")
	rootCmd.PersistentFlags().String("base-url", "", "Override the configured base URL")
}
