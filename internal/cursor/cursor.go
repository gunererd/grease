package cursor

// Cursor represents a 2D cursor position
type Cursor struct {
	Row int
	Col int
}

// New creates a new Cursor instance
func New() Cursor {
	return Cursor{
		Row: 0,
		Col: 0,
	}
}

// MoveLeft moves cursor left if possible
func (c *Cursor) MoveLeft() {
	if c.Col > 0 {
		c.Col--
	}
}

// MoveRight moves cursor right if within line bounds
func (c *Cursor) MoveRight(lineLength int) {
	if c.Col < lineLength {
		c.Col++
	}
}

// MoveUp moves cursor up and adjusts column if needed
func (c *Cursor) MoveUp(prevLineLength int) {
	if c.Row > 0 {
		c.Row--
		if c.Col > prevLineLength {
			c.Col = prevLineLength
		}
	}
}

// MoveDown moves cursor down and adjusts column if needed
func (c *Cursor) MoveDown(nextLineLength int) {
	c.Row++
	if c.Col > nextLineLength {
		c.Col = nextLineLength
	}
}

// SetPosition sets the cursor position
func (c *Cursor) SetPosition(row, col int) {
	c.Row = row
	c.Col = col
}

// GetPosition returns the current cursor position
func (c *Cursor) GetPosition() (row, col int) {
	return c.Row, c.Col
}
