# tada

A vim-inspired terminal-based todo list manager using the todo.txt format.

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Bubbles](https://github.com/charmbracelet/bubbles), and [Lipgloss](https://github.com/charmbracelet/lipgloss) from Charm.

## Features

- Interactive TUI with vim-inspired keybindings
- Todo.txt format support with automatic parsing and saving
- **Context-based organization** - todos automatically grouped by `@context` tags
- **Priority system** - color-coded visual badges for priorities (A-Z)
- **Archiving** - automatically archive completed todos older than 5 days
- **Multiple modes** - Normal, Insert, Command, and Visual modes
- **Leader key** (`Space`) for quick access to common commands
- Beautiful styling with color-coded priority badges and mode indicators
- Safe deletion with confirmation prompts
- Both `tada` and `td` commands available

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

3. **Configure your todo directory:**
   ```bash
   tada config set dir ~/.tada
   ```

   This sets the directory where your `todo.txt` and archive files will be stored.

4. **Start using tada:**
   ```bash
   tada    # If installed to PATH
   # or
   ./tada  # If not installed to PATH
   ```

   On first run with a new directory, an empty `todo.txt` will be created automatically.

### Configuration

A configuration file is **required** to use tada. On first run without configuration, you'll see:

```
Error: No todo directory configured.

To get started, set your todo directory:
  tada config set dir /path/to/your/todo/directory

Example:
  tada config set dir ~/.tada
```

**Configuration commands:**
```bash
tada config set dir PATH  # Set todo directory (required)
tada config get           # Show all configuration
tada config get dir       # Show todo directory location
tada config path          # Show config file path (~/.tada/config.yml)
```

**Directory structure:**

Once configured, your todo directory will contain:
```
~/.tada/
├── todo.txt                    # Your active todos
├── todo_archive_2024_11.txt    # November 2024 archive
├── todo_archive_2024_12.txt    # December 2024 archive
└── ...                         # Other monthly archives
```

All files (todo.txt and archives) are stored in the configured directory.

## Usage & Keybindings

### Vim-Inspired Modes

- **Normal mode** (default) - Navigate and view todos
- **Insert mode** (`i` or `Enter`) - Edit the selected todo
- **Command mode** (`:`) - Execute commands
- **Visual mode** (`v`) - Selection mode (ready for future operations)

Press `Esc` to return to Normal mode from any other mode.

### Normal Mode Keybindings

**Navigation:**
- `j/k` or `↑/↓` - Navigate todos (within and across context lists)
- `h/l` or `←/→` - Switch between context lists
- `q` or `Ctrl+C` - Quit

**Mode switching:**
- `i` or `Enter` - Enter Insert mode (edit selected todo)
- `:` - Enter Command mode
- `v` - Enter Visual mode

**Leader key** (`Space` + ...):
- `e` - Edit current task
- `a` or `n` - Add new task
- `d` or `x` - Delete current task (with confirmation)

### Command Mode

Type `:` to enter command mode, then use these commands:

- `:add <description>` - Add a new task
- `:edit <new description>` - Edit the selected task
- `:done` - Mark the selected task as complete
- `:delete` or `:del` - Delete the selected task (with confirmation)
- `:archive` - Archive completed todos older than 5 days
- `Esc` - Cancel and return to Normal mode

### Common Workflows

**Adding a task (quick method):**
1. Press `Space` then `a`
2. Type your task: `Buy groceries @Personal`
3. Press `Enter`

**Adding a prioritized task:**
1. Press `:`
2. Type: `add (A) Fix critical bug @Work`
3. Press `Enter` - task appears with red (A) badge at the top

**Editing a task:**
1. Navigate to task with `j/k`
2. Press `i` or `Enter`
3. Edit the prefilled text
4. Press `Enter` to save

**Deleting a task:**
1. Navigate to task with `j/k`
2. Press `Space` then `d`
3. Confirm with `d`, `x`, or `Enter` (or `Esc` to cancel)

**Marking as complete:**
1. Navigate to task with `j/k`
2. Press `:`
3. Type `done`
4. Press `Enter`

**Archiving old completed tasks:**
1. Press `:`
2. Type `archive`
3. Press `Enter`

All changes are automatically saved to `todo.txt` in your configured directory.

## Context-Based Organization

Todos are automatically grouped by their `@context` tags:

- Add contexts: `Buy milk @Personal` or `Review PR @Work`
- Multiple contexts: `Call dentist @Personal @Health` (appears in both lists)
- No context: Todos without `@context` go to "No Context" list
- Contexts are sorted alphabetically
- Navigate between contexts with `h/l` keys
- Navigate within contexts with `j/k` keys

## Priority System

Following the todo.txt format, priorities are indicated by `(A)` through `(Z)`:

**Visual indicators:**
- **(A)** - Red background (urgent)
- **(B)** - Orange background (high priority)
- **(C)** - Yellow background (medium-high priority)
- **(D-F)** - Orange text (high)
- **(G-M)** - Blue text (medium)
- **(N-Z)** - Gray text (low)
- No indicator for unprioritized tasks

**Behavior:**
- Todos are automatically sorted by priority within each context (A → Z, then unprioritized)
- Priority must be at the start of the task description: `(A) Fix bug @Work`

## Archiving

Completed todos older than 5 days can be archived:

1. Press `:` and type `archive`
2. Completed todos are moved to monthly archive files: `todo_archive_YYYY_MM.txt`
3. Archives are stored in your configured todo directory
4. Example: A task completed in November 2024 goes to `todo_archive_2024_11.txt`

## Todo.txt Format

The app follows the [todo.txt standard](https://github.com/todotxt/todo.txt). Example:

```
x 2025-09-25 2025-09-24 Review blog post @Work
(A) 2025-09-26 Call dentist for appointment @Personal +Health
2025-09-27 Buy groceries for the weekend @Personal
(B) Finish quarterly report @Work +Q4
```

**Format details:**
- `x` at start = completed
- `(A)` through `(Z)` = priority
- First date = completion date (if completed)
- Second date = creation date
- `@Context` = context tags (used for grouping)
- `+Project` = project tags

You can edit `todo.txt` directly with any text editor. Changes are reflected when you restart `tada`.

## Theming

The app features a beautiful default color scheme:

- Bordered, styled header with app title
- Distinct colors for active/inactive context lists
- Highlighted cursor and selected items
- Color-coded priority badges
- Mode indicators (Normal: Blue, Insert: Green, Command: Orange, Visual: Purple)
- Styled help text with visual separators

**Future customization:** The theming system is architected to support custom colors from a config file. All colors are centralized in `internal/tui/theme.go`.

## Development

### Project Structure
```
tada/
├── cmd/
│   ├── root.go          # Main command
│   └── config.go        # Config commands
├── internal/
│   ├── config/
│   │   └── config.go    # Configuration management
│   ├── todo/
│   │   └── todo.go      # Todo.txt parser and archiving
│   └── tui/
│       ├── model.go     # Bubble Tea TUI model
│       └── theme.go     # Color scheme and styling
├── main.go              # Entry point
└── hooks/
    └── install.sh       # Git hooks installer
```

### Setup

1. Clone the repository
2. Install dependencies: `go mod download`
3. Install git hooks: `make install-hooks` (recommended)
4. Build: `make build` or `go build`

### Git Hooks

Pre-commit hooks ensure code quality:

```bash
# Install hooks (recommended)
make install-hooks

# Or manually
./hooks/install.sh
```

The pre-commit hook runs:
1. Code formatting checks (gofmt)
2. Linting (golangci-lint or go vet)
3. Tests (go test)

To skip the hook: `git commit --no-verify`

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

### CI/CD

GitHub Actions workflows automatically:
- Build and test on Go 1.21, 1.22, and 1.23
- Run linting checks
- Generate coverage reports
- Build binaries for Linux, macOS, and Windows
- Upload build artifacts

See `.github/workflows/ci.yml` for details.

## Building from Source

```bash
# Build the binary
go build -o tada

# Run locally
./tada

# Or use the built-in alias
./td
```

The `td` command is a built-in alias (cobra alias) and works the same as `tada`.
