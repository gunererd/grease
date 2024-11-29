package types

// Cursor represents a cursor in the buffer
type Cursor interface {
	GetPosition() Position
	SetPosition(pos Position)
	GetID() int
	GetPriority() int
	SetPriority(priority int)
}
