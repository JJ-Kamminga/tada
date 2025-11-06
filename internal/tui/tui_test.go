package tui

import (
	"tada/internal/todo"
	"testing"
	"time"
)

func TestPriorityValue(t *testing.T) {
	tests := []struct {
		name     string
		priority string
		expected int
	}{
		{
			name:     "priority A should have value 0",
			priority: "A",
			expected: 0,
		},
		{
			name:     "priority B should have value 1",
			priority: "B",
			expected: 1,
		},
		{
			name:     "priority C should have value 2",
			priority: "C",
			expected: 2,
		},
		{
			name:     "priority Z should have value 25",
			priority: "Z",
			expected: 25,
		},
		{
			name:     "empty priority should have value 1000",
			priority: "",
			expected: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := priorityValue(tt.priority)
			if result != tt.expected {
				t.Errorf("priorityValue(%q) = %d, want %d", tt.priority, result, tt.expected)
			}
		})
	}
}

func TestSortTodosByPriority(t *testing.T) {
	tests := []struct {
		name     string
		todos    []TodoWithIndex
		expected []string // Expected order of priorities
	}{
		{
			name: "sort A, C, B should result in A, B, C",
			todos: []TodoWithIndex{
				{Item: todo.Item{Priority: "A", Description: "Task A"}},
				{Item: todo.Item{Priority: "C", Description: "Task C"}},
				{Item: todo.Item{Priority: "B", Description: "Task B"}},
			},
			expected: []string{"A", "B", "C"},
		},
		{
			name: "sort with empty priority last",
			todos: []TodoWithIndex{
				{Item: todo.Item{Priority: "", Description: "No priority"}},
				{Item: todo.Item{Priority: "B", Description: "Task B"}},
				{Item: todo.Item{Priority: "A", Description: "Task A"}},
			},
			expected: []string{"A", "B", ""},
		},
		{
			name: "already sorted list remains sorted",
			todos: []TodoWithIndex{
				{Item: todo.Item{Priority: "A", Description: "Task A"}},
				{Item: todo.Item{Priority: "B", Description: "Task B"}},
				{Item: todo.Item{Priority: "C", Description: "Task C"}},
			},
			expected: []string{"A", "B", "C"},
		},
		{
			name: "reverse sorted list gets sorted correctly",
			todos: []TodoWithIndex{
				{Item: todo.Item{Priority: "Z", Description: "Task Z"}},
				{Item: todo.Item{Priority: "B", Description: "Task B"}},
				{Item: todo.Item{Priority: "A", Description: "Task A"}},
			},
			expected: []string{"A", "B", "Z"},
		},
		{
			name: "single item remains unchanged",
			todos: []TodoWithIndex{
				{Item: todo.Item{Priority: "A", Description: "Task A"}},
			},
			expected: []string{"A"},
		},
		{
			name:     "empty list remains empty",
			todos:    []TodoWithIndex{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortTodosByPriority(tt.todos)

			if len(tt.todos) != len(tt.expected) {
				t.Fatalf("Expected %d todos, got %d", len(tt.expected), len(tt.todos))
			}

			for i, expectedPriority := range tt.expected {
				if tt.todos[i].Item.Priority != expectedPriority {
					t.Errorf("todos[%d].Priority = %q, want %q", i, tt.todos[i].Item.Priority, expectedPriority)
				}
			}
		})
	}
}

func TestGroupTodosByContext(t *testing.T) {
	tests := []struct {
		name             string
		todos            []todo.Item
		expectedContexts []string // Expected context names in order
		expectedCounts   map[string]int
	}{
		{
			name: "single context",
			todos: []todo.Item{
				{Description: "Task 1 @Work", Contexts: []string{"Work"}},
				{Description: "Task 2 @Work", Contexts: []string{"Work"}},
			},
			expectedContexts: []string{"Work"},
			expectedCounts: map[string]int{
				"Work": 2,
			},
		},
		{
			name: "multiple contexts",
			todos: []todo.Item{
				{Description: "Task 1 @Work", Contexts: []string{"Work"}},
				{Description: "Task 2 @Personal", Contexts: []string{"Personal"}},
				{Description: "Task 3 @Work", Contexts: []string{"Work"}},
			},
			expectedContexts: []string{"Personal", "Work"},
			expectedCounts: map[string]int{
				"Work":     2,
				"Personal": 1,
			},
		},
		{
			name: "no context items go to No Context",
			todos: []todo.Item{
				{Description: "Task without context", Contexts: []string{}},
				{Description: "Another task", Contexts: []string{}},
			},
			expectedContexts: []string{"No Context"},
			expectedCounts: map[string]int{
				"No Context": 2,
			},
		},
		{
			name: "mixed context and no context",
			todos: []todo.Item{
				{Description: "Task @Work", Contexts: []string{"Work"}},
				{Description: "Task without context", Contexts: []string{}},
			},
			expectedContexts: []string{"No Context", "Work"},
			expectedCounts: map[string]int{
				"No Context": 1,
				"Work":       1,
			},
		},
		{
			name: "todo with multiple contexts appears in both",
			todos: []todo.Item{
				{Description: "Task @Work @Office", Contexts: []string{"Work", "Office"}},
			},
			expectedContexts: []string{"Office", "Work"},
			expectedCounts: map[string]int{
				"Office": 1,
				"Work":   1,
			},
		},
		{
			name: "contexts sorted alphabetically",
			todos: []todo.Item{
				{Description: "Task @Zebra", Contexts: []string{"Zebra"}},
				{Description: "Task @Apple", Contexts: []string{"Apple"}},
				{Description: "Task @Mango", Contexts: []string{"Mango"}},
			},
			expectedContexts: []string{"Apple", "Mango", "Zebra"},
			expectedCounts: map[string]int{
				"Zebra": 1,
				"Apple": 1,
				"Mango": 1,
			},
		},
		{
			name: "old completed items are filtered out",
			todos: []todo.Item{
				{
					Description:    "Old completed @Work",
					Contexts:       []string{"Work"},
					Completed:      true,
					CompletionDate: time.Now().AddDate(0, 0, -10).Format("2006-01-02"),
				},
				{Description: "Active task @Work", Contexts: []string{"Work"}},
			},
			expectedContexts: []string{"Work"},
			expectedCounts: map[string]int{
				"Work": 1,
			},
		},
		{
			name:             "empty todo list",
			todos:            []todo.Item{},
			expectedContexts: []string{},
			expectedCounts:   map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := groupTodosByContext(tt.todos)

			// Check number of context groups
			if len(result) != len(tt.expectedContexts) {
				t.Errorf("Expected %d context groups, got %d", len(tt.expectedContexts), len(result))
			}

			// Check context names and order
			for i, expectedContext := range tt.expectedContexts {
				if i >= len(result) {
					t.Errorf("Missing context at index %d: expected %q", i, expectedContext)
					continue
				}
				if result[i].Context != expectedContext {
					t.Errorf("Context[%d] = %q, want %q", i, result[i].Context, expectedContext)
				}
			}

			// Check counts per context
			for _, contextList := range result {
				expectedCount, exists := tt.expectedCounts[contextList.Context]
				if !exists {
					t.Errorf("Unexpected context: %q", contextList.Context)
					continue
				}
				if len(contextList.Todos) != expectedCount {
					t.Errorf("Context %q has %d todos, want %d", contextList.Context, len(contextList.Todos), expectedCount)
				}
			}
		})
	}
}

func TestGroupTodosByContext_PrioritySorting(t *testing.T) {
	todos := []todo.Item{
		{Description: "Task C @Work", Priority: "C", Contexts: []string{"Work"}},
		{Description: "Task A @Work", Priority: "A", Contexts: []string{"Work"}},
		{Description: "Task B @Work", Priority: "B", Contexts: []string{"Work"}},
	}

	result := groupTodosByContext(todos)

	if len(result) != 1 {
		t.Fatalf("Expected 1 context group, got %d", len(result))
	}

	workContext := result[0]
	if workContext.Context != "Work" {
		t.Fatalf("Expected Work context, got %q", workContext.Context)
	}

	// Verify todos are sorted by priority
	expectedPriorities := []string{"A", "B", "C"}
	if len(workContext.Todos) != len(expectedPriorities) {
		t.Fatalf("Expected %d todos, got %d", len(expectedPriorities), len(workContext.Todos))
	}

	for i, expectedPriority := range expectedPriorities {
		if workContext.Todos[i].Item.Priority != expectedPriority {
			t.Errorf("Todo[%d] priority = %q, want %q", i, workContext.Todos[i].Item.Priority, expectedPriority)
		}
	}
}

func TestSortTodosByPriority_CompletionStatus(t *testing.T) {
	tests := []struct {
		name              string
		todos             []TodoWithIndex
		expectedCompleted []bool     // Expected completion status in order
		expectedPriority  []string   // Expected priority in order
	}{
		{
			name: "uncompleted tasks come before completed tasks",
			todos: []TodoWithIndex{
				{Item: todo.Item{Priority: "A", Description: "Completed A", Completed: true}},
				{Item: todo.Item{Priority: "B", Description: "Uncompleted B", Completed: false}},
				{Item: todo.Item{Priority: "C", Description: "Uncompleted C", Completed: false}},
			},
			expectedCompleted: []bool{false, false, true},
			expectedPriority:  []string{"B", "C", "A"},
		},
		{
			name: "completed tasks sorted by priority within their group",
			todos: []TodoWithIndex{
				{Item: todo.Item{Priority: "C", Description: "Completed C", Completed: true}},
				{Item: todo.Item{Priority: "A", Description: "Completed A", Completed: true}},
				{Item: todo.Item{Priority: "B", Description: "Completed B", Completed: true}},
			},
			expectedCompleted: []bool{true, true, true},
			expectedPriority:  []string{"A", "B", "C"},
		},
		{
			name: "uncompleted tasks sorted by priority within their group",
			todos: []TodoWithIndex{
				{Item: todo.Item{Priority: "C", Description: "Uncompleted C", Completed: false}},
				{Item: todo.Item{Priority: "A", Description: "Uncompleted A", Completed: false}},
				{Item: todo.Item{Priority: "B", Description: "Uncompleted B", Completed: false}},
			},
			expectedCompleted: []bool{false, false, false},
			expectedPriority:  []string{"A", "B", "C"},
		},
		{
			name: "mixed priorities with mixed completion status",
			todos: []TodoWithIndex{
				{Item: todo.Item{Priority: "C", Description: "Completed C", Completed: true}},
				{Item: todo.Item{Priority: "A", Description: "Uncompleted A", Completed: false}},
				{Item: todo.Item{Priority: "B", Description: "Completed B", Completed: true}},
				{Item: todo.Item{Priority: "", Description: "Uncompleted no priority", Completed: false}},
				{Item: todo.Item{Priority: "", Description: "Completed no priority", Completed: true}},
			},
			expectedCompleted: []bool{false, false, true, true, true},
			expectedPriority:  []string{"A", "", "B", "C", ""},
		},
		{
			name: "completed with higher priority still comes after uncompleted with lower priority",
			todos: []TodoWithIndex{
				{Item: todo.Item{Priority: "A", Description: "Completed A", Completed: true}},
				{Item: todo.Item{Priority: "Z", Description: "Uncompleted Z", Completed: false}},
			},
			expectedCompleted: []bool{false, true},
			expectedPriority:  []string{"Z", "A"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortTodosByPriority(tt.todos)

			if len(tt.todos) != len(tt.expectedCompleted) {
				t.Fatalf("Expected %d todos, got %d", len(tt.expectedCompleted), len(tt.todos))
			}

			for i := range tt.todos {
				if tt.todos[i].Item.Completed != tt.expectedCompleted[i] {
					t.Errorf("todos[%d].Completed = %v, want %v (description: %s)",
						i, tt.todos[i].Item.Completed, tt.expectedCompleted[i], tt.todos[i].Item.Description)
				}
				if tt.todos[i].Item.Priority != tt.expectedPriority[i] {
					t.Errorf("todos[%d].Priority = %q, want %q (description: %s)",
						i, tt.todos[i].Item.Priority, tt.expectedPriority[i], tt.todos[i].Item.Description)
				}
			}
		})
	}
}
