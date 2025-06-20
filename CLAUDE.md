# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Salthaven is a Go CLI tool that scans directories for markdown files containing today's date in their YAML frontmatter. The tool parses various date formats and lists notes created on the current day.

## Development Commands

### Building and Running
```bash
# Build the project
go build -o salthaven

# Run the today command (default: current directory or SALTHAVEN_FOLDER env var)
./salthaven today

# Run the onthisday command (finds notes with same month/day, any year)
./salthaven onthisday

# Run with specific folder (overrides env var)
./salthaven today /path/to/notes
./salthaven onthisday /path/to/notes

# Set default folder via environment variable
export SALTHAVEN_FOLDER=/path/to/your/notes
./salthaven today  # Will use /path/to/your/notes

# Or create a .env file in the current directory
echo "SALTHAVEN_FOLDER=/path/to/your/notes" > .env
./salthaven today  # Will use /path/to/your/notes

# Run directly without building
go run main.go today [folder_path]
go run main.go onthisday [folder_path]
```

### Testing and Quality
```bash
# Run tests (when added)
go test ./...

# Format code
go fmt ./...

# Lint code (requires golangci-lint)
golangci-lint run

# Check for vulnerabilities
go mod tidy && go list -json -m all | nancy sleuth
```

## Architecture

### Core Components

- **main.go**: Command router handling subcommand dispatch
- **cmd/today/**: Today subcommand package for finding notes with today's exact date
- **cmd/onthisday/**: OnThisDay subcommand package for finding notes with same month/day (any year)
- **Shared functionality across commands**:
  - **parseYAMLDate()**: Parses dates from YAML frontmatter with multiple format support
  - **scanMarkdownNotes*()**: Recursively walks directories to find markdown files
  - **readFileContent()**: File reading utility
  - **Execute()**: Main entry point for each command

### Key Design Patterns

- **Subcommand Architecture**: Modular command structure for extensibility
- **Date Format Flexibility**: Supports ISO 8601, US/European formats, and natural language dates
- **Error Tolerance**: Continues processing when individual files fail to parse
- **Recursive Directory Walking**: Uses filepath.WalkDir for efficient directory traversal

### Date Format Support
The tool recognizes these date formats in YAML frontmatter:
- `2006-01-02` (YYYY-MM-DD)
- `2006-01-02T15:04:05` (ISO 8601 variants)
- `01/02/2006` and `02/01/2006` (US/European)
- `January 2, 2006` and `Jan 2, 2006` (Named months)

## Environment Variables

- **SALTHAVEN_FOLDER**: Default folder path to scan for markdown files
  - Priority: 1. Command line argument, 2. Environment variable, 3. .env file, 4. Current directory
  - Can be set via: `export SALTHAVEN_FOLDER=/path/to/notes` or in `.env` file
  - .env file format: `SALTHAVEN_FOLDER=/path/to/notes`
  - Supports quoted values and comments (lines starting with #)

## Module Information

- **Module**: `github.com/travis-mark/salthaven`
- **Go Version**: 1.23.4
- **Dependencies**: Standard library only