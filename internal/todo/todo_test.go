package todo

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestIsCompletedOlderThanDays(t *testing.T) {
	tests := []struct {
		name          string
		item          Item
		days          int
		expectedOlder bool
	}{
		{
			name: "completed item from 10 days ago should be older than 5 days",
			item: Item{
				Completed:      true,
				CompletionDate: time.Now().AddDate(0, 0, -10).Format("2006-01-02"),
			},
			days:          5,
			expectedOlder: true,
		},
		{
			name: "completed item from 2 days ago should not be older than 5 days",
			item: Item{
				Completed:      true,
				CompletionDate: time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
			},
			days:          5,
			expectedOlder: false,
		},
		{
			name: "completed item from exactly 5 days ago should be older than 5 days (boundary)",
			item: Item{
				Completed:      true,
				CompletionDate: time.Now().AddDate(0, 0, -5).Format("2006-01-02"),
			},
			days:          5,
			expectedOlder: true,
		},
		{
			name: "completed item from 6 days ago should be older than 5 days",
			item: Item{
				Completed:      true,
				CompletionDate: time.Now().AddDate(0, 0, -6).Format("2006-01-02"),
			},
			days:          5,
			expectedOlder: true,
		},
		{
			name: "uncompleted item should return false",
			item: Item{
				Completed:      false,
				CompletionDate: time.Now().AddDate(0, 0, -10).Format("2006-01-02"),
			},
			days:          5,
			expectedOlder: false,
		},
		{
			name: "completed item without completion date should return false",
			item: Item{
				Completed:      true,
				CompletionDate: "",
			},
			days:          5,
			expectedOlder: false,
		},
		{
			name: "completed item with invalid date format should return false",
			item: Item{
				Completed:      true,
				CompletionDate: "invalid-date",
			},
			days:          5,
			expectedOlder: false,
		},
		{
			name: "completed item from 30 days ago should be older than 10 days",
			item: Item{
				Completed:      true,
				CompletionDate: time.Now().AddDate(0, 0, -30).Format("2006-01-02"),
			},
			days:          10,
			expectedOlder: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.item.IsCompletedOlderThanDays(tt.days)
			if result != tt.expectedOlder {
				t.Errorf("IsCompletedOlderThanDays() = %v, want %v", result, tt.expectedOlder)
			}
		})
	}
}

func TestShouldBeVisible(t *testing.T) {
	tests := []struct {
		name            string
		item            Item
		expectedVisible bool
	}{
		{
			name: "uncompleted item should be visible",
			item: Item{
				Completed: false,
			},
			expectedVisible: true,
		},
		{
			name: "recently completed item should be visible",
			item: Item{
				Completed:      true,
				CompletionDate: time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
			},
			expectedVisible: true,
		},
		{
			name: "completed item from 10 days ago should not be visible",
			item: Item{
				Completed:      true,
				CompletionDate: time.Now().AddDate(0, 0, -10).Format("2006-01-02"),
			},
			expectedVisible: false,
		},
		{
			name: "completed item from exactly 5 days ago should not be visible (boundary)",
			item: Item{
				Completed:      true,
				CompletionDate: time.Now().AddDate(0, 0, -5).Format("2006-01-02"),
			},
			expectedVisible: false,
		},
		{
			name: "completed item without date should be visible",
			item: Item{
				Completed:      true,
				CompletionDate: "",
			},
			expectedVisible: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.item.ShouldBeVisible()
			if result != tt.expectedVisible {
				t.Errorf("ShouldBeVisible() = %v, want %v", result, tt.expectedVisible)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected Item
	}{
		{
			name: "simple task with no metadata",
			line: "Buy groceries",
			expected: Item{
				Raw:         "Buy groceries",
				Completed:   false,
				Description: "Buy groceries",
				Contexts:    []string{},
				Projects:    []string{},
			},
		},
		{
			name: "task with priority",
			line: "(A) Call dentist",
			expected: Item{
				Raw:         "(A) Call dentist",
				Completed:   false,
				Priority:    "A",
				Description: "Call dentist",
				Contexts:    []string{},
				Projects:    []string{},
			},
		},
		{
			name: "task with priority and creation date",
			line: "(B) 2025-09-26 Finish report",
			expected: Item{
				Raw:          "(B) 2025-09-26 Finish report",
				Completed:    false,
				Priority:     "B",
				CreationDate: "2025-09-26",
				Description:  "Finish report",
				Contexts:     []string{},
				Projects:     []string{},
			},
		},
		{
			name: "task with context",
			line: "Buy milk @Grocery",
			expected: Item{
				Raw:         "Buy milk @Grocery",
				Completed:   false,
				Description: "Buy milk @Grocery",
				Contexts:    []string{"Grocery"},
				Projects:    []string{},
			},
		},
		{
			name: "task with multiple contexts",
			line: "Team meeting @Work @Office",
			expected: Item{
				Raw:         "Team meeting @Work @Office",
				Completed:   false,
				Description: "Team meeting @Work @Office",
				Contexts:    []string{"Work", "Office"},
				Projects:    []string{},
			},
		},
		{
			name: "task with project",
			line: "Write documentation +WebsiteRedesign",
			expected: Item{
				Raw:         "Write documentation +WebsiteRedesign",
				Completed:   false,
				Description: "Write documentation +WebsiteRedesign",
				Contexts:    []string{},
				Projects:    []string{"WebsiteRedesign"},
			},
		},
		{
			name: "task with context and project",
			line: "Review code @Work +ProjectX",
			expected: Item{
				Raw:         "Review code @Work +ProjectX",
				Completed:   false,
				Description: "Review code @Work +ProjectX",
				Contexts:    []string{"Work"},
				Projects:    []string{"ProjectX"},
			},
		},
		{
			name: "completed task",
			line: "x Buy milk",
			expected: Item{
				Raw:         "x Buy milk",
				Completed:   true,
				Description: "Buy milk",
				Contexts:    []string{},
				Projects:    []string{},
			},
		},
		{
			name: "completed task with dates",
			line: "x 2025-09-25 2025-09-24 Review blog post",
			expected: Item{
				Raw:            "x 2025-09-25 2025-09-24 Review blog post",
				Completed:      true,
				CompletionDate: "2025-09-25",
				CreationDate:   "2025-09-24",
				Description:    "Review blog post",
				Contexts:       []string{},
				Projects:       []string{},
			},
		},
		{
			name: "completed task with dates, context and project",
			line: "x 2025-09-20 2025-09-18 Setup development environment @Work +DevOps",
			expected: Item{
				Raw:            "x 2025-09-20 2025-09-18 Setup development environment @Work +DevOps",
				Completed:      true,
				CompletionDate: "2025-09-20",
				CreationDate:   "2025-09-18",
				Description:    "Setup development environment @Work +DevOps",
				Contexts:       []string{"Work"},
				Projects:       []string{"DevOps"},
			},
		},
		{
			name: "task with pri: format",
			line: "Learn Go concurrency patterns @Learning +Programming pri:B",
			expected: Item{
				Raw:         "Learn Go concurrency patterns @Learning +Programming pri:B",
				Completed:   false,
				Priority:    "B",
				Description: "Learn Go concurrency patterns @Learning +Programming pri:B",
				Contexts:    []string{"Learning"},
				Projects:    []string{"Programming"},
			},
		},
		{
			name: "empty line",
			line: "",
			expected: Item{
				Raw:      "",
				Contexts: []string{},
				Projects: []string{},
			},
		},
		{
			name: "whitespace only",
			line: "   ",
			expected: Item{
				Raw:      "   ",
				Contexts: []string{},
				Projects: []string{},
			},
		},
		{
			name: "completed task with only completion marker",
			line: "x",
			expected: Item{
				Raw:       "x",
				Completed: true,
				Contexts:  []string{},
				Projects:  []string{},
			},
		},
		{
			name: "task with priority only",
			line: "(A)",
			expected: Item{
				Raw:      "(A)",
				Priority: "A",
				Contexts: []string{},
				Projects: []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Parse(tt.line)

			if result.Raw != tt.expected.Raw {
				t.Errorf("Raw = %q, want %q", result.Raw, tt.expected.Raw)
			}
			if result.Completed != tt.expected.Completed {
				t.Errorf("Completed = %v, want %v", result.Completed, tt.expected.Completed)
			}
			if result.Priority != tt.expected.Priority {
				t.Errorf("Priority = %q, want %q", result.Priority, tt.expected.Priority)
			}
			if result.CompletionDate != tt.expected.CompletionDate {
				t.Errorf("CompletionDate = %q, want %q", result.CompletionDate, tt.expected.CompletionDate)
			}
			if result.CreationDate != tt.expected.CreationDate {
				t.Errorf("CreationDate = %q, want %q", result.CreationDate, tt.expected.CreationDate)
			}
			if result.Description != tt.expected.Description {
				t.Errorf("Description = %q, want %q", result.Description, tt.expected.Description)
			}

			if len(result.Contexts) != len(tt.expected.Contexts) {
				t.Errorf("Contexts length = %d, want %d", len(result.Contexts), len(tt.expected.Contexts))
			} else {
				for i := range result.Contexts {
					if result.Contexts[i] != tt.expected.Contexts[i] {
						t.Errorf("Contexts[%d] = %q, want %q", i, result.Contexts[i], tt.expected.Contexts[i])
					}
				}
			}

			if len(result.Projects) != len(tt.expected.Projects) {
				t.Errorf("Projects length = %d, want %d", len(result.Projects), len(tt.expected.Projects))
			} else {
				for i := range result.Projects {
					if result.Projects[i] != tt.expected.Projects[i] {
						t.Errorf("Projects[%d] = %q, want %q", i, result.Projects[i], tt.expected.Projects[i])
					}
				}
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		item     Item
		expected string
	}{
		{
			name: "simple item returns raw string",
			item: Item{
				Raw: "Buy groceries @Store",
			},
			expected: "Buy groceries @Store",
		},
		{
			name: "completed item returns raw string",
			item: Item{
				Raw: "x 2025-09-25 Complete task",
			},
			expected: "x 2025-09-25 Complete task",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.item.String()
			if result != tt.expected {
				t.Errorf("String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestLoadFromFile(t *testing.T) {
	// Create a temporary file with test data
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_todo.txt")

	testContent := `(A) 2025-09-26 Call dentist @Personal +Health
2025-09-27 Buy groceries @Personal
(B) Finish quarterly report @Work +Q4
x 2025-09-20 2025-09-18 Setup development environment @Work
Learn Go concurrency patterns @Learning +Programming`

	err := os.WriteFile(tmpFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test loading
	items, err := LoadFromFile(tmpFile)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}

	expectedCount := 5
	if len(items) != expectedCount {
		t.Errorf("LoadFromFile() returned %d items, want %d", len(items), expectedCount)
	}

	// Verify first item
	if items[0].Priority != "A" {
		t.Errorf("items[0].Priority = %q, want %q", items[0].Priority, "A")
	}
	if items[0].CreationDate != "2025-09-26" {
		t.Errorf("items[0].CreationDate = %q, want %q", items[0].CreationDate, "2025-09-26")
	}

	// Verify completed item
	if !items[3].Completed {
		t.Error("items[3].Completed should be true")
	}
	if items[3].CompletionDate != "2025-09-20" {
		t.Errorf("items[3].CompletionDate = %q, want %q", items[3].CompletionDate, "2025-09-20")
	}
}

func TestLoadFromFile_NonExistentFile(t *testing.T) {
	_, err := LoadFromFile("/nonexistent/file.txt")
	if err == nil {
		t.Error("LoadFromFile() should return error for nonexistent file")
	}
}

func TestSaveToFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_save.txt")

	items := []Item{
		{Raw: "(A) Task one @Work"},
		{Raw: "x 2025-09-25 Task two"},
		{Raw: "Task three +Project"},
	}

	err := SaveToFile(tmpFile, items)
	if err != nil {
		t.Fatalf("SaveToFile() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Fatal("SaveToFile() did not create file")
	}

	// Load it back and verify content
	loadedItems, err := LoadFromFile(tmpFile)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}

	if len(loadedItems) != len(items) {
		t.Errorf("Loaded %d items, want %d", len(loadedItems), len(items))
	}

	for i := range items {
		if loadedItems[i].Raw != items[i].Raw {
			t.Errorf("loadedItems[%d].Raw = %q, want %q", i, loadedItems[i].Raw, items[i].Raw)
		}
	}
}

func TestSaveToFile_InvalidPath(t *testing.T) {
	items := []Item{{Raw: "Test"}}
	err := SaveToFile("/nonexistent/dir/file.txt", items)
	if err == nil {
		t.Error("SaveToFile() should return error for invalid path")
	}
}

func TestArchiveOldCompletedTodos(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test items
	oldDate := time.Now().AddDate(0, 0, -10).Format("2006-01-02")
	recentDate := time.Now().AddDate(0, 0, -2).Format("2006-01-02")

	items := []Item{
		{
			Raw:            "x " + oldDate + " 2025-09-01 Old completed task",
			Completed:      true,
			CompletionDate: oldDate,
			Description:    "Old completed task",
		},
		{
			Raw:            "x " + recentDate + " 2025-09-20 Recent completed task",
			Completed:      true,
			CompletionDate: recentDate,
			Description:    "Recent completed task",
		},
		{
			Raw:         "Active task @Work",
			Completed:   false,
			Description: "Active task @Work",
		},
	}

	// Archive old items
	remaining, err := ArchiveOldCompletedTodos(items, tmpDir)
	if err != nil {
		t.Fatalf("ArchiveOldCompletedTodos() error = %v", err)
	}

	// Should have 2 remaining items (recent completed + active)
	if len(remaining) != 2 {
		t.Errorf("ArchiveOldCompletedTodos() returned %d items, want 2", len(remaining))
	}

	// Verify the old item was archived
	hasOldItem := false
	for _, item := range remaining {
		if item.Description == "Old completed task" {
			hasOldItem = true
		}
	}
	if hasOldItem {
		t.Error("Old completed task should have been archived")
	}

	// Verify archive file was created
	monthKey := time.Now().AddDate(0, 0, -10).Format("2006_01")
	archiveFile := filepath.Join(tmpDir, "todo_archive_"+monthKey+".txt")

	if _, err := os.Stat(archiveFile); os.IsNotExist(err) {
		t.Error("Archive file was not created")
	}

	// Verify archive file content
	archivedItems, err := LoadFromFile(archiveFile)
	if err != nil {
		t.Fatalf("Failed to load archive file: %v", err)
	}

	if len(archivedItems) != 1 {
		t.Errorf("Archive file contains %d items, want 1", len(archivedItems))
	}
}

func TestArchiveOldCompletedTodos_NoOldItems(t *testing.T) {
	tmpDir := t.TempDir()

	items := []Item{
		{
			Raw:       "Active task",
			Completed: false,
		},
	}

	remaining, err := ArchiveOldCompletedTodos(items, tmpDir)
	if err != nil {
		t.Fatalf("ArchiveOldCompletedTodos() error = %v", err)
	}

	if len(remaining) != 1 {
		t.Errorf("ArchiveOldCompletedTodos() returned %d items, want 1", len(remaining))
	}

	// No archive files should be created
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read temp dir: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("Expected no files in archive dir, found %d", len(files))
	}
}
