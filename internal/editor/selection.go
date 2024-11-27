package editor

// Selection represents a range of selected characters in visual mode
type Selection struct {
	StartRow, StartCol int // Starting position (row, column)
	EndRow, EndCol     int // Ending position (row, column)
}

// NewSelection creates a new selection starting at the given position
func NewSelection(row, col int) *Selection {
	return &Selection{
		StartRow: row,
		StartCol: col,
		EndRow:   row,
		EndCol:   col,
	}
}

// Contains checks if the given position is within the selection range
func (s *Selection) Contains(row, col int) bool {
	if s == nil {
		return false
	}

	// Get ordered positions
	startRow, startCol, endRow, endCol := s.GetOrderedRange()

	// If single line selection
	if startRow == endRow {
		return row == startRow && col >= startCol && col <= endCol
	}

	// Multi-line selection
	switch {
	case row == startRow:
		return col >= startCol
	case row == endRow:
		return col <= endCol
	case row > startRow && row < endRow:
		return true
	default:
		return false
	}
}

// UpdateEnd updates the end position of the selection
func (s *Selection) UpdateEnd(row, col int) {
	if s != nil {
		s.EndRow = row
		s.EndCol = col
	}
}

// Clear clears the selection
func (s *Selection) Clear() {
	if s != nil {
		s.StartRow = 0
		s.StartCol = 0
		s.EndRow = 0
		s.EndCol = 0
	}
}

// GetOrderedRange returns the selection range in order (top-to-bottom, left-to-right)
func (s *Selection) GetOrderedRange() (startRow, startCol, endRow, endCol int) {
	if s == nil {
		return 0, 0, 0, 0
	}

	if s.StartRow < s.EndRow || (s.StartRow == s.EndRow && s.StartCol <= s.EndCol) {
		return s.StartRow, s.StartCol, s.EndRow, s.EndCol
	}
	return s.EndRow, s.EndCol, s.StartRow, s.StartCol
}

// NumSelected returns the number of selected characters
func (s *Selection) NumSelected() int {
	if s == nil {
		return 0
	}
	startRow, startCol, endRow, endCol := s.GetOrderedRange()
	if startRow == endRow {
		return endCol - startCol + 1
	}
	numSelected := endCol - startCol + 1
	for row := startRow + 1; row < endRow; row++ {
		numSelected += 100 // assuming 100 characters per line
	}
	numSelected += endCol + 1
	return numSelected
}
