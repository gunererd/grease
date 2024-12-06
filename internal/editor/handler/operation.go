package handler

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/state"
	"github.com/gunererd/grease/internal/types"
)

// Operation defines what actions can be performed between two positions in a buffer
type Operation interface {
	Execute(e types.Editor, from, to types.Position) (tea.Model, tea.Cmd)
}

// DeleteOperation implements deletion of text between two positions
type DeleteOperation struct{}

func NewDeleteOperation() *DeleteOperation {
	return &DeleteOperation{}
}

func (d *DeleteOperation) Execute(e types.Editor, from, to types.Position) (tea.Model, tea.Cmd) {
	buf := e.Buffer()
	if from.Line() == to.Line() {
		// Single line deletion
		line, _ := buf.GetLine(from.Line())
		newLine := line[:from.Column()] + line[to.Column():]
		buf.ReplaceLine(from.Line(), newLine)
	} else {
		// Multi-line deletion
		firstLine, _ := buf.GetLine(from.Line())
		lastLine, _ := buf.GetLine(to.Line())

		// Combine first and last line
		newLine := firstLine[:from.Column()] + lastLine[to.Column():]
		buf.ReplaceLine(from.Line(), newLine)

		// Remove lines in between
		for i := from.Line() + 1; i <= to.Line(); i++ {
			buf.RemoveLine(from.Line() + 1)
		}
	}
	return e, nil
}

// ChangeOperation implements change operation (delete + enter insert mode)
type ChangeOperation struct {
	*DeleteOperation
}

func NewChangeOperation() *ChangeOperation {
	return &ChangeOperation{
		DeleteOperation: NewDeleteOperation(),
	}
}

func (c *ChangeOperation) Execute(e types.Editor, from, to types.Position) (tea.Model, tea.Cmd) {
	model, cmd := c.DeleteOperation.Execute(e, from, to)
	e.SetMode(state.InsertMode)
	return model, cmd
}

// YankOperation implements copying of text between two positions
type YankOperation struct{}

func NewYankOperation() *YankOperation {
	return &YankOperation{}
}

func (y *YankOperation) Execute(e types.Editor, from, to types.Position) (tea.Model, tea.Cmd) {
	buf := e.Buffer()
	var yankedText string

	if from.Line() == to.Line() {
		// Handle single line yank
		line, _ := buf.GetLine(from.Line())
		yankedText = line[from.Column():to.Column()]
	} else {
		// Handle multi-line yank
		for i := from.Line(); i <= to.Line(); i++ {
			line, _ := buf.GetLine(i)
			if i == from.Line() {
				yankedText += line[from.Column():] + "\n"
			} else if i == to.Line() {
				yankedText += line[:to.Column()]
			} else {
				yankedText += line + "\n"
			}
		}
	}

	defaultRegister.Set(yankedText)
	log.Println("Yanked text:", yankedText)
	return e, nil
}

// PasteOperation implements pasting of text after or before cursor
type PasteOperation struct {
	before bool // If true, paste before cursor
}

func NewPasteOperation(before bool) *PasteOperation {
	return &PasteOperation{before: before}
}

func (p *PasteOperation) Execute(e types.Editor, from, to types.Position) (tea.Model, tea.Cmd) {
	buf := e.Buffer()
	text := defaultRegister.Get()
	cursor, _ := buf.GetPrimaryCursor()

	if text == "" {
		return e, nil
	}

	// Split text into lines
	lines := strings.Split(text, "\n")

	if len(lines) == 1 {
		// Single line paste
		line, _ := buf.GetLine(from.Line())
		insertPos := from.Column()
		if !p.before {
			insertPos++
		}
		newLine := line[:insertPos] + text + line[insertPos:]
		buf.ReplaceLine(from.Line(), newLine)

		// Move cursor to end of pasted text
		newCol := insertPos + len(text)
		buf.MoveCursor(cursor.ID(), from.Line(), newCol)
	} else {
		// Multi-line paste
		currentLine, _ := buf.GetLine(from.Line())
		insertPos := from.Column()
		if !p.before {
			insertPos++
		}

		// Handle first line
		firstLine := currentLine[:insertPos] + lines[0]
		buf.ReplaceLine(from.Line(), firstLine)

		// Insert middle lines
		for i := 1; i < len(lines)-1; i++ {
			buf.InsertLine(from.Line()+i, lines[i])
		}

		// Handle last line
		lastLine := lines[len(lines)-1] + currentLine[insertPos:]
		buf.InsertLine(from.Line()+len(lines)-1, lastLine)

		// Move cursor to end of pasted text
		buf.MoveCursor(cursor.ID(), from.Line()+len(lines)-1, len(lines[len(lines)-1]))
	}

	return e, nil
}
