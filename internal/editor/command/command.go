package command

import "github.com/gunererd/grease/internal/editor/types"

type Command interface {
	// Execute runs the command and returns the modified editor
	Execute(e types.Editor) types.Editor

	// Name returns the command name for command mode
	Name() string
}
