package clipboard

import (
	"strings"

	"github.com/gunererd/grease/internal/buffer"
	"github.com/gunererd/grease/internal/register"
	"github.com/gunererd/grease/internal/types"
)

type PasteCommand struct {
	before bool
}

func NewPasteCommand(before bool) *PasteCommand {
	return &PasteCommand{
		before: before,
	}
}

func (c *PasteCommand) Execute(lines []string, pos types.Position, register *register.Register) ([]string, types.Position) {
	text := register.Get()
	newLines, newPos := insertText(lines, pos, text, c.before)
	return newLines, newPos
}

func (c *PasteCommand) Name() string {
	return "paste"
}

func insertText(lines []string, pos types.Position, text string, before bool) ([]string, types.Position) {
	if text == "" {
		return lines, pos
	}

	// Copy lines to avoid modifying original slice
	result := make([]string, len(lines))
	copy(result, lines)

	textLines := strings.Split(text, "\n")
	insertPos := calculateInsertPosition(pos, before)

	if len(textLines) == 1 {
		return insertSingleLine(result, pos, textLines[0], insertPos)
	}
	return insertMultiLine(result, pos, textLines, insertPos)
}

func calculateInsertPosition(pos types.Position, before bool) int {
	if before {
		return pos.Column()
	}
	return pos.Column() + 1
}

func insertSingleLine(lines []string, pos types.Position, text string, insertPos int) ([]string, types.Position) {
	// Handle position beyond buffer
	if pos.Line() >= len(lines) {
		lines = append(lines, text)
		return lines, buffer.NewPosition(len(lines)-1, len(text)-1)
	}

	currentLine := lines[pos.Line()]

	// Handle empty line
	if len(currentLine) == 0 {
		lines[pos.Line()] = text
		return lines, buffer.NewPosition(pos.Line(), len(text)-1)
	}

	// Insert into existing line
	if insertPos > len(currentLine) {
		insertPos = len(currentLine)
	}
	newLine := currentLine[:insertPos] + text + currentLine[insertPos:]
	lines[pos.Line()] = newLine
	return lines, buffer.NewPosition(pos.Line(), insertPos+len(text)-1)
}

func insertMultiLine(lines []string, pos types.Position, textLines []string, insertPos int) ([]string, types.Position) {
	// Handle position beyond buffer
	if pos.Line() >= len(lines) {
		lines = append(lines, textLines...)
		return lines, buffer.NewPosition(len(lines)-1, len(textLines[len(textLines)-1])-1)
	}

	newLines := make([]string, 0, len(lines)+len(textLines)-1)

	// Handle first line
	firstLine := handleFirstLine(lines[pos.Line()], textLines[0], insertPos)
	newLines = append(newLines, lines[:pos.Line()]...)
	newLines = append(newLines, firstLine)

	// Add middle lines
	newLines = append(newLines, textLines[1:len(textLines)-1]...)

	// Handle last line
	lastLine := handleLastLine(lines[pos.Line()], textLines[len(textLines)-1], insertPos)
	newLines = append(newLines, lastLine)
	newLines = append(newLines, lines[pos.Line()+1:]...)

	return newLines, buffer.NewPosition(pos.Line()+len(textLines)-1, len(textLines[len(textLines)-1])-1)
}

func handleFirstLine(currentLine, textLine string, insertPos int) string {
	if len(currentLine) == 0 {
		return textLine
	}
	if insertPos > len(currentLine) {
		insertPos = len(currentLine)
	}
	return currentLine[:insertPos] + textLine
}

func handleLastLine(currentLine, textLine string, insertPos int) string {
	if len(currentLine) == 0 {
		return textLine
	}
	return textLine + currentLine[insertPos:]
}
