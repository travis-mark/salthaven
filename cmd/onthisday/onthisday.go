package onthisday

import (
	"fmt"
	"os"
	"time"

	"github.com/travis-mark/salthaven/internal/markdown"
)

// Execute runs the onthisday command
func Execute(folderPath string) error {
	// Check if folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return fmt.Errorf("folder does not exist: %s", folderPath)
	}

	today := time.Now()
	fmt.Printf("Scanning folder '%s' for markdown notes dated %s (any year)...\n\n", 
		folderPath, today.Format("January 2"))

	// Scan for notes on this day using the same day matcher
	notes, err := markdown.ScanMarkdownNotes(folderPath, markdown.SameDayMatcher, today)
	if err != nil {
		return fmt.Errorf("error scanning folder: %v", err)
	}

	// Display results
	if len(notes) == 0 {
		fmt.Printf("No markdown notes with date %s were found.\n", today.Format("January 2"))
	} else {
		fmt.Printf("Found %d markdown note(s) dated %s:\n", len(notes), today.Format("January 2"))
		for i, note := range notes {
			fmt.Printf("%d. %s\n", i+1, note)
		}
	}

	return nil
}