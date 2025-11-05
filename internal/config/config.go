package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	TodoDir string `yaml:"todo_dir"`
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".tada", "config.yml"), nil
}

// Load reads the config file and returns a Config struct
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// If config doesn't exist, return empty config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save writes the config to disk
func Save(cfg *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// GetTodoDir returns the todo directory from config
// Returns an error if config is not set
func GetTodoDir() (string, error) {
	cfg, err := Load()
	if err != nil {
		return "", err
	}

	if cfg.TodoDir == "" {
		return "", fmt.Errorf("todo directory not configured")
	}

	return cfg.TodoDir, nil
}

// GetTodoFilePath returns the full path to the todo.txt file
func GetTodoFilePath() (string, error) {
	dir, err := GetTodoDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "todo.txt"), nil
}
