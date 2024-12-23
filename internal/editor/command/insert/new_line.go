package insert

import (
	"github.com/gunererd/grease/internal/editor/state"
	"github.com/gunererd/grease/internal/editor/types"
)

type NewLineCommand struct {
	before bool
}

func NewNewLineCommand(before bool) *NewLineCommand {
	return &NewLineCommand{
		before: before,
	}
}

func (c *NewLineCommand) Execute(e types.Editor) types.Editor {
	buf := e.Buffer()
	cursor, err := buf.GetPrimaryCursor()
	if err != nil {
		return e
	}

	currentLine := cursor.GetPosition().Line()
	insertLine := currentLine
	if !c.before {
		insertLine++
	}

	if err := buf.InsertLine(insertLine, ""); err != nil {
		return e
	}

	buf.MoveCursor(cursor.ID(), insertLine, 0)
	e.SetMode(state.InsertMode)
	e.HandleCursorMovement()

	return e
}

func (c *NewLineCommand) Name() string {
	if c.before {
		return "new_line_before"
	}
	return "new_line_after"
}
