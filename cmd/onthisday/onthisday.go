package onthisday

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// parseYAMLDate extracts and parses the date from YAML frontmatter
func parseYAMLDate(content string) (time.Time, error) {
	// Check if content starts with YAML frontmatter
	if !strings.HasPrefix(content, "---") {
		return time.Time{}, fmt.Errorf("no YAML frontmatter found")
	}

	lines := strings.Split(content, "\n")
	inFrontmatter := false

	for i, line := range lines {
		line = strings.TrimSpace(line)

		if i == 0 && line == "---" {
			inFrontmatter = true
			continue
		}

		if inFrontmatter && line == "---" {
			break
		}

		if inFrontmatter {
			// Look for date property (supports various formats)
			dateRegex := regexp.MustCompile(`^date\s*:\s*(.+)$`)
			matches := dateRegex.FindStringSubmatch(line)

			if len(matches) > 1 {
				dateStr := strings.TrimSpace(matches[1])
				// Remove quotes if present
				dateStr = strings.Trim(dateStr, `"'`)

				// Try different date formats
				formats := []string{
					"2006-01-02",          // YYYY-MM-DD
					"2006-01-02T15:04",    // ISO 8601 without seconds
					"2006-01-02T15:04:05", // ISO 8601 without timezone
					"2006-01-02 15:04:05", // YYYY-MM-DD HH:MM:SS
					"01/02/2006",          // MM/DD/YYYY
					"02/01/2006",          // DD/MM/YYYY
					"January 2, 2006",     // Month DD, YYYY
					"Jan 2, 2006",         // Mon DD, YYYY
					"{{date}}T{{time}}",   // Template format
				}

				for _, format := range formats {
					if date, err := time.Parse(format, dateStr); err == nil {
						return date, nil
					}
				}

				return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
			}
		}
	}

	return time.Time{}, fmt.Errorf("no date property found in YAML frontmatter")
}

// readFileContent reads the entire content of a file
func readFileContent(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// scanMarkdownNotesOnThisDay scans the specified folder for markdown notes with the same month and day (any year)
func scanMarkdownNotesOnThisDay(folderPath string) ([]string, error) {
	var onThisDayNotes []string
	today := time.Now()

	err := filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Only process markdown files
		if !strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			return nil
		}

		// Read file content
		content, err := readFileContent(path)
		if err != nil {
			fmt.Printf("Warning: Could not read file %s: %v\n", path, err)
			return nil // Continue processing other files
		}

		// Parse date from YAML frontmatter
		fileDate, err := parseYAMLDate(content)
		if err != nil {
			fmt.Printf("Warning: Could not parse date from %s: %v\n", path, err)
			return nil // Continue processing other files
		}

		// Check if the month and day match today (any year)
		if fileDate.Month() == today.Month() && fileDate.Day() == today.Day() {
			onThisDayNotes = append(onThisDayNotes, path)
		}

		return nil
	})

	return onThisDayNotes, err
}

// Execute runs the onthisday command
func Execute(folderPath string) error {
	// Check if folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return fmt.Errorf("folder does not exist: %s", folderPath)
	}

	today := time.Now()
	fmt.Printf("Scanning folder '%s' for markdown notes dated %s (any year)...\n\n", 
		folderPath, today.Format("January 2"))

	// Scan for notes on this day
	notes, err := scanMarkdownNotesOnThisDay(folderPath)
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