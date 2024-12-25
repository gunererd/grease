package types

import tea "github.com/charmbracelet/bubbletea"

type Handler interface {
	Handle(msg tea.KeyMsg) (tea.Cmd, error)
}
