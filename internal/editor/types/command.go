package types

type Command interface {
	// Execute runs the command and returns the modified editor
	Execute(e Editor) Editor

	// Name returns the command name for command mode
	Name() string
	Explain() string
}
