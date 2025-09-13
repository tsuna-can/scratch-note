package config

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	ScratchNoteDir string `yaml:"scratch_note_dir"`
	Editor         string `yaml:"editor"`
}

// LoadConfig loads configuration from the specified file path
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// GetDefaultConfig returns the default configuration
func GetDefaultConfig() *Config {
	return &Config{
		ScratchNoteDir: "~/scratch-notes",
		Editor:         "vi",
	}
}

// ExpandPath expands ~ to home directory in path
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path // Return original if can't get home dir
		}
		return filepath.Join(homeDir, path[2:])
	}
	return path
}

// CreateDefaultConfig creates a default configuration file
func CreateDefaultConfig(configPath string) error {
	config := GetDefaultConfig()
	
	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		return err
	}
	
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	
	return os.WriteFile(configPath, data, 0644)
}