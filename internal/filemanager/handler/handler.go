package handler

import (
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	eTypes "github.com/gunererd/grease/internal/editor/types"
	"github.com/gunererd/grease/internal/filemanager/types"
)

type Handler struct {
	dirManager types.DirectoryManager
	editor     eTypes.Editor
	loadDir    func(string) error
	logger     types.Logger
}

func New(
	dirManager types.DirectoryManager,
	editor eTypes.Editor,
	loadDir func(string) error,
	logger types.Logger,
) *Handler {
	return &Handler{
		dirManager: dirManager,
		editor:     editor,
		loadDir:    loadDir,
		logger:     logger,
	}
}

func (h *Handler) Handle(msg tea.KeyMsg) (tea.Cmd, error) {
	switch msg.String() {
	case "enter":
		h.logger.Println("Enter key pressed")
		cursor, err := h.editor.Buffer().GetPrimaryCursor()
		if err != nil {
			return nil, err
		}

		line := cursor.GetPosition().Line()
		content, err := h.editor.Buffer().GetLine(line)
		if err != nil {
			return nil, err
		}

		if content[len(content)-1] == '/' {
			dirName := content[:len(content)-1]
			newPath := filepath.Join(h.dirManager.CurrentPath(), dirName)

			if err := h.dirManager.ChangeDirectory(newPath); err != nil {
				return nil, err
			}

			return nil, h.loadDir(newPath)
		}

	case "-":
		if h.dirManager.CurrentPath() != "/" {
			h.logger.Println("Back key pressed")
			parentDir := filepath.Dir(h.dirManager.CurrentPath())

			if err := h.dirManager.ChangeDirectory(parentDir); err != nil {
				return nil, err
			}

			return nil, h.loadDir(parentDir)
		}
	}

	return nil, nil
}
