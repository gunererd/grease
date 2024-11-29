package types

// Position represents a position in the buffer
type Position interface {
	String() string
	Before(other Position) bool
	Equal(other Position) bool
	Add(line, col int) Position
	Line() int
	Column() int
}
