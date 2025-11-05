# tada

A vim-inspired terminal-based todo list manager using the todo.txt format.

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Bubbles](https://github.com/charmbracelet/bubbles), and [Lipgloss](https://github.com/charmbracelet/lipgloss) from Charm.

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
- ✅ **Context-based organization** - todos automatically grouped by `@context` tags (sorted alphabetically)
- ✅ **Priority system** - todos sorted by priority (A-Z) with color-coded visual badges
- ✅ Multiple context lists displayed simultaneously
- ✅ Example todo list pre-populated
- ✅ Both `tada` and `td` commands work
- ✅ Default storage at `~/.tada/todo.txt` (customizable with `-f` flag)
- ✅ Built with Bubbles components for professional text input with cursor, selection, and clipboard support
- ✅ **Beautiful styling** with Lipgloss - customizable color themes (config file support coming soon)
- ✅ **Safe deletion** with confirmation prompts to prevent accidental data loss

**Vim-Inspired Modes:**
- ✅ **Normal mode** (default) - Navigation with j/k or arrow keys, h/l to switch contexts
- ✅ **Command mode** (`:`) - Execute commands to manage todos with interactive text input
- ✅ **Insert mode** (`i` or `Enter`) - Edit the currently selected todo (prefilled with existing text)
- ✅ **Visual mode** (`v`) - Ready for selection operations
- ✅ **Leader key** (`Space`) - Quick access to common commands (edit, add, delete)
- ✅ Mode switching with `Esc` to return to normal mode

**Keybindings (Normal mode):**
- `j/k` or up/down arrows: Navigate todos (up/down within and across lists)
- `h/l` or left/right arrows: Switch between context lists
- `i` or `Enter`: Edit the currently selected todo
- `v`: Enter visual mode
- `:`: Enter command mode
- `Space`: Leader key (opens leader command menu)
- `q` or `Ctrl+C`: Quit

**Leader Key Bindings (Space + ...):**
- `e`: Edit current task (same as `i` or `Enter`)
- `a` or `n`: Add new task (opens command mode with `:add ` prefilled)
- `d` or `x`: Delete current task (with confirmation prompt)

**Command Mode Commands:**
- `:add <task description>` - Add a new task
- `:edit <new description>` - Edit the currently selected task
- `:done` - Mark the currently selected task as complete
- `:delete` or `:del` - Delete the currently selected task
- `Esc` - Cancel and return to normal mode

## Usage Examples

**Editing a task (Insert Mode):**
1. Navigate to the task using `j/k` keys
2. Press `i` or `Enter` to enter insert mode (or `Space` then `e` for leader key)
3. The input will be prefilled with the current task description
4. Edit the text as needed
5. Press `Enter` to save changes

**Adding a new task (Leader Key - Quick):**
1. Press `Space` then `a` (or `n`)
2. The command mode opens with `:add ` already typed
3. Type your task description (e.g., `Buy groceries @Personal`)
4. Press `Enter` to create the task

**Adding a prioritized task:**
1. Press `:` to enter command mode
2. Type `add (A) Fix critical bug @Work` (priority must be at the start)
3. Press `Enter` to create a high-priority task with a red (A) badge
4. The task will appear at the top of the @Work context list

**Deleting a task (Leader Key):**
1. Navigate to the task using `j/k` keys
2. Press `Space` then `d` (or `x`)
3. A confirmation prompt appears showing the task to be deleted
4. Press `d`, `x`, or `Enter` to confirm deletion
5. Press `Esc` to cancel (any other key also cancels)

**Editing via Command Mode:**
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

## Installation

### Prerequisites
- Go 1.21 or higher ([install Go](https://go.dev/doc/install))
- Git

### Quick Start

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd tada
   ```

2. **Run the install script:**
   ```bash
   ./install.sh
   ```

   The script will:
   - Install dependencies
   - Build the application
   - Optionally install `tada` to your PATH

3. **Configure your todo file location:**

   **Option A: Use the default location** (recommended for first-time users)
   ```bash
   mkdir -p ~/.tada
   touch ~/.tada/todo.txt
   ```

   **Option B: Use a custom location**
   ```bash
   # Set your preferred location (tada will create the file if needed)
   tada config set file ~/Documents/todo.txt
   # Or any other location you prefer
   ```

4. **Start using tada:**
   ```bash
   tada    # If installed to PATH
   # or
   ./tada  # If not installed to PATH
   ```

### Configuring Todo File Location

If you want to store your todo file in a custom location, you can configure it once and `tada` will always use that location:

```bash
# Set your custom todo file location
tada config set file /path/to/your/todo.txt

# Or use a path relative to your home directory
tada config set file ~/Documents/todo.txt
```

Once configured, you can use `tada` without any flags. The configuration is stored in `~/.tada/config.yml`.

**Configuration commands:**
```bash
tada config get          # Show all configuration
tada config get file     # Show todo file location
tada config set file PATH  # Set todo file location
tada config path         # Show config file path
```

**Note:** The `-f` flag still works and takes priority over the config file if you need to temporarily use a different file.

### Creating a Short Alias

If you prefer the `td` command, add this to your shell configuration (~/.bashrc, ~/.zshrc, etc.):

```bash
alias td='tada'
```

### Getting Started

1. Run `tada` to start the app
2. Press `Space` then `a` (or `n`) to add your first task
3. Type your task, for example: `Buy groceries @Personal`
4. Press `Enter` to save
5. Use `j/k` to navigate and explore the keybindings!

Your todos are automatically saved to your todo.txt file. You can edit it directly with any text editor, and changes will be reflected when you restart `tada`.

## Context-Based Organization

Todos are automatically grouped by their `@context` tags:
- Add contexts to your todos: `Buy milk @Personal` or `Review PR @Work`
- Multiple contexts are supported: `Call dentist @Personal @Health`
- Todos appear in all their associated context lists
- Todos without contexts go to "No Context"
- **Contexts are sorted alphabetically** for easy navigation
- Use `h/l` keys to switch between context lists
- Use `j/k` keys to navigate within and across lists

## Priority System

Todos support priority levels following the todo.txt format:
- **Priority syntax:** `(A)` through `(Z)`, where `(A)` is highest priority
- **Automatic sorting:** Within each context, todos are sorted by priority (A → Z, then unprioritized)
- **Visual indicators:** Priority badges with color-coding:
  - **(A)** - Red background (urgent)
  - **(B)** - Orange background (high priority)
  - **(C)** - Yellow background (medium-high priority)
  - **(D-F)** - Orange (high)
  - **(G-M)** - Blue (medium)
  - **(N-Z)** - Gray (low)
  - No priority indicator for unprioritized tasks

**Example:** `:add (A) Fix critical bug @Work` creates a high-priority task

## Theming

The app features a beautiful default color scheme with:
- Bordered, styled header with app title
- Distinct colors for active/inactive context lists
- Highlighted cursor and selected items
- **Color-coded priority badges** (Red for A, Orange for B, Yellow for C, etc.)
- Color-coded mode indicators (Normal: Blue, Insert: Green, Command: Orange, Visual: Purple)
- Styled help text with visual separators
- Professional text input components

**Future Customization:** The theming system is architected to support loading custom colors from a config file. All colors are centralized in `internal/tui/theme.go`, making it easy to add configuration file support in the future.

## Development

### Setup

1. Clone the repository
2. Install dependencies: `go mod download`
3. Install git hooks: `make install-hooks`

### Git Hooks

Pre-commit hooks are available to ensure code quality:

```bash
# Install hooks (recommended)
make install-hooks

# Or manually
./hooks/install.sh
```

The pre-commit hook runs:
1. **Code formatting** checks (gofmt)
2. **Linting** (golangci-lint or go vet)
3. **Tests** (go test)

To skip the hook when needed: `git commit --no-verify`

### Make Commands

```bash
make build                 # Build the application
make test                  # Run tests
make test-verbose          # Run tests with verbose output
make test-coverage         # Run tests with coverage report
make test-coverage-detail  # Generate HTML coverage report
make lint                  # Run linter (golangci-lint or go vet)
make fmt                   # Format code with gofmt
make install-hooks         # Install git hooks
make uninstall-hooks       # Uninstall git hooks
make clean                 # Remove build artifacts
```

### CI/CD

GitHub Actions workflows automatically:
- Build and test on Go 1.21, 1.22, and 1.23
- Run linting checks
- Generate coverage reports
- Build binaries for Linux, macOS, and Windows
- Upload build artifacts

See `.github/workflows/ci.yml` for details.

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -cover

# Run specific package
go test ./internal/todo -v
```

Current test coverage:
- `internal/todo`: 88.9%
- `internal/tui`: 7.9% (utility functions)
- Overall: 22.9%

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
