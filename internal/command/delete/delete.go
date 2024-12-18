package delete

import (
	"github.com/gunererd/grease/internal/command/motion"
	"github.com/gunererd/grease/internal/types"
)

type DeleteCommand struct {
	motion motion.Motion
}

func NewDeleteCommand(motion motion.Motion) *DeleteCommand {
	return &DeleteCommand{
		motion: motion,
	}
}

func (c *DeleteCommand) Execute(lines []string, pos types.Position) ([]string, types.Position) {
	targetPos := c.motion.Calculate(lines, pos)

	// Ensure from position is before to position
	from, to := pos, targetPos
	if from.Line() > to.Line() || (from.Line() == to.Line() && from.Column() > to.Column()) {
		from, to = to, from
	}

	if from.Line() == to.Line() {
		return deleteSingleLine(lines, from, to)
	}
	return deleteMultiLine(lines, from, to)
}

func (c *DeleteCommand) Name() string {
	return "delete"
}

func deleteSingleLine(lines []string, from, to types.Position) ([]string, types.Position) {
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

func deleteMultiLine(lines []string, from, to types.Position) ([]string, types.Position) {
	if from.Line() >= len(lines) {
		return lines, from
	}

	// Handle first line
	firstLine := lines[from.Line()][:from.Column()]

	// Handle last line
	var lastLineSuffix string
	if to.Line() < len(lines) {
		lastLine := lines[to.Line()]
		if to.Column() < len(lastLine) {
			lastLineSuffix = lastLine[to.Column():]
		}
	}

	newLine := firstLine + lastLineSuffix

	// Create new slice with deleted lines removed
	newLines := make([]string, 0, len(lines)-(to.Line()-from.Line()))
	newLines = append(newLines, lines[:from.Line()]...)
	newLines = append(newLines, newLine)
	if to.Line()+1 < len(lines) {
		newLines = append(newLines, lines[to.Line()+1:]...)
	}

	return newLines, from
}
