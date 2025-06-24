package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/travis-mark/salthaven/cmd/list"
	"github.com/travis-mark/salthaven/cmd/serve"
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

func usage() {
	fmt.Println("Usage: salthaven <command> [folder_path] [options] [args...]")
	fmt.Println("Commands:")
	fmt.Println("  list           List markdown notes matching today's date")
	fmt.Println("  serve          Serve a web page with today's entries")
	fmt.Println("Options:")
	fmt.Println("  -v, --verbose  Enable verbose output (show warnings)")
	fmt.Println("  -p, --port     Port number for serve command (default: 8080)")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "list":
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

		if err := list.Execute(folderPath, verbose); err != nil {
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
		usage()
	}
}
