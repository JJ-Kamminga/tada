# Claude Development Notes

This document contains notes and context for Claude Code when working on the tada project.

## Project Overview

tada is a vim-inspired terminal-based todo list manager written in Go. It uses:
- **Cobra** for CLI command structure
- **Bubble Tea** for the interactive TUI
- **lipgloss** for styling
- **todo.txt format** for data storage

## Architecture

### Package Structure

```
tada/
├── cmd/              # Cobra commands
├── internal/
│   ├── todo/        # Todo.txt parsing and file I/O
│   └── tui/         # Bubble Tea TUI implementation
└── main.go          # Application entry point
```

### Key Components

1. **cmd/root.go**: Main command setup
   - Manages file path resolution (defaults to ~/.tada/todo.txt)
   - Initializes the Bubble Tea program

2. **internal/todo/todo.go**: Todo.txt parser
   - `Item` struct represents a todo item
   - `Parse()` function parses todo.txt format
   - `LoadFromFile()` and `SaveToFile()` for persistence

3. **internal/tui/model.go**: TUI model
   - Implements Bubble Tea's Model interface
   - Manages application state and modes
   - Handles keyboard input and rendering

## Current Implementation Status

### Completed
- ✅ Basic Cobra CLI setup with aliases
- ✅ Todo.txt format parser
- ✅ Interactive TUI with Bubble Tea
- ✅ Mode system (normal, command, insert, visual)
- ✅ Mode switching keybindings
- ✅ Basic navigation (j/k, arrow keys)
- ✅ Visual display of todos with completion status

### Not Yet Implemented
- ❌ Insert mode functionality (adding/editing todos)
- ❌ Command mode commands
- ❌ Visual mode selection
- ❌ Todo item manipulation (add, edit, delete, toggle completion)
- ❌ Saving changes back to file
- ❌ Filtering and sorting
- ❌ Search functionality

## Development Guidelines

### Code Style
- Follow standard Go conventions
- Keep packages focused and small
- Use clear, descriptive names
- Comment exported functions and types

### Todo.txt Format Reference
From https://github.com/todotxt/todo.txt:

Format: `x (priority) completion-date creation-date description @context +project`

Examples:
```
x 2025-09-25 2025-09-24 Review blog post @Work pri:A
(A) 2025-09-26 Call dentist @Personal +Health
```

- Completion marker: `x` at start
- Priority: `(A)` through `(Z)` or inline `pri:A`
- Dates: ISO format YYYY-MM-DD
- Contexts: `@Context`
- Projects: `+Project`

### TUI Implementation Notes
- Bubble Tea uses the Elm Architecture (Model-Update-View)
- `tea.WithAltScreen()` provides full-screen experience
- The model is immutable - Update() returns a new model
- Messages drive all state changes

## Future Enhancements

Potential features to add:
1. Full insert mode with text input for new todos
2. Command mode with vim-like commands (`:w`, `:q`, `:filter`, etc.)
3. Visual mode for multi-select operations
4. Due date support and highlighting
5. Recurring tasks
6. Tab completion
7. Undo/redo functionality
8. Export to different formats
9. Sync with cloud services
10. Configuration file support

## Testing

Currently no tests are implemented. Future testing should include:
- Unit tests for todo.txt parser
- Integration tests for file I/O
- TUI interaction tests (if feasible)

## Build and Deploy

Build: `go build -o tada`
The binary is self-contained and can be copied anywhere in the PATH.
