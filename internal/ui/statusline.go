package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gunererd/grease/internal/types"
)

// StatusLine represents the editor's status line
type StatusLine struct {
	styles *StatusLineStyle
}

// NewStatusLine creates a new StatusLine
func NewStatusLine() types.StatusLine {
	return &StatusLine{
		styles: NewStatusLineStyle(),
	}
}

// Render renders the status line with the given editor state
func (s *StatusLine) Render(mode string, cursor types.Cursor, bufferLineCount int, viewX, viewY int, width int) string {
	pos := cursor.GetPosition()
	progress := int(float64(pos.Line()+1) / float64(bufferLineCount) * 100)

	// Build the status line components using the styles
	modeIndicator := s.styles.GetModeStyle(mode).Render(mode)
	bufferPos := s.styles.GetPositionStyle().Render(fmt.Sprintf("Buf[%d,%d]", pos.Line()+1, pos.Column()+1))
	viewPos := s.styles.GetPositionStyle().Render(fmt.Sprintf("View[%d,%d]", viewY+1, viewX+1))
	progressIndicator := s.styles.GetProgressStyle().Render(fmt.Sprintf("%d%%", progress))

	// Calculate the total width of all fixed components
	fixedWidth := lipgloss.Width(modeIndicator) +
		lipgloss.Width(bufferPos) +
		lipgloss.Width(viewPos) +
		lipgloss.Width(progressIndicator) +
		3 // for spaces between components

	// Create a flexible space that fills the remaining width
	flexSpace := strings.Repeat(" ", max(0, width-fixedWidth))

	// Join all components with the flexible space before the progress
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		modeIndicator,
		" ",
		bufferPos,
		" ",
		viewPos,
		flexSpace,
		progressIndicator,
	)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
