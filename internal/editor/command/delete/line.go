package delete

import (
	"github.com/gunererd/grease/internal/editor/command"
	"github.com/gunererd/grease/internal/editor/types"
)

type DeleteLineCommand struct {
	cursorID int
}

func NewDeleteLineCommand(cursorID int) command.Command {
	return &DeleteLineCommand{
		cursorID: cursorID,
	}
}

func (c *DeleteLineCommand) Execute(e types.Editor) types.Editor {
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
		buf.MoveCursor(c.cursorID, line, 0)
	}

	return e
}

func (c *DeleteLineCommand) Name() string {
	return "delete_line"
}

type DeleteToEndCommandOfLine struct {
	cursorID int
}

func NewDeleteToEndOfLineCommand(cursorID int) command.Command {
	return &DeleteToEndCommandOfLine{
		cursorID: cursorID,
	}
}

func (c *DeleteToEndCommandOfLine) Execute(e types.Editor) types.Editor {
	buf := e.Buffer()
	cursor, err := buf.GetCursor(c.cursorID)
	if err != nil {
		return e
	}

	pos := cursor.GetPosition()
	line, _ := buf.GetLine(pos.Line())
	newLine := line[:pos.Column()]
	buf.ReplaceLine(pos.Line(), newLine)

	return e
}

func (c *DeleteToEndCommandOfLine) Name() string {
	return "delete_to_end"
}
