package serve

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/travis-mark/salthaven/internal/markdown"
)

// NoteEntry represents a markdown note with its metadata
type NoteEntry struct {
	Path     string
	FullPath string
	Date     time.Time
	Title    string
	Content  string
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>On This Day - {{.FormattedDate}}</title>
    <style>
        :root {
            --bg-primary: #f9f9f9;
            --bg-secondary: white;
            --text-primary: #333;
            --text-secondary: #7f8c8d;
            --text-tertiary: #95a5a6;
            --text-accent: #2c3e50;
            --text-content: #34495e;
            --border-color: #eee;
            --shadow: rgba(0,0,0,0.1);
        }

        @media (prefers-color-scheme: dark) {
            :root {
                --bg-primary: #1a1a1a;
                --bg-secondary: #2d2d2d;
                --text-primary: #e0e0e0;
                --text-secondary: #b0b0b0;
                --text-tertiary: #888;
                --text-accent: #64b5f6;
                --text-content: #d0d0d0;
                --border-color: #444;
                --shadow: rgba(0,0,0,0.3);
            }
        }

        [data-theme="dark"] {
            --bg-primary: #1a1a1a;
            --bg-secondary: #2d2d2d;
            --text-primary: #e0e0e0;
            --text-secondary: #b0b0b0;
            --text-tertiary: #888;
            --text-accent: #64b5f6;
            --text-content: #d0d0d0;
            --border-color: #444;
            --shadow: rgba(0,0,0,0.3);
        }

        [data-theme="light"] {
            --bg-primary: #f9f9f9;
            --bg-secondary: white;
            --text-primary: #333;
            --text-secondary: #7f8c8d;
            --text-tertiary: #95a5a6;
            --text-accent: #2c3e50;
            --text-content: #34495e;
            --border-color: #eee;
            --shadow: rgba(0,0,0,0.1);
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            line-height: 1.6;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            color: var(--text-primary);
            background-color: var(--bg-primary);
            transition: background-color 0.3s ease, color 0.3s ease;
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
            padding: 20px;
            background: var(--bg-secondary);
            border-radius: 8px;
            box-shadow: 0 2px 4px var(--shadow);
            position: relative;
            transition: background-color 0.3s ease, box-shadow 0.3s ease;
        }
        .header h1 {
            color: var(--text-accent);
            margin: 0;
            transition: color 0.3s ease;
        }
        .header p {
            color: var(--text-secondary);
            margin: 10px 0 0 0;
            transition: color 0.3s ease;
        }
        .theme-toggle {
            position: absolute;
            top: 20px;
            right: 20px;
            background: none;
            border: 2px solid var(--text-tertiary);
            border-radius: 50%;
            width: 40px;
            height: 40px;
            cursor: pointer;
            font-size: 18px;
            display: flex;
            align-items: center;
            justify-content: center;
            transition: all 0.3s ease;
            color: var(--text-tertiary);
        }
        .theme-toggle:hover {
            border-color: var(--text-accent);
            color: var(--text-accent);
            transform: scale(1.1);
        }
        .note {
            background: var(--bg-secondary);
            margin: 20px 0;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px var(--shadow);
            transition: background-color 0.3s ease, box-shadow 0.3s ease;
        }
        .note-header {
            border-bottom: 1px solid var(--border-color);
            padding-bottom: 10px;
            margin-bottom: 15px;
            transition: border-color 0.3s ease;
        }
        .note-title {
            font-size: 1.3em;
            font-weight: bold;
            color: var(--text-accent);
            margin: 0;
            transition: color 0.3s ease;
        }
        .note-title-link {
            color: var(--text-accent);
            text-decoration: none;
            transition: all 0.3s ease;
        }
        .note-title-link:hover {
            text-decoration: underline;
            opacity: 0.8;
        }
        .note-date {
            color: var(--text-secondary);
            font-size: 0.9em;
            margin: 5px 0;
            transition: color 0.3s ease;
        }
        .note-path {
            color: var(--text-tertiary);
            font-size: 0.8em;
            font-family: monospace;
            transition: color 0.3s ease;
        }
        .note-content {
            white-space: pre-wrap;
            color: var(--text-content);
            transition: color 0.3s ease;
        }
        .no-notes {
            text-align: center;
            color: var(--text-secondary);
            font-style: italic;
            padding: 40px;
            background: var(--bg-secondary);
            border-radius: 8px;
            box-shadow: 0 2px 4px var(--shadow);
            transition: all 0.3s ease;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            padding: 20px;
            color: var(--text-tertiary);
            font-size: 0.9em;
            transition: color 0.3s ease;
        }
        .footer a {
            color: var(--text-accent);
            text-decoration: none;
            transition: color 0.3s ease;
        }
        .footer a:hover {
            text-decoration: underline;
        }
        
        /* Checkbox styling */
        .checkbox-item {
            display: flex;
            align-items: flex-start;
            margin: 2px 0;
            line-height: 1.5;
        }
        .checkbox {
            width: 12px;
            height: 12px;
            margin-right: 6px;
            margin-top: 3px;
            border: 1px solid var(--text-tertiary);
            border-radius: 2px;
            background: var(--bg-secondary);
            flex-shrink: 0;
            display: flex;
            align-items: center;
            justify-content: center;
            transition: background-color 0.2s ease;
        }
        .checkbox.checked {
            background: var(--text-accent);
            border-color: var(--text-accent);
        }
        .checkbox.checked::after {
            content: '✓';
            color: var(--bg-secondary);
            font-size: 9px;
            font-weight: bold;
            line-height: 1;
        }
        .checkbox-text {
            flex: 1;
            color: var(--text-content);
        }
    </style>
</head>
<body>
    <div class="header">
        <button class="theme-toggle" onclick="toggleTheme()" title="Toggle dark/light mode">
            <span class="theme-icon">🌙</span>
        </button>
        <h1>On This Day</h1>
        <p>{{.FormattedDate}} • {{.Count}} {{if eq .Count 1}}entry{{else}}entries{{end}} found</p>
    </div>

    {{if .Notes}}
        {{range .Notes}}
        <div class="note">
            <div class="note-header">
                {{if .Title}}
                <h2 class="note-title">
                    <a href="obsidian://open?path={{.FullPath}}" class="note-title-link">{{.Title}}</a>
                </h2>
                {{end}}
                <div class="note-date">{{.Date.Format "January 2, 2006"}}</div>
                <div class="note-path">{{.Path}}</div>
            </div>
            <div class="note-content">{{.Content}}</div>
        </div>
        {{end}}
    {{else}}
        <div class="no-notes">
            No notes found for this day
        </div>
    {{end}}

    <div class="footer">
        Generated by Salthaven • <a href="javascript:location.reload()">Refresh</a>
    </div>

    <script>
        // Theme management
        function getStoredTheme() {
            return localStorage.getItem('theme');
        }

        function setStoredTheme(theme) {
            localStorage.setItem('theme', theme);
        }

        function getPreferredTheme() {
            const storedTheme = getStoredTheme();
            if (storedTheme) {
                return storedTheme;
            }
            return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
        }

        function setTheme(theme) {
            document.documentElement.setAttribute('data-theme', theme);
            const themeIcon = document.querySelector('.theme-icon');
            if (themeIcon) {
                themeIcon.textContent = theme === 'dark' ? '☀️' : '🌙';
            }
        }

        function toggleTheme() {
            const currentTheme = document.documentElement.getAttribute('data-theme');
            const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
            setTheme(newTheme);
            setStoredTheme(newTheme);
        }

        // Initialize theme on page load
        document.addEventListener('DOMContentLoaded', function() {
            const preferredTheme = getPreferredTheme();
            setTheme(preferredTheme);
        });

        // Listen for system theme changes
        window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', function(e) {
            if (!getStoredTheme()) {
                setTheme(e.matches ? 'dark' : 'light');
            }
        });

        // Checkbox functionality
        function convertCheckboxes() {
            const noteContents = document.querySelectorAll('.note-content');
            
            noteContents.forEach(function(content) {
                let html = content.innerHTML;
                
                // Convert unchecked checkboxes: - [ ] or * [ ] 
                html = html.replace(/^(\s*)([-*])\s+\[\s\]\s+(.+)$/gm, function(match, indent, bullet, text) {
                    return indent + '<div class="checkbox-item">' +
                           '<div class="checkbox"></div>' +
                           '<span class="checkbox-text">' + text + '</span>' +
                           '</div>';
                });
                
                // Convert checked checkboxes: - [x] or * [x]
                html = html.replace(/^(\s*)([-*])\s+\[[xX]\]\s+(.+)$/gm, function(match, indent, bullet, text) {
                    return indent + '<div class="checkbox-item checked">' +
                           '<div class="checkbox checked"></div>' +
                           '<span class="checkbox-text">' + text + '</span>' +
                           '</div>';
                });
                
                content.innerHTML = html;
            });
        }

        // Initialize checkboxes after theme is set
        document.addEventListener('DOMContentLoaded', function() {
            const preferredTheme = getPreferredTheme();
            setTheme(preferredTheme);
            
            // Convert checkboxes after a short delay to ensure content is rendered
            setTimeout(convertCheckboxes, 100);
        });
    </script>
</body>
</html>`

// extractTitleFromContent extracts title from markdown content
func extractTitleFromContent(content string) string {
	lines := splitLines(content)
	inFrontmatter := false

	for i, line := range lines {
		line = trimSpace(line)

		if i == 0 && line == "---" {
			inFrontmatter = true
			continue
		}

		if inFrontmatter && line == "---" {
			inFrontmatter = false
			continue
		}

		if inFrontmatter {
			if startsWithTitle(line) {
				return extractTitleValue(line)
			}
		} else {
			// Look for first markdown header
			if startsWithHash(line) {
				return trimLeadingHashes(line)
			}
		}
	}

	return ""
}

// Helper functions to avoid regex
func splitLines(s string) []string {
	var lines []string
	var current string
	for i, c := range s {
		if c == '\n' {
			lines = append(lines, current)
			current = ""
		} else if c == '\r' && i+1 < len(s) && s[i+1] == '\n' {
			lines = append(lines, current)
			current = ""
		} else if c != '\r' {
			current += string(c)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

func trimSpace(s string) string {
	start := 0
	end := len(s)

	for start < end && isSpace(s[start]) {
		start++
	}

	for end > start && isSpace(s[end-1]) {
		end--
	}

	return s[start:end]
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func startsWithTitle(line string) bool {
	return len(line) >= 6 && line[:6] == "title:"
}

func extractTitleValue(line string) string {
	if len(line) <= 6 {
		return ""
	}
	value := trimSpace(line[6:])
	if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'')) {
		return value[1 : len(value)-1]
	}
	return value
}

func startsWithHash(line string) bool {
	return len(line) > 0 && line[0] == '#'
}

func trimLeadingHashes(line string) string {
	i := 0
	for i < len(line) && line[i] == '#' {
		i++
	}
	if i < len(line) && line[i] == ' ' {
		i++
	}
	return line[i:]
}

// getContentWithoutFrontmatter removes YAML frontmatter from content
func getContentWithoutFrontmatter(content string) string {
	lines := splitLines(content)
	if len(lines) == 0 {
		return content
	}

	if trimSpace(lines[0]) != "---" {
		return content
	}

	for i := 1; i < len(lines); i++ {
		if trimSpace(lines[i]) == "---" {
			if i+1 < len(lines) {
				var result string
				for j := i + 1; j < len(lines); j++ {
					if j > i+1 {
						result += "\n"
					}
					result += lines[j]
				}
				return result
			}
			return ""
		}
	}

	return content
}

// PageData represents the data passed to the HTML template
type PageData struct {
	Notes         []NoteEntry
	FormattedDate string
	Count         int
}

// Execute runs the serve command
func Execute(folderPath string, verbose bool, port int) error {
	// Check if folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return fmt.Errorf("folder does not exist: %s", folderPath)
	}

	// Set up HTTP handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Get notes for today using the same logic as onthisday
		today := time.Now()
		notePaths, err := markdown.ScanMarkdownNotes(folderPath, markdown.SameDayMatcher, today, verbose)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning folder: %v", err), http.StatusInternalServerError)
			return
		}

		// Process each note to extract metadata and content
		var notes []NoteEntry
		for _, path := range notePaths {
			content, err := markdown.ReadFileContent(path)
			if err != nil {
				if verbose {
					fmt.Printf("Warning: Could not read file %s: %v\n", path, err)
				}
				continue
			}

			fileDate, err := markdown.ParseYAMLDate(content)
			if err != nil {
				if verbose {
					fmt.Printf("Warning: Could not parse date from %s: %v\n", path, err)
				}
				continue
			}

			title := extractTitleFromContent(content)
			cleanContent := getContentWithoutFrontmatter(content)

			// Get relative path for display
			relPath, err := filepath.Rel(folderPath, path)
			if err != nil {
				relPath = path
			}

			// Get absolute path for Obsidian link
			fullPath, err := filepath.Abs(path)
			if err != nil {
				fullPath = path
			}

			notes = append(notes, NoteEntry{
				Path:     relPath,
				FullPath: fullPath,
				Date:     fileDate,
				Title:    title,
				Content:  cleanContent,
			})
		}

		// Sort notes by date, newest first
		sort.Slice(notes, func(i, j int) bool {
			return notes[i].Date.After(notes[j].Date)
		})

		// Prepare template data
		data := PageData{
			Notes:         notes,
			FormattedDate: today.Format("Monday, January 2"),
			Count:         len(notes),
		}

		// Parse and execute template
		tmpl, err := template.New("onthisday").Parse(htmlTemplate)
		if err != nil {
			http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, fmt.Sprintf("Template execution error: %v", err), http.StatusInternalServerError)
			return
		}
	})

	// Start server
	addr := ":" + strconv.Itoa(port)
	fmt.Printf("Starting server on http://localhost%s\n", addr)
	fmt.Printf("Serving notes from: %s\n", folderPath)
	fmt.Printf("Press Ctrl+C to stop\n")

	return http.ListenAndServe(addr, nil)
}
