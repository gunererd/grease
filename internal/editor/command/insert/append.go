package insert

import (
	"github.com/gunererd/grease/internal/editor/state"
	"github.com/gunererd/grease/internal/editor/types"
)

type AppendCommand struct {
	endOfLine bool
	cursor    types.Cursor
}

func NewAppendCommand(endOfLine bool, cursor types.Cursor) *AppendCommand {
	return &AppendCommand{
		endOfLine: endOfLine,
		cursor:    cursor,
	}
}

func (c *AppendCommand) Execute(e types.Editor) types.Editor {
	buf := e.Buffer()

	pos := c.cursor.GetPosition()
	line, _ := buf.GetLine(pos.Line())

	if c.endOfLine {
		buf.MoveCursor(c.cursor.ID(), pos.Line(), len(line))
	} else {
		if pos.Column() < len(line) {
			buf.MoveCursor(c.cursor.ID(), pos.Line(), pos.Column()+1)
		}
	}

	e.SetMode(state.InsertMode)
	e.HandleCursorMovement()
	return e
}

func (c *AppendCommand) Name() string {
	if c.endOfLine {
		return "append_end_of_line"
	}
	return "append"
}
