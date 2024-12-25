package directory

import (
	"os"
	"sort"
	"strings"

	"github.com/gunererd/grease/internal/filemanager/entry"
	"github.com/gunererd/grease/internal/filemanager/types"
)

type Manager struct {
	currentPath string
	logger      types.Logger
}

func NewDirectoryManager(initialPath string, logger types.Logger) types.DirectoryManager {
	return &Manager{
		currentPath: initialPath,
		logger:      logger,
	}
}

func (m *Manager) ReadDirectory() ([]types.Entry, error) {
	entries, err := os.ReadDir(m.currentPath)
	if err != nil {
		m.logger.Println("Failed to read directory:", err)
		return nil, err
	}

	result := make([]types.Entry, 0, len(entries))
	for _, e := range entries {
		entryType := types.File
		name := e.Name()

		if e.IsDir() {
			entryType = types.Directory
			name = name + "/"
		}

		result = append(result, entry.New(name, entryType))
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Type() == result[j].Type() {
			return result[i].Name() < result[j].Name()
		}
		return result[i].Type() == types.Directory
	})

	return result, nil
}

func (m *Manager) ChangeDirectory(path string) error {
	m.currentPath = path
	return nil
}

func (m *Manager) CurrentPath() string {
	return m.currentPath
}

func (m *Manager) GetDirectoryContent() ([]byte, error) {
	entries, err := m.ReadDirectory()
	if err != nil {
		m.logger.Println("Failed to read directory:", err)
		return nil, err
	}

	var sb strings.Builder
	for i, entry := range entries {
		sb.WriteString(entry.Name())
		if i < len(entries)-1 {
			sb.WriteString("\n")
		}
	}

	return []byte(sb.String()), nil
}
