package today

import (
	"fmt"
	"os"
	"time"

	"github.com/travis-mark/salthaven/internal/markdown"
)

// Execute runs the today command
func Execute(folderPath string) error {
	// Check if folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return fmt.Errorf("folder does not exist: %s", folderPath)
	}

	fmt.Printf("Scanning folder '%s' for markdown notes dated today...\n\n", folderPath)

	// Scan for today's notes using the exact date matcher
	today := time.Now()
	notes, err := markdown.ScanMarkdownNotes(folderPath, markdown.ExactDateMatcher, today)
	if err != nil {
		return fmt.Errorf("error scanning folder: %v", err)
	}

	// Display results
	if len(notes) == 0 {
		fmt.Println("No markdown notes with today's date were found.")
	} else {
		fmt.Printf("Found %d markdown note(s) dated today:\n", len(notes))
		for i, note := range notes {
			fmt.Printf("%d. %s\n", i+1, note)
		}
	}

	return nil
}