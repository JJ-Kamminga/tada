package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"tada/internal/config"
	"tada/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var todoFile string

var rootCmd = &cobra.Command{
	Use:     "tada",
	Aliases: []string{"td"},
	Short:   "A vim-inspired todo list manager",
	Long:    `tada is a terminal-based todo list manager using the todo.txt format with vim-inspired keybindings.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the todo file path (priority: flag > config > default)
		var err error
		todoFile, err = config.GetTodoFilePath(todoFile)
		if err != nil {
			fmt.Println("Error getting todo file path:", err)
			os.Exit(1)
		}

		// Ensure the directory exists
		dir := filepath.Dir(todoFile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Println("Error creating todo directory:", err)
			os.Exit(1)
		}

		// If todo.txt doesn't exist, copy the example one
		if _, err := os.Stat(todoFile); os.IsNotExist(err) {
			// Try to copy from current directory
			if _, err := os.Stat("todo.txt"); err == nil {
				data, err := os.ReadFile("todo.txt")
				if err == nil {
					_ = os.WriteFile(todoFile, data, 0644) // Best effort copy of example file
				}
			}
		}

		// Start the TUI
		m := tui.NewModel(todoFile)
		p := tea.NewProgram(m, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&todoFile, "file", "f", "", "Path to todo.txt file (default: ~/.tada/todo.txt)")
}
