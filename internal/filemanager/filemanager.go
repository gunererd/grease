package filemanager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	eTypes "github.com/gunererd/grease/internal/editor/types"
	"github.com/gunererd/grease/internal/filemanager/handler"
	"github.com/gunererd/grease/internal/filemanager/hook"
	"github.com/gunererd/grease/internal/filemanager/types"
)

type Filemanager struct {
	dirManager types.DirectoryManager
	opManager  types.OperationManager
	view       types.View
	handler    types.Handler
	editor     eTypes.Editor
	logger     types.Logger
}

func New(
	dirManager types.DirectoryManager,
	opManager types.OperationManager,
	view types.View,
	editor eTypes.Editor,
	logger types.Logger,
) types.FileManager {
	fm := &Filemanager{
		dirManager: dirManager,
		opManager:  opManager,
		view:       view,
		editor:     editor,
		logger:     logger,
	}

	fm.handler = handler.New(dirManager, editor, fm.LoadDirectory, logger)
	editor.AddHook(hook.NewFileOperationHook(dirManager, opManager))

	return fm
}

func (fm *Filemanager) DirectoryManager() types.DirectoryManager {
	return fm.dirManager
}

func (fm *Filemanager) OperationManager() types.OperationManager {
	return fm.opManager
}

func (fm *Filemanager) Editor() eTypes.Editor {
	return fm.editor
}

// Implement tea.Model interface
func (fm *Filemanager) Init() tea.Cmd {
	return nil
}

func (fm *Filemanager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Let handler process the input first
		if cmd, err := fm.handler.Handle(msg); err != nil {
			return fm, cmd
		}
		// If handler didn't handle it, pass to editor
		if _, cmd := fm.editor.Update(msg); cmd != nil {
			return fm, cmd
		}
	default:
		if _, cmd := fm.editor.Update(msg); cmd != nil {
			return fm, cmd
		}
	}
	return fm, nil
}

func (fm *Filemanager) View() string {
	return fm.view.Render()
}

func (fm *Filemanager) LoadDirectory(path string) error {
	resolvedPath, err := resolvePath(path)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	if err := fm.dirManager.ChangeDirectory(resolvedPath); err != nil {
		return fmt.Errorf("failed to change directory: %w", err)
	}

	entries, err := fm.dirManager.ReadDirectory()
	if err != nil {
		return err
	}

	var sb strings.Builder
	for i, entry := range entries {
		sb.WriteString(entry.Name())
		if i < len(entries)-1 {
			sb.WriteString("\n")
		}
	}

	return fm.editor.Buffer().LoadFromReader(strings.NewReader(sb.String()))
}

// resolvePath sanitizes and resolves the given path
func resolvePath(path string) (string, error) {
	if path == "" {
		pwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get working directory: %w", err)
		}
		return pwd, nil
	}

	// Expand home directory if needed
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		path = strings.Replace(path, "~", homeDir, 1)
	}

	// Clean and make absolute
	absPath, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Verify directory exists and is accessible
	info, err := os.Stat(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to access path: %w", err)
	}

	if !info.IsDir() {
		return "", fmt.Errorf("path is not a directory: %s", absPath)
	}

	return absPath, nil
}

func (fm *Filemanager) Logger() types.Logger {
	return fm.logger
}
