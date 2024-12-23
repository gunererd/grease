package types

import (
	tea "github.com/charmbracelet/bubbletea"
)

// ModeHandler defines the interface for handling key events in different editor modes
type ModeHandler interface {
	Handle(msg tea.KeyMsg, editor Editor) (Editor, tea.Cmd)
}
