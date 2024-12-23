package insert

import (
	"log"

	"github.com/gunererd/grease/internal/editor/state"
	"github.com/gunererd/grease/internal/editor/types"
)

type InsertCommand struct {
	startOfLine bool
	cursor      types.Cursor
}

func NewInsertCommand(startOfLine bool, cursor types.Cursor) *InsertCommand {
	return &InsertCommand{
		startOfLine: startOfLine,
		cursor:      cursor,
	}
}

func (c *InsertCommand) Explain() {
	log.Printf("type:<InsertCommand>, startOfLine:<%t>, cursor:<%d>, pos:<%v>\n", c.startOfLine, c.cursor.ID(), c.cursor.GetPosition())
}

func (c *InsertCommand) Execute(e types.Editor) types.Editor {
	buf := e.Buffer()

	pos := c.cursor.GetPosition()
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
		buf.MoveCursor(c.cursor.ID(), pos.Line(), col)
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
