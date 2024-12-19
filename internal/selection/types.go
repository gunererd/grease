package selection

import "github.com/gunererd/grease/internal/types"

// Mode represents different selection modes
type Mode int

const (
	Character Mode = iota
	Line
	Block
)

// Selection represents a selected region in the buffer
type Selection struct {
	id          int
	highlightID int
	anchor      types.Position // Where selection started
	head        types.Position // Where selection currently ends
	mode        Mode
}

// Add these methods to the Selection struct
func (s *Selection) GetAnchor() types.Position {
	return s.anchor
}

func (s *Selection) GetHead() types.Position {
	return s.head
}

// Manager defines the interface for managing selections
type Manager interface {
	// StartSelection starts a new selection at the given position
	StartSelection(pos types.Position) *Selection

	// UpdateSelection updates an existing selection
	UpdateSelection(sel *Selection, newHead types.Position)

	// ClearSelection removes a specific selection
	ClearSelection(sel *Selection)

	// ClearAllSelections removes all selections
	ClearAllSelections()

	// GetSelections returns all active selections
	GetSelections() []*Selection

	// GetSelectedText returns the text content of a selection
	GetSelectedText(sel *Selection, buf types.Buffer) string
}
