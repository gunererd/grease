package handler

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/editor/keytree"
	"github.com/gunererd/grease/internal/editor/register"
	"github.com/gunererd/grease/internal/editor/state"
	"github.com/gunererd/grease/internal/editor/types"
)

type CommandMode struct {
	buffer   string
	executor *CommandExecutor
	logger   types.Logger
}

func NewCommandMode(
	kt *keytree.KeyTree,
	register *register.Register,
	hlm types.HighlightManager,
	logger types.Logger,
) *CommandMode {
	return &CommandMode{
		buffer: "",
		// executor: executor,
		logger: logger,
	}
}

func (h *CommandMode) Handle(msg tea.KeyMsg, e types.Editor) (types.Editor, tea.Cmd) {

	switch msg.String() {
	case "esc", "ctrl+c":
		h.buffer = ""
		e.SetMode(state.NormalMode)
	case "enter":
		e = h.executeCommand(e)
		h.buffer = ""
		e.SetMode(state.NormalMode)
	case "backspace":
		if len(h.buffer) > 0 {
			h.buffer = h.buffer[:len(h.buffer)-1]
		}
	default:
		h.buffer += string(msg.Runes)
	}

	h.logger.Printf("Buffer: %s", h.buffer)

	return e, nil
}

func (h *CommandMode) executeCommand(e types.Editor) types.Editor {
	// Parse command from buffer
	cmd := parseCommand(h.buffer)
	if cmd == nil {
		h.logger.Printf("Invalid command: %s", h.buffer)
		return e
	}

	return h.executor.Execute(cmd, e)
}

func (h *CommandMode) GetBuffer() string {
	return h.buffer
}

func parseCommand(input string) types.Command {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil
	}

	switch parts[0] {
	case "w", "write":
		return CreateWriteCommand()
	}
	return nil
}
