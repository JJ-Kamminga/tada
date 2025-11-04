package todo

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
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
	defer file.Close()

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
func SaveToFile(filename string, items []Item) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, item := range items {
		_, err := fmt.Fprintln(writer, item.String())
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}
