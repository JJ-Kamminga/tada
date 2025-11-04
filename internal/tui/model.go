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

// TodoWithIndex wraps a todo item with its index in the main todos slice
type TodoWithIndex struct {
	Item  todo.Item
	Index int
}

// ContextList represents a group of todos for a specific context
type ContextList struct {
	Context string
	Todos   []TodoWithIndex
}

// groupTodosByContext groups todos by their contexts
func groupTodosByContext(todos []todo.Item) []ContextList {
	contextMap := make(map[string][]TodoWithIndex)

	// Group todos by context
	for i, item := range todos {
		todoWithIdx := TodoWithIndex{Item: item, Index: i}
		if len(item.Contexts) == 0 {
			// No context, put in "No Context" list
			contextMap["No Context"] = append(contextMap["No Context"], todoWithIdx)
		} else {
			// Add to each context it belongs to
			for _, context := range item.Contexts {
				contextMap[context] = append(contextMap[context], todoWithIdx)
			}
		}
	}

	// Convert map to sorted list
	var lists []ContextList

	// Add "No Context" first if it exists
	if items, ok := contextMap["No Context"]; ok {
		lists = append(lists, ContextList{Context: "No Context", Todos: items})
		delete(contextMap, "No Context")
	}

	// Add other contexts in alphabetical order
	contexts := make([]string, 0, len(contextMap))
	for context := range contextMap {
		contexts = append(contexts, context)
	}

	// Simple sort
	for i := 0; i < len(contexts); i++ {
		for j := i + 1; j < len(contexts); j++ {
			if contexts[i] > contexts[j] {
				contexts[i], contexts[j] = contexts[j], contexts[i]
			}
		}
	}

	for _, context := range contexts {
		lists = append(lists, ContextList{Context: context, Todos: contextMap[context]})
	}

	return lists
}

// Model represents the application state
type Model struct {
	todos         []todo.Item
	contextLists  []ContextList   // Grouped todos by context
	listCursor    int             // Which context list is selected
	itemCursor    int             // Which item in the current list is selected
	mode          Mode
	filename      string
	width         int
	height        int
	commandInput  textinput.Model // Text input for command mode
	insertInput   textinput.Model // Text input for insert mode
	editingIndex  int             // Index of the todo being edited in insert mode (-1 if adding new)
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
	insInput.Placeholder = "edit todo description..."
	insInput.Prompt = "> "
	insInput.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("green")).Bold(true)
	insInput.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	insInput.CharLimit = 500

	return Model{
		todos:        todos,
		contextLists: groupTodosByContext(todos),
		listCursor:   0,
		itemCursor:   0,
		mode:         ModeNormal,
		filename:     filename,
		commandInput: cmdInput,
		insertInput:  insInput,
		editingIndex: -1,
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

		// Get current todo to edit
		_, idx := m.getCurrentTodo()
		if idx != -1 {
			// Prefill with current todo description
			m.editingIndex = idx
			m.insertInput.SetValue(m.todos[idx].Raw)
		} else {
			// No todo selected, will add new one
			m.editingIndex = -1
			m.insertInput.Reset()
		}

		m.insertInput.Focus()
		return m, textinput.Blink
	case "v":
		m.mode = ModeVisual
	case "q":
		return m, tea.Quit
	case "up", "k":
		// Move up within current list
		if m.itemCursor > 0 {
			m.itemCursor--
		} else if m.listCursor > 0 {
			// Move to previous list
			m.listCursor--
			if len(m.contextLists) > 0 && m.listCursor < len(m.contextLists) {
				m.itemCursor = len(m.contextLists[m.listCursor].Todos) - 1
			}
		}
	case "down", "j":
		// Move down within current list
		if len(m.contextLists) > 0 && m.listCursor < len(m.contextLists) {
			if m.itemCursor < len(m.contextLists[m.listCursor].Todos)-1 {
				m.itemCursor++
			} else if m.listCursor < len(m.contextLists)-1 {
				// Move to next list
				m.listCursor++
				m.itemCursor = 0
			}
		}
	case "left", "h":
		// Move to previous list
		if m.listCursor > 0 {
			m.listCursor--
			// Adjust item cursor if needed
			if len(m.contextLists) > 0 && m.itemCursor >= len(m.contextLists[m.listCursor].Todos) {
				m.itemCursor = len(m.contextLists[m.listCursor].Todos) - 1
			}
		}
	case "right", "l":
		// Move to next list
		if len(m.contextLists) > 0 && m.listCursor < len(m.contextLists)-1 {
			m.listCursor++
			// Adjust item cursor if needed
			if m.itemCursor >= len(m.contextLists[m.listCursor].Todos) {
				m.itemCursor = len(m.contextLists[m.listCursor].Todos) - 1
			}
		}
	}

	return m, nil
}

// getCurrentTodo returns the currently selected todo item and its index in the todos slice
func (m Model) getCurrentTodo() (*todo.Item, int) {
	if len(m.contextLists) == 0 || m.listCursor >= len(m.contextLists) {
		return nil, -1
	}

	currentList := m.contextLists[m.listCursor]
	if m.itemCursor >= len(currentList.Todos) {
		return nil, -1
	}

	todoWithIdx := currentList.Todos[m.itemCursor]
	idx := todoWithIdx.Index

	// Validate index is still valid
	if idx < 0 || idx >= len(m.todos) {
		return nil, -1
	}

	return &m.todos[idx], idx
}

// refreshContextLists rebuilds the context lists after todos change
func (m *Model) refreshContextLists() {
	m.contextLists = groupTodosByContext(m.todos)

	// Ensure cursors are still valid
	if m.listCursor >= len(m.contextLists) {
		m.listCursor = len(m.contextLists) - 1
	}
	if m.listCursor < 0 {
		m.listCursor = 0
	}

	if len(m.contextLists) > 0 && m.listCursor < len(m.contextLists) {
		if m.itemCursor >= len(m.contextLists[m.listCursor].Todos) {
			m.itemCursor = len(m.contextLists[m.listCursor].Todos) - 1
		}
		if m.itemCursor < 0 {
			m.itemCursor = 0
		}
	}
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

	// Parse the new todo to extract contexts
	newItem := todo.Parse(description)

	m.todos = append(m.todos, newItem)

	// Save to file
	if err := todo.SaveToFile(m.filename, m.todos); err != nil {
		// TODO: Handle error (could add error message to model)
		return m, nil
	}

	// Refresh context lists
	m.refreshContextLists()

	// Return to normal mode
	m.mode = ModeNormal
	m.commandInput.Blur()

	return m, nil
}

// cmdEdit edits the current task
func (m Model) cmdEdit(newDescription string) (Model, tea.Cmd) {
	if newDescription == "" {
		return m, nil
	}

	// Get current todo
	_, idx := m.getCurrentTodo()
	if idx == -1 {
		return m, nil
	}

	// Parse the new description to extract contexts
	updatedItem := todo.Parse(newDescription)

	// Update the item in todos
	m.todos[idx] = updatedItem

	// Save to file
	if err := todo.SaveToFile(m.filename, m.todos); err != nil {
		return m, nil
	}

	// Refresh context lists
	m.refreshContextLists()

	// Return to normal mode
	m.mode = ModeNormal
	m.commandInput.Blur()

	return m, nil
}

// cmdDone marks the current task as complete
func (m Model) cmdDone(args string) (Model, tea.Cmd) {
	// Get current todo
	_, idx := m.getCurrentTodo()
	if idx == -1 {
		return m, nil
	}

	// Mark as completed
	m.todos[idx].Completed = true

	// Update raw string to include 'x' marker
	raw := m.todos[idx].Raw
	if !strings.HasPrefix(raw, "x ") {
		m.todos[idx].Raw = "x " + raw
	}

	// Save to file
	if err := todo.SaveToFile(m.filename, m.todos); err != nil {
		return m, nil
	}

	// Refresh context lists
	m.refreshContextLists()

	// Return to normal mode
	m.mode = ModeNormal
	m.commandInput.Blur()

	return m, nil
}

// cmdDelete deletes the current task
func (m Model) cmdDelete(args string) (Model, tea.Cmd) {
	// Get current todo
	_, idx := m.getCurrentTodo()
	if idx == -1 {
		return m, nil
	}

	// Remove the item
	m.todos = append(m.todos[:idx], m.todos[idx+1:]...)

	// Save to file
	if err := todo.SaveToFile(m.filename, m.todos); err != nil {
		return m, nil
	}

	// Refresh context lists
	m.refreshContextLists()

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
		m.editingIndex = -1
		return m, nil
	case "enter":
		description := m.insertInput.Value()
		if description != "" {
			if m.editingIndex >= 0 && m.editingIndex < len(m.todos) {
				// Edit existing todo
				updatedItem := todo.Parse(description)
				m.todos[m.editingIndex] = updatedItem
			} else {
				// Add new todo
				newItem := todo.Parse(description)
				m.todos = append(m.todos, newItem)
			}

			// Save to file
			if err := todo.SaveToFile(m.filename, m.todos); err == nil {
				// Refresh context lists
				m.refreshContextLists()
			}
		}

		// Return to normal mode
		m.mode = ModeNormal
		m.insertInput.Blur()
		m.editingIndex = -1
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

	// Todo lists by context
	if len(m.todos) == 0 {
		s += "  No todos yet. Press 'i' to add one!\n"
	} else {
		// Render each context list
		for listIdx, contextList := range m.contextLists {
			// Context header
			contextHeaderStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("magenta")).
				Underline(true)

			if listIdx == m.listCursor {
				contextHeaderStyle = contextHeaderStyle.Foreground(lipgloss.Color("cyan"))
			}

			s += contextHeaderStyle.Render(fmt.Sprintf("@%s (%d)", contextList.Context, len(contextList.Todos))) + "\n"

			// Render todos in this context
			for itemIdx, todoWithIdx := range contextList.Todos {
				cursor := "  "
				if listIdx == m.listCursor && itemIdx == m.itemCursor {
					cursor = "> "
				}

				// Style the item
				itemStyle := lipgloss.NewStyle()
				if todoWithIdx.Item.Completed {
					itemStyle = itemStyle.Foreground(lipgloss.Color("240")).Strikethrough(true)
				}

				s += fmt.Sprintf("%s%s\n", cursor, itemStyle.Render(todoWithIdx.Item.Description))
			}
			s += "\n" // Space between lists
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
		help = "i/enter: edit todo | v: visual | :: command | j/k: up/down | h/l: prev/next list | q: quit"
	case ModeInsert:
		help = "enter: save changes | esc: cancel"
	case ModeCommand:
		help = "add <task> | edit <new text> | done | delete/del | enter: execute | esc: cancel"
	case ModeVisual:
		help = "esc: back to normal mode"
	}

	s += "\n" + helpStyle.Render(help)

	return s
}
