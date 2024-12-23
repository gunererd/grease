package types

type Hook interface {
	OnBeforeCommand(cmd Command, e Editor)
	OnAfterCommand(cmd Command, e Editor)
}

type HookManager interface {
	AddHook(h Hook)
	RemoveHook(h Hook)
	GetHooks() []Hook
	ExecuteBeforeHooks(cmd Command, e Editor)
	ExecuteAfterHooks(cmd Command, e Editor)
}
