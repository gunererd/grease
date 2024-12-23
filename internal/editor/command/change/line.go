package change

import (
	"log"

	"github.com/gunererd/grease/internal/editor/state"
	"github.com/gunererd/grease/internal/editor/types"
)

type ChangeLineCommand struct {
	cursor types.Cursor
}

func NewChangeLineCommand(cursor types.Cursor) types.Command {
	return &ChangeLineCommand{
		cursor: cursor,
	}
}

func (c *ChangeLineCommand) Explain() {
	log.Printf("type:<ChangeLineCommand>, cursor:<%d>, pos:<%v>\n", c.cursor.ID(), c.cursor.GetPosition())
}

func (c *ChangeLineCommand) Execute(e types.Editor) types.Editor {
	buf := e.Buffer()

	line := c.cursor.GetPosition().Line()
	buf.RemoveLine(line)

	if line >= buf.LineCount() && line > 0 {
		buf.MoveCursor(c.cursor.ID(), line-1, 0)
	} else {
		buf.InsertLine(line, "")
		buf.MoveCursor(c.cursor.ID(), line, 0)
	}

	e.SetMode(state.InsertMode)
	return e
}

func (c *ChangeLineCommand) Name() string {
	return "change_line"
}

type ChangeToEndOfLineCommand struct {
	cursor types.Cursor
}

func NewChangeToEndOfLineCommand(cursor types.Cursor) types.Command {
	return &ChangeToEndOfLineCommand{
		cursor: cursor,
	}
}

func (c *ChangeToEndOfLineCommand) Explain() {
	log.Printf("type:<ChangeToEndOfLineCommand>, cursor:<%d>, pos:<%v>\n", c.cursor.ID(), c.cursor.GetPosition())
}

func (c *ChangeToEndOfLineCommand) Execute(e types.Editor) types.Editor {
	buf := e.Buffer()

	pos := c.cursor.GetPosition()
	line, _ := buf.GetLine(pos.Line())
	newLine := line[:pos.Column()]
	buf.ReplaceLine(pos.Line(), newLine)

	e.SetMode(state.InsertMode)
	return e
}

func (c *ChangeToEndOfLineCommand) Name() string {
	return "change_to_end"
}
