package types

// HistoryManagerInterface defines the interface for managing history
// This allows for different implementations of history management
// to be used interchangeably.
type HistoryManager interface {
	Undo(e Editor) Editor
	Redo(e Editor) Editor
	Push(entry HistoryEntry)
	CanUndo() bool
	CanRedo() bool
}

// HistoryEntry represents a single operation in the history
// This is a simplified version for the interface
type HistoryEntry struct {
	OperationType string
	BeforeLines   map[int]string
	AfterLines    map[int]string
	CursorBefore  Position
	CursorAfter   Position
}
