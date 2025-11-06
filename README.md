# tada

Manage a todo.txt To Do list, in your terminal, with vim-like keybinds!

The app follows the [todo.txt standard](https://github.com/todotxt/todo.txt). `todo.txt` is a text file, you can edit it with any text editor. Changes are reflected when you restart `tada`.

It has awesome styles thanks to [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Bubbles](https://github.com/charmbracelet/bubbles), and [Lipgloss](https://github.com/charmbracelet/lipgloss) from Charm.

This is my personal pet project, but I did give it a license, see the paragraph License at the bottom. I am not planning to accept feature requests, though bug reports and PR's for those are welcome.

## Installation from source

This is the only way to install tada.

### Prerequisites

- Go 1.21 or higher ([install Go](https://go.dev/doc/install))
- Git

### Steps

1. Git clone this repository.
2. Run the install script `install.sh`.
3. Configure your todo directory: `tada config set dir ~/.tada`
4. Start using tada:

   ```bash
   tada    # If installed to PATH
   # or
   ./tada  # If not installed to PATH

### Configuration

A configuration file is **required** to use tada.

**Configuration commands:**

```bash
tada config set dir PATH  # Set todo directory (required)
tada config get           # Show all configuration
tada config get dir       # Show todo directory location
tada config path          # Show config file path (~/.tada/config.yml)
```

**Directory structure:**

Once configured, your todo directory will contain:

```bash
~/.tada/
├── todo.txt                    # Your active todos
├── todo_archive_2024_11.txt    # November 2024 archive
├── todo_archive_2024_12.txt    # December 2024 archive
└── ...                         # Other monthly archives
```

All files (todo.txt and archives) are stored in the configured directory.

## Usage & Keybindings

All keybinds are displayed in the app.

All commands can be viewed from command mode by typing `/`.

## Archiving

Completed todos older than 5 days can be archived:

1. Press `:` and type `archive`
2. Completed todos are moved to monthly archive files: `todo_archive_YYYY_MM.txt`
3. Archives are stored in your configured todo directory
4. Example: A task completed in November 2024 goes to `todo_archive_2024_11.txt`

## Todo.txt Format

 Example:

```tada
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

If you really really want to, you can do this.

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

## License

This project is licensed under the Mozilla Public License Version 2.0, which included in the file LICENSE.md.

In simpler terms: the code is open source, but this is my project, and I reserve the right to use the name 'tada'.
