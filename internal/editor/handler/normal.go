package handler

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/keytree"
	"github.com/gunererd/grease/internal/state"
	"github.com/gunererd/grease/internal/types"
)

type NormalMode struct {
	keytree *keytree.KeyTree
}

func NewNormalMode(kt *keytree.KeyTree) *NormalMode {

	// Vim style Jump to beginning of buffer
	kt.Add([]string{"g", "g"}, keytree.KeyAction{
		Execute: func(e types.Editor) (tea.Model, tea.Cmd) {
			cursor, _ := e.Buffer().GetPrimaryCursor()
			e.Buffer().MoveCursor(cursor.ID(), 0, 0)
			e.HandleCursorMovement()
			return e, nil
		},
	})

	return &NormalMode{
		keytree: kt,
	}
}

func (h *NormalMode) Handle(msg tea.KeyMsg, e types.Editor) (tea.Model, tea.Cmd) {

	// Handle key sequences
	if handled, model, cmd := h.keytree.Handle(msg.String(), e); handled {
		return model, cmd
	}

	cursor, err := e.Buffer().GetPrimaryCursor()
	if err != nil {
		return e, nil
	}

	switch msg.String() {
	case "ctrl+c":
		return e, tea.Quit
	case "h":
		e.Buffer().MoveCursorRelative(cursor.ID(), 0, -1)
		e.HandleCursorMovement()
	case "l":
		e.Buffer().MoveCursorRelative(cursor.ID(), 0, 1)
		e.HandleCursorMovement()
	case "j":
		e.Buffer().MoveCursorRelative(cursor.ID(), 1, 0)
		e.HandleCursorMovement()
	case "k":
		e.Buffer().MoveCursorRelative(cursor.ID(), -1, 0)
		e.HandleCursorMovement()
	case "i":
		e.SetMode(state.InsertMode)
	case "v":
		e.SetMode(state.VisualMode)
	case ":":
		e.SetMode(state.CommandMode)
	case "q":
		return e, tea.Quit
	case "G":
		// Vim style end of buffer
		cursor, _ := e.Buffer().GetPrimaryCursor()
		e.Buffer().MoveCursor(cursor.ID(), e.Buffer().LineCount()-1, 0)
		e.HandleCursorMovement()
		return e, nil
	case "$":
		// Vim style jump to end of line
		cursor, _ := e.Buffer().GetPrimaryCursor()
		line := cursor.GetPosition().Line()
		lineLength, _ := e.Buffer().LineLen(line)
		e.Buffer().MoveCursor(cursor.ID(), line, lineLength)
		e.HandleCursorMovement()
		return e, nil
	case "^":
		// Vim style jump to beginning of line
		cursor, _ := e.Buffer().GetPrimaryCursor()
		line := cursor.GetPosition().Line()
		e.Buffer().MoveCursor(cursor.ID(), line, 0)
		e.HandleCursorMovement()
		return e, nil
	}
	return e, nil
}
