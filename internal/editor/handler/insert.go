package handler

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/state"
	"github.com/gunererd/grease/internal/types"
)

type InsertMode struct{}

func NewInsertMode() *InsertMode {
	return &InsertMode{}
}

func (h *InsertMode) Handle(msg tea.KeyMsg, e types.Editor) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		e.SetMode(state.NormalMode)
	case tea.KeyRunes:
		e.Buffer().Insert(string(msg.Runes))
		e.HandleCursorMovement()
	case tea.KeySpace:
		e.Buffer().Insert(" ")
		e.HandleCursorMovement()
	case tea.KeyEnter:
		e.Buffer().Insert("\n")
		e.HandleCursorMovement()
	case tea.KeyBackspace:
		e.Buffer().Delete(-1)
		e.HandleCursorMovement()
	}
	return e, nil
}
