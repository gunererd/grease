package history

import (
	"github.com/gunererd/grease/internal/editor/types"
)

type HistoryManager struct {
	undoStack []types.HistoryEntry
	redoStack []types.HistoryEntry
	maxSize   int
}

func New(maxSize int) types.HistoryManager {
	return &HistoryManager{
		undoStack: make([]types.HistoryEntry, 0),
		redoStack: make([]types.HistoryEntry, 0),
		maxSize:   maxSize,
	}
}

func (h *HistoryManager) Push(entry types.HistoryEntry) {
	// Clear redo stack when new operation is performed
	h.redoStack = h.redoStack[:0]

	h.undoStack = append(h.undoStack, entry)

	// Remove oldest entries if we exceed maxSize
	if h.maxSize > 0 && len(h.undoStack) > h.maxSize {
		h.undoStack = h.undoStack[1:]
	}
}

func (h *HistoryManager) CanUndo() bool {
	return len(h.undoStack) > 0
}

func (h *HistoryManager) CanRedo() bool {
	return len(h.redoStack) > 0
}

func (h *HistoryManager) Undo(e types.Editor) types.Editor {
	if !h.CanUndo() {
		return e
	}

	// Pop last entry from undo stack
	entry := h.undoStack[len(h.undoStack)-1]
	h.undoStack = h.undoStack[:len(h.undoStack)-1]

	// Restore buffer lines
	buf := e.Buffer()
	for lineNum, content := range entry.BeforeLines {
		buf.ReplaceLine(lineNum, content)
	}

	// Restore cursor position
	cursor, _ := buf.GetPrimaryCursor()
	cursor.SetPosition(entry.CursorBefore)

	// Add to redo stack
	h.redoStack = append(h.redoStack, entry)

	return e
}

func (h *HistoryManager) Redo(e types.Editor) types.Editor {
	if !h.CanRedo() {
		return e
	}

	// Pop last entry from redo stack
	entry := h.redoStack[len(h.redoStack)-1]
	h.redoStack = h.redoStack[:len(h.redoStack)-1]

	// Restore buffer lines
	buf := e.Buffer()
	for lineNum, content := range entry.AfterLines {
		buf.ReplaceLine(lineNum, content)
	}

	// Restore cursor position
	cursor, _ := buf.GetPrimaryCursor()
	cursor.SetPosition(entry.CursorAfter)

	// Add back to undo stack
	h.undoStack = append(h.undoStack, entry)

	return e
}
