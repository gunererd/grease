package filemanager

import (
	"bytes"
	"os"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/editor/types"
)

type FileManager struct {
	currentPath string
	editor      types.Editor
	operations  []Operation
	handler     *FileManagerHandler
}

func New(initialPath string, editor types.Editor) *FileManager {
	fm := &FileManager{
		currentPath: initialPath,
		editor:      editor,
	}
	fm.handler = NewFileManagerHandler(fm)
	return fm
}

func (fm *FileManager) GetHandler() types.ModeHandler {
	return fm.handler
}

func (fm *FileManager) ReadDirectory() ([]Entry, error) {
	entries, err := os.ReadDir(fm.currentPath)
	if err != nil {
		return nil, err
	}

	result := make([]Entry, 0, len(entries))
	for _, entry := range entries {
		entryType := File
		name := entry.Name()

		if entry.IsDir() {
			entryType = Directory
			name = name + "/"
		}

		result = append(result, Entry{
			Name: name,
			Type: entryType,
		})
	}

	// Sort entries: directories first, then files, both alphabetically
	sort.Slice(result, func(i, j int) bool {
		if result[i].Type == result[j].Type {
			return result[i].Name < result[j].Name
		}
		return result[i].Type == Directory
	})

	return result, nil
}

func (fm *FileManager) AsSource() types.Source {
	return NewDirectorySource(fm)
}

func (fm *FileManager) LoadDirectory() error {
	source := fm.AsSource()

	if err := fm.editor.IO().SetSource(source); err != nil {
		return err
	}

	content, err := source.Read()
	if err != nil {
		return err
	}

	return fm.editor.Buffer().LoadFromReader(bytes.NewReader(content))
}

// Implement tea.Model interface
func (fm *FileManager) Init() tea.Cmd {
	return nil
}

func (fm *FileManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Filemanager handles first
		if _, cmd := fm.handler.Handle(msg, fm.editor); cmd != nil {
			return fm, cmd
		}

		// If filemanager didn't handle it, pass to editor
		if _, cmd := fm.editor.Update(msg); cmd != nil {
			return fm, cmd
		}

		return fm, nil
	default:
		// Pass other messages to editor but maintain FileManager as model
		if _, cmd := fm.editor.Update(msg); cmd != nil {
			return fm, cmd
		}
		return fm, nil
	}
}

func (fm *FileManager) View() string {
	return fm.editor.View()
}
