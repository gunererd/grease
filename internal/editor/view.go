package editor

import (
	"strings"

	"github.com/gunererd/grease/internal/buffer"
	"github.com/gunererd/grease/internal/navigator"
	"github.com/gunererd/grease/internal/ui"
)

// View handles all UI-related concerns
type View struct {
	Width        int
	Height       int
	ScrollOffset int
	ViewHeight   int
	navigator    *navigator.Navigator
	state        *State // Reference to editor state
}

// NewView creates a new View instance
func NewView(n *navigator.Navigator, s *State) *View {
	return &View{
		navigator:  n,
		ViewHeight: 10, // Default height, will be updated when window size is received
		state:      s,
	}
}

// UpdateSize updates the view dimensions
func (v *View) UpdateSize(width, height int) {
	v.Width = width
	v.Height = height
	v.ViewHeight = height - 1 // Reserve 1 line for statusline
}

// EnsureVisible adjusts scroll offset to keep cursor in view
func (v *View) EnsureVisible(cursorRow, totalEntries int) {
	if cursorRow < v.ScrollOffset {
		v.ScrollOffset = cursorRow
	}
	if cursorRow >= v.ScrollOffset+v.ViewHeight {
		v.ScrollOffset = cursorRow - v.ViewHeight + 1
	}
}

// GetVisibleRange returns the range of entries that should be visible
func (v *View) GetVisibleRange(totalEntries int) (start, end int) {
	start = v.ScrollOffset
	end = min(start+v.ViewHeight, totalEntries)
	return start, end
}

// RenderEntry renders a single entry with appropriate styling
func (v *View) RenderEntry(entry buffer.Entry, row int) string {
	line := v.navigator.Buffer.GetLine(row)
	if entry.IsDir {
		line += "/"
	}

	// Convert line to runes for proper character handling
	runes := []rune(line)
	var result strings.Builder

	// Handle empty line case
	if len(runes) == 0 {
		return line
	}

	// Get cursor position for this line
	cursorCol := -1
	if row == v.navigator.Cursor.Row {
		cursorCol = v.navigator.Cursor.Col
		if cursorCol >= len(runes) {
			cursorCol = len(runes) - 1
		}
	}

	// Render each character with appropriate style
	for col, r := range runes {
		char := string(r)
		switch {
		case col == cursorCol:
			char = ui.CursorStyle.Render(char)
		case v.navigator.Buffer.IsLineModified(row):
			char = ui.ModifiedStyle.Render(char)
		case v.state.IsSelected(row, col):
			char = ui.SelectedStyle.Render(char)
		}
		result.WriteString(char)
	}

	return result.String()
}

// RenderStatusLine renders the status line at the bottom
func (v *View) RenderStatusLine(mode, path string, cursorRow int) string {
	return ui.RenderStatusLine(mode, path, v.Width, cursorRow, 0)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
