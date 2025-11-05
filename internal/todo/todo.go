package todo

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

// Item represents a todo item following the todo.txt format
type Item struct {
	Raw            string
	Completed      bool
	Priority       string
	CompletionDate string
	CreationDate   string
	Description    string
	Contexts       []string
	Projects       []string
}

// Parse parses a todo.txt line into an Item
func Parse(line string) Item {
	item := Item{
		Raw:      line,
		Contexts: []string{},
		Projects: []string{},
	}

	if strings.TrimSpace(line) == "" {
		return item
	}

	parts := strings.Fields(line)
	if len(parts) == 0 {
		return item
	}

	idx := 0

	// Check for completion marker
	if parts[idx] == "x" {
		item.Completed = true
		idx++
		if idx >= len(parts) {
			return item
		}
	}

	// Check for priority at the start (e.g., "(A)")
	if !item.Completed && len(parts[idx]) == 3 && parts[idx][0] == '(' && parts[idx][2] == ')' {
		item.Priority = string(parts[idx][1])
		idx++
		if idx >= len(parts) {
			return item
		}
	}

	// Date regex pattern
	datePattern := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

	// Check for completion date (only if completed)
	if item.Completed && datePattern.MatchString(parts[idx]) {
		item.CompletionDate = parts[idx]
		idx++
		if idx >= len(parts) {
			return item
		}
	}

	// Check for creation date
	if datePattern.MatchString(parts[idx]) {
		item.CreationDate = parts[idx]
		idx++
		if idx >= len(parts) {
			return item
		}
	}

	// Rest is description with contexts and projects
	descParts := parts[idx:]
	for _, part := range descParts {
		if strings.HasPrefix(part, "@") {
			item.Contexts = append(item.Contexts, part[1:])
		} else if strings.HasPrefix(part, "+") {
			item.Projects = append(item.Projects, part[1:])
		} else if strings.HasPrefix(part, "pri:") && len(part) == 5 {
			item.Priority = string(part[4])
		}
	}

	item.Description = strings.Join(descParts, " ")

	return item
}

// String returns the formatted todo.txt string
func (i Item) String() string {
	return i.Raw
}

// LoadFromFile loads todos from a file
func LoadFromFile(filename string) ([]Item, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close() //nolint:errcheck // Read-only operation, close error not critical

	var items []Item
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		items = append(items, Parse(line))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// SaveToFile saves todos to a file
func SaveToFile(filename string, items []Item) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	writer := bufio.NewWriter(file)
	for _, item := range items {
		_, err := fmt.Fprintln(writer, item.String())
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

// IsCompletedOlderThanDays checks if a completed todo is older than the specified number of days
func (i Item) IsCompletedOlderThanDays(days int) bool {
	if !i.Completed || i.CompletionDate == "" {
		return false
	}

	completionTime, err := time.Parse("2006-01-02", i.CompletionDate)
	if err != nil {
		return false
	}

	cutoffDate := time.Now().AddDate(0, 0, -days)
	return completionTime.Before(cutoffDate)
}

// ShouldBeVisible returns true if a todo should be visible in the main view
// Completed todos older than 5 days should not be visible
func (i Item) ShouldBeVisible() bool {
	if !i.Completed {
		return true
	}
	return !i.IsCompletedOlderThanDays(5)
}

// ArchiveOldCompletedTodos moves completed todos older than 5 days to archive files
// Returns the remaining todos (without archived items) and any error
func ArchiveOldCompletedTodos(todos []Item, archiveDir string) ([]Item, error) {
	// Group old completed todos by month
	archiveByMonth := make(map[string][]Item)
	var remainingTodos []Item

	for _, item := range todos {
		if item.IsCompletedOlderThanDays(5) {
			// Parse completion date to get year and month
			if item.CompletionDate != "" {
				completionTime, err := time.Parse("2006-01-02", item.CompletionDate)
				if err == nil {
					monthKey := completionTime.Format("2006_01")
					archiveByMonth[monthKey] = append(archiveByMonth[monthKey], item)
					continue
				}
			}
			// If we can't parse the date, keep it in the main list
			remainingTodos = append(remainingTodos, item)
		} else {
			remainingTodos = append(remainingTodos, item)
		}
	}

	// Write each month's items to the appropriate archive file
	for monthKey, items := range archiveByMonth {
		archiveFilename := fmt.Sprintf("%s/todo_archive_%s.txt", archiveDir, monthKey)

		// Open file in append mode, create if doesn't exist
		file, err := os.OpenFile(archiveFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open archive file %s: %w", archiveFilename, err)
		}

		writer := bufio.NewWriter(file)
		for _, item := range items {
			_, err := fmt.Fprintln(writer, item.String())
			if err != nil {
				_ = file.Close() // Best effort close on error path
				return nil, fmt.Errorf("failed to write to archive file %s: %w", archiveFilename, err)
			}
		}

		if err := writer.Flush(); err != nil {
			_ = file.Close() // Best effort close on error path
			return nil, fmt.Errorf("failed to flush archive file %s: %w", archiveFilename, err)
		}

		if err := file.Close(); err != nil {
			return nil, fmt.Errorf("failed to close archive file %s: %w", archiveFilename, err)
		}
	}

	return remainingTodos, nil
}
