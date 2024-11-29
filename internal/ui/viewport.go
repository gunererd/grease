package ui

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gunererd/grease/internal/buffer"
	"github.com/gunererd/grease/internal/state"
)

// Viewport represents a view into a portion of the buffer
type Viewport struct {
	width       int
	height      int
	offset      buffer.Position // Top-left position of viewport in buffer
	scrollOff   int             // Number of lines to keep visible above/below cursor
	cursor      buffer.Position // Current cursor position
	showCursor  bool            // Controls cursor blinking state
	cursorStyle *buffer.CursorStyle
	mode        state.Mode
}

// NewViewport creates a new viewport with the given dimensions
func NewViewport(width, height int) *Viewport {
	return &Viewport{
		width:       width,
		height:      height,
		offset:      buffer.Position{},
		scrollOff:   5,
		cursor:      buffer.Position{},
		showCursor:  true,
		cursorStyle: buffer.NewCursorStyle(),
		mode:        state.NormalMode,
	}
}

// SetSize sets the viewport dimensions
func (v *Viewport) SetSize(width, height int) {
	v.width = width
	v.height = height
}

// GetSize returns the viewport dimensions
func (v *Viewport) GetSize() (width, height int) {
	return v.width, v.height
}

// GetOffset returns the viewport's offset in the buffer
func (v *Viewport) GetOffset() buffer.Position {
	return v.offset
}

// SetScrollOff sets the number of lines to keep visible above/below cursor
func (v *Viewport) SetScrollOff(lines int) {
	v.scrollOff = lines
}

// ScrollTo scrolls the viewport to ensure the target position is visible
func (v *Viewport) ScrollTo(pos buffer.Position) {
	// Vertical scrolling
	if pos.Line < v.offset.Line+v.scrollOff {
		v.offset.Line = int(math.Max(0, float64(pos.Line-v.scrollOff)))
	} else if pos.Line >= v.offset.Line+v.height-v.scrollOff {
		v.offset.Line = pos.Line - v.height + v.scrollOff + 1
	}

	// Horizontal scrolling
	if pos.Column < v.offset.Column {
		v.offset.Column = int(math.Max(0, float64(pos.Column)))
	} else if pos.Column >= v.offset.Column+v.width {
		v.offset.Column = pos.Column - v.width + 1
	}
}

// SetCursor sets the cursor position
func (v *Viewport) SetCursor(pos buffer.Position) {
	v.cursor = pos
	v.ScrollTo(pos)
}

// GetCursor returns the current cursor position
func (v *Viewport) GetCursor() buffer.Position {
	return v.cursor
}

// ToggleCursor toggles the cursor visibility state for blinking effect
func (v *Viewport) ToggleCursor() {
	v.showCursor = !v.showCursor
}

// IsPositionVisible returns true if the position is within the viewport
func (v *Viewport) IsPositionVisible(pos buffer.Position) bool {
	return pos.Line >= v.offset.Line &&
		pos.Line < v.offset.Line+v.height &&
		pos.Column >= v.offset.Column &&
		pos.Column < v.offset.Column+v.width
}

// GetVisibleLines returns the range of visible line numbers
func (v *Viewport) GetVisibleLines() (start, end int) {
	return v.offset.Line, v.offset.Line + v.height
}

// GetVisibleColumns returns the range of visible column numbers
func (v *Viewport) GetVisibleColumns() (start, end int) {
	return v.offset.Column, v.offset.Column + v.width
}

// SetMode sets the current editor mode
func (v *Viewport) SetMode(mode state.Mode) {
	v.mode = mode
}

// renderLine processes and formats a single line of content
func (v *Viewport) renderLine(content string, lineNum int) string {
	// Create base content
	visibleContent := content
	if content == "" {
		visibleContent = strings.Repeat(" ", v.width)
	} else {
		// Calculate visible portion of the line
		start, end := v.GetVisibleColumns()
		if start >= len(content) {
			visibleContent = strings.Repeat(" ", v.width)
		} else {
			// Ensure we don't go past the end of the content
			if end > len(content) {
				end = len(content)
			}

			// Get visible portion of the line
			if start < len(content) {
				if end <= len(content) {
					visibleContent = content[start:end]
				} else {
					visibleContent = content[start:]
				}
			}

			// Pad to viewport width
			if len(visibleContent) < v.width {
				visibleContent += strings.Repeat(" ", v.width-len(visibleContent))
			} else if len(visibleContent) > v.width {
				visibleContent = visibleContent[:v.width]
			}
		}
	}

	// Add cursor if needed
	if v.showCursor && v.cursor.Line == lineNum {
		cursorCol := v.cursor.Column - v.offset.Column
		if cursorCol >= 0 && cursorCol < v.width {
			// Split the line at cursor position
			before := visibleContent[:cursorCol]
			cursor := " "
			if cursorCol < len(visibleContent) {
				cursor = string(visibleContent[cursorCol])
			}
			after := ""
			if cursorCol+1 < len(visibleContent) {
				after = visibleContent[cursorCol+1:]
			}

			// Get cursor style based on mode
			var cursorStyle lipgloss.Style
			switch v.mode {
			case state.InsertMode:
				cursorStyle = v.cursorStyle.GetInsertStyle()
			case state.CommandMode:
				cursorStyle = v.cursorStyle.GetCommandStyle()
			default:
				cursorStyle = v.cursorStyle.GetNormalStyle()
			}

			// Style the cursor
			styledCursor := cursorStyle.Render(cursor)

			// Combine all parts
			visibleContent = before + styledCursor + after
		}
	}

	return visibleContent
}

// createEmptyLine creates an empty line with proper formatting
func (v *Viewport) createEmptyLine() string {
	return strings.Repeat(" ", v.width)
}

// View returns the visible portion of the buffer content
func (v *Viewport) View(buf *buffer.Buffer) []string {
	start, end := v.GetVisibleLines()
	result := make([]string, 0, v.height)

	// Render visible lines
	for line := start; line < end && line < buf.LineCount(); line++ {
		content, err := buf.GetLine(line)
		if err != nil {
			continue
		}
		result = append(result, v.renderLine(content, line))
	}

	// Fill remaining space with empty lines
	emptyLine := v.createEmptyLine()
	for len(result) < v.height {
		result = append(result, emptyLine)
	}

	return result
}

// CenterOn centers the viewport on the given position
func (v *Viewport) CenterOn(pos buffer.Position) {
	v.offset.Line = int(math.Max(0, float64(pos.Line-v.height/2)))
	v.offset.Column = int(math.Max(0, float64(pos.Column-v.width/2)))
}

// GetRelativePosition converts a buffer position to viewport coordinates
func (v *Viewport) GetRelativePosition(pos buffer.Position) (x, y int) {
	return pos.Column - v.offset.Column, pos.Line - v.offset.Line
}

// GetAbsolutePosition converts viewport coordinates to buffer position
func (v *Viewport) GetAbsolutePosition(x, y int) buffer.Position {
	return buffer.Position{
		Line:   y + v.offset.Line,
		Column: x + v.offset.Column,
	}
}

// ScrollUp scrolls the viewport up by the specified number of lines
func (v *Viewport) ScrollUp(lines int) {
	v.offset.Line = int(math.Max(0, float64(v.offset.Line-lines)))
}

// ScrollDown scrolls the viewport down by the specified number of lines
func (v *Viewport) ScrollDown(lines int) {
	v.offset.Line += lines
}

// ScrollLeft scrolls the viewport left by the specified number of columns
func (v *Viewport) ScrollLeft(cols int) {
	v.offset.Column = int(math.Max(0, float64(v.offset.Column-cols)))
}

// ScrollRight scrolls the viewport right by the specified number of columns
func (v *Viewport) ScrollRight(cols int) {
	v.offset.Column += cols
}
