package handler

// import (
// 	"github.com/gunererd/grease/internal/types"
// )

// type operationManager struct {
// 	historyManager types.HistoryManager
// 	factories      map[types.OperationType]func() types.Operation
// }

// func NewOperationManager(hm types.HistoryManager) types.OperationManager {
// 	om := &operationManager{
// 		historyManager: hm,
// 		factories:      make(map[types.OperationType]func() types.Operation),
// 	}

// 	om.factories[types.OpYank] = NewYankOperation
// 	om.factories[types.OpDelete] = NewDeleteOperation
// 	om.factories[types.OpChange] = NewChangeOperation
// 	om.factories[types.OpPaste] = func() types.Operation { return NewPasteOperation(false) }

// 	return om
// }

// func (om *operationManager) Execute(opType types.OperationType, e types.Editor, from, to types.Position) types.Editor {
// 	factory := om.factories[opType]
// 	if factory == nil {
// 		return nil
// 	}

// 	op := factory()
// 	if opType != types.OpYank { // Yank operations don't need history tracking
// 		op = NewHistoryAwareOperation(op, om.historyManager)
// 	}

// 	return op.Execute(e, from, to)
// }

// func (om *operationManager) CreateHistoryAwareOperation(op types.Operation) types.Operation {
// 	return NewHistoryAwareOperation(op, om.historyManager)
// }

// func (om *operationManager) GetOperationFactory(opType types.OperationType) func() types.Operation {
// 	return om.factories[opType]
// }
