package editor

// State manages the editor's state and history
type State struct {
	Mode      Mode
	Selection *Selection
}

// NewState creates a new State instance
func NewState() *State {
	return &State{
		Mode:      NormalMode,
		Selection: nil,
	}
}

// SetMode changes the editor mode
func (s *State) SetMode(mode Mode) {
	s.Mode = mode
}

func (s *State) GetMode() Mode {
	return s.Mode
}

// StartSelection starts a new selection from the given row
func (s *State) StartSelection(row int) {
	s.Selection = NewSelection(row)
}

// UpdateSelection updates the selection end point
func (s *State) UpdateSelection(row int) {
	if s.Selection != nil {
		s.Selection.UpdateEnd(row)
	}
}

// ClearSelection clears the current selection
func (s *State) ClearSelection() {
	s.Selection = nil
}

// IsSelected checks if a given row is currently selected
func (s *State) IsSelected(row int) bool {
	if s.Selection == nil {
		return false
	}
	return s.Selection.Contains(row)
}
