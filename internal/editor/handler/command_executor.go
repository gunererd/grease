package handler

import (
	"github.com/gunererd/grease/internal/editor/hook"
	"github.com/gunererd/grease/internal/editor/types"
)

// CommandExecutor handles all command executions and hooks
type CommandExecutor struct {
	hookManager types.HookManager
}

func NewCommandExecutor(history types.HistoryManager, hookManager types.HookManager) *CommandExecutor {
	ce := &CommandExecutor{
		hookManager: hookManager,
	}

	historyHook := hook.NewHistoryHook(history)
	ce.AddHook(historyHook)

	return ce
}

func (ce *CommandExecutor) AddHook(h types.Hook) {
	ce.hookManager.AddHook(h)
}

func (ce *CommandExecutor) Execute(cmd types.Command, e types.Editor) types.Editor {
	ce.hookManager.ExecuteBeforeHooks(cmd, e)

	e = cmd.Execute(e)
	cmd.Explain()

	ce.hookManager.ExecuteAfterHooks(cmd, e)

	return e
}
