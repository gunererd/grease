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

func (c *Cursor) Move(deltaRow, deltaCol int) {
	c.Row += deltaRow
	c.Col += deltaCol
}

func (c *Cursor) SetPosition(row, col int) {
	c.Row = row
	c.Col = col
}
