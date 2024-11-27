package model

type Change struct {
	LineIndex int // Index of the line that was changed
	OldText   string
	NewText   string
	Position  int // Position in the line where change occurred
}

type History struct {
	undoStack []Change
	redoStack []Change
	maxSize   int
}

func NewHistory(maxSize int) *History {
	return &History{
		undoStack: make([]Change, 0),
		redoStack: make([]Change, 0),
		maxSize:   maxSize,
	}
}

func (h *History) Push(change Change) {
	h.undoStack = append(h.undoStack, change)
	if len(h.undoStack) > h.maxSize {
		h.undoStack = h.undoStack[1:]
	}
	// Clear redo stack when new change is made
	h.redoStack = make([]Change, 0)
}

func (h *History) Undo() (Change, bool) {
	if len(h.undoStack) == 0 {
		return Change{}, false
	}

	change := h.undoStack[len(h.undoStack)-1]
	h.undoStack = h.undoStack[:len(h.undoStack)-1]
	h.redoStack = append(h.redoStack, change)
	return change, true
}

func (h *History) Redo() (Change, bool) {
	if len(h.redoStack) == 0 {
		return Change{}, false
	}

	change := h.redoStack[len(h.redoStack)-1]
	h.redoStack = h.redoStack[:len(h.redoStack)-1]
	h.undoStack = append(h.undoStack, change)
	return change, true
}
