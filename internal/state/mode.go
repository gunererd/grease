package state

type Mode int

const (
	NormalMode Mode = iota
	InsertMode
	CommandMode
)
