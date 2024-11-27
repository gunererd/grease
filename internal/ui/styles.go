package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Cursor styles
	CursorStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#458588")).
			Foreground(lipgloss.Color("#000000"))

	// Selected line style (for visual mode)
	SelectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#d79921")).
			Foreground(lipgloss.Color("#000000"))
)
