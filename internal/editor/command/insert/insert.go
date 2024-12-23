package insert

import (
	"github.com/gunererd/grease/internal/editor/state"
	"github.com/gunererd/grease/internal/editor/types"
)

type InsertCommand struct {
	startOfLine bool
	cursorID    int
}

func NewInsertCommand(startOfLine bool, cursorID int) *InsertCommand {
	return &InsertCommand{
		startOfLine: startOfLine,
		cursorID:    cursorID,
	}
}

func (c *InsertCommand) Execute(e types.Editor) types.Editor {
	buf := e.Buffer()
	cursor, err := buf.GetCursor(c.cursorID)
	if err != nil {
		return e
	}

	pos := cursor.GetPosition()
	if c.startOfLine {
		// For 'I', move to first non-whitespace character
		line, _ := buf.GetLine(pos.Line())
		col := 0
		for i, ch := range line {
			if ch != ' ' && ch != '\t' {
				col = i
				break
			}
		}
		buf.MoveCursor(c.cursorID, pos.Line(), col)
	}

	e.SetMode(state.InsertMode)
	e.HandleCursorMovement()
	return e
}

func (c *InsertCommand) Name() string {
	if c.startOfLine {
		return "insert_start_of_line"
	}
	return "insert"
}
