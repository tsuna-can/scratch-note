package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"scratch-note/config"
	"scratch-note/utils"
)

// CommandType represents the type of command to execute
type CommandType int

const (
	CommandTypeCreate CommandType = iota
	CommandTypeConfig
	CommandTypeHelp
)

// Command represents a parsed command
type Command struct {
	Type  CommandType
	Title string
}

// EditorLauncher interface for launching editors
type EditorLauncher interface {
	Launch(filePath string) error
}

// RealEditor implements EditorLauncher for real editor execution
type RealEditor struct {
	EditorName string
}

func (r *RealEditor) Launch(filePath string) error {
	cmd := exec.Command(r.EditorName, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	err := cmd.Run()
	if err != nil {
		return &EditorError{Editor: r.EditorName, Err: err.Error()}
	}
	return nil
}

// EditorError represents an error when launching the editor
type EditorError struct {
	Editor string
	Err    string
}

func (e *EditorError) Error() string {
	return fmt.Sprintf("Editor '%s' not found: %s", e.Editor, e.Err)
}

// ParseArgs parses command line arguments
func ParseArgs(args []string) (Command, error) {
	if len(args) == 1 {
		return Command{Type: CommandTypeCreate, Title: ""}, nil
	}

	if len(args) == 2 {
		switch args[1] {
		case "--config":
			return Command{Type: CommandTypeConfig}, nil
		case "--help", "-h":
			return Command{Type: CommandTypeHelp}, nil
		default:
			// Check if it's an unknown flag
			if args[1][0] == '-' {
				return Command{}, fmt.Errorf("unknown flag: %s", args[1])
			}
			return Command{Type: CommandTypeCreate, Title: args[1]}, nil
		}
	}

	if len(args) > 2 {
		return Command{}, fmt.Errorf("too many arguments")
	}

	return Command{}, fmt.Errorf("unknown command")
}

// CreateScratchNote creates a new scratch note file and opens it in editor
func CreateScratchNote(title, directory string, t time.Time, editor EditorLauncher) (string, error) {
	// Check if directory exists
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return "", fmt.Errorf("scratch-note directory does not exist: %s", directory)
	}

	// Generate filename
	filename := utils.GenerateFileName(title, t)
	filePath := filepath.Join(directory, filename)

	// Create empty file
	err := os.WriteFile(filePath, []byte(""), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}

	// Launch editor
	err = editor.Launch(filePath)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func main() {
	cmd, err := ParseArgs(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		printUsage()
		os.Exit(1)
	}

	switch cmd.Type {
	case CommandTypeHelp:
		printUsage()
	case CommandTypeConfig:
		handleConfigCommand()
	case CommandTypeCreate:
		handleCreateCommand(cmd.Title)
	}
}

func printUsage() {
	fmt.Println("scratch-note - A simple terminal-based note-taking tool")
	fmt.Println("")
	fmt.Println("USAGE:")
	fmt.Println("  scratch-note                    Create new timestamped note")
	fmt.Println("  scratch-note \"title\"            Create note with custom title")
	fmt.Println("  scratch-note --config           Edit configuration file")
	fmt.Println("  scratch-note --help             Show this help message")
	fmt.Println("")
	fmt.Println("EXAMPLES:")
	fmt.Println("  scratch-note                    # Creates: 2025-08-16_143045.md")
	fmt.Println("  scratch-note \"meeting notes\"    # Creates: 2025-08-16_143045_meeting-notes.md")
	fmt.Println("")
	fmt.Println("CONFIGURATION:")
	fmt.Println("  Config file: ~/.config/scratch-note/config.yaml")
	fmt.Println("  Run 'scratch-note --config' to create or edit configuration")
	fmt.Println("")
	fmt.Println("For more information, visit: https://github.com/your-repo/scratch-note")
}

// getConfigPath returns the path to the configuration file
func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home dir is not available
		return "config.yaml"
	}
	return filepath.Join(homeDir, ".config", "scratch-note", "config.yaml")
}

func handleConfigCommand() {
	configPath := getConfigPath()
	
	// Check if config file exists, create if not
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("Config file not found. Creating default config at: %s\n", configPath)
		err := config.CreateDefaultConfig(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to create config file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Default config file created successfully.")
	}
	
	// Load config to get editor
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to load config file: %v\n", err)
		os.Exit(1)
	}
	
	// Use configured editor or default to vi
	editorName := cfg.Editor
	if editorName == "" {
		editorName = "vi"
	}
	
	// Launch editor to edit config
	editor := &RealEditor{EditorName: editorName}
	err = editor.Launch(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func handleCreateCommand(title string) {
	configPath := getConfigPath()
	
	// Load config or use defaults
	var cfg *config.Config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Config file doesn't exist, prompt user to create one
		fmt.Printf("Error: Config file not found. Run 'scratch-note --config' to create one.\n")
		os.Exit(1)
	}
	
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid config file format: %v\n", err)
		os.Exit(1)
	}
	
	// Expand path and check if directory exists
	notesDir := config.ExpandPath(cfg.ScratchNoteDir)
	if _, err := os.Stat(notesDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: scratch-note directory does not exist: %s\n", notesDir)
		os.Exit(1)
	}
	
	// Use configured editor or default to vi
	editorName := cfg.Editor
	if editorName == "" {
		editorName = "vi"
	}
	
	// Create scratch note
	editor := &RealEditor{EditorName: editorName}
	filePath, err := CreateScratchNote(title, notesDir, time.Now(), editor)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Created scratch-note: %s\n", filePath)
}