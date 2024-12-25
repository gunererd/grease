package handler

import (
	"github.com/gunererd/grease/internal/editor/hook"
	"github.com/gunererd/grease/internal/editor/types"
)

// CommandExecutor handles all command executions and hooks
type CommandExecutor struct {
	hookManager types.HookManager
	logger      types.Logger
}

func NewCommandExecutor(history types.HistoryManager, hookManager types.HookManager, logger types.Logger) *CommandExecutor {
	ce := &CommandExecutor{
		hookManager: hookManager,
		logger:      logger,
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
	ce.logger.Println(cmd.Explain())

	ce.hookManager.ExecuteAfterHooks(cmd, e)

	return e
}
