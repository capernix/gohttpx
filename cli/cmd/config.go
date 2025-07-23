package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type Config struct {
	BaseURL   string `json:"base_url"`
	AuthToken string `json:"auth_token"`
}

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "gohttpx", "config.json"), nil
}

func LoadConfig() Config {
	cfg, err := Load()
	if err != nil {
		return Config{BaseURL: "http://localhost"}
	}
	return *cfg
}

func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return &Config{BaseURL: "http://localhost:8080"}, nil
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	return enc.Encode(cfg)
}

func Print(cfg *Config) {
	fmt.Println("CLI Configuration:")
	fmt.Printf("Base URL:   %s\n", cfg.BaseURL)
	if cfg.AuthToken != "" {
		fmt.Println("Auth Token: [set]")
	} else {
		fmt.Println("Auth Token: (none)")
	}
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
	Long:  "Configure API endpoint and other settings for the GoHTTPx CLI",
}

var configSetUrlCmd = &cobra.Command{
	Use:   "set-url",
	Short: "Set the API base URL",
	Long:  "Set the base URL for the GoHTTPx API",
	Example: `  gohttpx-cli config set-url http://localhost:8080
  gohttpx-cli config set-url https://api.example.com`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		baseURL := args[0]

		cfg := LoadConfig()
		cfg.BaseURL = baseURL

		if err := Save(&cfg); err != nil {
			PrintError(fmt.Sprintf("Failed to save configuration: %v", err))
			return
		}

		PrintSuccess(fmt.Sprintf("Base URL set to %s", baseURL))
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  "Display the current CLI configuration settings",
	Example: `  gohttpx-cli config show
  gohttpx-cli config show --json`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonOutput, _ := cmd.Flags().GetBool("json")
		cfg := LoadConfig()

		if jsonOutput {
			output, _ := json.MarshalIndent(cfg, "", "  ")
			fmt.Println(string(output))
		} else {
			Print(&cfg)
		}
	},
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset configuration to defaults",
	Long:  "Reset all configuration settings to their default values",
	Example: `  gohttpx-cli config reset
  gohttpx-cli config reset --confirm`,
	Run: func(cmd *cobra.Command, args []string) {
		skipConfirm, _ := cmd.Flags().GetBool("confirm")

		if !skipConfirm {
			if !ConfirmAction("Are you sure you want to reset all configuration to defaults?") {
				fmt.Println("Reset cancelled")
				return
			}
		}

		defaultCfg := Config{BaseURL: "http://localhost"}
		if err := Save(&defaultCfg); err != nil {
			PrintError(fmt.Sprintf("Failed to reset configuration: %v", err))
			return
		}

		PrintSuccess("Configuration reset to defaults")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configSetUrlCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configResetCmd)

	configResetCmd.Flags().BoolP("confirm", "y", false, "Skip confirmation prompt")
}
