package buffer

import "testing"

func TestNewCursor(t *testing.T) {
	pos := Position{Line: 1, Column: 2}
	cursor := NewCursor(pos, 1, 50)

	if cursor.pos != pos {
		t.Errorf("Expected cursor position %v, got %v", pos, cursor.pos)
	}

	if cursor.id != 1 {
		t.Errorf("Expected cursor id 1, got %d", cursor.id)
	}

	if cursor.priority != 50 {
		t.Errorf("Expected cursor priority 50, got %d", cursor.priority)
	}
}

func TestCursorGetPosition(t *testing.T) {
	pos := Position{Line: 1, Column: 2}
	cursor := NewCursor(pos, 1, 50)

	gotPos := cursor.GetPosition()
	if gotPos != pos {
		t.Errorf("Expected position %v, got %v", pos, gotPos)
	}
}

func TestCursorSetPosition(t *testing.T) {
	cursor := NewCursor(Position{Line: 0, Column: 0}, 1, 50)
	newPos := Position{Line: 1, Column: 2}

	cursor.SetPosition(newPos)
	if cursor.pos != newPos {
		t.Errorf("Expected position %v, got %v", newPos, cursor.pos)
	}
}

func TestCursorGetID(t *testing.T) {
	cursor := NewCursor(Position{Line: 0, Column: 0}, 42, 50)

	if id := cursor.GetID(); id != 42 {
		t.Errorf("Expected cursor ID 42, got %d", id)
	}
}

func TestCursorGetSetPriority(t *testing.T) {
	cursor := NewCursor(Position{Line: 0, Column: 0}, 1, 50)

	if priority := cursor.GetPriority(); priority != 50 {
		t.Errorf("Expected priority 50, got %d", priority)
	}

	cursor.SetPriority(75)
	if priority := cursor.GetPriority(); priority != 75 {
		t.Errorf("Expected priority 75, got %d", priority)
	}
}
