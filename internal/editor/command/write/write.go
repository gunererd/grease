package write

import "github.com/gunererd/grease/internal/editor/types"

type WriteCommand struct{}

func NewWriteCommand() *WriteCommand {
	return &WriteCommand{}
}

func (c *WriteCommand) Execute(e types.Editor) types.Editor {
	// TODO: Implement write command
	return e
}

func (c *WriteCommand) Name() string {
	return "write"
}

func (c *WriteCommand) Explain() string {
	return "Save current buffer to file"
}
