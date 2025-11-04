package tui

import (
	"fmt"
	"strings"
	"tada/internal/todo"

	"github.com/charmbracelet/bubbles/textinput"
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
	todos        []todo.Item
	cursor       int
	mode         Mode
	filename     string
	width        int
	height       int
	commandInput textinput.Model // Text input for command mode
	insertInput  textinput.Model // Text input for insert mode
}

// NewModel creates a new TUI model
func NewModel(filename string) Model {
	todos, err := todo.LoadFromFile(filename)
	if err != nil {
		// If file doesn't exist, start with empty list
		todos = []todo.Item{}
	}

	// Initialize command input
	cmdInput := textinput.New()
	cmdInput.Placeholder = "enter command..."
	cmdInput.Prompt = ":"
	cmdInput.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("51")).Bold(true)
	cmdInput.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	cmdInput.CharLimit = 200

	// Initialize insert input
	insInput := textinput.New()
	insInput.Placeholder = "enter todo description..."
	insInput.Prompt = "> "
	insInput.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("green")).Bold(true)
	insInput.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	insInput.CharLimit = 500

	return Model{
		todos:        todos,
		cursor:       0,
		mode:         ModeNormal,
		filename:     filename,
		commandInput: cmdInput,
		insertInput:  insInput,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	// Update textinput components for cursor blink and other messages
	if m.mode == ModeCommand {
		m.commandInput, cmd = m.commandInput.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.mode == ModeInsert {
		m.insertInput, cmd = m.insertInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
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
		m.commandInput.Reset()
		m.commandInput.Focus()
		return m, textinput.Blink
	case "i", "enter":
		m.mode = ModeInsert
		m.insertInput.Reset()
		m.insertInput.Focus()
		return m, textinput.Blink
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

// executeCommand parses and executes a command
func (m Model) executeCommand() (Model, tea.Cmd) {
	cmdLine := m.commandInput.Value()
	parts := strings.Fields(cmdLine)
	if len(parts) == 0 {
		return m, nil
	}

	cmd := parts[0]
	args := strings.Join(parts[1:], " ")

	switch cmd {
	case "add":
		return m.cmdAdd(args)
	case "edit":
		return m.cmdEdit(args)
	case "done":
		return m.cmdDone(args)
	case "delete", "del":
		return m.cmdDelete(args)
	}

	return m, nil
}

// cmdAdd adds a new task
func (m Model) cmdAdd(description string) (Model, tea.Cmd) {
	if description == "" {
		return m, nil
	}

	// Create new todo item
	newItem := todo.Item{
		Raw:         description,
		Description: description,
		Completed:   false,
		Contexts:    []string{},
		Projects:    []string{},
	}

	m.todos = append(m.todos, newItem)

	// Save to file
	if err := todo.SaveToFile(m.filename, m.todos); err != nil {
		// TODO: Handle error (could add error message to model)
		return m, nil
	}

	// Return to normal mode
	m.mode = ModeNormal
	m.commandInput.Blur()

	// Move cursor to the new item
	m.cursor = len(m.todos) - 1

	return m, nil
}

// cmdEdit edits the current task
func (m Model) cmdEdit(newDescription string) (Model, tea.Cmd) {
	if len(m.todos) == 0 || newDescription == "" {
		return m, nil
	}

	// Update the description
	m.todos[m.cursor].Description = newDescription
	m.todos[m.cursor].Raw = newDescription

	// Save to file
	if err := todo.SaveToFile(m.filename, m.todos); err != nil {
		return m, nil
	}

	// Return to normal mode
	m.mode = ModeNormal
	m.commandInput.Blur()

	return m, nil
}

// cmdDone marks the current task as complete
func (m Model) cmdDone(args string) (Model, tea.Cmd) {
	if len(m.todos) == 0 {
		return m, nil
	}

	// Mark as completed
	m.todos[m.cursor].Completed = true

	// Update raw string to include 'x' marker
	raw := m.todos[m.cursor].Raw
	if !strings.HasPrefix(raw, "x ") {
		m.todos[m.cursor].Raw = "x " + raw
	}

	// Save to file
	if err := todo.SaveToFile(m.filename, m.todos); err != nil {
		return m, nil
	}

	// Return to normal mode
	m.mode = ModeNormal
	m.commandInput.Blur()

	return m, nil
}

// cmdDelete deletes the current task
func (m Model) cmdDelete(args string) (Model, tea.Cmd) {
	if len(m.todos) == 0 {
		return m, nil
	}

	// Remove the item at cursor
	m.todos = append(m.todos[:m.cursor], m.todos[m.cursor+1:]...)

	// Adjust cursor if needed
	if m.cursor >= len(m.todos) && m.cursor > 0 {
		m.cursor--
	}

	// Save to file
	if err := todo.SaveToFile(m.filename, m.todos); err != nil {
		return m, nil
	}

	// Return to normal mode
	m.mode = ModeNormal
	m.commandInput.Blur()

	return m, nil
}

// handleCommandMode handles key presses in command mode
func (m Model) handleCommandMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = ModeNormal
		m.commandInput.Blur()
		return m, nil
	case "enter":
		// Execute the command
		return m.executeCommand()
	}

	// Let the textinput handle the key
	var cmd tea.Cmd
	m.commandInput, cmd = m.commandInput.Update(msg)
	return m, cmd
}

// handleInsertMode handles key presses in insert mode
func (m Model) handleInsertMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = ModeNormal
		m.insertInput.Blur()
		return m, nil
	case "enter":
		// Add the todo
		description := m.insertInput.Value()
		if description != "" {
			newItem := todo.Item{
				Raw:         description,
				Description: description,
				Completed:   false,
				Contexts:    []string{},
				Projects:    []string{},
			}
			m.todos = append(m.todos, newItem)

			// Save to file
			if err := todo.SaveToFile(m.filename, m.todos); err == nil {
				m.cursor = len(m.todos) - 1
			}
		}

		// Return to normal mode
		m.mode = ModeNormal
		m.insertInput.Blur()
		return m, nil
	}

	// Let the textinput handle the key
	var cmd tea.Cmd
	m.insertInput, cmd = m.insertInput.Update(msg)
	return m, cmd
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

	// Command/Insert input prompt
	if m.mode == ModeCommand {
		s += "\n" + m.commandInput.View()
	} else if m.mode == ModeInsert {
		s += "\n" + m.insertInput.View()
	}

	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Padding(1, 1)

	help := ""
	switch m.mode {
	case ModeNormal:
		help = "i/enter: insert | v: visual | :: command | j/k: navigate | q: quit"
	case ModeInsert:
		help = "enter: add todo | esc: cancel"
	case ModeCommand:
		help = "add <task> | edit <new text> | done | delete/del | enter: execute | esc: cancel"
	case ModeVisual:
		help = "esc: back to normal mode"
	}

	s += "\n" + helpStyle.Render(help)

	return s
}
