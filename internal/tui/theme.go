package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme holds all color and style configurations
// This structure makes it easy to load colors from a config file in the future
type Theme struct {
	// Colors
	Primary          lipgloss.Color
	Secondary        lipgloss.Color
	Accent           lipgloss.Color
	Success          lipgloss.Color
	Warning          lipgloss.Color
	Danger           lipgloss.Color
	Muted            lipgloss.Color
	Background       lipgloss.Color
	Foreground       lipgloss.Color
	Border           lipgloss.Color
	SelectedBorder   lipgloss.Color
	CompletedText    lipgloss.Color

	// Mode-specific colors
	NormalModeColor  lipgloss.Color
	InsertModeColor  lipgloss.Color
	CommandModeColor lipgloss.Color
	VisualModeColor  lipgloss.Color
}

// DefaultTheme returns the default color scheme
// In the future, this could be replaced with LoadThemeFromConfig()
func DefaultTheme() Theme {
	return Theme{
		Primary:          lipgloss.Color("39"),  // Bright blue
		Secondary:        lipgloss.Color("170"), // Purple
		Accent:           lipgloss.Color("205"), // Pink
		Success:          lipgloss.Color("42"),  // Green
		Warning:          lipgloss.Color("214"), // Orange
		Danger:           lipgloss.Color("196"), // Red
		Muted:           lipgloss.Color("241"),  // Gray
		Background:       lipgloss.Color("235"), // Dark gray
		Foreground:       lipgloss.Color("255"), // White
		Border:           lipgloss.Color("240"), // Border gray
		SelectedBorder:   lipgloss.Color("39"),  // Bright blue
		CompletedText:    lipgloss.Color("240"), // Gray for completed items

		NormalModeColor:  lipgloss.Color("39"),  // Blue
		InsertModeColor:  lipgloss.Color("42"),  // Green
		CommandModeColor: lipgloss.Color("214"), // Orange
		VisualModeColor:  lipgloss.Color("170"), // Purple
	}
}

// Styles holds all the styled components
type Styles struct {
	Theme Theme

	// Header
	AppTitle       lipgloss.Style

	// Context headers
	ContextHeader         lipgloss.Style
	ContextHeaderActive   lipgloss.Style

	// Todo items
	TodoNormal     lipgloss.Style
	TodoCompleted  lipgloss.Style
	TodoCursor     lipgloss.Style

	// Mode indicator
	ModeNormal     lipgloss.Style
	ModeInsert     lipgloss.Style
	ModeCommand    lipgloss.Style
	ModeVisual     lipgloss.Style

	// Help text
	HelpText       lipgloss.Style

	// Input prompts
	CommandPrompt  lipgloss.Style
	InsertPrompt   lipgloss.Style
	InputText      lipgloss.Style

	// Layout
	ContentBox     lipgloss.Style
	ContextListBox lipgloss.Style
}

// NewStyles creates a new Styles instance with the given theme
func NewStyles(theme Theme) Styles {
	return Styles{
		Theme: theme,

		// Header - bold, centered, with gradient effect
		AppTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Primary).
			Background(theme.Background).
			Padding(0, 2).
			MarginBottom(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Primary),

		// Context headers
		ContextHeader: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Secondary).
			Background(theme.Background).
			Padding(0, 1).
			MarginTop(0).
			MarginBottom(0).
			Border(lipgloss.Border{
				Top:    "─",
				Bottom: "─",
				Left:   "",
				Right:  "",
			}).
			BorderForeground(theme.Border),

		ContextHeaderActive: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Accent).
			Background(theme.Background).
			Padding(0, 1).
			MarginTop(0).
			MarginBottom(0).
			Border(lipgloss.Border{
				Top:    "━",
				Bottom: "━",
				Left:   "",
				Right:  "",
			}).
			BorderForeground(theme.SelectedBorder),

		// Todo items
		TodoNormal: lipgloss.NewStyle().
			Foreground(theme.Foreground).
			Padding(0, 1),

		TodoCompleted: lipgloss.NewStyle().
			Foreground(theme.CompletedText).
			Strikethrough(true).
			Padding(0, 1),

		TodoCursor: lipgloss.NewStyle().
			Foreground(theme.Accent).
			Bold(true),

		// Mode indicators with colored backgrounds
		ModeNormal: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("0")).
			Background(theme.NormalModeColor).
			Padding(0, 2),

		ModeInsert: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("0")).
			Background(theme.InsertModeColor).
			Padding(0, 2),

		ModeCommand: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("0")).
			Background(theme.CommandModeColor).
			Padding(0, 2),

		ModeVisual: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("0")).
			Background(theme.VisualModeColor).
			Padding(0, 2),

		// Help text
		HelpText: lipgloss.NewStyle().
			Foreground(theme.Muted).
			Padding(1, 2).
			Border(lipgloss.NormalBorder()).
			BorderForeground(theme.Border).
			BorderTop(true).
			MarginTop(1),

		// Input prompts
		CommandPrompt: lipgloss.NewStyle().
			Foreground(theme.CommandModeColor).
			Bold(true),

		InsertPrompt: lipgloss.NewStyle().
			Foreground(theme.InsertModeColor).
			Bold(true),

		InputText: lipgloss.NewStyle().
			Foreground(theme.Foreground),

		// Layout
		ContentBox: lipgloss.NewStyle().
			Padding(0, 2),

		ContextListBox: lipgloss.NewStyle().
			Padding(0, 1).
			MarginBottom(1),
	}
}
