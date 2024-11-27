package navigator

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/buffer"
	"github.com/gunererd/grease/internal/cursor"
)

// Navigator handles navigation through a hierarchical structure
type Navigator struct {
	Buffer *buffer.Buffer
	Cursor cursor.Cursor
}

// New creates a new Navigator instance
func New() *Navigator {
	return &Navigator{
		Buffer: buffer.NewBuffer(),
		Cursor: cursor.New(),
	}
}

// MoveCursor moves the cursor by the given row offset
func (n *Navigator) MoveCursor(offset int) bool {
	newRow := n.Cursor.Row + offset
	if newRow >= 0 && newRow < n.Buffer.NumEntries() {
		n.Cursor.Row = newRow
		// Adjust column if needed
		if entry, ok := n.Buffer.GetEntry(newRow); ok {
			lineLength := len(entry.Name)
			if n.Cursor.Col > lineLength {
				n.Cursor.Col = lineLength
			}
		}
		return true
	}
	return false
}

// MoveCursorLeft moves cursor left
func (n *Navigator) MoveCursorLeft() {
	n.Cursor.MoveLeft()
}

// MoveCursorRight moves cursor right within current entry bounds
func (n *Navigator) MoveCursorRight() {
	if entry, ok := n.Buffer.GetEntry(n.Cursor.Row); ok {
		lineLength := len(entry.Name)
		n.Cursor.MoveRight(lineLength)
	}
}

// SetCursor sets the cursor to a specific row
func (n *Navigator) SetCursor(row int) {
	if row >= 0 && row < n.Buffer.NumEntries() {
		n.Cursor.Row = row
	}
}

// GetCurrentEntry returns the entry at the cursor position
func (n *Navigator) GetCurrentEntry() (buffer.Entry, bool) {
	return n.Buffer.GetEntry(n.Cursor.Row)
}

// ReadDirectory reads the contents of the specified directory
func (n *Navigator) ReadDirectory(path string) tea.Cmd {
	return func() tea.Msg {
		err := n.Buffer.ReadDirectory(path)
		return ReadDirectoryMsg{err: err}
	}
}

// GetParentDir attempts to navigate to the parent directory
func (n *Navigator) GetParentDir() (string, error) {
	return n.Buffer.GetParentDir()
}

// GetCurrentDir returns the current directory path
func (n *Navigator) GetCurrentDir() string {
	return n.Buffer.GetCurrentDir()
}

// NumEntries returns the total number of entries
func (n *Navigator) NumEntries() int {
	return n.Buffer.NumEntries()
}

// GetEntry returns the entry at the specified index
func (n *Navigator) GetEntry(idx int) (buffer.Entry, bool) {
	return n.Buffer.GetEntry(idx)
}

// GetInput returns the current input buffer
func (n *Navigator) GetInput() string {
	return n.Buffer.GetInput()
}

// ReadDirectoryMsg is sent when a directory has been read
type ReadDirectoryMsg struct {
	err error
}
