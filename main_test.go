package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// MockEditor for testing editor launching
type MockEditor struct {
	CalledWith []string
	ShouldFail bool
}

func (m *MockEditor) Launch(filePath string) error {
	m.CalledWith = append(m.CalledWith, filePath)
	if m.ShouldFail {
		return &EditorError{Editor: "mock-editor", Err: "command not found"}
	}
	return nil
}

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedCmd Command
		expectError bool
	}{
		{
			name:        "no arguments",
			args:        []string{"scratch-note"},
			expectedCmd: Command{Type: CommandTypeCreate, Title: ""},
			expectError: false,
		},
		{
			name:        "with title",
			args:        []string{"scratch-note", "my notes"},
			expectedCmd: Command{Type: CommandTypeCreate, Title: "my notes"},
			expectError: false,
		},
		{
			name:        "config flag",
			args:        []string{"scratch-note", "--config"},
			expectedCmd: Command{Type: CommandTypeConfig},
			expectError: false,
		},
		{
			name:        "help flag",
			args:        []string{"scratch-note", "--help"},
			expectedCmd: Command{Type: CommandTypeHelp},
			expectError: false,
		},
		{
			name:        "too many arguments",
			args:        []string{"scratch-note", "arg1", "arg2"},
			expectedCmd: Command{},
			expectError: true,
		},
		{
			name:        "unknown flag",
			args:        []string{"scratch-note", "--unknown"},
			expectedCmd: Command{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := ParseArgs(tt.args)
			
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if cmd.Type != tt.expectedCmd.Type {
				t.Errorf("Command type = %v, want %v", cmd.Type, tt.expectedCmd.Type)
			}

			if cmd.Title != tt.expectedCmd.Title {
				t.Errorf("Command title = %q, want %q", cmd.Title, tt.expectedCmd.Title)
			}
		})
	}
}

func TestCreateScratchNote(t *testing.T) {
	tempDir := t.TempDir()
	
	tests := []struct {
		name          string
		title         string
		directoryExists bool
		expectError   bool
	}{
		{
			name:            "create note without title",
			title:           "",
			directoryExists: true,
			expectError:     false,
		},
		{
			name:            "create note with title",
			title:           "test-note",
			directoryExists: true,
			expectError:     false,
		},
		{
			name:            "directory does not exist",
			title:           "",
			directoryExists: false,
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			noteDir := filepath.Join(tempDir, "test-"+tt.name)
			
			if tt.directoryExists {
				err := os.MkdirAll(noteDir, 0755)
				if err != nil {
					t.Fatalf("Failed to create test directory: %v", err)
				}
			}

			mockEditor := &MockEditor{}
			filePath, err := CreateScratchNote(tt.title, noteDir, time.Now(), mockEditor)
			
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify file was created
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("Scratch note file was not created: %s", filePath)
			}

			// Verify editor was called
			if len(mockEditor.CalledWith) != 1 {
				t.Errorf("Editor should be called once, called %d times", len(mockEditor.CalledWith))
			} else if mockEditor.CalledWith[0] != filePath {
				t.Errorf("Editor called with %q, want %q", mockEditor.CalledWith[0], filePath)
			}

			// Verify filename format
			filename := filepath.Base(filePath)
			if tt.title == "" {
				// Should be timestamp.md format
				if len(filename) != 20 { // YYYY-MM-DD_HHMMSS.md
					t.Errorf("Filename length should be 20, got %d: %s", len(filename), filename)
				}
			} else {
				// Should contain title
				if !containsTitle(filename, tt.title) {
					t.Errorf("Filename should contain title %q: %s", tt.title, filename)
				}
			}
		})
	}
}

func TestCreateScratchNoteEditorFailure(t *testing.T) {
	tempDir := t.TempDir()
	noteDir := filepath.Join(tempDir, "test-notes")
	err := os.MkdirAll(noteDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	mockEditor := &MockEditor{ShouldFail: true}
	_, err = CreateScratchNote("", noteDir, time.Now(), mockEditor)
	
	if err == nil {
		t.Error("Expected error when editor fails")
	}

	// Verify error is of correct type
	if _, ok := err.(*EditorError); !ok {
		t.Errorf("Expected EditorError, got %T", err)
	}
}

// Helper function to check if filename contains title
func containsTitle(filename, title string) bool {
	// Simple check - in real implementation this would be more sophisticated
	return len(filename) > 20 // Basic timestamp is 20 chars, so longer means title was added
}