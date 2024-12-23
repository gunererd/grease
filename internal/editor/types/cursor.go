package types

// Cursor represents a cursor in the buffer
type Cursor interface {
	GetPosition() Position
	SetPosition(pos Position)
	ID() int
	GetPriority() int
	SetPriority(priority int)
}
