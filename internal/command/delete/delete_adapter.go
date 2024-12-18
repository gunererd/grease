package delete

import (
	"github.com/gunererd/grease/internal/command"
	"github.com/gunererd/grease/internal/command/motion"
	"github.com/gunererd/grease/internal/types"
)

type DeleteCommandAdapter struct {
	cmd *DeleteCommand
}

func NewDeleteCommandAdapter(motion motion.Motion) command.Command {
	return &DeleteCommandAdapter{
		cmd: NewDeleteCommand(motion),
	}
}

func (a *DeleteCommandAdapter) Execute(e types.Editor) types.Editor {
	buf := e.Buffer()
	cursor, _ := buf.GetPrimaryCursor()
	curPos := cursor.GetPosition()

	lines := bufferToLines(buf)
	newLines, newPos := a.cmd.Execute(lines, curPos)

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

func (a *DeleteCommandAdapter) Name() string {
	return a.cmd.Name()
}

func bufferToLines(buf types.Buffer) []string {
	lines := make([]string, buf.LineCount())
	for i := 0; i < buf.LineCount(); i++ {
		line, _ := buf.GetLine(i)
		lines[i] = line
	}
	return lines
}
