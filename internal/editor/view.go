package editor

import (
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
}

// NewView creates a new View instance
func NewView(n *navigator.Navigator) *View {
	return &View{
		navigator:  n,
		ViewHeight: 10, // Default height, will be updated when window size is received
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
func (v *View) RenderEntry(entry buffer.Entry, idx int, isSelected bool) string {
	line := v.navigator.Buffer.GetLine(idx)
	if entry.IsDir {
		line += "/"
	}

	// First apply selection style if selected
	if isSelected {
		line = ui.SelectedStyle.Render(line)
		return line
	}

	// Then handle cursor if this is the cursor row
	if idx == v.navigator.Cursor.Row && len(line) > 0 {
		runes := []rune(line)
		cursorCol := v.navigator.Cursor.Col
		if cursorCol >= len(runes) {
			cursorCol = len(runes) - 1
		}
		if cursorCol < 0 {
			cursorCol = 0
		}

		// Split the line into three parts: before cursor, cursor char, and after cursor
		before := string(runes[:cursorCol])
		cursor := string(runes[cursorCol])
		after := ""
		if cursorCol < len(runes)-1 {
			after = string(runes[cursorCol+1:])
		}

		// Apply styles
		if v.navigator.Buffer.IsLineModified(idx) {
			line = ui.ModifiedStyle.Render(before) +
				ui.CursorStyle.Render(cursor) +
				ui.ModifiedStyle.Render(after)
		} else {
			line = before + ui.CursorStyle.Render(cursor) + after
		}
		return line
	}

	// Finally, show modified style if line is modified
	if v.navigator.Buffer.IsLineModified(idx) {
		line = ui.ModifiedStyle.Render(line)
	}

	return line
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
