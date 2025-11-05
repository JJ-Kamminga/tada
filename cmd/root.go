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

var rootCmd = &cobra.Command{
	Use:   "tada",
	Short: "A vim-inspired todo list manager",
	Long:  `tada is a terminal-based todo list manager using the todo.txt format with vim-inspired keybindings.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the todo directory from config
		todoDir, err := config.GetTodoDir()
		if err != nil {
			fmt.Println("Error: No todo directory configured.")
			fmt.Println()
			fmt.Println("To get started, set your todo directory:")
			fmt.Println("  tada config set dir /path/to/your/todo/directory")
			fmt.Println()
			fmt.Println("Example:")
			fmt.Println("  tada config set dir ~/.tada")
			os.Exit(1)
		}

		// Ensure the directory exists
		if err := os.MkdirAll(todoDir, 0755); err != nil {
			fmt.Println("Error creating todo directory:", err)
			os.Exit(1)
		}

		// Get the full path to todo.txt
		todoFile := filepath.Join(todoDir, "todo.txt")

		// If todo.txt doesn't exist, create an empty one
		if _, err := os.Stat(todoFile); os.IsNotExist(err) {
			if err := os.WriteFile(todoFile, []byte(""), 0644); err != nil {
				fmt.Println("Error creating todo.txt:", err)
				os.Exit(1)
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
