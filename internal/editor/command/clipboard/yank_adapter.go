package clipboard

import (
	"log"

	"github.com/gunererd/grease/internal/editor/command/motion"
	"github.com/gunererd/grease/internal/editor/register"
	"github.com/gunererd/grease/internal/editor/types"
)

// YankCommandAdapter wraps YankCommand to implement the Command interface
type YankCommandAdapter struct {
	cmd      *YankCommand
	register *register.Register
	yanked   string
}

func NewYankCommandAdapter(motion motion.Motion, register *register.Register) types.Command {
	return &YankCommandAdapter{
		cmd:      NewYankCommand(motion),
		register: register,
	}
}

func (a *YankCommandAdapter) Execute(e types.Editor) types.Editor {
	buf := e.Buffer()
	cursor, _ := buf.GetPrimaryCursor()
	curPos := cursor.GetPosition()

	lines := bufferToLines(buf)
	_, _ = a.cmd.Execute(lines, curPos, a.register)
	a.yanked = a.register.Get()
	return e
}

func (a *YankCommandAdapter) Explain() {
	log.Printf("type:<YankCommandAdapter>, cmd:<%s>, yanked:<%q>\n", a.cmd.Name(), a.yanked)
}
func (a *YankCommandAdapter) Name() string {
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
