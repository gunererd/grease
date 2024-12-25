package types

type OperationQueue interface {
	Push(op Operation)
	Clear()
	IsEmpty() bool
	GetOperationDescriptions() []string
	Execute() error
	Operations() []Operation
}

type OperationType int

const (
	Delete OperationType = iota
	Rename
	Move
	Create
)

type Operation interface {
	Type() OperationType
	Source() string
	Target() string
}

type OperationManager interface {
	QueueOperation(op Operation)
	ExecuteOperations() error
	GetPendingOperations() []Operation
	Clear()
}

type OperationExecutor interface {
	Execute(op Operation) error
	ValidateOperation(op Operation) error
}
