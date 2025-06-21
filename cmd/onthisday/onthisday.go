package onthisday

import (
	"fmt"
	"os"
	"time"

	"github.com/travis-mark/salthaven/internal/markdown"
)

// Execute runs the onthisday command
func Execute(folderPath string, verbose bool) error {
	// Check if folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return fmt.Errorf("folder does not exist: %s", folderPath)
	}
	// Scan for notes on this day using the same day matcher
	today := time.Now()
	notes, err := markdown.ScanMarkdownNotes(folderPath, markdown.SameDayMatcher, today, verbose)
	if err != nil {
		return fmt.Errorf("error scanning folder: %v", err)
	}
	// Check for results
	if len(notes) == 0 {
		return fmt.Errorf("no notes found")
	}
	// Display results
	for _, note := range notes {
		fmt.Printf("%s\n", note)
	}

	return nil
}
