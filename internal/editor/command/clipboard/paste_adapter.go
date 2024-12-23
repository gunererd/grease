package clipboard

import (
	"github.com/gunererd/grease/internal/editor/command"
	"github.com/gunererd/grease/internal/editor/register"
	"github.com/gunererd/grease/internal/editor/types"
)

// PasteCommandAdapter wraps PasteCommand to implement the Command interface
type PasteCommandAdapter struct {
	cmd      *PasteCommand
	register *register.Register
	cursor   types.Cursor
}

func NewPasteCommandAdapter(cursor types.Cursor, register *register.Register, before bool) command.Command {
	return &PasteCommandAdapter{
		cmd:      NewPasteCommand(before),
		register: register,
		cursor:   cursor,
	}
}

func (a *PasteCommandAdapter) Execute(e types.Editor) types.Editor {
	buf := e.Buffer()
	curPos := a.cursor.GetPosition()

	lines := bufferToLines(buf)
	newLines, newPos := a.cmd.Execute(lines, curPos, a.register)

	// Clear existing lines
	for i := buf.LineCount() - 1; i >= 0; i-- {
		buf.ReplaceLine(i, "")
	}

	// Add new lines
	for i, line := range newLines {
		if i < buf.LineCount() {
			buf.ReplaceLine(i, line)
		} else {
			buf.InsertLine(i, line)
		}
	}

	a.cursor.SetPosition(newPos)
	return e
}

func (a *PasteCommandAdapter) Name() string {
	return a.cmd.Name()
}
