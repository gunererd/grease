package hook

import (
	"github.com/gunererd/grease/internal/editor/types"
)

type HistoryHook struct {
	history      types.HistoryManager
	beforeLines  map[int]string
	cursorBefore types.Position
}

func NewHistoryHook(history types.HistoryManager) *HistoryHook {
	return &HistoryHook{
		history: history,
	}
}

func (h *HistoryHook) OnBeforeCommand(cmd types.Command, e types.Editor) {
	// Only track history for modifying commands
	if !isModifyingCommand(cmd.Name()) {
		return
	}

	buf := e.Buffer()
	cursor, _ := buf.GetPrimaryCursor()
	h.cursorBefore = cursor.GetPosition()

	h.beforeLines = make(map[int]string)
	for i := 0; i < buf.LineCount(); i++ {
		line, _ := buf.GetLine(i)
		h.beforeLines[i] = line
	}
}

func (h *HistoryHook) OnAfterCommand(cmd types.Command, e types.Editor) {
	// Only track history for modifying commands
	if !isModifyingCommand(cmd.Name()) {
		return
	}

	buf := e.Buffer()
	cursor, _ := buf.GetPrimaryCursor()
	cursorAfter := cursor.GetPosition()

	afterLines := make(map[int]string)
	for i := 0; i < buf.LineCount(); i++ {
		line, _ := buf.GetLine(i)
		afterLines[i] = line
	}

	entry := types.HistoryEntry{
		OperationType: cmd.Name(),
		BeforeLines:   h.beforeLines,
		AfterLines:    afterLines,
		CursorBefore:  h.cursorBefore,
		CursorAfter:   cursorAfter,
	}

	h.history.Push(entry)
}

func isModifyingCommand(cmdName string) bool {
	modifyingCommands := map[string]bool{
		"delete_line":          true,
		"change_line":          true,
		"delete":               true,
		"change":               true,
		"paste":                true,
		"insert":               true,
		"append":               true,
		"delete_to_end":        true,
		"change_to_end":        true,
		"new_line":             true,
		"append_end_of_line":   true,
		"insert_start_of_line": true,
	}
	return modifyingCommands[cmdName]
}
