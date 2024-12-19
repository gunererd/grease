package motion

import (
	"github.com/gunererd/grease/internal/buffer"
	"github.com/gunererd/grease/internal/types"
)

type HalfPageDownMotion struct {
	viewport types.Viewport
}

func NewHalfPageDownMotion(viewport types.Viewport) *HalfPageDownMotion {
	return &HalfPageDownMotion{viewport: viewport}
}

func (m *HalfPageDownMotion) Calculate(lines []string, pos types.Position) types.Position {
	// Don't scroll past buffer end
	if pos.Line() >= len(lines)-1 {
		return pos
	}

	// Calculate if we can scroll a full half page
	lastLine := len(lines) - 1
	_, endLine := m.viewport.VisibleLines()

	// If we're already showing the last line, don't scroll
	if endLine >= lastLine {
		return buffer.NewPosition(lastLine, pos.Column())
	}

	m.viewport.ScrollHalfPageDown(len(lines))
	startLine, _ := m.viewport.VisibleLines()

	// Ensure we don't go past buffer end
	targetLine := startLine + m.viewport.ScrollOff()
	if targetLine >= len(lines) {
		targetLine = len(lines) - 1
	}

	return buffer.NewPosition(targetLine, pos.Column())
}

func (m *HalfPageDownMotion) Name() string {
	return "half_page_down"
}

type HalfPageUpMotion struct {
	viewport types.Viewport
}

func NewHalfPageUpMotion(viewport types.Viewport) *HalfPageUpMotion {
	return &HalfPageUpMotion{viewport: viewport}
}

func (m *HalfPageUpMotion) Calculate(lines []string, pos types.Position) types.Position {
	// If viewport is already at top, move cursor to top
	if m.viewport.Offset().Line() <= 0 {
		return buffer.NewPosition(0, pos.Column())
	}

	if pos.Line() <= 0 {
		return pos
	}

	m.viewport.ScrollHalfPageUp()
	_, endLine := m.viewport.VisibleLines()

	// Ensure we don't go past buffer start
	targetLine := endLine - m.viewport.ScrollOff() - 1
	if targetLine < 0 {
		targetLine = 0
	}

	return buffer.NewPosition(targetLine, pos.Column())
}

func (m *HalfPageUpMotion) Name() string {
	return "half_page_up"
}
