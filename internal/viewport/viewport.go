package viewport

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gunererd/grease/internal/buffer"
)

// Viewport represents a view into a portion of the buffer
type Viewport struct {
	width     int
	height    int
	offset    buffer.Position // Top-left position of viewport in buffer
	scrollOff int             // Number of lines to keep visible above/below cursor
	cursor    buffer.Position // Current cursor position
}

// New creates a new viewport with the given dimensions
func New(width, height int) *Viewport {
	return &Viewport{
		width:     width,
		height:    height,
		offset:    buffer.Position{Line: 0, Column: 0},
		scrollOff: 5, // Default scroll offset
		cursor:    buffer.Position{Line: 0, Column: 0},
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
	// Vertical scrolling with scroll-off
	if pos.Line < v.offset.Line+v.scrollOff {
		v.offset.Line = int(math.Max(0, float64(pos.Line-v.scrollOff)))
	} else if pos.Line >= v.offset.Line+v.height-v.scrollOff {
		v.offset.Line = pos.Line - v.height + v.scrollOff + 1
	}

	// Horizontal scrolling with padding
	padding := 5 // Keep 5 characters visible on either side when possible
	if pos.Column < v.offset.Column+padding {
		v.offset.Column = int(math.Max(0, float64(pos.Column-padding)))
	} else if pos.Column >= v.offset.Column+v.width-padding {
		v.offset.Column = pos.Column - v.width + padding + 1
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

// renderLine processes and formats a single line of content
func (v *Viewport) renderLine(content string, lineNum int, maxLineNumWidth int, showLineNumbers bool) string {
	// Handle horizontal scrolling
	runes := []rune(content)
	colStart := v.offset.Column
	colEnd := colStart + v.width

	// Adjust for line numbers if showing them
	if showLineNumbers {
		colEnd -= maxLineNumWidth + 1 // +1 for spacing
	}

	lineContent := ""
	if colStart < len(runes) {
		if colEnd > len(runes) {
			colEnd = len(runes)
		}
		lineContent = string(runes[colStart:colEnd])
	}

	// Add line numbers if requested
	if showLineNumbers {
		lineNum := fmt.Sprintf("%*d ", maxLineNumWidth, lineNum+1)
		lineContent = lineNum + lineContent
	}

	// If this is the cursor line and column, style it
	if lineNum == v.cursor.Line {
		// Convert lineContent to runes for proper unicode handling
		contentRunes := []rune(lineContent)
		cursorCol := v.cursor.Column - colStart
		if showLineNumbers {
			cursorCol += maxLineNumWidth + 1
		}

		if cursorCol >= 0 && cursorCol < len(contentRunes) {
			// Split the line into before, cursor, and after parts
			before := string(contentRunes[:cursorCol])
			cursorChar := string(contentRunes[cursorCol])
			after := ""
			if cursorCol+1 < len(contentRunes) {
				after = string(contentRunes[cursorCol+1:])
			}

			// Style the cursor character with inverse colors
			styledCursor := lipgloss.NewStyle().Reverse(true).Render(cursorChar)
			lineContent = before + styledCursor + after
		}
	}

	// Pad line to viewport width
	if lipgloss.Width(lineContent) < v.width {
		lineContent += strings.Repeat(" ", v.width-lipgloss.Width(lineContent))
	}

	return lineContent
}

// createEmptyLine creates an empty line with proper formatting
func (v *Viewport) createEmptyLine(maxLineNumWidth int, showLineNumbers bool) string {
	if showLineNumbers {
		return strings.Repeat(" ", maxLineNumWidth) + " " + strings.Repeat(" ", v.width-maxLineNumWidth-1)
	}
	return strings.Repeat(" ", v.width)
}

// View returns the visible portion of the buffer content with optional line numbers
func (v *Viewport) View(buf *buffer.Buffer, showLineNumbers bool) []string {
	start, end := v.GetVisibleLines()
	result := make([]string, 0, v.height)
	maxLineNumWidth := len(string(end))

	// Render visible lines
	for line := start; line < end && line < buf.LineCount(); line++ {
		content, err := buf.GetLine(line)
		if err != nil {
			continue
		}
		result = append(result, v.renderLine(content, line, maxLineNumWidth, showLineNumbers))
	}

	// Fill remaining space with empty lines
	emptyLine := v.createEmptyLine(maxLineNumWidth, showLineNumbers)
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
	v.offset.Line -= lines
	if v.offset.Line < 0 {
		v.offset.Line = 0
	}
}

// ScrollDown scrolls the viewport down by the specified number of lines
func (v *Viewport) ScrollDown(lines int) {
	v.offset.Line += lines
}

// ScrollLeft scrolls the viewport left by the specified number of columns
func (v *Viewport) ScrollLeft(cols int) {
	v.offset.Column -= cols
	if v.offset.Column < 0 {
		v.offset.Column = 0
	}
}

// ScrollRight scrolls the viewport right by the specified number of columns
func (v *Viewport) ScrollRight(cols int) {
	v.offset.Column += cols
}
