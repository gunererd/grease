package handler

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/types"
)

// ModeHandler defines the interface for handling key events in different editor modes
type ModeHandler interface {
	Handle(msg tea.KeyMsg, e types.Editor) (tea.Model, tea.Cmd)
}
