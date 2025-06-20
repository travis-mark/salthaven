package markdown

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// ParseYAMLDate extracts and parses the date from YAML frontmatter
func ParseYAMLDate(content string) (time.Time, error) {
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

// ReadFileContent reads the entire content of a file
func ReadFileContent(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// DateMatcher is a function type that determines if a file date matches the criteria
type DateMatcher func(fileDate, referenceDate time.Time) bool

// ScanMarkdownNotes scans the specified folder for markdown notes matching the date criteria
func ScanMarkdownNotes(folderPath string, matcher DateMatcher, referenceDate time.Time) ([]string, error) {
	var matchingNotes []string

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
		content, err := ReadFileContent(path)
		if err != nil {
			fmt.Printf("Warning: Could not read file %s: %v\n", path, err)
			return nil // Continue processing other files
		}

		// Parse date from YAML frontmatter
		fileDate, err := ParseYAMLDate(content)
		if err != nil {
			fmt.Printf("Warning: Could not parse date from %s: %v\n", path, err)
			return nil // Continue processing other files
		}

		// Check if the date matches using the provided matcher
		if matcher(fileDate, referenceDate) {
			matchingNotes = append(matchingNotes, path)
		}

		return nil
	})

	return matchingNotes, err
}

// ExactDateMatcher returns true if the file date exactly matches the reference date (same year, month, day)
func ExactDateMatcher(fileDate, referenceDate time.Time) bool {
	return fileDate.Year() == referenceDate.Year() && 
		   fileDate.Month() == referenceDate.Month() && 
		   fileDate.Day() == referenceDate.Day()
}

// SameDayMatcher returns true if the file date has the same month and day (any year)
func SameDayMatcher(fileDate, referenceDate time.Time) bool {
	return fileDate.Month() == referenceDate.Month() && 
		   fileDate.Day() == referenceDate.Day()
}