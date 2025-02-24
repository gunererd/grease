package types

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/editor/state"
)

// Editor defines the interface for editor operations needed by handlers
type Editor interface {
	Buffer() Buffer
	Viewport() Viewport
	Width() int
	Height() int
	SetMode(mode state.Mode)
	HandleCursorMovement()
	UpdateViewport(width, height int)
	HighlightManager() HighlightManager
	HistoryManager() HistoryManager
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	Init() tea.Cmd
	View() string
	IO() IOManager
	AddHook(h Hook)
	RemoveHook(h Hook)
	GetHooks() []Hook
	Logger() Logger
}
