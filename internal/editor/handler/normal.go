package handler

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/state"
	"github.com/gunererd/grease/internal/types"
)

type NormalMode struct{}

func NewNormalMode() *NormalMode {
	return &NormalMode{}
}

func (h *NormalMode) Handle(msg tea.KeyMsg, e types.Editor) (tea.Model, tea.Cmd) {
	cursor, err := e.Buffer().GetPrimaryCursor()
	if err != nil {
		return e, nil
	}

	switch msg.String() {
	case "ctrl+c":
		return e, tea.Quit
	case "h":
		e.Buffer().MoveCursor(cursor.ID(), 0, -1)
		e.HandleCursorMovement()
	case "l":
		e.Buffer().MoveCursor(cursor.ID(), 0, 1)
		e.HandleCursorMovement()
	case "j":
		e.Buffer().MoveCursor(cursor.ID(), 1, 0)
		e.HandleCursorMovement()
	case "k":
		e.Buffer().MoveCursor(cursor.ID(), -1, 0)
		e.HandleCursorMovement()
	case "i":
		e.SetMode(state.InsertMode)
	case "v":
		e.SetMode(state.VisualMode)
	case ":":
		e.SetMode(state.CommandMode)
	case "q":
		return e, tea.Quit
	case "z":
		// Center viewport on cursor
		cursor, _ := e.Buffer().GetPrimaryCursor()
		e.Viewport().CenterOn(cursor.GetPosition())
	}

	return e, nil
}
