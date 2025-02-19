package ui

import (
	"math"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gunererd/grease/internal/editor/buffer"
	"github.com/gunererd/grease/internal/editor/highlight"
	"github.com/gunererd/grease/internal/editor/state"
	"github.com/gunererd/grease/internal/editor/types"
)

// Add at the top with other type definitions
type CursorInfo struct {
	show    bool
	pos     types.Position
	style   *buffer.CursorStyle
	primary bool
}

// Viewport represents a view into a portion of the buffer
type Viewport struct {
	width            int
	height           int
	offset           types.Position // Top-left position of viewport in buffer
	scrollOff        int            // Number of lines to keep visible above/below cursor
	cursors          []CursorInfo   // New field for multiple cursors
	cursor           types.Position // Keep temporarily for backwards compatibility
	showCursor       bool           // Keep temporarily for backwards compatibility
	cursorStyle      *buffer.CursorStyle
	mode             state.Mode
	highlightManager types.HighlightManager
}

// NewViewport creates a new viewport with the given dimensions
func NewViewport(width, height int) types.Viewport {
	cursorStyle := buffer.NewCursorStyle()
	return &Viewport{
		width:     width,
		height:    height,
		offset:    buffer.NewPosition(0, 0),
		scrollOff: 5,
		cursors: []CursorInfo{
			{
				show:    true,
				pos:     buffer.NewPosition(0, 0),
				style:   cursorStyle,
				primary: true,
			},
		},
		cursor:           buffer.NewPosition(0, 0), // Keep temporarily
		showCursor:       true,                     // Keep temporarily
		cursorStyle:      cursorStyle,
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

func (vp *Viewport) ScrollOff() int {
	return vp.scrollOff
}

func (vp *Viewport) ScrollTo(pos types.Position, bufferLineCount int) {
	// Vertical scrolling
	if pos.Line() < vp.offset.Line()+vp.scrollOff {
		// When scrolling up, respect scrollOff
		line := int(math.Max(0, float64(pos.Line()-vp.scrollOff)))
		vp.offset = buffer.NewPosition(line, vp.offset.Column())
	} else if pos.Line() >= vp.offset.Line()+vp.height-vp.scrollOff {
		// When scrolling down, check if we're near buffer end
		line := pos.Line() - vp.height + vp.scrollOff + 1

		// If this would show past buffer end, adjust to show last line at bottom
		if line+vp.height > bufferLineCount {
			line = bufferLineCount - vp.height
			if line < 0 {
				line = 0
			}
		}

		vp.offset = buffer.NewPosition(line, vp.offset.Column())
	}

	// Horizontal scrolling (unchanged)
	if pos.Column() < vp.offset.Column() {
		col := int(math.Max(0, float64(pos.Column())))
		vp.offset = buffer.NewPosition(vp.offset.Line(), col)
	} else if pos.Column() >= vp.offset.Column()+vp.width {
		col := pos.Column() - vp.width + 1
		vp.offset = buffer.NewPosition(vp.offset.Line(), col)
	}
}

// SetCursor sets the cursor position
func (vp *Viewport) SetCursor(pos types.Position, bufferLineCount int) {
	vp.cursor = pos
	vp.ScrollTo(pos, bufferLineCount)
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
func (vp *Viewport) collectStyleRanges(lineNumber int, contentLength int) ([]StyleRange, []ViewportCursor) {
	ranges := []StyleRange{}
	viewportOffset := vp.offset.Column()

	// Get highlight styles
	if highlights := vp.getHighlightRanges(lineNumber, viewportOffset, contentLength); len(highlights) > 0 {
		ranges = append(ranges, highlights...)
	}

	// Add cursor styles for all visible cursors on this line
	var cursors []ViewportCursor
	for _, cursorInfo := range vp.cursors {
		if cursorInfo.pos.Line() == lineNumber {
			cursorCol := cursorInfo.pos.Column() - viewportOffset
			if cursorCol >= 0 && cursorCol < contentLength {
				cursors = append(cursors, ViewportCursor{
					position: cursorCol,
					style:    vp.getCursorStyle(),
				})
			}
		}
	}

	return ranges, cursors
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
func (vp *Viewport) applyStyles(lineContent string, highlightRanges []StyleRange, cursors []ViewportCursor) string {
	var result strings.Builder
	result.Grow(len(lineContent) * 2)

	lastPos := 0
	cursorIndex := 0

	// Sort cursors by position
	sort.Slice(cursors, func(i, j int) bool {
		return cursors[i].position < cursors[j].position
	})

	for _, r := range highlightRanges {
		// Handle cursors before this range
		for cursorIndex < len(cursors) && cursors[cursorIndex].position < r.start {
			if cursors[cursorIndex].position > lastPos {
				result.WriteString(lineContent[lastPos:cursors[cursorIndex].position])
			}
			result.WriteString(cursors[cursorIndex].style.Render(string(lineContent[cursors[cursorIndex].position])))
			lastPos = cursors[cursorIndex].position + 1
			cursorIndex++
		}

		// Add unstyled text before this range
		if r.start > lastPos {
			result.WriteString(lineContent[lastPos:r.start])
		}

		// Handle cursors within this range
		rangeStart := r.start
		for cursorIndex < len(cursors) && cursors[cursorIndex].position < r.end {
			if cursors[cursorIndex].position > rangeStart {
				result.WriteString(r.style.Render(lineContent[rangeStart:cursors[cursorIndex].position]))
			}
			result.WriteString(cursors[cursorIndex].style.Render(string(lineContent[cursors[cursorIndex].position])))
			rangeStart = cursors[cursorIndex].position + 1
			cursorIndex++
		}

		// Add remaining styled text in this range
		if rangeStart < r.end {
			result.WriteString(r.style.Render(lineContent[rangeStart:r.end]))
		}
		lastPos = r.end
	}

	// Handle remaining cursors
	for cursorIndex < len(cursors) {
		if cursors[cursorIndex].position > lastPos {
			result.WriteString(lineContent[lastPos:cursors[cursorIndex].position])
		}
		result.WriteString(cursors[cursorIndex].style.Render(string(lineContent[cursors[cursorIndex].position])))
		lastPos = cursors[cursorIndex].position + 1
		cursorIndex++
	}

	// Add remaining unstyled text
	if lastPos < len(lineContent) {
		result.WriteString(lineContent[lastPos:])
	}

	return result.String()
}

// renderLine processes and formats a single line of content
func (vp *Viewport) renderLine(content string, lineNumber int) string {
	// Ensure empty lines have at least one space for cursor rendering
	if len(content) == 0 {
		content = " "
	}

	visibleContent := vp.prepareVisibleContent(content)
	highlightRanges, cursors := vp.collectStyleRanges(lineNumber, len(content))

	if len(highlightRanges) == 0 && len(cursors) == 0 {
		return visibleContent
	}

	mergedHighlights := vp.mergeHighlightRanges(highlightRanges)
	return vp.applyStyles(visibleContent, mergedHighlights, cursors)
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

func (vp *Viewport) ScrollDown(lines int, bufferLineCount int) {
	// Calculate maximum allowed scroll position that shows last line at bottom of viewport
	maxLine := bufferLineCount - vp.height
	if maxLine < 0 {
		maxLine = 0
	}

	// Calculate new scroll position
	newLine := vp.offset.Line() + lines

	// If scrolling would show past buffer end, adjust to show last line at bottom
	if newLine+vp.height > bufferLineCount {
		newLine = maxLine
	}

	vp.offset = buffer.NewPosition(newLine, vp.offset.Column())
}

// ScrollLeft scrolls the viewport left by the specified number of columns
func (vp *Viewport) ScrollLeft(cols int) {
	vp.offset = buffer.NewPosition(vp.offset.Line(), int(math.Max(0, float64(vp.offset.Column()-cols))))
}

// ScrollRight scrolls the viewport right by the specified number of columns
func (vp *Viewport) ScrollRight(cols int) {
	vp.offset = buffer.NewPosition(vp.offset.Line(), vp.offset.Column()+cols)
}

func (vp *Viewport) SyncCursors(bufferCursors []types.Cursor, bufferLineCount int) {
	vp.cursors = make([]CursorInfo, len(bufferCursors))
	for i, cursor := range bufferCursors {
		vp.cursors[i] = CursorInfo{
			show:    true,
			pos:     cursor.GetPosition(),
			style:   vp.cursorStyle,
			primary: i == 0, // First cursor is primary, handle it in buffer later
		}
	}

	// Ensure at least primary cursor is visible
	if len(vp.cursors) > 0 {
		vp.ScrollTo(vp.cursors[0].pos, bufferLineCount)
		// Update legacy cursor field temporarily
		vp.cursor = vp.cursors[0].pos
	}
}

func (vp *Viewport) GetVisibleCursors() []CursorInfo {
	visible := make([]CursorInfo, 0)
	for _, cursor := range vp.cursors {
		if vp.IsPositionVisible(cursor.pos) {
			visible = append(visible, cursor)
		}
	}
	return visible
}

func (vp *Viewport) ScrollHalfPageUp() {
	// If we're already at the top, just move cursor to top
	cursor := vp.cursors[0]
	if vp.offset.Line() <= 0 {
		cursor.pos = buffer.NewPosition(0, cursor.pos.Column())
		vp.cursors[0] = cursor
		return
	}

	vp.ScrollUp(vp.height / 2)
	if !vp.IsPositionVisible(vp.cursors[0].pos) {
		_, endLine := vp.VisibleLines()
		// Move cursor to last visible line minus scrollOff
		cursor.pos = buffer.NewPosition(endLine-vp.scrollOff-1, cursor.pos.Column())
		vp.cursors[0] = cursor
	}
}

func (vp *Viewport) ScrollHalfPageDown(bufferLineCount int) {
	vp.ScrollDown(vp.height/2, bufferLineCount)
	cursor := vp.cursors[0]
	if !vp.IsPositionVisible(cursor.pos) {
		startLine, _ := vp.VisibleLines()
		cursor.pos = buffer.NewPosition(startLine+vp.scrollOff, cursor.pos.Column())
		vp.cursors[0] = cursor
	}
}
