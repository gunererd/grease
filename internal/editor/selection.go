package editor

// Selection represents a range of selected entries in visual mode
type Selection struct {
	Start int // Starting entry index
	End   int // Ending entry index (inclusive)
}

// NewSelection creates a new selection starting at the given position
func NewSelection(pos int) *Selection {
	return &Selection{
		Start: pos,
		End:   pos,
	}
}

// Contains checks if the given index is within the selection range
func (s *Selection) Contains(idx int) bool {
	if s == nil {
		return false
	}
	if s.Start <= s.End {
		return idx >= s.Start && idx <= s.End
	}
	return idx >= s.End && idx <= s.Start
}

// UpdateEnd updates the end position of the selection
func (s *Selection) UpdateEnd(pos int) {
	if s != nil {
		s.End = pos
	}
}

// Clear clears the selection
func (s *Selection) Clear() {
	if s != nil {
		s.Start = 0
		s.End = 0
	}
}

// GetRange returns the start and end indices of the selection in ascending order
func (s *Selection) GetRange() (start, end int) {
	if s == nil {
		return 0, 0
	}
	if s.Start <= s.End {
		return s.Start, s.End
	}
	return s.End, s.Start
}

// NumSelected returns the number of selected entries
func (s *Selection) NumSelected() int {
	if s == nil {
		return 0
	}
	start, end := s.GetRange()
	return end - start + 1
}
