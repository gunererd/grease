package styles

import "github.com/charmbracelet/lipgloss"

// CursorStyle provides styling functions for the cursor
type CursorStyle struct{}

// NewCursorStyle creates a new CursorStyle provider
func NewCursorStyle() *CursorStyle {
	return &CursorStyle{}
}

// GetNormalStyle returns the cursor style for normal mode
func (s *CursorStyle) GetNormalStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color("7")).
		Foreground(lipgloss.Color("0"))
}

// GetInsertStyle returns the cursor style for insert mode
func (s *CursorStyle) GetInsertStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color("2")). // Green background
		Foreground(lipgloss.Color("0"))  // Black text
}

// GetCommandStyle returns the cursor style for command mode
func (s *CursorStyle) GetCommandStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color("5")). // Magenta background
		Foreground(lipgloss.Color("0"))  // Black text
}
