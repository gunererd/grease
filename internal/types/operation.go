package types

// Operation defines what actions can be performed between two positions in a buffer
type Operation interface {
	Execute(e Editor, from, to Position) Editor
}
