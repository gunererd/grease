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

	AddCursor() (Cursor, error)
	RemoveCursor(id int)
	GetPrimaryCursor() (Cursor, error)
	GetCursor(id int) (Cursor, error)
	GetCursors() []Cursor
	ClearCursors()
	MoveCursorRelative(cursorID int, lineOffset, columnOffset int) error
	MoveCursor(cursorID int, lineOffset, columnOffset int) error
	NextWordPosition(pos Position, bigWord bool) Position
	NextWordEndPosition(pos Position, bigWord bool) Position
	PrevWordPosition(pos Position, bigWord bool) Position

	Insert(text string) error
	Delete(count int) error
	ReplaceLine(line int, content string) error
	InsertLine(line int, content string) error
	RemoveLine(line int) error
}
