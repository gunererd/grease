package selection

import (
	"strings"

	"github.com/gunererd/grease/internal/editor/highlight"
	"github.com/gunererd/grease/internal/editor/types"
)

type manager struct {
	highlightManager types.HighlightManager
	selections       map[int]*Selection
	nextID           int
}

func NewManager(hm types.HighlightManager) Manager {
	return &manager{
		highlightManager: hm,
		selections:       make(map[int]*Selection),
		nextID:           1,
	}
}

func (m *manager) StartSelection(pos types.Position) *Selection {
	sel := &Selection{
		id:     m.nextID,
		anchor: pos,
		head:   pos,
		mode:   Character,
	}
	m.nextID++

	// Create highlight for this selection
	sel.highlightID = m.highlightManager.Add(
		highlight.CreateVisualHighlight(pos, pos),
	)

	m.selections[sel.id] = sel
	return sel
}

func (m *manager) UpdateSelection(sel *Selection, newHead types.Position) {
	if sel == nil {
		return
	}

	sel.head = newHead
	m.highlightManager.Update(
		sel.highlightID,
		highlight.CreateVisualHighlight(sel.anchor, sel.head),
	)
}

func (m *manager) ClearSelection(sel *Selection) {
	if sel == nil {
		return
	}

	m.highlightManager.Remove(sel.highlightID)
	delete(m.selections, sel.id)
}

func (m *manager) ClearAllSelections() {
	for _, sel := range m.selections {
		m.highlightManager.Remove(sel.highlightID)
	}
	m.selections = make(map[int]*Selection)
}

func (m *manager) GetSelections() []*Selection {
	selections := make([]*Selection, 0, len(m.selections))
	for _, sel := range m.selections {
		selections = append(selections, sel)
	}
	return selections
}

func (m *manager) GetSelectedText(sel *Selection, buf types.Buffer) string {
	if sel == nil {
		return ""
	}

	// Reference the extractSingleLine and extractMultiLine functions
	// from clipboard/yank.go for implementation
	lines := bufferToLines(buf)

	if sel.anchor.Line() == sel.head.Line() {
		return extractSingleLine(lines, sel.anchor, sel.head)
	}
	return extractMultiLine(lines, sel.anchor, sel.head)
}

// Helper function to convert buffer to lines
func bufferToLines(buf types.Buffer) []string {
	lines := make([]string, buf.LineCount())
	for i := 0; i < buf.LineCount(); i++ {
		line, _ := buf.GetLine(i)
		lines[i] = line
	}
	return lines
}

func extractSingleLine(lines []string, from, to types.Position) string {
	if from.Line() >= len(lines) {
		return ""
	}

	line := lines[from.Line()]

	// Ensure positions are in correct order
	startCol, endCol := from.Column(), to.Column()
	if startCol > endCol {
		startCol, endCol = endCol, startCol
	}

	// Bounds checking
	if startCol >= len(line) {
		return ""
	}
	if endCol >= len(line) {
		endCol = len(line) - 1
	}

	return line[startCol : endCol+1]
}

func extractMultiLine(lines []string, from, to types.Position) string {
	// Ensure positions are in correct order
	start, end := from, to
	if start.Line() > end.Line() || (start.Line() == end.Line() && start.Column() > end.Column()) {
		start, end = end, start
	}

	if start.Line() >= len(lines) {
		return ""
	}

	var result strings.Builder

	// Handle first line
	firstLine := lines[start.Line()]
	if start.Column() < len(firstLine) {
		result.WriteString(firstLine[start.Column():])
		result.WriteString("\n")
	}

	// Handle middle lines
	for i := start.Line() + 1; i < end.Line() && i < len(lines); i++ {
		result.WriteString(lines[i])
		result.WriteString("\n")
	}

	// Handle last line
	if end.Line() < len(lines) {
		lastLine := lines[end.Line()]
		endCol := end.Column()
		if endCol >= len(lastLine) {
			endCol = len(lastLine) - 1
		}
		if endCol >= 0 {
			result.WriteString(lastLine[:endCol+1])
		}
	}

	return result.String()
}
