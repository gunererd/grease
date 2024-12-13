package clipboard

import (
	"github.com/gunererd/grease/internal/command"
	"github.com/gunererd/grease/internal/command/motion"
	"github.com/gunererd/grease/internal/register"
	"github.com/gunererd/grease/internal/types"
)

// PasteCommandAdapter wraps PasteCommand to implement the Command interface
type PasteCommandAdapter struct {
	cmd      *PasteCommand
	register *register.Register
}

func NewPasteCommandAdapter(motion motion.Motion, register *register.Register, before bool) command.Command {
	return &PasteCommandAdapter{
		cmd:      NewPasteCommand(before),
		register: register,
	}
}

func (a *PasteCommandAdapter) Execute(e types.Editor) types.Editor {
	buf := e.Buffer()
	cursor, _ := buf.GetPrimaryCursor()
	curPos := cursor.GetPosition()

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

	cursor.SetPosition(newPos)
	return e
}

func (a *PasteCommandAdapter) Name() string {
	return a.cmd.Name()
}
