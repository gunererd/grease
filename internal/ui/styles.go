package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Cursor styles
	CursorStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#ebdbb2")).
			Foreground(lipgloss.Color("#282828"))

	// Selected line style (for visual mode)
	SelectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#504945"))

	// Style for showing unsaved changes
	ModifiedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#fe8019")) // Orange color for modified text
)
