package change

import (
	"github.com/gunererd/grease/internal/command"
	"github.com/gunererd/grease/internal/state"
	"github.com/gunererd/grease/internal/types"
)

type ChangeLineCommand struct {
	cursorID int
}

func NewChangeLineCommand(cursorID int) command.Command {
	return &ChangeLineCommand{
		cursorID: cursorID,
	}
}

func (c *ChangeLineCommand) Execute(e types.Editor) types.Editor {
	buf := e.Buffer()
	cursor, err := buf.GetCursor(c.cursorID)
	if err != nil {
		return e
	}

	line := cursor.GetPosition().Line()
	buf.RemoveLine(line)

	if line >= buf.LineCount() && line > 0 {
		buf.MoveCursor(c.cursorID, line-1, 0)
	} else {
		buf.InsertLine(line, "")
		buf.MoveCursor(c.cursorID, line, 0)
	}

	e.SetMode(state.InsertMode)
	return e
}

func (c *ChangeLineCommand) Name() string {
	return "change_line"
}

type ChangeToEndOfLineCommand struct {
	cursorID int
}

func NewChangeToEndOfLineCommand(cursorID int) command.Command {
	return &ChangeToEndOfLineCommand{
		cursorID: cursorID,
	}
}

func (c *ChangeToEndOfLineCommand) Execute(e types.Editor) types.Editor {
	buf := e.Buffer()
	cursor, err := buf.GetCursor(c.cursorID)
	if err != nil {
		return e
	}

	pos := cursor.GetPosition()
	line, _ := buf.GetLine(pos.Line())
	newLine := line[:pos.Column()]
	buf.ReplaceLine(pos.Line(), newLine)

	e.SetMode(state.InsertMode)
	return e
}

func (c *ChangeToEndOfLineCommand) Name() string {
	return "change_to_end"
}
