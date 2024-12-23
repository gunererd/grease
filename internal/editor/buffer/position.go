package buffer

import (
	"fmt"

	"github.com/gunererd/grease/internal/editor/types"
)

// Position represents a position in the buffer
type Position struct {
	line   int
	column int
}

func NewPosition(line, column int) Position {
	return Position{
		line:   line,
		column: column,
	}
}

// String returns a string representation of the position
func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.line, p.column)
}

// Before returns true if this position is before other position
func (p Position) Before(other types.Position) bool {
	if p.line < other.Line() {
		return true
	}
	if p.line == other.Line() {
		return p.column < other.Column()
	}
	return false
}

// Equal returns true if this position is equal to other position
func (p Position) Equal(other types.Position) bool {
	return p.line == other.Line() && p.column == other.Column()
}

func (p Position) Line() int {
	return p.line
}

func (p Position) Column() int {
	return p.column
}

// Add returns a new position offset by the given line and column
func (p Position) Add(line, col int) types.Position {
	return Position{
		line:   p.line + line,
		column: p.column + col,
	}
}

func (p Position) Set(line, col int) types.Position {
	return Position{
		line:   line,
		column: col,
	}
}
