package buffer

import "github.com/gunererd/grease/internal/types"

// Cursor represents a cursor in the buffer
type Cursor struct {
	pos      types.Position
	id       int
	priority int // Higher priority cursors take precedence in overlapping operations
}

// NewCursor creates a new cursor at the given position
func NewCursor(pos Position, id, priority int) *Cursor {
	return &Cursor{
		pos:      pos,
		id:       id,
		priority: priority,
	}
}

// GetPosition returns the current cursor position
func (c *Cursor) GetPosition() types.Position {
	return c.pos
}

// SetPosition sets the cursor position
func (c *Cursor) SetPosition(pos types.Position) {
	c.pos = pos
}

// GetID returns the cursor's unique identifier
func (c *Cursor) GetID() int {
	return c.id
}

// GetPriority returns the cursor's priority
func (c *Cursor) GetPriority() int {
	return c.priority
}

// SetPriority sets the cursor's priority
func (c *Cursor) SetPriority(priority int) {
	c.priority = priority
}
