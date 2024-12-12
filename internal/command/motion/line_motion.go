package motion

import (
	"github.com/gunererd/grease/internal/buffer"
	"github.com/gunererd/grease/internal/types"
)

// StartOfLineMotion moves cursor to first character of line
type StartOfLineMotion struct{}

func NewStartOfLineMotion() *StartOfLineMotion {
	return &StartOfLineMotion{}
}

func (m *StartOfLineMotion) Calculate(lines []string, pos types.Position) types.Position {
	if pos.Line() >= len(lines) {
		return pos
	}
	return buffer.NewPosition(pos.Line(), 0)
}

func (m *StartOfLineMotion) Name() string {
	return "goto_line_start"
}

// EndOfLineMotion moves cursor to last character of line
type EndOfLineMotion struct{}

func NewEndOfLineMotion() *EndOfLineMotion {
	return &EndOfLineMotion{}
}

func (m *EndOfLineMotion) Calculate(lines []string, pos types.Position) types.Position {
	if pos.Line() >= len(lines) {
		return pos
	}
	lineLen := len([]rune(lines[pos.Line()]))
	if lineLen == 0 {
		return buffer.NewPosition(pos.Line(), 0)
	}
	return buffer.NewPosition(pos.Line(), lineLen-1)
}

func (m *EndOfLineMotion) Name() string {
	return "goto_line_end"
}
