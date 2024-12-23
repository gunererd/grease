package motion

import (
	"github.com/gunererd/grease/internal/editor/buffer"
	"github.com/gunererd/grease/internal/editor/types"
)

// StartOfBufferMotion moves cursor to first line of buffer
type StartOfBufferMotion struct{}

func NewStartOfBufferMotion() *StartOfBufferMotion {
	return &StartOfBufferMotion{}
}

func (m *StartOfBufferMotion) Calculate(lines []string, pos types.Position) types.Position {
	if len(lines) == 0 {
		return pos
	}
	return buffer.NewPosition(0, 0)
}

func (m *StartOfBufferMotion) Name() string {
	return "move_to_buffer_start"
}

// EndOfBufferMotion moves cursor to last line of buffer
type EndOfBufferMotion struct{}

func NewEndOfBufferMotion() *EndOfBufferMotion {
	return &EndOfBufferMotion{}
}

func (m *EndOfBufferMotion) Calculate(lines []string, pos types.Position) types.Position {
	if len(lines) == 0 {
		return pos
	}
	lastLine := len(lines) - 1
	return buffer.NewPosition(lastLine, 0)
}

func (m *EndOfBufferMotion) Name() string {
	return "move_to_buffer_end"
}
