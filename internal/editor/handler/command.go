package handler

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/state"
	"github.com/gunererd/grease/internal/types"
)

type CommandMode struct{}

func NewCommandMode() *CommandMode {
	return &CommandMode{}
}

func (h *CommandMode) Handle(msg tea.KeyMsg, e types.Editor) (types.Editor, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		e.SetMode(state.NormalMode)
	}
	return e, nil
}
