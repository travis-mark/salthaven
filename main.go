package main

import (
	"fmt"
	"log"
	"os"

	"github.com/travis-mark/salthaven/cmd/today"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: salthaven <command> [args...]")
		fmt.Println("Commands:")
		fmt.Println("  today [folder_path]  - Find markdown notes with today's date")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "today":
		// Default folder to scan (current directory)
		folderPath := "."

		// Check if folder path is provided as command line argument
		if len(os.Args) > 2 {
			folderPath = os.Args[2]
		}

		if err := today.Execute(folderPath); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Available commands:")
		fmt.Println("  today [folder_path]  - Find markdown notes with today's date")
		os.Exit(1)
	}
}