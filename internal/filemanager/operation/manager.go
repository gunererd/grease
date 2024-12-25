package operation

import (
	"github.com/gunererd/grease/internal/filemanager/types"
)

type Manager struct {
	queue      *OperationQueue
	executor   types.OperationExecutor
	dirManager types.DirectoryManager
	logger     types.Logger
}

func NewOperationManager(dirManager types.DirectoryManager, logger types.Logger) types.OperationManager {
	executor := NewExecutor(dirManager)
	return &Manager{
		queue:      NewOperationQueue(executor),
		executor:   executor,
		dirManager: dirManager,
		logger:     logger,
	}
}

func (m *Manager) QueueOperation(op types.Operation) {
	m.queue.Push(op)
}

func (m *Manager) ExecuteOperations() error {
	if err := m.queue.Execute(); err != nil {
		return err
	}

	// After executing operations, reload the directory
	return nil
}

func (m *Manager) GetPendingOperations() []types.Operation {
	return m.queue.Operations()
}

func (m *Manager) Clear() {
	m.queue.Clear()
}
