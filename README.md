# tada

A vim-inspired terminal-based todo list manager using the todo.txt format.

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Bubbles](https://github.com/charmbracelet/bubbles) from Charm.

## Project Structure
```
tada/
├── cmd/
│   └── root.go          # Cobra command setup
├── internal/
│   ├── todo/
│   │   └── todo.go      # Todo.txt parser
│   └── tui/
│       └── model.go     # Bubble Tea TUI model
├── main.go              # Entry point
├── todo.txt             # Example todo list
└── requirements.md
```

## Features Implemented

**Core Functionality:**
- ✅ Interactive TUI that takes over the terminal
- ✅ Todo.txt format parser (reads/writes todo.txt files)
- ✅ Example todo list pre-populated
- ✅ Both `tada` and `td` commands work
- ✅ Default storage at `~/.tada/todo.txt` (customizable with `-f` flag)
- ✅ Built with Bubbles components for professional text input with cursor, selection, and clipboard support

**Vim-Inspired Modes:**
- ✅ **Normal mode** (default) - Navigation with j/k or arrow keys
- ✅ **Command mode** (`:`) - Execute commands to manage todos with interactive text input
- ✅ **Insert mode** (`i` or `Enter`) - Add new todos with interactive text input
- ✅ **Visual mode** (`v`) - Ready for selection operations
- ✅ Mode switching with `Esc` to return to normal mode

**Keybindings (Normal mode):**
- `j/k` or arrow keys: Navigate todos
- `i` or `Enter`: Enter insert mode
- `v`: Enter visual mode
- `:`: Enter command mode
- `q` or `Ctrl+C`: Quit

**Command Mode Commands:**
- `:add <task description>` - Add a new task
- `:edit <new description>` - Edit the currently selected task
- `:done` - Mark the currently selected task as complete
- `:delete` or `:del` - Delete the currently selected task
- `Esc` - Cancel and return to normal mode

## Usage Examples

**Adding a new task (Insert Mode):**
1. Press `i` or `Enter` to enter insert mode
2. Type your task description (e.g., `Buy groceries @Personal`)
3. Press `Enter` to create the task

**Adding a new task (Command Mode):**
1. Press `:` to enter command mode
2. Type `add Buy groceries @Personal`
3. Press `Enter` to create the task

**Editing an existing task:**
1. Navigate to the task using `j/k` keys
2. Press `:` to enter command mode
3. Type `edit Buy groceries and milk @Personal`
4. Press `Enter` to update the task

**Marking a task as complete:**
1. Navigate to the task using `j/k` keys
2. Press `:` to enter command mode
3. Type `done`
4. Press `Enter` to mark it complete

**Deleting a task:**
1. Navigate to the task using `j/k` keys
2. Press `:` to enter command mode
3. Type `delete` (or `del`)
4. Press `Enter` to remove the task

All changes are automatically saved to your todo.txt file.

## Building

```bash
go build -o tada
cp tada td  # Create the 'td' alias
```

## Running the App

```bash
./tada                    # Run with default todo.txt location
./td                      # Same using the 'td' alias
./tada -f custom.txt      # Use a custom todo file
./tada --help             # Show help
```

The app will display your todo list with visual indicators for completed items (strikethrough) and the current mode in the footer. Command mode is fully functional for managing your todos!

## Todo.txt Format

The app follows the [todo.txt standard](https://github.com/todotxt/todo.txt). Example todo items:

```
x 2025-09-25 2025-09-24 Review blog post @Work pri:A
(A) 2025-09-26 Call dentist for appointment @Personal +Health
2025-09-27 Buy groceries for the weekend @Personal
(B) Finish quarterly report @Work +Q4
```

Format:
- `x` at start = completed
- `(A)` = priority
- First date = completion date (if completed)
- Second date = creation date
- `@Context` = context tags
- `+Project` = project tags
