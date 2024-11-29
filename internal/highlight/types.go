package highlight

import "github.com/gunererd/grease/internal/types"

// Type represents different types of highlights
type Type int

const (
	Visual Type = iota
	Search
)

// Highlight represents a highlighted region in the buffer
type Highlight struct {
	// Unique identifier for this highlight
	ID int
	// Starting position of the highlight
	StartPosition types.Position
	// Ending position of the highlight (inclusive)
	EndPosition types.Position
	// Type of the highlight
	Type Type
	// Priority determines which highlight is shown when multiple highlights overlap
	// Higher priority highlights take precedence
	Priority int
}

// Region represents a region of text in the buffer
type Region struct {
	Start types.Position
	End   types.Position
}

// highlight implements types.Highlight interface
type highlight struct {
	id            int
	startPosition types.Position
	endPosition   types.Position
	highlightType types.HighlightType
	priority      int
}

func (h *highlight) GetID() int {
	return h.id
}

func (h *highlight) GetType() types.HighlightType {
	return h.highlightType
}

func (h *highlight) GetStartPosition() types.Position {
	return h.startPosition
}

func (h *highlight) GetEndPosition() types.Position {
	return h.endPosition
}

func (h *highlight) GetPriority() int {
	return h.priority
}

// Contains checks if a position is within this highlight
func (h *highlight) Contains(pos types.Position) bool {
	// If highlight is on a single line
	if h.startPosition.Line() == h.endPosition.Line() {
		return pos.Line() == h.startPosition.Line() &&
			pos.Column() >= h.startPosition.Column() &&
			pos.Column() <= h.endPosition.Column()
	}

	// If position is on first line
	if pos.Line() == h.startPosition.Line() {
		return pos.Column() >= h.startPosition.Column()
	}

	// If position is on last line
	if pos.Line() == h.endPosition.Line() {
		return pos.Column() <= h.endPosition.Column()
	}

	// If position is between first and last line
	return pos.Line() > h.startPosition.Line() && pos.Line() < h.endPosition.Line()
}

// NewHighlight creates a new highlight
func NewHighlight(start, end types.Position, highlightType types.HighlightType, priority int) types.Highlight {
	return &highlight{
		startPosition: start,
		endPosition:   end,
		highlightType: highlightType,
		priority:      priority,
	}
}

// CreateVisualHighlight creates a new highlight for visual mode selection
func CreateVisualHighlight(start, end types.Position) types.Highlight {
	return NewHighlight(start, end, types.VisualHighlight, 100) // Visual mode has highest priority
}

// CreateSearchHighlight creates a new highlight for search results
func CreateSearchHighlight(start, end types.Position) types.Highlight {
	return NewHighlight(start, end, types.SearchHighlight, 50)
}

// Overlaps checks if this highlight overlaps with another highlight
func (h *Highlight) Overlaps(other *Highlight) bool {
	// If highlights are on different lines
	if h.EndPosition.Line() < other.StartPosition.Line() ||
		h.StartPosition.Line() > other.EndPosition.Line() {
		return false
	}

	// If highlights are on the same line
	if h.StartPosition.Line() == h.EndPosition.Line() &&
		other.StartPosition.Line() == other.EndPosition.Line() {
		return !(h.EndPosition.Column() < other.StartPosition.Column() ||
			h.StartPosition.Column() > other.EndPosition.Column())
	}

	// If highlights span multiple lines
	return true
}
