package handler

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/keytree"
	"github.com/gunererd/grease/internal/state"
	"github.com/gunererd/grease/internal/types"
)

type NormalMode struct {
	keytree *keytree.KeyTree
	history types.HistoryManager
}

func NewNormalMode(kt *keytree.KeyTree, history types.HistoryManager, om types.OperationManager) *NormalMode {

	// Vim style Jump to beginning of buffer
	kt.Add([]string{"g", "g"}, keytree.KeyAction{
		Execute: func(e types.Editor) types.Editor {
			cursor, _ := e.Buffer().GetPrimaryCursor()
			e.Buffer().MoveCursor(cursor.ID(), 0, 0)
			e.HandleCursorMovement()
			return e
		},
	})

	// Undo command
	kt.Add([]string{"u"}, keytree.KeyAction{
		Execute: func(e types.Editor) types.Editor {
			return history.Undo(e)
		},
	})

	// Redo command
	kt.Add([]string{"ctrl+r"}, keytree.KeyAction{
		Execute: func(e types.Editor) types.Editor {
			return history.Redo(e)
		},
	})

	// Word motion commands - change
	kt.Add([]string{"c", "w"}, keytree.KeyAction{
		Execute: CreateWordMotionCommand(false, NewHistoryAwareOperation(NewChangeOperation(), history)),
	})
	kt.Add([]string{"c", "W"}, keytree.KeyAction{
		Execute: CreateWordMotionCommand(true, NewHistoryAwareOperation(NewChangeOperation(), history)),
	})

	kt.Add([]string{"c", "e"}, keytree.KeyAction{
		Execute: CreateWordEndMotionCommand(false, NewHistoryAwareOperation(NewChangeOperation(), history)),
	})

	kt.Add([]string{"c", "E"}, keytree.KeyAction{
		Execute: CreateWordEndMotionCommand(true, NewHistoryAwareOperation(NewChangeOperation(), history)),
	})

	kt.Add([]string{"c", "b"}, keytree.KeyAction{
		Execute: CreateWordBackMotionCommand(false, NewHistoryAwareOperation(NewChangeOperation(), history)),
	})
	kt.Add([]string{"c", "B"}, keytree.KeyAction{
		Execute: CreateWordBackMotionCommand(true, NewHistoryAwareOperation(NewChangeOperation(), history)),
	})

	// Word motion commands - delete
	kt.Add([]string{"d", "w"}, keytree.KeyAction{
		Execute: CreateWordMotionCommand(false, NewHistoryAwareOperation(NewDeleteOperation(), history)),
	})
	kt.Add([]string{"d", "W"}, keytree.KeyAction{
		Execute: CreateWordMotionCommand(true, NewHistoryAwareOperation(NewDeleteOperation(), history)),
	})

	kt.Add([]string{"d", "e"}, keytree.KeyAction{
		Execute: CreateWordEndMotionCommand(false, NewHistoryAwareOperation(NewDeleteOperation(), history)),
	})

	kt.Add([]string{"d", "E"}, keytree.KeyAction{
		Execute: CreateWordEndMotionCommand(true, NewHistoryAwareOperation(NewDeleteOperation(), history)),
	})

	kt.Add([]string{"d", "b"}, keytree.KeyAction{
		Execute: CreateWordBackMotionCommand(false, NewHistoryAwareOperation(NewDeleteOperation(), history)),
	})
	kt.Add([]string{"d", "B"}, keytree.KeyAction{
		Execute: CreateWordBackMotionCommand(true, NewHistoryAwareOperation(NewDeleteOperation(), history)),
	})

	// Word motion commands - yank
	kt.Add([]string{"y", "w"}, keytree.KeyAction{
		Execute: CreateWordMotionCommand(false, NewYankOperation()),
	})
	kt.Add([]string{"y", "W"}, keytree.KeyAction{
		Execute: CreateWordMotionCommand(true, NewYankOperation()),
	})
	kt.Add([]string{"y", "e"}, keytree.KeyAction{
		Execute: CreateWordEndMotionCommand(false, NewYankOperation()),
	})
	kt.Add([]string{"y", "E"}, keytree.KeyAction{
		Execute: CreateWordEndMotionCommand(true, NewYankOperation()),
	})
	kt.Add([]string{"y", "b"}, keytree.KeyAction{
		Execute: CreateWordBackMotionCommand(false, NewYankOperation()),
	})
	kt.Add([]string{"y", "B"}, keytree.KeyAction{
		Execute: CreateWordBackMotionCommand(true, NewYankOperation()),
	})

	return &NormalMode{
		keytree: kt,
		history: history,
	}
}

func (h *NormalMode) Handle(msg tea.KeyMsg, e types.Editor) (types.Editor, tea.Cmd) {

	// Handle key sequences
	if handled, model := h.keytree.Handle(msg.String(), e); handled {
		e.HandleCursorMovement()
		return model, nil
	}

	cursor, err := e.Buffer().GetPrimaryCursor()
	if err != nil {
		return e, nil
	}

	switch msg.String() {
	case "ctrl+c":
		return e, tea.Quit
	case "h":
		e.Buffer().MoveCursorRelative(cursor.ID(), 0, -1)
		e.HandleCursorMovement()
	case "l":
		e.Buffer().MoveCursorRelative(cursor.ID(), 0, 1)
		e.HandleCursorMovement()
	case "j":
		e.Buffer().MoveCursorRelative(cursor.ID(), 1, 0)
		e.HandleCursorMovement()
	case "k":
		e.Buffer().MoveCursorRelative(cursor.ID(), -1, 0)
		e.HandleCursorMovement()
	case "i":
		e.SetMode(state.InsertMode)
	case "v":
		e.SetMode(state.VisualMode)
	case ":":
		e.SetMode(state.CommandMode)
	case "q":
		return e, tea.Quit
	case "G":
		// Vim style end of buffer
		cursor, _ := e.Buffer().GetPrimaryCursor()
		e.Buffer().MoveCursor(cursor.ID(), e.Buffer().LineCount()-1, 0)
		e.HandleCursorMovement()
		return e, nil
	case "$":
		// Vim style jump to end of line
		cursor, _ := e.Buffer().GetPrimaryCursor()
		line := cursor.GetPosition().Line()
		lineLength, _ := e.Buffer().LineLen(line)
		e.Buffer().MoveCursor(cursor.ID(), line, lineLength)
		e.HandleCursorMovement()
		return e, nil
	case "^", "0":
		// Vim style jump to beginning of line
		cursor, _ := e.Buffer().GetPrimaryCursor()
		line := cursor.GetPosition().Line()
		e.Buffer().MoveCursor(cursor.ID(), line, 0)
		e.HandleCursorMovement()
		return e, nil
	case "w":
		model := CreateWordMotionCommand(false, nil)(e)
		e.HandleCursorMovement()
		return model, nil
	case "W":
		model := CreateWordMotionCommand(true, nil)(e)
		e.HandleCursorMovement()
		return model, nil
	case "e":
		model := CreateWordEndMotionCommand(false, nil)(e)
		e.HandleCursorMovement()
		return model, nil
	case "E":
		model := CreateWordEndMotionCommand(true, nil)(e)
		e.HandleCursorMovement()
		return model, nil
	case "b":
		model := CreateWordBackMotionCommand(false, nil)(e)
		e.HandleCursorMovement()
		return model, nil
	case "B":
		model := CreateWordBackMotionCommand(true, nil)(e)
		e.HandleCursorMovement()
		return model, nil
	case "p":
		model := NewHistoryAwareOperation(NewPasteOperation(false), e.HistoryManager()).Execute(e, cursor.GetPosition(), cursor.GetPosition())
		e.HandleCursorMovement()
		return model, nil
	case "P":
		model := NewHistoryAwareOperation(NewPasteOperation(true), e.HistoryManager()).Execute(e, cursor.GetPosition(), cursor.GetPosition())
		e.HandleCursorMovement()
		return model, nil
	}
	return e, nil
}
