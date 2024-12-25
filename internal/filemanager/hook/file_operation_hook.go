package hook

import (
	"path/filepath"

	eTypes "github.com/gunererd/grease/internal/editor/types"
	"github.com/gunererd/grease/internal/filemanager/operation"
	types "github.com/gunererd/grease/internal/filemanager/types"
)

type FileOperationHook struct {
	fm        types.FileManager
	opManager types.OperationManager
	clipboard map[string]string // stores original paths of deleted/moved files
}

func NewFileOperationHook(fm types.FileManager, opManager types.OperationManager) *FileOperationHook {
	return &FileOperationHook{
		fm:        fm,
		opManager: opManager,
		clipboard: make(map[string]string),
	}
}

func (foh *FileOperationHook) OnBeforeCommand(cmd eTypes.Command, e eTypes.Editor) {
	// We don't need to do anything before command execution
}

func (foh *FileOperationHook) OnAfterCommand(cmd eTypes.Command, e eTypes.Editor) {
	switch cmd.Name() {
	case "delete_line":
		cursor, err := e.Buffer().GetPrimaryCursor()
		if err != nil {
			return
		}

		// Get the line that was deleted
		line := cursor.GetPosition().Line()
		content, err := e.Buffer().GetLine(line)
		if err != nil {
			return
		}

		// Store in clipboard for potential move operations
		foh.clipboard[content] = filepath.Join(foh.fm.DirectoryManager().CurrentPath(), content)

		// Queue delete operation
		foh.opManager.QueueOperation(operation.New(
			types.Delete,
			filepath.Join(foh.fm.DirectoryManager().CurrentPath(), content),
			"",
		))

	case "change_line":
		cursor, err := e.Buffer().GetPrimaryCursor()
		if err != nil {
			return
		}

		// Get the line before and after change
		line := cursor.GetPosition().Line()
		newContent, err := e.Buffer().GetLine(line)
		if err != nil {
			return
		}

		// Get original content from history
		history := e.HistoryManager()
		if len(history.UndoStack()) > 0 {
			lastEntry := history.UndoStack()[len(history.UndoStack())-1]
			oldContent := lastEntry.BeforeLines[line]

			// Queue rename operation
			foh.opManager.QueueOperation(operation.New(
				types.Rename,
				filepath.Join(foh.fm.DirectoryManager().CurrentPath(), oldContent),
				filepath.Join(foh.fm.DirectoryManager().CurrentPath(), newContent),
			))
		}

	case "paste":
		cursor, err := e.Buffer().GetPrimaryCursor()
		if err != nil {
			return
		}

		// Get the line where paste occurred
		line := cursor.GetPosition().Line()
		content, err := e.Buffer().GetLine(line)
		if err != nil {
			return
		}

		// If we have this content in clipboard, it's a move operation
		if originalPath, exists := foh.clipboard[content]; exists {
			foh.opManager.QueueOperation(operation.New(
				types.Move,
				originalPath,
				foh.fm.DirectoryManager().CurrentPath(),
			))
			delete(foh.clipboard, content)
		}
	}
}
