package operation

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gunererd/grease/internal/filemanager/types"
)

type Executor struct {
	dirManager types.DirectoryManager
}

func NewExecutor(dirManager types.DirectoryManager) types.OperationExecutor {
	return &Executor{
		dirManager: dirManager,
	}
}

func (e *Executor) Execute(op types.Operation) error {
	if err := e.ValidateOperation(op); err != nil {
		return fmt.Errorf("operation validation failed: %w", err)
	}

	switch op.Type() {
	case types.Delete:
		return os.Remove(op.Source())
	case types.Rename:
		return os.Rename(op.Source(), op.Target())
	case types.Move:
		target := filepath.Join(op.Target(), filepath.Base(op.Source()))
		return os.Rename(op.Source(), target)
	case types.Create:
		if op.Source()[len(op.Source())-1] == '/' {
			return os.MkdirAll(op.Source(), 0755)
		}
		f, err := os.Create(op.Source())
		if err != nil {
			return err
		}
		return f.Close()
	default:
		return fmt.Errorf("unknown operation type: %v", op.Type())
	}
}

func (e *Executor) ValidateOperation(op types.Operation) error {
	switch op.Type() {
	case types.Delete:
		if _, err := os.Stat(op.Source()); err != nil {
			return fmt.Errorf("source does not exist: %w", err)
		}
	case types.Rename, types.Move:
		if _, err := os.Stat(op.Source()); err != nil {
			return fmt.Errorf("source does not exist: %w", err)
		}
		if _, err := os.Stat(op.Target()); err == nil {
			return fmt.Errorf("target already exists")
		}
	case types.Create:
		if _, err := os.Stat(op.Source()); err == nil {
			return fmt.Errorf("file/directory already exists")
		}
	}
	return nil
}
