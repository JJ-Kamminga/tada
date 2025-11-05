package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"tada/internal/config"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage tada configuration",
	Long:  `Configure tada settings like todo file location.`,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long:  `Set a configuration value. Available keys: file`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		cfg, err := config.Load()
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}

		switch key {
		case "file":
			// Expand home directory if present
			if value[:2] == "~/" {
				home, err := os.UserHomeDir()
				if err != nil {
					fmt.Println("Error getting home directory:", err)
					os.Exit(1)
				}
				value = filepath.Join(home, value[2:])
			}

			// Convert to absolute path
			absPath, err := filepath.Abs(value)
			if err != nil {
				fmt.Println("Error resolving path:", err)
				os.Exit(1)
			}

			cfg.TodoFile = absPath
			fmt.Printf("Set todo file location to: %s\n", absPath)
		default:
			fmt.Printf("Unknown config key: %s\n", key)
			fmt.Println("Available keys: file")
			os.Exit(1)
		}

		if err := config.Save(cfg); err != nil {
			fmt.Println("Error saving config:", err)
			os.Exit(1)
		}

		configPath, _ := config.GetConfigPath()
		fmt.Printf("Configuration saved to: %s\n", configPath)
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a configuration value",
	Long:  `Get a configuration value. If no key is specified, shows all config. Available keys: file`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}

		if len(args) == 0 {
			// Show all config
			configPath, _ := config.GetConfigPath()
			fmt.Printf("Configuration file: %s\n\n", configPath)
			if cfg.TodoFile != "" {
				fmt.Printf("file: %s\n", cfg.TodoFile)
			} else {
				home, _ := os.UserHomeDir()
				defaultPath := filepath.Join(home, ".tada", "todo.txt")
				fmt.Printf("file: %s (default)\n", defaultPath)
			}
		} else {
			key := args[0]
			switch key {
			case "file":
				if cfg.TodoFile != "" {
					fmt.Println(cfg.TodoFile)
				} else {
					home, _ := os.UserHomeDir()
					defaultPath := filepath.Join(home, ".tada", "todo.txt")
					fmt.Printf("%s (default)\n", defaultPath)
				}
			default:
				fmt.Printf("Unknown config key: %s\n", key)
				fmt.Println("Available keys: file")
				os.Exit(1)
			}
		}
	},
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show the config file path",
	Long:  `Display the path to the tada configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, err := config.GetConfigPath()
		if err != nil {
			fmt.Println("Error getting config path:", err)
			os.Exit(1)
		}
		fmt.Println(configPath)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configPathCmd)
}
