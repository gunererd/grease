package buffer

import "fmt"

// Position represents a position in the buffer
type Position struct {
	Line   int
	Column int
}

// String returns a string representation of the position
func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line+1, p.Column+1)
}

// Before returns true if this position is before other position
func (p Position) Before(other Position) bool {
	if p.Line < other.Line {
		return true
	}
	if p.Line == other.Line {
		return p.Column < other.Column
	}
	return false
}

// Equal returns true if this position is equal to other position
func (p Position) Equal(other Position) bool {
	return p.Line == other.Line && p.Column == other.Column
}

// Add returns a new position offset by the given line and column
func (p Position) Add(line, col int) Position {
	return Position{
		Line:   p.Line + line,
		Column: p.Column + col,
	}
}
