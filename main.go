package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/travis-mark/salthaven/cmd/onthisday"
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
		fmt.Println("Usage: salthaven <command> [args...]")
		fmt.Println("Commands:")
		fmt.Println("  today [folder_path]     - Find markdown notes with today's date")
		fmt.Println("  onthisday [folder_path] - Find markdown notes with today's month/day (any year)")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "today":
		folderPath := getDefaultFolderPath()

		// Check if folder path is provided as command line argument
		if len(os.Args) > 2 {
			folderPath = os.Args[2]
		}

		if err := today.Execute(folderPath); err != nil {
			log.Fatal(err)
		}
	case "onthisday":
		folderPath := getDefaultFolderPath()

		// Check if folder path is provided as command line argument
		if len(os.Args) > 2 {
			folderPath = os.Args[2]
		}

		if err := onthisday.Execute(folderPath); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Available commands:")
		fmt.Println("  today [folder_path]     - Find markdown notes with today's date")
		fmt.Println("  onthisday [folder_path] - Find markdown notes with today's month/day (any year)")
		os.Exit(1)
	}
}
