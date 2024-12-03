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
func (vp *Viewport) SetSize(width, height int) {
	vp.width = width
	vp.height = height
}

// Size returns the viewport dimensions
func (vp *Viewport) Size() (width, height int) {
	return vp.width, vp.height
}

// Offset returns the viewport's offset in the buffer
func (vp *Viewport) Offset() types.Position {
	return vp.offset
}

// SetScrollOff sets the number of lines to keep visible above/below cursor
func (vp *Viewport) SetScrollOff(lines int) {
	vp.scrollOff = lines
}

// ScrollTo scrolls the viewport to ensure the target position is visible
func (vp *Viewport) ScrollTo(pos types.Position) {
	// Vertical scrolling
	if pos.Line() < vp.offset.Line()+vp.scrollOff {
		line := int(math.Max(0, float64(pos.Line()-vp.scrollOff)))
		vp.offset = buffer.NewPosition(line, vp.offset.Column())
	} else if pos.Line() >= vp.offset.Line()+vp.height-vp.scrollOff {
		line := pos.Line() - vp.height + vp.scrollOff + 1
		vp.offset = buffer.NewPosition(line, vp.offset.Column())
	}

	// Horizontal scrolling
	if pos.Column() < vp.offset.Column() {
		col := int(math.Max(0, float64(pos.Column())))
		vp.offset = buffer.NewPosition(vp.offset.Line(), col)
	} else if pos.Column() >= vp.offset.Column()+vp.width {
		col := pos.Column() - vp.width + 1
		vp.offset = buffer.NewPosition(vp.offset.Line(), col)
	}
}

// SetCursor sets the cursor position
func (vp *Viewport) SetCursor(pos types.Position) {
	vp.cursor = pos
	vp.ScrollTo(pos)
}

// Cursor returns the current cursor position
func (vp *Viewport) Cursor() types.Position {
	return vp.cursor
}

// ToggleCursor toggles the cursor visibility.
func (vp *Viewport) ToggleCursor() {
	vp.showCursor = true
}

// IsPositionVisible returns true if the position is within the viewport
func (vp *Viewport) IsPositionVisible(pos types.Position) bool {
	isLineVisible := pos.Line() >= vp.offset.Line() && pos.Line() < vp.offset.Line()+vp.height
	isColumnVisible := pos.Column() >= vp.offset.Column() && pos.Column() < vp.offset.Column()+vp.width
	return isLineVisible && isColumnVisible
}

// VisibleLines returns the range of visible line numbers
func (vp *Viewport) VisibleLines() (start, end int) {
	return vp.offset.Line(), vp.offset.Line() + vp.height
}

// VisibleColumns returns the range of visible column numbers
func (vp *Viewport) VisibleColumns() (start, end int) {
	return vp.offset.Column(), vp.offset.Column() + vp.width
}

// SetMode sets the current editor mode
func (vp *Viewport) SetMode(mode state.Mode) {
	vp.mode = mode
	// Disable cursor blinking in visual mode since we show selection
	if mode == state.VisualMode {
		vp.showCursor = true // Keep cursor always visible in visual mode
	}
}

// SetHighlightManager sets the highlight manager for the viewport
func (vp *Viewport) SetHighlightManager(hm types.HighlightManager) {
	vp.highlightManager = hm
}

// StyleRange represents a range of text with a specific style
type StyleRange struct {
	start, end int
	style      lipgloss.Style
}

// ViewportCursor represents a cursor position and its visual style in the viewport
type ViewportCursor struct {
	position int
	style    lipgloss.Style
}

// prepareVisibleContent handles content preparation and padding
func (vp *Viewport) prepareVisibleContent(content string) string {
	if content == "" {
		return strings.Repeat(" ", vp.width)
	}

	// Calculate visible portion of the line
	startCol, endCol := vp.VisibleColumns()
	if startCol >= len(content) {
		return strings.Repeat(" ", vp.width)
	}

	// Ensure we don't go past the end of the content
	if endCol > len(content) {
		endCol = len(content)
	}

	// Get visible portion of the line
	visibleContent := content
	if startCol < len(content) {
		if endCol <= len(content) {
			visibleContent = content[startCol:endCol]
		} else {
			visibleContent = content[startCol:]
		}
	}

	// Pad to viewport width
	if len(visibleContent) < vp.width {
		visibleContent += strings.Repeat(" ", vp.width-len(visibleContent))
	} else if len(visibleContent) > vp.width {
		visibleContent = visibleContent[:vp.width]
	}

	return visibleContent
}

// getCursorStyle returns the appropriate cursor style for current mode
func (vp *Viewport) getCursorStyle() lipgloss.Style {
	switch vp.mode {
	case state.InsertMode:
		return vp.cursorStyle.GetInsertStyle()
	case state.CommandMode:
		return vp.cursorStyle.GetCommandStyle()
	case state.VisualMode:
		return vp.cursorStyle.GetVisualStyle()
	default:
		return vp.cursorStyle.GetNormalStyle()
	}
}

// getHighlightBounds calculates viewport-relative bounds for a highlight
func (vp *Viewport) getHighlightBounds(h types.Highlight, lineNumber, viewportOffset, contentLength int) (start, end int) {
	startPos := h.GetStartPosition()
	endPos := h.GetEndPosition()

	// Ensure start position is before end position
	if endPos.Line() < startPos.Line() || (endPos.Line() == startPos.Line() && endPos.Column() < startPos.Column()) {
		startPos, endPos = endPos, startPos
	}

	// Calculate bounds
	start = 0
	end = contentLength

	if lineNumber == startPos.Line() {
		start = startPos.Column() - viewportOffset
	}
	if lineNumber == endPos.Line() {
		end = endPos.Column() - viewportOffset + 1
	}

	// Skip if completely outside visible area
	if end <= 0 || start >= contentLength {
		return 0, 0
	}

	// Clip to visible area
	if start < 0 {
		start = 0
	}
	if end > contentLength {
		end = contentLength
	}

	return start, end
}

// getHighlightRanges returns styled ranges for highlights in a line
func (vp *Viewport) getHighlightRanges(lineNumber, viewportOffset, contentLength int) []StyleRange {
	if vp.highlightManager == nil {
		return nil
	}

	highlights := vp.highlightManager.GetForLine(lineNumber)
	if len(highlights) == 0 {
		return nil
	}

	ranges := make([]StyleRange, 0, len(highlights))
	style := highlight.NewStyle()

	for _, h := range highlights {
		start, end := vp.getHighlightBounds(h, lineNumber, viewportOffset, contentLength)
		if start >= end {
			continue
		}
		ranges = append(ranges, StyleRange{
			start: start,
			end:   end,
			style: style.GetStyle(h.GetType()),
		})
	}
	return ranges
}

// collectStyleRanges gathers cursor and highlight ranges for a line
func (vp *Viewport) collectStyleRanges(lineNumber int, contentLength int) ([]StyleRange, *ViewportCursor) {
	ranges := []StyleRange{}
	viewportOffset := vp.offset.Column()
	cursorCol := vp.cursor.Column() - viewportOffset

	// Get highlight styles
	if highlights := vp.getHighlightRanges(lineNumber, viewportOffset, contentLength); len(highlights) > 0 {
		ranges = append(ranges, highlights...)
	}

	// Add cursor style if on this line
	var cursor *ViewportCursor
	if vp.showCursor && vp.cursor.Line() == lineNumber && cursorCol >= 0 && cursorCol < contentLength {
		cursor = &ViewportCursor{
			position: cursorCol,
			style:    vp.getCursorStyle(),
		}
	}

	return ranges, cursor
}

// mergeHighlightRanges handles merging overlapping highlights
func (vp *Viewport) mergeHighlightRanges(ranges []StyleRange) []StyleRange {
	if len(ranges) == 0 {
		return nil
	}

	// Sort highlight ranges by start position
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].start < ranges[j].start
	})

	// Merge overlapping ranges
	merged := make([]StyleRange, 0, len(ranges))
	current := ranges[0]
	for i := 1; i < len(ranges); i++ {
		if ranges[i].start <= current.end {
			if ranges[i].end > current.end {
				current.end = ranges[i].end
			}
		} else {
			merged = append(merged, current)
			current = ranges[i]
		}
	}
	merged = append(merged, current)

	return merged
}

// applyStyles handles the actual style application
func (vp *Viewport) applyStyles(lineContent string, highlightRanges []StyleRange, cursor *ViewportCursor) string {
	var result strings.Builder
	result.Grow(len(lineContent) * 2)

	lastPos := 0
	for _, r := range highlightRanges {
		// Ensure valid bounds
		if r.start < 0 {
			r.start = 0
		}
		if r.end > len(lineContent) {
			r.end = len(lineContent)
		}
		if r.start >= r.end || r.start >= len(lineContent) {
			continue
		}

		// Add unstyled text before this range
		if r.start > lastPos {
			result.WriteString(lineContent[lastPos:r.start])
		}

		// If cursor is in this range, split the styling
		if cursor != nil && cursor.position >= r.start && cursor.position < r.end {
			// Write highlighted text before cursor
			if cursor.position > r.start {
				result.WriteString(r.style.Render(lineContent[r.start:cursor.position]))
			}
			// Write cursor
			result.WriteString(cursor.style.Render(string(lineContent[cursor.position])))
			// Write highlighted text after cursor
			if cursor.position+1 < r.end {
				result.WriteString(r.style.Render(lineContent[cursor.position+1:r.end]))
			}
		} else {
			// Add styled text without cursor
			result.WriteString(r.style.Render(lineContent[r.start:r.end]))
		}
		lastPos = r.end
	}

	// Handle cursor if it's after all highlights
	if cursor != nil && (len(highlightRanges) == 0 || cursor.position >= lastPos) {
		// Add any unstyled text before cursor
		if cursor.position > lastPos {
			result.WriteString(lineContent[lastPos:cursor.position])
		}
		// Add cursor
		result.WriteString(cursor.style.Render(string(lineContent[cursor.position])))
		lastPos = cursor.position + 1
	}

	// Add remaining unstyled text
	if lastPos < len(lineContent) {
		result.WriteString(lineContent[lastPos:])
	}

	return result.String()
}

// renderLine processes and formats a single line of content
func (vp *Viewport) renderLine(content string, lineNumber int) string {
	visibleContent := vp.prepareVisibleContent(content)
	highlightRanges, cursor := vp.collectStyleRanges(lineNumber, len(visibleContent))

	if len(highlightRanges) == 0 && cursor == nil {
		return visibleContent
	}

	mergedHighlights := vp.mergeHighlightRanges(highlightRanges)
	return vp.applyStyles(visibleContent, mergedHighlights, cursor)
}

// createEmptyLine creates an empty line with proper formatting
func (vp *Viewport) createEmptyLine() string {
	return strings.Repeat(" ", vp.width)
}

// Render returns the visible portion of the buffer content
func (vp *Viewport) Render(buf types.Buffer) []string {
	startLine, endLine := vp.VisibleLines()
	renderedLines := make([]string, 0, vp.height)

	// Render visible lines
	for line := startLine; line < endLine && line < buf.LineCount(); line++ {
		content, err := buf.GetLine(line)
		if err != nil {
			continue
		}
		renderedLines = append(renderedLines, vp.renderLine(content, line))
	}

	// Fill remaining space with empty lines
	emptyLine := vp.createEmptyLine()
	for len(renderedLines) < vp.height {
		renderedLines = append(renderedLines, emptyLine)
	}

	return renderedLines
}

// CenterOn centers the viewport on the given position
func (vp *Viewport) CenterOn(pos types.Position) {
	line := int(math.Max(0, float64(pos.Line()-vp.height/2)))
	column := int(math.Max(0, float64(pos.Column()-vp.width/2)))
	vp.offset = buffer.NewPosition(line, column)
}

func (vp *Viewport) BufferToViewportPosition(pos types.Position) (x, y int) {
	return pos.Column() - vp.offset.Column(), pos.Line() - vp.offset.Line()
}

func (vp *Viewport) ViewportToBufferPosition(x, y int) types.Position {
	return buffer.NewPosition(
		y+vp.offset.Line(),
		x+vp.offset.Column(),
	)
}

// ScrollUp scrolls the viewport up by the specified number of lines
func (vp *Viewport) ScrollUp(lines int) {
	line := int(math.Max(0, float64(vp.offset.Line()-lines)))
	vp.offset = buffer.NewPosition(line, vp.offset.Column())
}

// ScrollDown scrolls the viewport down by the specified number of lines
func (vp *Viewport) ScrollDown(lines int) {
	vp.offset = buffer.NewPosition(vp.offset.Line()+lines, vp.offset.Column())
}

// ScrollLeft scrolls the viewport left by the specified number of columns
func (vp *Viewport) ScrollLeft(cols int) {
	vp.offset = buffer.NewPosition(vp.offset.Line(), int(math.Max(0, float64(vp.offset.Column()-cols))))
}

// ScrollRight scrolls the viewport right by the specified number of columns
func (vp *Viewport) ScrollRight(cols int) {
	vp.offset = buffer.NewPosition(vp.offset.Line(), vp.offset.Column()+cols)
}
