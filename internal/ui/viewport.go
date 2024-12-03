package ui

import (
	"math"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gunererd/grease/internal/buffer"
	"github.com/gunererd/grease/internal/highlight"
	"github.com/gunererd/grease/internal/state"
	"github.com/gunererd/grease/internal/types"
)

// Viewport represents a view into a portion of the buffer
type Viewport struct {
	width            int
	height           int
	offset           types.Position // Top-left position of viewport in buffer
	scrollOff        int            // Number of lines to keep visible above/below cursor
	cursor           types.Position // Current cursor position
	showCursor       bool
	cursorStyle      *buffer.CursorStyle
	mode             state.Mode
	highlightManager types.HighlightManager
}

// NewViewport creates a new viewport with the given dimensions
func NewViewport(width, height int) types.Viewport {
	return &Viewport{
		width:            width,
		height:           height,
		offset:           buffer.NewPosition(0, 0),
		scrollOff:        5,
		cursor:           buffer.NewPosition(0, 0),
		showCursor:       true,
		cursorStyle:      buffer.NewCursorStyle(),
		mode:             state.NormalMode,
		highlightManager: nil,
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
func (v *Viewport) GetOffset() types.Position {
	return v.offset
}

// SetScrollOff sets the number of lines to keep visible above/below cursor
func (v *Viewport) SetScrollOff(lines int) {
	v.scrollOff = lines
}

// ScrollTo scrolls the viewport to ensure the target position is visible
func (v *Viewport) ScrollTo(pos types.Position) {
	// Vertical scrolling
	if pos.Line() < v.offset.Line()+v.scrollOff {
		line := int(math.Max(0, float64(pos.Line()-v.scrollOff)))
		v.offset = buffer.NewPosition(line, v.offset.Column())
	} else if pos.Line() >= v.offset.Line()+v.height-v.scrollOff {
		line := pos.Line() - v.height + v.scrollOff + 1
		v.offset = buffer.NewPosition(line, v.offset.Column())
	}

	// Horizontal scrolling
	if pos.Column() < v.offset.Column() {
		col := int(math.Max(0, float64(pos.Column())))
		v.offset = buffer.NewPosition(v.offset.Line(), col)
	} else if pos.Column() >= v.offset.Column()+v.width {
		col := pos.Column() - v.width + 1
		v.offset = buffer.NewPosition(v.offset.Line(), col)
	}
}

// SetCursor sets the cursor position
func (v *Viewport) SetCursor(pos types.Position) {
	v.cursor = pos
	v.ScrollTo(pos)
}

// GetCursor returns the current cursor position
func (v *Viewport) GetCursor() types.Position {
	return v.cursor
}

// ToggleCursor toggles the cursor visibility.
func (v *Viewport) ToggleCursor() {
	v.showCursor = true
}

// IsPositionVisible returns true if the position is within the viewport
func (v *Viewport) IsPositionVisible(pos types.Position) bool {
	return pos.Line() >= v.offset.Line() &&
		pos.Line() < v.offset.Line()+v.height &&
		pos.Column() >= v.offset.Column() &&
		pos.Column() < v.offset.Column()+v.width
}

// GetVisibleLines returns the range of visible line numbers
func (v *Viewport) GetVisibleLines() (start, end int) {
	return v.offset.Line(), v.offset.Line() + v.height
}

// GetVisibleColumns returns the range of visible column numbers
func (v *Viewport) GetVisibleColumns() (start, end int) {
	return v.offset.Column(), v.offset.Column() + v.width
}

// SetMode sets the current editor mode
func (v *Viewport) SetMode(mode state.Mode) {
	v.mode = mode
	// Disable cursor blinking in visual mode since we show selection
	if mode == state.VisualMode {
		v.showCursor = true // Keep cursor always visible in visual mode
	}
}

// SetHighlightManager sets the highlight manager for the viewport
func (v *Viewport) SetHighlightManager(hm types.HighlightManager) {
	v.highlightManager = hm
}

// styledRange represents a range of text with a specific style
type styledRange struct {
	start, end int
	style      lipgloss.Style
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

	// Get all styled ranges for this line
	ranges := make([]styledRange, 0, 8) // Pre-allocate for common case

	// Add cursor style if on this line
	currentCol := v.offset.Column()
	cursorCol := v.cursor.Column() - currentCol
	if v.showCursor && v.cursor.Line() == lineNum && cursorCol >= 0 && cursorCol < len(visibleContent) {
		var cursorStyle lipgloss.Style
		switch v.mode {
		case state.InsertMode:
			cursorStyle = v.cursorStyle.GetInsertStyle()
		case state.CommandMode:
			cursorStyle = v.cursorStyle.GetCommandStyle()
		default:
			cursorStyle = v.cursorStyle.GetNormalStyle()
		}
		ranges = append(ranges, styledRange{
			start: cursorCol,
			end:   cursorCol + 1,
			style: cursorStyle,
		})
	}

	// Add highlight styles
	if v.highlightManager != nil {
		highlights := v.highlightManager.GetForLine(lineNum)
		if len(highlights) > 0 {
			style := highlight.NewStyle()
			for _, h := range highlights {
				start := h.GetStartPosition()
				end := h.GetEndPosition()

				// Ensure start position is before end position
				if end.Line() < start.Line() || (end.Line() == start.Line() && end.Column() < start.Column()) {
					start, end = end, start
				}

				// Convert buffer positions to viewport-relative positions
				startCol := 0
				endCol := len(visibleContent)

				// If this is the start line, start from the start column
				if lineNum == start.Line() {
					startCol = start.Column() - currentCol
				}

				// If this is the end line, end at the end column
				if lineNum == end.Line() {
					endCol = end.Column() - currentCol + 1
				}

				// Skip if completely outside visible area
				if endCol <= 0 || startCol >= len(visibleContent) {
					continue
				}

				// Clip to visible area
				if startCol < 0 {
					startCol = 0
				}
				if endCol > len(content) {
					endCol = len(content)
				}

				// Skip invalid ranges
				if startCol >= endCol {
					continue
				}

				ranges = append(ranges, styledRange{
					start: startCol,
					end:   endCol,
					style: style.GetStyle(h.GetType()),
				})
			}
		}
	}

	// If no styles to apply, return as is
	if len(ranges) == 0 {
		return visibleContent
	}

	// Sort ranges by start position and merge overlapping ranges
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].start < ranges[j].start
	})

	// Merge overlapping ranges
	merged := make([]styledRange, 0, len(ranges))
	current := ranges[0]
	for i := 1; i < len(ranges); i++ {
		if ranges[i].start <= current.end {
			// Ranges overlap, extend current range if needed
			if ranges[i].end > current.end {
				current.end = ranges[i].end
			}
		} else {
			// No overlap, add current range and start new one
			merged = append(merged, current)
			current = ranges[i]
		}
	}
	merged = append(merged, current)

	// Apply styles efficiently
	var result strings.Builder
	result.Grow(len(visibleContent) * 2) // Estimate space needed including style codes

	lastPos := 0
	for _, r := range merged {
		// Ensure valid bounds
		if r.start < 0 {
			r.start = 0
		}
		if r.end > len(visibleContent) {
			r.end = len(visibleContent)
		}
		if r.start >= r.end || r.start >= len(visibleContent) {
			continue
		}

		// Add unstyled text before this range
		if r.start > lastPos {
			result.WriteString(visibleContent[lastPos:r.start])
		}
		// Add styled text
		result.WriteString(r.style.Render(visibleContent[r.start:r.end]))
		lastPos = r.end
	}
	// Add remaining unstyled text
	if lastPos < len(visibleContent) {
		result.WriteString(visibleContent[lastPos:])
	}

	return result.String()
}

// createEmptyLine creates an empty line with proper formatting
func (v *Viewport) createEmptyLine() string {
	return strings.Repeat(" ", v.width)
}

// View returns the visible portion of the buffer content
func (v *Viewport) View(buf types.Buffer) []string {
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
func (v *Viewport) CenterOn(pos types.Position) {
	line := int(math.Max(0, float64(pos.Line()-v.height/2)))
	column := int(math.Max(0, float64(pos.Column()-v.width/2)))
	v.offset = buffer.NewPosition(line, column)
}

// GetRelativePosition converts a buffer position to viewport coordinates
func (v *Viewport) GetRelativePosition(pos types.Position) (x, y int) {
	return pos.Column() - v.offset.Column(), pos.Line() - v.offset.Line()
}

// GetAbsolutePosition converts viewport coordinates to buffer position
func (v *Viewport) GetAbsolutePosition(x, y int) types.Position {
	return buffer.NewPosition(
		y+v.offset.Line(),
		x+v.offset.Column(),
	)
}

// ScrollUp scrolls the viewport up by the specified number of lines
func (v *Viewport) ScrollUp(lines int) {
	line := int(math.Max(0, float64(v.offset.Line()-lines)))
	v.offset = buffer.NewPosition(line, v.offset.Column())
}

// ScrollDown scrolls the viewport down by the specified number of lines
func (v *Viewport) ScrollDown(lines int) {
	v.offset = buffer.NewPosition(v.offset.Line()+lines, v.offset.Column())
}

// ScrollLeft scrolls the viewport left by the specified number of columns
func (v *Viewport) ScrollLeft(cols int) {
	v.offset = buffer.NewPosition(v.offset.Line(), int(math.Max(0, float64(v.offset.Column()-cols))))
}

// ScrollRight scrolls the viewport right by the specified number of columns
func (v *Viewport) ScrollRight(cols int) {
	v.offset = buffer.NewPosition(v.offset.Line(), v.offset.Column()+cols)
}
