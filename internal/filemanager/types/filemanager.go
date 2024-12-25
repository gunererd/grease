package types

import (
	tea "github.com/charmbracelet/bubbletea"
	eTypes "github.com/gunererd/grease/internal/editor/types"
)

type FileManager interface {
	LoadDirectory(path string) error
	DirectoryManager() DirectoryManager
	OperationManager() OperationManager
	Editor() eTypes.Editor
	Logger() Logger

	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	Init() tea.Cmd
	View() string
}

type Source interface {
	Path() string
	Read() ([]byte, error)
	Close() error
}
