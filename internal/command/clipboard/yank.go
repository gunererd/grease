package clipboard

import (
	"github.com/gunererd/grease/internal/command/motion"
	"github.com/gunererd/grease/internal/register"
	"github.com/gunererd/grease/internal/types"
)

type YankCommand struct {
	motion motion.Motion
}

func NewYankCommand(motion motion.Motion) *YankCommand {
	return &YankCommand{
		motion: motion,
	}
}

func (c *YankCommand) Execute(lines []string, pos types.Position, register *register.Register) ([]string, types.Position) {
	targetPos := c.motion.Calculate(lines, pos)
	yankedText := extractText(lines, pos, targetPos)
	register.Set(yankedText)
	return lines, pos
}

func (c *YankCommand) Name() string {
	return "yank"
}

func extractText(lines []string, from, to types.Position) string {
	if len(lines) == 0 {
		return ""
	}

	// Ensure from position is before to position
	if from.Line() > to.Line() || (from.Line() == to.Line() && from.Column() > to.Column()) {
		from, to = to, from
	}

	if from.Line() == to.Line() {
		return extractSingleLine(lines, from, to)
	}
	return extractMultiLine(lines, from, to)
}

func extractSingleLine(lines []string, from, to types.Position) string {
	line := lines[from.Line()]
	if from.Column() >= len(line) {
		return ""
	}
	endCol := to.Column()
	if endCol >= len(line)-1 {
		endCol = len(line)
		return line[from.Column():endCol]
	}
	return line[from.Column():endCol]
}

func extractMultiLine(lines []string, from, to types.Position) string {
	var text string
	for i := from.Line(); i <= to.Line()+1; i++ {
		if i >= len(lines) {
			break
		}

		line := lines[i]
		switch {
		case i == from.Line():
			if from.Column() < len(line) {
				text += line[from.Column():] + "\n"
			}
		case i == to.Line():
			endCol := to.Column()
			if endCol > len(line) {
				endCol = len(line)
			}
			text += line[:endCol+1]
		default:
			text += line + "\n"
		}
	}
	return text
}
