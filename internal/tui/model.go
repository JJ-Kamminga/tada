package tui

import (
	"fmt"
	"tada/internal/todo"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Mode represents the current mode of the app
type Mode int

const (
	ModeNormal Mode = iota
	ModeCommand
	ModeInsert
	ModeVisual
)

func (m Mode) String() string {
	switch m {
	case ModeNormal:
		return "NORMAL"
	case ModeCommand:
		return "COMMAND"
	case ModeInsert:
		return "INSERT"
	case ModeVisual:
		return "VISUAL"
	default:
		return "UNKNOWN"
	}
}

// Model represents the application state
type Model struct {
	todos    []todo.Item
	cursor   int
	mode     Mode
	filename string
	width    int
	height   int
}

// NewModel creates a new TUI model
func NewModel(filename string) Model {
	todos, err := todo.LoadFromFile(filename)
	if err != nil {
		// If file doesn't exist, start with empty list
		todos = []todo.Item{}
	}

	return Model{
		todos:    todos,
		cursor:   0,
		mode:     ModeNormal,
		filename: filename,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	return m, nil
}

// handleKeyPress handles key presses based on current mode
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global quit keys
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}

	switch m.mode {
	case ModeNormal:
		return m.handleNormalMode(msg)
	case ModeCommand:
		return m.handleCommandMode(msg)
	case ModeInsert:
		return m.handleInsertMode(msg)
	case ModeVisual:
		return m.handleVisualMode(msg)
	}

	return m, nil
}

// handleNormalMode handles key presses in normal mode
func (m Model) handleNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case ":":
		m.mode = ModeCommand
	case "i", "enter":
		m.mode = ModeInsert
	case "v":
		m.mode = ModeVisual
	case "q":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.todos)-1 {
			m.cursor++
		}
	}

	return m, nil
}

// handleCommandMode handles key presses in command mode
func (m Model) handleCommandMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = ModeNormal
	}

	return m, nil
}

// handleInsertMode handles key presses in insert mode
func (m Model) handleInsertMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = ModeNormal
	}

	return m, nil
}

// handleVisualMode handles key presses in visual mode
func (m Model) handleVisualMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = ModeNormal
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	var s string

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("cyan")).
		Padding(0, 1)

	s += headerStyle.Render("TADA - Todo List") + "\n\n"

	// Todo list
	if len(m.todos) == 0 {
		s += "  No todos yet. Press 'i' to add one!\n"
	} else {
		for i, item := range m.todos {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			// Style the item
			itemStyle := lipgloss.NewStyle()
			if item.Completed {
				itemStyle = itemStyle.Foreground(lipgloss.Color("240")).Strikethrough(true)
			}

			s += fmt.Sprintf("%s %s\n", cursor, itemStyle.Render(item.Description))
		}
	}

	// Footer with mode indicator
	s += "\n"
	modeStyle := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		Background(lipgloss.Color("blue")).
		Foreground(lipgloss.Color("white"))

	switch m.mode {
	case ModeInsert:
		modeStyle = modeStyle.Background(lipgloss.Color("green"))
	case ModeCommand:
		modeStyle = modeStyle.Background(lipgloss.Color("yellow")).Foreground(lipgloss.Color("black"))
	case ModeVisual:
		modeStyle = modeStyle.Background(lipgloss.Color("magenta"))
	}

	s += modeStyle.Render(m.mode.String())

	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Padding(1, 1)

	help := ""
	switch m.mode {
	case ModeNormal:
		help = "i/enter: insert | v: visual | :: command | j/k: navigate | q: quit"
	case ModeInsert:
		help = "esc: back to normal mode"
	case ModeCommand:
		help = "esc: back to normal mode"
	case ModeVisual:
		help = "esc: back to normal mode"
	}

	s += "\n" + helpStyle.Render(help)

	return s
}
