package buffer

import (
	"testing"
)

func TestNew(t *testing.T) {
	b := New()
	if len(b.lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(b.lines))
	}
	if len(b.lines[0]) != 0 {
		t.Errorf("Expected empty line, got %v", b.lines[0])
	}

	cursor, err := b.GetPrimaryCursor()
	if err != nil {
		t.Errorf("Expected primary cursor, got error: %v", err)
	}
	if cursor.GetPriority() != 100 {
		t.Errorf("Expected primary cursor priority 100, got %d", cursor.GetPriority())
	}
}

func TestNewFromString(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		wantLines     int
		wantFirstLine string
		wantLastLine  string
	}{
		{
			name:          "empty string",
			content:       "",
			wantLines:     1,
			wantFirstLine: "",
			wantLastLine:  "",
		},
		{
			name:          "single line",
			content:       "hello",
			wantLines:     1,
			wantFirstLine: "hello",
			wantLastLine:  "hello",
		},
		{
			name:          "multiple lines",
			content:       "hello\nworld",
			wantLines:     2,
			wantFirstLine: "hello",
			wantLastLine:  "world",
		},
		{
			name:          "trailing newline",
			content:       "hello\n",
			wantLines:     2,
			wantFirstLine: "hello",
			wantLastLine:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewFromString(tt.content)
			if len(b.lines) != tt.wantLines {
				t.Errorf("Expected %d lines, got %d", tt.wantLines, len(b.lines))
			}
			if string(b.lines[0]) != tt.wantFirstLine {
				t.Errorf("Expected first line %q, got %q", tt.wantFirstLine, string(b.lines[0]))
			}
			if string(b.lines[len(b.lines)-1]) != tt.wantLastLine {
				t.Errorf("Expected last line %q, got %q", tt.wantLastLine, string(b.lines[len(b.lines)-1]))
			}

			cursor, err := b.GetPrimaryCursor()
			if err != nil {
				t.Errorf("Expected primary cursor, got error: %v", err)
			}
			if cursor.GetPriority() != 100 {
				t.Errorf("Expected primary cursor priority 100, got %d", cursor.GetPriority())
			}
		})
	}
}

func TestBuffer_AddCursor(t *testing.T) {
	b := NewFromString("hello\nworld")

	tests := []struct {
		name        string
		pos         Position
		priority    int
		wantErr     bool
		wantCursors int
	}{
		{
			name:        "valid position",
			pos:         Position{Line: 0, Column: 0},
			priority:    50,
			wantErr:     false,
			wantCursors: 2, // including primary cursor
		},
		{
			name:        "invalid line",
			pos:         Position{Line: -1, Column: 0},
			priority:    50,
			wantErr:     true,
			wantCursors: 2,
		},
		{
			name:        "invalid column",
			pos:         Position{Line: 0, Column: 10},
			priority:    50,
			wantErr:     true,
			wantCursors: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := b.AddCursor(tt.pos, tt.priority)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCursor() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(b.cursors) != tt.wantCursors {
				t.Errorf("Expected %d cursors, got %d", tt.wantCursors, len(b.cursors))
			}
		})
	}
}

func TestBuffer_Insert(t *testing.T) {
	tests := []struct {
		name        string
		initial     string
		cursorPos   []Position
		insert      string
		want        string
		wantCursors []Position
	}{
		{
			name:    "single cursor insert",
			initial: "hello\nworld",
			cursorPos: []Position{
				{Line: 0, Column: 5},
			},
			insert: "!",
			want:   "hello!\nworld",
			wantCursors: []Position{
				{Line: 0, Column: 6},
			},
		},
		{
			name:    "multiple cursor insert",
			initial: "hello\nworld",
			cursorPos: []Position{
				{Line: 0, Column: 5},
				{Line: 1, Column: 5},
			},
			insert: "!",
			want:   "hello!\nworld!",
			wantCursors: []Position{
				{Line: 0, Column: 6},
				{Line: 1, Column: 6},
			},
		},
		{
			name:    "newline insert",
			initial: "hello\nworld",
			cursorPos: []Position{
				{Line: 0, Column: 5},
			},
			insert: "\n",
			want:   "hello\n\nworld",
			wantCursors: []Position{
				{Line: 1, Column: 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewFromString(tt.initial)
			
			// Remove primary cursor and add test cursors
			b.cursors = nil
			for i, pos := range tt.cursorPos {
				b.AddCursor(pos, 50-i) // Decreasing priority
			}

			if err := b.Insert(tt.insert); err != nil {
				t.Errorf("Insert() error = %v", err)
			}

			if got := b.Get(); got != tt.want {
				t.Errorf("After insert, got %q, want %q", got, tt.want)
			}

			for i, wantPos := range tt.wantCursors {
				if got := b.cursors[i].GetPosition(); got != wantPos {
					t.Errorf("Cursor %d position = %v, want %v", i, got, wantPos)
				}
			}
		})
	}
}

func TestBuffer_Delete(t *testing.T) {
	tests := []struct {
		name        string
		initial     string
		cursorPos   []Position
		count       int
		want        string
		wantCursors []Position
	}{
		{
			name:    "single cursor delete",
			initial: "hello!\nworld!",
			cursorPos: []Position{
				{Line: 0, Column: 6},
			},
			count: 1,
			want:  "hello\nworld!",
			wantCursors: []Position{
				{Line: 0, Column: 5},
			},
		},
		{
			name:    "multiple cursor delete",
			initial: "hello!\nworld!",
			cursorPos: []Position{
				{Line: 0, Column: 6},
				{Line: 1, Column: 6},
			},
			count: 1,
			want:  "hello\nworld",
			wantCursors: []Position{
				{Line: 0, Column: 5},
				{Line: 1, Column: 5},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewFromString(tt.initial)
			
			// Remove primary cursor and add test cursors
			b.cursors = nil
			for i, pos := range tt.cursorPos {
				b.AddCursor(pos, 50-i) // Decreasing priority
			}

			if err := b.Delete(tt.count); err != nil {
				t.Errorf("Delete() error = %v", err)
			}

			if got := b.Get(); got != tt.want {
				t.Errorf("After delete, got %q, want %q", got, tt.want)
			}

			for i, wantPos := range tt.wantCursors {
				if got := b.cursors[i].GetPosition(); got != wantPos {
					t.Errorf("Cursor %d position = %v, want %v", i, got, wantPos)
				}
			}
		})
	}
}

func TestBuffer_MoveCursor(t *testing.T) {
	b := NewFromString("hello\nworld")
	cursor, _ := b.AddCursor(Position{Line: 0, Column: 0}, 50)

	tests := []struct {
		name         string
		lineOffset   int
		columnOffset int
		wantPos      Position
		wantErr      bool
	}{
		{
			name:         "move right",
			lineOffset:   0,
			columnOffset: 1,
			wantPos:      Position{Line: 0, Column: 1},
			wantErr:      false,
		},
		{
			name:         "move down",
			lineOffset:   1,
			columnOffset: 0,
			wantPos:      Position{Line: 1, Column: 1},
			wantErr:      false,
		},
		{
			name:         "invalid move up",
			lineOffset:   -2,
			columnOffset: 0,
			wantPos:      Position{Line: 1, Column: 1},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := b.MoveCursor(cursor.GetID(), tt.lineOffset, tt.columnOffset)
			if (err != nil) != tt.wantErr {
				t.Errorf("MoveCursor() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				if got := cursor.GetPosition(); got != tt.wantPos {
					t.Errorf("Cursor position = %v, want %v", got, tt.wantPos)
				}
			}
		})
	}
}
