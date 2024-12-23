package filemanager

type OperationType int

const (
	Delete OperationType = iota
	Rename
	Move
)

type Operation struct {
	Type   OperationType
	Source string
	Target string
}
