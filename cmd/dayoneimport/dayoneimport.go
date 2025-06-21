package dayoneimport

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// DayOneEntry represents a single Day One journal entry
type DayOneEntry struct {
	Date     time.Time
	Weather  string
	Location string
	Title    string
	Content  string
}

// ParseDayOneExport parses a Day One export file and returns entries
func ParseDayOneExport(filePath string, verbose bool) ([]DayOneEntry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	var entries []DayOneEntry
	var currentEntry DayOneEntry
	var contentLines []string
	inContent := false
	
	scanner := bufio.NewScanner(file)
	lineNum := 0
	
	// Regex patterns for parsing
	dateRegex := regexp.MustCompile(`^\s*Date:\s*(.+)$`)
	weatherRegex := regexp.MustCompile(`^\s*Weather:\s*(.+)$`)
	locationRegex := regexp.MustCompile(`^\s*Location:\s*(.+)$`)
	
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		
		// Check for new entry start (Date line)
		if matches := dateRegex.FindStringSubmatch(line); matches != nil {
			// Save previous entry if we have one
			if !currentEntry.Date.IsZero() {
				currentEntry.Content = strings.Join(contentLines, "\n")
				entries = append(entries, currentEntry)
			}
			
			// Start new entry
			dateStr := strings.TrimSpace(matches[1])
			parsedDate, err := parseDayOneDate(dateStr)
			if err != nil {
				if verbose {
					fmt.Printf("Warning: Could not parse date '%s' at line %d: %v\n", dateStr, lineNum, err)
				}
				continue
			}
			
			currentEntry = DayOneEntry{Date: parsedDate}
			contentLines = []string{}
			inContent = false
			continue
		}
		
		// Check for weather
		if matches := weatherRegex.FindStringSubmatch(line); matches != nil {
			currentEntry.Weather = strings.TrimSpace(matches[1])
			continue
		}
		
		// Check for location
		if matches := locationRegex.FindStringSubmatch(line); matches != nil {
			currentEntry.Location = strings.TrimSpace(matches[1])
			continue
		}
		
		// Empty line after metadata starts content
		if strings.TrimSpace(line) == "" && !inContent {
			inContent = true
			continue
		}
		
		// Collect content lines
		if inContent {
			// First non-empty line after metadata is the title
			if currentEntry.Title == "" && strings.TrimSpace(line) != "" {
				currentEntry.Title = strings.TrimSpace(line)
			} else {
				contentLines = append(contentLines, line)
			}
		}
	}
	
	// Don't forget the last entry
	if !currentEntry.Date.IsZero() {
		currentEntry.Content = strings.Join(contentLines, "\n")
		entries = append(entries, currentEntry)
	}
	
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}
	
	return entries, nil
}

// parseDayOneDate parses Day One's date format
func parseDayOneDate(dateStr string) (time.Time, error) {
	// Day One format: "January 1, 2016 at 21:01:41 EST"
	formats := []string{
		"January 2, 2006 at 15:04:05 MST",
		"January 2, 2006 at 15:04:05 EST",
		"January 2, 2006 at 15:04:05 CST", 
		"January 2, 2006 at 15:04:05 PST",
		"January 2, 2006 at 15:04:05 EDT",
		"January 2, 2006 at 15:04:05 CDT",
		"January 2, 2006 at 15:04:05 PDT",
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}
	
	return time.Time{}, fmt.Errorf("could not parse date format: %s", dateStr)
}

// ConvertToMarkdown converts a Day One entry to markdown format
func (entry DayOneEntry) ConvertToMarkdown() string {
	var sb strings.Builder
	
	// YAML frontmatter
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("date: %s\n", entry.Date.Format("2006-01-02")))
	if entry.Weather != "" {
		sb.WriteString(fmt.Sprintf("weather: %s\n", entry.Weather))
	}
	if entry.Location != "" {
		sb.WriteString(fmt.Sprintf("location: %s\n", entry.Location))
	}
	sb.WriteString("source: Day One\n")
	sb.WriteString("---\n\n")
	
	// Title
	if entry.Title != "" {
		sb.WriteString(fmt.Sprintf("# %s\n\n", entry.Title))
	}
	
	// Content
	sb.WriteString(entry.Content)
	sb.WriteString("\n")
	
	return sb.String()
}

// Execute runs the dayoneimport command
func Execute(folderPath string, verbose bool) error {
	// Check if output folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return fmt.Errorf("output folder does not exist: %s", folderPath)
	}

	// Look for Journal.txt in Desktop export folder
	journalPath := "/Users/travis/Desktop/06-21-2025_15-17-day-one-export/Journal.txt"
	if _, err := os.Stat(journalPath); os.IsNotExist(err) {
		return fmt.Errorf("Day One export file not found: %s", journalPath)
	}

	fmt.Printf("Importing Day One entries from %s to %s...\n\n", journalPath, folderPath)

	// Parse the Day One export
	entries, err := ParseDayOneExport(journalPath, verbose)
	if err != nil {
		return fmt.Errorf("failed to parse Day One export: %v", err)
	}

	fmt.Printf("Found %d Day One entries\n", len(entries))

	// Convert and write entries
	successCount := 0
	for _, entry := range entries {
		filename := fmt.Sprintf("%s-dayone.md", entry.Date.Format("2006-01-02"))
		outputPath := filepath.Join(folderPath, filename)
		
		markdown := entry.ConvertToMarkdown()
		
		if err := os.WriteFile(outputPath, []byte(markdown), 0644); err != nil {
			if verbose {
				fmt.Printf("Warning: Could not write %s: %v\n", outputPath, err)
			}
			continue
		}
		
		successCount++
		if verbose {
			fmt.Printf("Created: %s\n", outputPath)
		}
	}

	fmt.Printf("\nSuccessfully imported %d entries to %s\n", successCount, folderPath)
	return nil
}