package delete

import (
	"log"

	"github.com/gunererd/grease/internal/editor/types"
)

type DeleteLineCommand struct {
	cursor types.Cursor
}

func NewDeleteLineCommand(cursor types.Cursor) types.Command {
	return &DeleteLineCommand{
		cursor: cursor,
	}
}

func (c *DeleteLineCommand) Explain() {
	log.Printf("type:<DeleteLineCommand>, cursor:<%d>, pos:<%v>\n", c.cursor.ID(), c.cursor.GetPosition())
}

func (c *DeleteLineCommand) Execute(e types.Editor) types.Editor {
	buf := e.Buffer()

	line := c.cursor.GetPosition().Line()
	buf.RemoveLine(line)

	if line >= buf.LineCount() && line > 0 {
		buf.MoveCursor(c.cursor.ID(), line-1, 0)
	} else {
		buf.MoveCursor(c.cursor.ID(), line, 0)
	}

	return e
}

func (c *DeleteLineCommand) Name() string {
	return "delete_line"
}

type DeleteToEndCommandOfLine struct {
	cursor types.Cursor
}

func NewDeleteToEndOfLineCommand(cursor types.Cursor) types.Command {
	return &DeleteToEndCommandOfLine{
		cursor: cursor,
	}
}

func (c *DeleteToEndCommandOfLine) Explain() {
	log.Printf("type:<DeleteToEndCommandOfLine>, cursor:<%d>, pos:<%v>\n", c.cursor.ID(), c.cursor.GetPosition())
}

func (c *DeleteToEndCommandOfLine) Execute(e types.Editor) types.Editor {
	buf := e.Buffer()

	pos := c.cursor.GetPosition()
	line, _ := buf.GetLine(pos.Line())
	newLine := line[:pos.Column()]
	buf.ReplaceLine(pos.Line(), newLine)

	return e
}

func (c *DeleteToEndCommandOfLine) Name() string {
	return "delete_to_end"
}
