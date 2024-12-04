package handler

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/keytree"
	"github.com/gunererd/grease/internal/state"
	"github.com/gunererd/grease/internal/types"
)

type NormalMode struct {
	keytree *keytree.KeyTree
}

func NewNormalMode(kt *keytree.KeyTree) *NormalMode {

	// Vim style Jump to beginning of buffer
	kt.Add([]string{"g", "g"}, keytree.KeyAction{
		Execute: func(e types.Editor) (tea.Model, tea.Cmd) {
			cursor, _ := e.Buffer().GetPrimaryCursor()
			e.Buffer().MoveCursor(cursor.ID(), 0, 0)
			e.HandleCursorMovement()
			return e, nil
		},
	})

	// Word motion commands - change
	kt.Add([]string{"c", "w"}, keytree.KeyAction{
		Execute: CreateWordMotionCommand(false, NewChangeOperation()),
	})
	kt.Add([]string{"c", "W"}, keytree.KeyAction{
		Execute: CreateWordMotionCommand(true, NewChangeOperation()),
	})

	kt.Add([]string{"c", "e"}, keytree.KeyAction{
		Execute: CreateWordEndMotionCommand(false, NewChangeOperation()),
	})

	kt.Add([]string{"c", "E"}, keytree.KeyAction{
		Execute: CreateWordEndMotionCommand(true, NewChangeOperation()),
	})

	kt.Add([]string{"c", "b"}, keytree.KeyAction{
		Execute: CreateWordBackMotionCommand(false, NewChangeOperation()),
	})
	kt.Add([]string{"c", "B"}, keytree.KeyAction{
		Execute: CreateWordBackMotionCommand(true, NewChangeOperation()),
	})

	// Word motion commands - delete
	kt.Add([]string{"d", "w"}, keytree.KeyAction{
		Execute: CreateWordMotionCommand(false, NewDeleteOperation()),
	})
	kt.Add([]string{"d", "W"}, keytree.KeyAction{
		Execute: CreateWordMotionCommand(true, NewDeleteOperation()),
	})

	kt.Add([]string{"d", "e"}, keytree.KeyAction{
		Execute: CreateWordEndMotionCommand(false, NewDeleteOperation()),
	})

	kt.Add([]string{"d", "E"}, keytree.KeyAction{
		Execute: CreateWordEndMotionCommand(true, NewDeleteOperation()),
	})

	kt.Add([]string{"d", "b"}, keytree.KeyAction{
		Execute: CreateWordBackMotionCommand(false, NewDeleteOperation()),
	})
	kt.Add([]string{"d", "B"}, keytree.KeyAction{
		Execute: CreateWordBackMotionCommand(true, NewDeleteOperation()),
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
	}
}

func (h *NormalMode) Handle(msg tea.KeyMsg, e types.Editor) (tea.Model, tea.Cmd) {

	// Handle key sequences
	if handled, model, cmd := h.keytree.Handle(msg.String(), e); handled {
		e.HandleCursorMovement()
		return model, cmd
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
	case "^":
		// Vim style jump to beginning of line
		cursor, _ := e.Buffer().GetPrimaryCursor()
		line := cursor.GetPosition().Line()
		e.Buffer().MoveCursor(cursor.ID(), line, 0)
		e.HandleCursorMovement()
		return e, nil
	case "w":
		model, cmd := CreateWordMotionCommand(false, nil)(e)
		e.HandleCursorMovement()
		return model, cmd
	case "W":
		model, cmd := CreateWordMotionCommand(true, nil)(e)
		e.HandleCursorMovement()
		return model, cmd
	case "e":
		model, cmd := CreateWordEndMotionCommand(false, nil)(e)
		e.HandleCursorMovement()
		return model, cmd
	case "E":
		model, cmd := CreateWordEndMotionCommand(true, nil)(e)
		e.HandleCursorMovement()
		return model, cmd
	case "b":
		model, cmd := CreateWordBackMotionCommand(false, nil)(e)
		e.HandleCursorMovement()
		return model, cmd
	case "B":
		model, cmd := CreateWordBackMotionCommand(true, nil)(e)
		e.HandleCursorMovement()
		return model, cmd
	case "p":
		model, cmd := NewPasteOperation(false).Execute(e, cursor.GetPosition(), cursor.GetPosition())
		e.HandleCursorMovement()
		return model, cmd
	case "P":
		model, cmd := NewPasteOperation(true).Execute(e, cursor.GetPosition(), cursor.GetPosition())
		e.HandleCursorMovement()
		return model, cmd
	}
	return e, nil
}
