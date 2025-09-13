package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name           string
		configContent  string
		expectError    bool
		expectedConfig Config
	}{
		{
			name: "valid config",
			configContent: `scratch_note_dir: "~/scratch-notes"
editor: "nvim"`,
			expectError: false,
			expectedConfig: Config{
				ScratchNoteDir: "~/scratch-notes",
				Editor:         "nvim",
			},
		},
		{
			name: "config with only editor",
			configContent: `editor: "vim"`,
			expectError: false,
			expectedConfig: Config{
				ScratchNoteDir: "",
				Editor:         "vim",
			},
		},
		{
			name: "config with only directory",
			configContent: `scratch_note_dir: "/home/user/notes"`,
			expectError: false,
			expectedConfig: Config{
				ScratchNoteDir: "/home/user/notes",
				Editor:         "",
			},
		},
		{
			name:          "invalid yaml",
			configContent: `invalid: yaml: content: [`,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary config file
			tempDir := t.TempDir()
			configPath := filepath.Join(tempDir, "config.yaml")
			
			err := os.WriteFile(configPath, []byte(tt.configContent), 0644)
			if err != nil {
				t.Fatalf("Failed to write test config file: %v", err)
			}

			config, err := LoadConfig(configPath)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if config.ScratchNoteDir != tt.expectedConfig.ScratchNoteDir {
				t.Errorf("ScratchNoteDir = %q, want %q", config.ScratchNoteDir, tt.expectedConfig.ScratchNoteDir)
			}

			if config.Editor != tt.expectedConfig.Editor {
				t.Errorf("Editor = %q, want %q", config.Editor, tt.expectedConfig.Editor)
			}
		})
	}
}

func TestLoadConfigFileNotExists(t *testing.T) {
	tempDir := t.TempDir()
	nonExistentPath := filepath.Join(tempDir, "nonexistent.yaml")
	
	_, err := LoadConfig(nonExistentPath)
	if err == nil {
		t.Error("Expected error for non-existent config file")
	}
}

func TestGetDefaultConfig(t *testing.T) {
	config := GetDefaultConfig()
	
	if config.Editor != "vi" {
		t.Errorf("Default editor should be 'vi', got: %q", config.Editor)
	}
	
	if config.ScratchNoteDir == "" {
		t.Error("Default scratch note directory should not be empty")
	}
	
	// Check that default directory contains home directory reference
	if config.ScratchNoteDir != "~/scratch-notes" {
		t.Errorf("Default directory should be '~/scratch-notes', got: %q", config.ScratchNoteDir)
	}
}

func TestExpandPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "tilde expansion",
			input:    "~/scratch-notes",
			expected: filepath.Join(os.Getenv("HOME"), "scratch-notes"),
		},
		{
			name:     "absolute path",
			input:    "/absolute/path",
			expected: "/absolute/path",
		},
		{
			name:     "relative path",
			input:    "relative/path",
			expected: "relative/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExpandPath(tt.input)
			if result != tt.expected {
				t.Errorf("ExpandPath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCreateDefaultConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	
	err := CreateDefaultConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to create default config: %v", err)
	}
	
	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Default config file was not created")
	}
	
	// Verify content by loading it
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load created config: %v", err)
	}
	
	defaultConfig := GetDefaultConfig()
	if config.Editor != defaultConfig.Editor {
		t.Errorf("Created config editor = %q, want %q", config.Editor, defaultConfig.Editor)
	}
	
	if config.ScratchNoteDir != defaultConfig.ScratchNoteDir {
		t.Errorf("Created config directory = %q, want %q", config.ScratchNoteDir, defaultConfig.ScratchNoteDir)
	}
}