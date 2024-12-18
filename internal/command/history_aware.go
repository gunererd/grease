package command

import (
	"github.com/gunererd/grease/internal/types"
)

type HistoryAwareCommand struct {
	command Command
	history types.HistoryManager
}

func NewHistoryAwareCommand(cmd Command, history types.HistoryManager) Command {
	return &HistoryAwareCommand{
		command: cmd,
		history: history,
	}
}

func (h *HistoryAwareCommand) Execute(e types.Editor) types.Editor {
	// Capture state before command execution
	buf := e.Buffer()
	cursor, _ := buf.GetPrimaryCursor()
	cursorBefore := cursor.GetPosition()

	beforeLines := make(map[int]string)
	for i := 0; i < buf.LineCount(); i++ {
		line, _ := buf.GetLine(i)
		beforeLines[i] = line
	}

	e = h.command.Execute(e)

	// Capture state after command execution
	cursor, _ = buf.GetPrimaryCursor()
	cursorAfter := cursor.GetPosition()

	afterLines := make(map[int]string)
	for i := 0; i < buf.LineCount(); i++ {
		line, _ := buf.GetLine(i)
		afterLines[i] = line
	}

	entry := types.HistoryEntry{
		OperationType: h.command.Name(),
		BeforeLines:   beforeLines,
		AfterLines:    afterLines,
		CursorBefore:  cursorBefore,
		CursorAfter:   cursorAfter,
	}

	h.history.Push(entry)
	return e
}

func (h *HistoryAwareCommand) Name() string {
	return h.command.Name()
}
