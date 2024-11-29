package types

import (
	"io"
)

// Buffer represents the text content and provides operations to modify it
type Buffer interface {
	Get() string
	GetLine(line int) (string, error)
	LineCount() int
	LineLen(line int) (int, error)
	LoadFromReader(r io.Reader) error

	AddCursor(pos Position, priority int) (*Cursor, error)
	RemoveCursor(id int)
	GetPrimaryCursor() (*Cursor, error)
	MoveCursor(cursorID int, lineOffset, columnOffset int) error

	Insert(text string) error
	Delete(count int) error
}
