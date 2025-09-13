# Scratch Note

A simple terminal-based note-taking tool written in Go that creates timestamped markdown files and opens them in your preferred editor.

## Features

- 📝 Creates timestamped markdown files automatically
- ⚡ Quick note creation from the command line
- 🔧 Configurable editor and storage directory
- 🎯 Optional note titles for better organization
- 🖥️ Cross-platform support (Linux, macOS, Windows)

## Installation

### From Source

```bash
git clone https://github.com/yourusername/scratch-note.git
cd scratch-note
make install
```

### Build Manually

```bash
go build -o scratch-note .
sudo cp scratch-note /usr/local/bin/
```

## Usage

### Basic Commands

```bash
# Create a note with current timestamp
scratch-note

# Create a note with a title
scratch-note "shopping list"

# Edit configuration file
scratch-note --config
```

### File Naming Convention

- Basic format: `2025-08-16_143045.md` (YYYY-MM-DD_HHMMSS.md)
- With title: `2025-08-16_143045_shopping-list.md`

## Configuration

Configuration file is located at `~/.config/scratch-note/config.yaml`:

```yaml
scratch-note_dir: "~/scratch-notes"    # Directory to store notes
editor: "nvim"                         # Editor to use (default: vi)
```

### First Run

On first run, if no configuration file exists, you'll be prompted to create one:

```bash
scratch-note --config
```

## Directory Structure

```
~/.config/scratch-note/
└── config.yaml

~/scratch-notes/
├── 2025-08-16_143045.md
├── 2025-08-16_144520_meeting-notes.md
└── 2025-08-16_150030_todo.md
```

## Development

### Prerequisites

- Go 1.19 or later
- Make (optional, for using Makefile commands)

### Build

```bash
make build          # Build binary
make test           # Run tests
make test-coverage  # Run tests with coverage report
make cross-compile  # Build for multiple platforms
```

### Testing

The project follows Test-Driven Development (TDD) practices:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
make test

# Generate coverage report
make test-coverage
```

### Cross-Platform Builds

```bash
make cross-compile
```

This creates binaries for:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

## Error Handling

The tool provides clear error messages for common issues:

- **Missing directory**: `Error: scratch-note directory does not exist: /path/to/dir`
- **No config file**: `Error: Config file not found. Run 'scratch-note --config' to create one.`
- **Editor not found**: `Error: Editor 'nvim' not found in PATH`
- **Invalid config**: `Error: Invalid config file format`

## Success Messages

```
Created scratch-note: /path/to/scratch-notes/2025-08-16_143045.md
Created scratch-note: /path/to/scratch-notes/2025-08-16_143045_shopping-list.md
```

## Project Structure

```
scratch-note/
├── main.go                 # Main application logic
├── main_test.go           # Main application tests
├── config/
│   ├── config.go          # Configuration management
│   └── config_test.go     # Configuration tests
├── utils/
│   ├── file.go            # File operations utilities
│   └── file_test.go       # File utilities tests
├── integration_test.go    # End-to-end tests
├── Makefile              # Build and development commands
├── go.mod                # Go module definition
└── go.sum                # Go module checksums
```
