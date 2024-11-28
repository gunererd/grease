package styles

import "github.com/charmbracelet/lipgloss"

var (
	baseStatusStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1)

	// Mode indicator styles
	normalModeStyle = baseStatusStyle.
			Bold(true).
			Background(lipgloss.Color("#005f87")).
			Foreground(lipgloss.Color("#ffffff"))

	insertModeStyle = baseStatusStyle.
			Bold(true).
			Background(lipgloss.Color("#5f8700")).
			Foreground(lipgloss.Color("#ffffff"))

	commandModeStyle = baseStatusStyle.
				Bold(true).
				Background(lipgloss.Color("#af5f00")).
				Foreground(lipgloss.Color("#ffffff"))

	// Position and progress indicator styles
	positionStyle = baseStatusStyle.
			Background(lipgloss.Color("#303030")).
			Foreground(lipgloss.Color("#d0d0d0"))

	progressStyle = baseStatusStyle.
			Background(lipgloss.Color("#444444")).
			Foreground(lipgloss.Color("#d0d0d0"))
)

// StatusLineStyle provides styling functions for the status line
type StatusLineStyle struct{}

// NewStatusLineStyle creates a new StatusLineStyle provider
func NewStatusLineStyle() *StatusLineStyle {
	return &StatusLineStyle{}
}

// GetModeStyle returns the appropriate style for the given mode string
func (s *StatusLineStyle) GetModeStyle(mode string) lipgloss.Style {
	switch mode {
	case "NORMAL":
		return normalModeStyle
	case "INSERT":
		return insertModeStyle
	case "COMMAND":
		return commandModeStyle
	default:
		return baseStatusStyle
	}
}

// GetPositionStyle returns the style for position indicators
func (s *StatusLineStyle) GetPositionStyle() lipgloss.Style {
	return positionStyle
}

// GetProgressStyle returns the style for the progress indicator
func (s *StatusLineStyle) GetProgressStyle() lipgloss.Style {
	return progressStyle
}
