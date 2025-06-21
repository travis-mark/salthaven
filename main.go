package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/travis-mark/salthaven/cmd/dayoneimport"
	"github.com/travis-mark/salthaven/cmd/onthisday"
	"github.com/travis-mark/salthaven/cmd/serve"
	"github.com/travis-mark/salthaven/cmd/today"
)

// loadEnvFile loads environment variables from a .env file
func loadEnvFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return // Silently ignore if .env file doesn't exist
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE format
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Remove quotes if present
			if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'')) {
				value = value[1 : len(value)-1]
			}

			// Only set if not already set in environment
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
	}
}

// getDefaultFolderPath returns the default folder path to scan
// Priority: 1. SALTHAVEN_FOLDER env var, 2. .env file, 3. current directory
func getDefaultFolderPath() string {
	// Load .env file if it exists (only affects environment if var not already set)
	loadEnvFile(".env")

	if envPath := os.Getenv("SALTHAVEN_FOLDER"); envPath != "" {
		return envPath
	}
	return "."
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: salthaven <command> [options] [args...]")
		fmt.Println("Commands:")
		fmt.Println("  today [-v|--verbose] [folder_path]       - Find markdown notes with today's date")
		fmt.Println("  onthisday [-v|--verbose] [folder_path]   - Find markdown notes with today's month/day (any year)")
		fmt.Println("  dayoneimport [-v|--verbose] [folder_path] - Import Day One entries to markdown")
		fmt.Println("  serve [-v|--verbose] [-p|--port PORT] [folder_path] - Serve a web page with today's entries")
		fmt.Println("Options:")
		fmt.Println("  -v, --verbose    - Enable verbose output (show warnings)")
		fmt.Println("  -p, --port       - Port number for serve command (default: 8080)")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "today":
		folderPath := getDefaultFolderPath()
		verbose := false

		// Parse arguments
		for i := 2; i < len(os.Args); i++ {
			arg := os.Args[i]
			if arg == "-v" || arg == "--verbose" {
				verbose = true
			} else {
				folderPath = arg
			}
		}

		if err := today.Execute(folderPath, verbose); err != nil {
			log.Fatal(err)
		}
	case "onthisday":
		folderPath := getDefaultFolderPath()
		verbose := false

		// Parse arguments
		for i := 2; i < len(os.Args); i++ {
			arg := os.Args[i]
			if arg == "-v" || arg == "--verbose" {
				verbose = true
			} else {
				folderPath = arg
			}
		}

		if err := onthisday.Execute(folderPath, verbose); err != nil {
			log.Fatal(err)
		}
	case "dayoneimport":
		folderPath := getDefaultFolderPath()
		verbose := false

		// Parse arguments
		for i := 2; i < len(os.Args); i++ {
			arg := os.Args[i]
			if arg == "-v" || arg == "--verbose" {
				verbose = true
			} else {
				folderPath = arg
			}
		}

		if err := dayoneimport.Execute(folderPath, verbose); err != nil {
			log.Fatal(err)
		}
	case "serve":
		folderPath := getDefaultFolderPath()
		verbose := false
		port := 8080

		// Parse arguments
		for i := 2; i < len(os.Args); i++ {
			arg := os.Args[i]
			if arg == "-v" || arg == "--verbose" {
				verbose = true
			} else if arg == "-p" || arg == "--port" {
				if i+1 < len(os.Args) {
					if p, err := strconv.Atoi(os.Args[i+1]); err == nil {
						port = p
						i++ // Skip the port number argument
					}
				}
			} else {
				folderPath = arg
			}
		}

		if err := serve.Execute(folderPath, verbose, port); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Available commands:")
		fmt.Println("  today [-v|--verbose] [folder_path]       - Find markdown notes with today's date")
		fmt.Println("  onthisday [-v|--verbose] [folder_path]   - Find markdown notes with today's month/day (any year)")
		fmt.Println("  dayoneimport [-v|--verbose] [folder_path] - Import Day One entries to markdown")
		fmt.Println("  serve [-v|--verbose] [-p|--port PORT] [folder_path] - Serve a web page with today's entries")
		os.Exit(1)
	}
}
