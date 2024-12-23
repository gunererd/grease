package filemanager

import (
	"bytes"
	"os"
	"sort"

	"github.com/gunererd/grease/internal/editor/types"
)

type FileManager struct {
	currentPath string
	editor      types.Editor
	operations  []Operation
}

func New(initialPath string, editor types.Editor) *FileManager {
	return &FileManager{
		currentPath: initialPath,
		editor:      editor,
	}
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
