package ui

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	statusLineStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Width(100).
			Height(1)

	modeStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			MarginRight(1)

	// Mode-specific styles
	normalModeStyle = modeStyle.
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#98971a")) // Green

	insertModeStyle = modeStyle.
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#458588")) // Blue

	visualModeStyle = modeStyle.
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#d79921")) // Yellow

	pathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ebdbb2")).
			MarginLeft(1)

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ebdbb2"))
)

// RenderStatusLine creates a styled status line
func RenderStatusLine(mode, path string, width int, row, col int) string {
	// Adjust status line width to terminal width
	statusLineStyle = statusLineStyle.Width(width)

	var modeIndicator string
	switch mode {
	case "NORMAL":
		modeIndicator = normalModeStyle.Render(mode)
	case "INSERT":
		modeIndicator = insertModeStyle.Render(mode)
	case "VISUAL":
		modeIndicator = visualModeStyle.Render(mode)
	default:
		modeIndicator = modeStyle.Render(mode)
	}

	pathDisplay := filepath.Clean(path)
	if usr, err := user.Current(); err == nil {
		pathDisplay = strings.Replace(pathDisplay, usr.HomeDir, "~", 1)
	}

	leftContent := lipgloss.JoinHorizontal(
		lipgloss.Left,
		modeIndicator,
		pathStyle.Render(pathDisplay),
	)

	cursorPos := fmt.Sprintf("%d,%d  ", row+1, col+1)
	rightContent := cursorStyle.Render(cursorPos)

	// Calculate padding needed between left and right content
	padding := width - lipgloss.Width(leftContent) - lipgloss.Width(rightContent)
	if padding < 1 {
		padding = 1
	}

	content := lipgloss.JoinHorizontal(
		lipgloss.Left,
		leftContent,
		strings.Repeat(" ", padding),
		rightContent,
	)

	return statusLineStyle.Render(content)
}
