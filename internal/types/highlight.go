package types

// HighlightManager defines the interface for managing highlights
type HighlightManager interface {
	// Add adds a new highlight and returns its ID
	Add(h Highlight) int

	// Remove removes a highlight by its ID
	Remove(id int)

	// Clear removes all highlights of a given type
	Clear(highlightType HighlightType)

	// Get returns a highlight by its ID
	Get(id int) (Highlight, bool)

	// GetForLine returns all highlights that intersect with the given line
	GetForLine(line int) []Highlight

	// GetForPosition returns all highlights that contain the given position
	GetForPosition(pos Position) []Highlight

	// Update updates an existing highlight
	Update(id int, h Highlight) bool
}

// Highlight represents a highlighted region in the buffer
type Highlight interface {
	// GetID returns the highlight's unique identifier
	GetID() int

	// GetType returns the type of highlight
	GetType() HighlightType

	// GetStartPosition returns the starting position of the highlight
	GetStartPosition() Position

	// GetEndPosition returns the ending position of the highlight
	GetEndPosition() Position

	// GetPriority returns the highlight's priority
	GetPriority() int

	// Contains checks if a position is within this highlight
	Contains(pos Position) bool
}

// HighlightType represents different types of highlights
type HighlightType int

const (
	// VisualHighlight for visual mode selection
	VisualHighlight HighlightType = iota
	// SearchHighlight for search results
	SearchHighlight
)
