package command

import "github.com/gunererd/grease/internal/editor/types"

type CompositeCommand struct {
	commands []Command
	name     string
}

func NewCompositeCommand(name string, commands ...Command) *CompositeCommand {
	return &CompositeCommand{
		commands: commands,
		name:     name,
	}
}

func (c *CompositeCommand) Execute(e types.Editor) types.Editor {
	for _, cmd := range c.commands {
		e = cmd.Execute(e)
	}
	return e
}

func (c *CompositeCommand) Name() string {
	return c.name
}
