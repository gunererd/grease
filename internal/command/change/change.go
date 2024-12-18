package change

import (
	"github.com/gunererd/grease/internal/command/motion"
	"github.com/gunererd/grease/internal/types"
)

type ChangeCommand struct {
	motion motion.Motion
}

func NewChangeCommand(motion motion.Motion) *ChangeCommand {
	return &ChangeCommand{
		motion: motion,
	}
}

func (c *ChangeCommand) Execute(lines []string, pos types.Position) ([]string, types.Position) {
	targetPos := c.motion.Calculate(lines, pos)

	// Ensure from position is before to position
	from, to := pos, targetPos
	if from.Line() > to.Line() || (from.Line() == to.Line() && from.Column() > to.Column()) {
		from, to = to, from
	}

	if from.Line() == to.Line() {
		return changeSingleLine(lines, from, to)
	}
	return changeMultiLine(lines, from, to)
}

func (c *ChangeCommand) Name() string {
	return "change"
}

func changeSingleLine(lines []string, from, to types.Position) ([]string, types.Position) {
	if from.Line() >= len(lines) {
		return lines, from
	}

	line := lines[from.Line()]
	if from.Column() >= len(line) {
		return lines, from
	}

	endCol := to.Column()
	if endCol > len(line) {
		endCol = len(line)
	}

	newLine := line[:from.Column()] + line[endCol:]
	lines[from.Line()] = newLine
	return lines, from
}

func changeMultiLine(lines []string, from, to types.Position) ([]string, types.Position) {
	if from.Line() >= len(lines) {
		return lines, from
	}

	firstLine := lines[from.Line()][:from.Column()]

	// Create new slice with changed lines removed
	newLines := make([]string, 0, len(lines)-(to.Line()-from.Line()))
	newLines = append(newLines, lines[:from.Line()]...)
	newLines = append(newLines, firstLine)
	if to.Line()+1 < len(lines) {
		newLines = append(newLines, lines[to.Line()+1:]...)
	}

	return newLines, from
}
