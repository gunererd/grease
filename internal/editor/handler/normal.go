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
		Execute: createWordMotionCommand(nextWordStart, false, true),
	})
	kt.Add([]string{"c", "W"}, keytree.KeyAction{
		Execute: createWordMotionCommand(nextWordStart, true, true),
	})
	kt.Add([]string{"c", "e"}, keytree.KeyAction{
		Execute: createWordMotionCommand(nextWordEnd, false, true),
	})
	kt.Add([]string{"c", "E"}, keytree.KeyAction{
		Execute: createWordMotionCommand(nextWordEnd, true, true),
	})
	kt.Add([]string{"c", "b"}, keytree.KeyAction{
		Execute: createWordMotionCommand(prevWordStart, false, true),
	})
	kt.Add([]string{"c", "B"}, keytree.KeyAction{
		Execute: createWordMotionCommand(prevWordStart, true, true),
	})

	// Word motion commands - delete
	kt.Add([]string{"d", "w"}, keytree.KeyAction{
		Execute: createWordMotionCommand(nextWordStart, false, false),
	})
	kt.Add([]string{"d", "W"}, keytree.KeyAction{
		Execute: createWordMotionCommand(nextWordStart, true, false),
	})
	kt.Add([]string{"d", "e"}, keytree.KeyAction{
		Execute: createWordMotionCommand(nextWordEnd, false, false),
	})
	kt.Add([]string{"d", "E"}, keytree.KeyAction{
		Execute: createWordMotionCommand(nextWordEnd, true, false),
	})
	kt.Add([]string{"d", "b"}, keytree.KeyAction{
		Execute: createWordMotionCommand(prevWordStart, false, false),
	})
	kt.Add([]string{"d", "B"}, keytree.KeyAction{
		Execute: createWordMotionCommand(prevWordStart, true, false),
	})

	return &NormalMode{
		keytree: kt,
	}
}

func (h *NormalMode) Handle(msg tea.KeyMsg, e types.Editor) (tea.Model, tea.Cmd) {

	// Handle key sequences
	if handled, model, cmd := h.keytree.Handle(msg.String(), e); handled {
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
		cursor, _ := e.Buffer().GetPrimaryCursor()
		newPos := e.Buffer().NextWordPosition(cursor.GetPosition(), false)
		e.Buffer().MoveCursor(cursor.ID(), newPos.Line(), newPos.Column())
		e.HandleCursorMovement()
	case "W":
		cursor, _ := e.Buffer().GetPrimaryCursor()
		newPos := e.Buffer().NextWordPosition(cursor.GetPosition(), true)
		e.Buffer().MoveCursor(cursor.ID(), newPos.Line(), newPos.Column())
		e.HandleCursorMovement()
	case "e":
		cursor, _ := e.Buffer().GetPrimaryCursor()
		newPos := e.Buffer().NextWordEndPosition(cursor.GetPosition(), false)
		e.Buffer().MoveCursor(cursor.ID(), newPos.Line(), newPos.Column())
		e.HandleCursorMovement()
	case "E":
		cursor, _ := e.Buffer().GetPrimaryCursor()
		newPos := e.Buffer().NextWordEndPosition(cursor.GetPosition(), true)
		e.Buffer().MoveCursor(cursor.ID(), newPos.Line(), newPos.Column())
		e.HandleCursorMovement()
	case "b":
		cursor, _ := e.Buffer().GetPrimaryCursor()
		newPos := e.Buffer().PrevWordPosition(cursor.GetPosition(), false)
		e.Buffer().MoveCursor(cursor.ID(), newPos.Line(), newPos.Column())
		e.HandleCursorMovement()
	case "B":
		cursor, _ := e.Buffer().GetPrimaryCursor()
		newPos := e.Buffer().PrevWordPosition(cursor.GetPosition(), true)
		e.Buffer().MoveCursor(cursor.ID(), newPos.Line(), newPos.Column())
		e.HandleCursorMovement()
	}
	return e, nil
}
