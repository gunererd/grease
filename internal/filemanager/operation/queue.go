package operation

import (
	"github.com/gunererd/grease/internal/filemanager/types"
)

type OperationQueue struct {
	operations []types.Operation
	executor   types.OperationExecutor
}

func NewOperationQueue(executor types.OperationExecutor) *OperationQueue {
	return &OperationQueue{
		operations: make([]types.Operation, 0),
		executor:   executor,
	}
}

func (q *OperationQueue) Operations() []types.Operation {
	return q.operations
}

func (q *OperationQueue) Push(op types.Operation) {
	q.operations = append(q.operations, op)
}

func (q *OperationQueue) Clear() {
	q.operations = q.operations[:0]
}

func (q *OperationQueue) IsEmpty() bool {
	return len(q.operations) == 0
}

func (q *OperationQueue) Execute() error {
	for _, op := range q.operations {
		if err := q.executor.Execute(op); err != nil {
			return err
		}
	}
	q.Clear()
	return nil
}
