package handler

import (
	tea "github.com/charmbracelet/bubbletea"
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
		line, _ := buf.GetLine(from.Line())
		newLine := line[:from.Column()] + line[to.Column():]
		buf.ReplaceLine(from.Line(), newLine)
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
	// TODO: Add command to switch to insert mode
	return model, cmd
}

// YankOperation implements copying of text between two positions
type YankOperation struct {
	register string // Register to store yanked text
}

func NewYankOperation() *YankOperation {
	return &YankOperation{}
}

func (y *YankOperation) Execute(e types.Editor, from, to types.Position) (tea.Model, tea.Cmd) {
	buf := e.Buffer()
	if from.Line() == to.Line() {
		// Handle single line yank
		line, _ := buf.GetLine(from.Line())
		y.register = line[from.Column():to.Column()]
	} else {
		// Handle multi-line yank
		var yankedText string
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
		y.register = yankedText
	}
	return e, nil
}
