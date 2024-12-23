package types

type OperationType string

const (
	OpYank   OperationType = "yank"
	OpDelete OperationType = "delete"
	OpChange OperationType = "change"
	OpPaste  OperationType = "paste"
)

// Operation defines what actions can be performed between two positions in a buffer
type Operation interface {
	Execute(e Editor, from, to Position) Editor
}

type OperationManager interface {
	Execute(opType OperationType, e Editor, from, to Position) Editor

	// CreateHistoryAwareOperation wraps an operation with history tracking
	CreateHistoryAwareOperation(op Operation) Operation

	GetOperationFactory(opType OperationType) func() Operation
}
