package filemanager

import (
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/editor/types"
)

type FileManagerHandler struct {
	fm *FileManager
}

func NewFileManagerHandler(fm *FileManager) *FileManagerHandler {
	return &FileManagerHandler{
		fm: fm,
	}
}

func (h *FileManagerHandler) Handle(msg tea.KeyMsg, e types.Editor) (types.Editor, tea.Cmd) {
	switch msg.String() {
	case "enter":
		cursor, err := e.Buffer().GetPrimaryCursor()
		if err != nil {
			return e, nil
		}

		// Get current line content
		line := cursor.GetPosition().Line()
		content, err := e.Buffer().GetLine(line)
		if err != nil {
			return e, nil
		}

		// If it's a directory (ends with "/"), enter ij
		// TODO: Instead of checking the last character, check if the content is a directory
		//       by checking by using os.Stat. Get content name from editor.Buffer()
		if content[len(content)-1] == '/' {
			dirName := content[:len(content)-1]
			newPath := filepath.Join(h.fm.currentPath, dirName)
			h.fm.currentPath = newPath

			// Update buffer with new directory content
			if err := h.fm.LoadDirectory(); err != nil {
				return e, nil
			}
		}

	case "-":
		// Go to parent directory
		if h.fm.currentPath != "/" {
			h.fm.currentPath = filepath.Dir(h.fm.currentPath)

			// Update buffer with new directory contents
			if err := h.fm.LoadDirectory(); err != nil {
				return e, nil
			}
		}
	}

	return e, nil
}
