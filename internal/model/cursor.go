package model

type Cursor struct {
	Row int
	Col int
}

func NewCursor() Cursor {
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

func (c *Cursor) SetPosition(row, col int) {
	c.Row = row
	c.Col = col
}
