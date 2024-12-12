package handler

import (
	"github.com/gunererd/grease/internal/types"
)

// HistoryManager manages the undo/redo stacks
type HistoryManager struct {
	undoStack []types.HistoryEntry
	redoStack []types.HistoryEntry
	maxSize   int
}

// NewHistoryManager creates a new HistoryManager with a specified maximum size
func NewHistoryManager(maxSize int) types.HistoryManager {
	return &HistoryManager{
		undoStack: make([]types.HistoryEntry, 0),
		redoStack: make([]types.HistoryEntry, 0),
		maxSize:   maxSize,
	}
}

// Push adds a new entry to the undo stack and clears the redo stack
func (h *HistoryManager) Push(entry types.HistoryEntry) {
	// Clear redo stack when new operation is performed
	h.redoStack = h.redoStack[:0]

	// Add new entry to undo stack
	h.undoStack = append(h.undoStack, entry)

	// Remove oldest entries if we exceed maxSize
	if h.maxSize > 0 && len(h.undoStack) > h.maxSize {
		h.undoStack = h.undoStack[1:]
	}
}

// CanUndo returns true if there are operations that can be undone
func (h *HistoryManager) CanUndo() bool {
	return len(h.undoStack) > 0
}

// CanRedo returns true if there are operations that can be redone
func (h *HistoryManager) CanRedo() bool {
	return len(h.redoStack) > 0
}

// GetLastEntry returns the most recent entry from the undo stack without removing it
func (h *HistoryManager) GetLastEntry() (types.HistoryEntry, bool) {
	if !h.CanUndo() {
		return types.HistoryEntry{}, false
	}
	return h.undoStack[len(h.undoStack)-1], true
}

// PopUndo removes and returns the most recent entry from the undo stack
func (h *HistoryManager) PopUndo() (types.HistoryEntry, bool) {
	if !h.CanUndo() {
		return types.HistoryEntry{}, false
	}
	entry := h.undoStack[len(h.undoStack)-1]
	h.undoStack = h.undoStack[:len(h.undoStack)-1]
	return entry, true
}

// PopRedo removes and returns the most recent entry from the redo stack
func (h *HistoryManager) PopRedo() (types.HistoryEntry, bool) {
	if !h.CanRedo() {
		return types.HistoryEntry{}, false
	}
	entry := h.redoStack[len(h.redoStack)-1]
	h.redoStack = h.redoStack[:len(h.redoStack)-1]
	return entry, true
}

// PushRedo adds an entry to the redo stack
func (h *HistoryManager) PushRedo(entry types.HistoryEntry) {
	h.redoStack = append(h.redoStack, entry)
	if h.maxSize > 0 && len(h.redoStack) > h.maxSize {
		h.redoStack = h.redoStack[1:]
	}
}

// Undo restores the buffer and cursor to the state before the last operation
func (h *HistoryManager) Undo(e types.Editor) types.Editor {
	if !h.CanUndo() {
		return e
	}

	entry, _ := h.PopUndo()
	buf := e.Buffer()

	// Restore buffer lines
	for lineNum, content := range entry.BeforeLines {
		buf.ReplaceLine(lineNum, content)
	}

	// Restore cursor position
	cursor, _ := buf.GetPrimaryCursor()
	cursor.SetPosition(entry.CursorBefore)

	// Push to redo stack
	h.PushRedo(entry)

	return e
}

// Redo reapplies the last undone operation
func (h *HistoryManager) Redo(e types.Editor) types.Editor {
	if !h.CanRedo() {
		return e
	}

	entry, _ := h.PopRedo()
	buf := e.Buffer()

	// Restore buffer lines
	for lineNum, content := range entry.AfterLines {
		buf.ReplaceLine(lineNum, content)
	}

	// Restore cursor position
	cursor, _ := buf.GetPrimaryCursor()
	cursor.SetPosition(entry.CursorAfter)

	// Push back to undo stack
	h.undoStack = append(h.undoStack, entry)

	return e
}

// HistoryAwareOperation wraps an Operation with history tracking
type HistoryAwareOperation struct {
	operation types.Operation
	history   types.HistoryManager
}

// NewHistoryAwareOperation creates a new history-aware wrapper around an operation
func NewHistoryAwareOperation(op types.Operation, history types.HistoryManager) *HistoryAwareOperation {
	return &HistoryAwareOperation{
		operation: op,
		history:   history,
	}
}

// getOperationType returns a string identifier for the operation type
func getOperationType(op types.Operation) string {
	switch op.(type) {
	case *DeleteOperation:
		return "delete"
	case *ChangeOperation:
		return "change"
	case *PasteOperation:
		return "paste"
	default:
		return "unknown"
	}
}

// captureLines captures the content of lines between from and to positions
func captureLines(buf types.Buffer, from, to types.Position) map[int]string {
	lines := make(map[int]string)

	// If single line operation
	if from.Line() == to.Line() {
		line, err := buf.GetLine(from.Line())
		if err == nil {
			lines[from.Line()] = line
		}
		return lines
	}

	// Multi-line operation
	for i := from.Line(); i <= to.Line(); i++ {
		line, err := buf.GetLine(i)
		if err == nil {
			lines[i] = line
		}
	}
	return lines
}

// Execute implements the Operation interface with history tracking
func (h *HistoryAwareOperation) Execute(e types.Editor, from, to types.Position) types.Editor {
	// Capture state before operation
	buf := e.Buffer()
	cursor, _ := buf.GetPrimaryCursor()
	cursorBefore := cursor.GetPosition()
	beforeLines := captureLines(buf, from, to)

	// Execute the underlying operation
	model := h.operation.Execute(e, from, to)

	// Capture state after operation
	cursor, _ = buf.GetPrimaryCursor()
	cursorAfter := cursor.GetPosition()
	afterLines := captureLines(buf, from, to)

	// Create and push history entry
	entry := types.HistoryEntry{
		OperationType: getOperationType(h.operation),
		BeforeLines:   beforeLines,
		AfterLines:    afterLines,
		CursorBefore:  cursorBefore,
		CursorAfter:   cursorAfter,
	}
	h.history.Push(entry)

	return model
}
