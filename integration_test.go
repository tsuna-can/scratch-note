package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"scratch-note/config"
	"scratch-note/utils"
)

func TestIntegrationFullWorkflow(t *testing.T) {
	// Setup temporary environment
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", "scratch-note")
	notesDir := filepath.Join(tempDir, "scratch-notes")
	configPath := filepath.Join(configDir, "config.yaml")

	// Create directories
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	err = os.MkdirAll(notesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create notes directory: %v", err)
	}

	// Create config file
	err = config.CreateDefaultConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Test 1: Create note without title
	mockEditor := &MockEditor{}
	filePath1, err := CreateScratchNote("", notesDir, time.Now(), mockEditor)
	if err != nil {
		t.Fatalf("Failed to create scratch note: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(filePath1); os.IsNotExist(err) {
		t.Errorf("Note file was not created: %s", filePath1)
	}

	// Verify filename format
	filename1 := filepath.Base(filePath1)
	if !strings.HasSuffix(filename1, ".md") {
		t.Errorf("Note file should have .md extension: %s", filename1)
	}

	// Test 2: Create note with title
	title := "integration-test-note"
	filePath2, err := CreateScratchNote(title, notesDir, time.Now(), mockEditor)
	if err != nil {
		t.Fatalf("Failed to create titled scratch note: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(filePath2); os.IsNotExist(err) {
		t.Errorf("Titled note file was not created: %s", filePath2)
	}

	// Verify filename contains title
	filename2 := filepath.Base(filePath2)
	if !strings.Contains(filename2, "integration-test-note") {
		t.Errorf("Note filename should contain title: %s", filename2)
	}

	// Test 3: Command parsing integration
	commands := []struct {
		args []string
		cmd  Command
	}{
		{[]string{"scratch-note"}, Command{Type: CommandTypeCreate, Title: ""}},
		{[]string{"scratch-note", "test title"}, Command{Type: CommandTypeCreate, Title: "test title"}},
		{[]string{"scratch-note", "--config"}, Command{Type: CommandTypeConfig}},
		{[]string{"scratch-note", "--help"}, Command{Type: CommandTypeHelp}},
	}

	for _, tc := range commands {
		parsed, err := ParseArgs(tc.args)
		if err != nil {
			t.Errorf("Failed to parse args %v: %v", tc.args, err)
			continue
		}

		if parsed.Type != tc.cmd.Type {
			t.Errorf("Args %v: got type %v, want %v", tc.args, parsed.Type, tc.cmd.Type)
		}

		if parsed.Title != tc.cmd.Title {
			t.Errorf("Args %v: got title %q, want %q", tc.args, parsed.Title, tc.cmd.Title)
		}
	}

	// Test 4: Config loading and path expansion
	loadedConfig, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify config structure
	if loadedConfig.Editor == "" {
		t.Error("Config should have editor set")
	}

	if loadedConfig.ScratchNoteDir == "" {
		t.Error("Config should have scratch note directory set")
	}

	// Test 5: Path expansion
	homeDir, _ := os.UserHomeDir()
	testPath := "~/test-path"
	expanded := config.ExpandPath(testPath)
	expected := filepath.Join(homeDir, "test-path")

	if expanded != expected {
		t.Errorf("Path expansion failed: got %q, want %q", expanded, expected)
	}
}

func TestIntegrationErrorHandling(t *testing.T) {
	tempDir := t.TempDir()

	// Test: Directory doesn't exist
	nonExistentDir := filepath.Join(tempDir, "nonexistent")
	mockEditor := &MockEditor{}

	_, err := CreateScratchNote("", nonExistentDir, time.Now(), mockEditor)
	if err == nil {
		t.Error("Expected error when directory doesn't exist")
	}

	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("Error should mention directory doesn't exist: %v", err)
	}

	// Test: Editor failure
	notesDir := filepath.Join(tempDir, "notes")
	err = os.MkdirAll(notesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create notes directory: %v", err)
	}

	failingEditor := &MockEditor{ShouldFail: true}
	_, err = CreateScratchNote("", notesDir, time.Now(), failingEditor)
	if err == nil {
		t.Error("Expected error when editor fails")
	}

	if _, ok := err.(*EditorError); !ok {
		t.Errorf("Expected EditorError, got %T", err)
	}
}

func TestIntegrationAllPackages(t *testing.T) {
	// This test ensures all packages work together correctly
	tempDir := t.TempDir()
	
	// Test config package
	configPath := filepath.Join(tempDir, "config.yaml")
	err := config.CreateDefaultConfig(configPath)
	if err != nil {
		t.Fatalf("Config package error: %v", err)
	}

	_, err = config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Config loading error: %v", err)
	}

	// Test utils package
	testTime := time.Date(2025, 8, 16, 14, 30, 45, 0, time.UTC)
	filename := utils.GenerateFileName("test", testTime)
	expected := "2025-08-16_143045_test.md"
	
	if filename != expected {
		t.Errorf("Utils integration failed: got %q, want %q", filename, expected)
	}

	t.Logf("Integration test completed - all packages work together")
}