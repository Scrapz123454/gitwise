package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aymenhmaidiwastaken/gitwise/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show or initialize configuration",
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		data, err := yaml.Marshal(cfg)
		if err != nil {
			return err
		}
		fmt.Print(string(data))
		return nil
	},
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a default configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		home, _ := os.UserHomeDir()
		configDir := filepath.Join(home, ".config", "gitwise")
		configPath := filepath.Join(configDir, "config.yaml")

		if _, err := os.Stat(configPath); err == nil {
			return fmt.Errorf("config file already exists at %s", configPath)
		}

		if err := os.MkdirAll(configDir, 0o755); err != nil {
			return err
		}

		cfg := config.DefaultConfig()
		data, err := yaml.Marshal(cfg)
		if err != nil {
			return err
		}

		if err := os.WriteFile(configPath, data, 0o644); err != nil {
			return err
		}

		fmt.Printf("Created config at %s\n", configPath)
		return nil
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configInitCmd)
	rootCmd.AddCommand(configCmd)
}
