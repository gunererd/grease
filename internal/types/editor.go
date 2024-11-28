package types

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/buffer"
	"github.com/gunererd/grease/internal/state"
	"github.com/gunererd/grease/internal/ui"
)

// Editor defines the interface for editor operations needed by handlers
type Editor interface {
	Buffer() *buffer.Buffer
	Viewport() *ui.Viewport
	Width() int
	Height() int
	SetMode(mode state.Mode)
	HandleCursorMovement()
	UpdateViewport(width, height int)
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	Init() tea.Cmd
	View() string
}
