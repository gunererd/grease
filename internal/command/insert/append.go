package insert

import (
	"github.com/gunererd/grease/internal/state"
	"github.com/gunererd/grease/internal/types"
)

type AppendCommand struct {
	endOfLine bool
	cursorID  int
}

func NewAppendCommand(endOfLine bool, cursorID int) *AppendCommand {
	return &AppendCommand{
		endOfLine: endOfLine,
		cursorID:  cursorID,
	}
}

func (c *AppendCommand) Execute(e types.Editor) types.Editor {
	buf := e.Buffer()
	cursor, err := buf.GetCursor(c.cursorID)
	if err != nil {
		return e
	}

	pos := cursor.GetPosition()
	line, _ := buf.GetLine(pos.Line())

	if c.endOfLine {
		buf.MoveCursor(c.cursorID, pos.Line(), len(line))
	} else {
		if pos.Column() < len(line) {
			buf.MoveCursor(c.cursorID, pos.Line(), pos.Column()+1)
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
