package handler

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/command/motion"
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
		Execute: motion.CreateBasicMotionCommand(motion.NewStartOfBufferMotion()),
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
		return motion.CreateBasicMotionCommand(motion.NewLeftMotion())(e), nil
	case "l":
		return motion.CreateBasicMotionCommand(motion.NewRightMotion())(e), nil
	case "j":
		return motion.CreateBasicMotionCommand(motion.NewDownMotion())(e), nil
	case "k":
		return motion.CreateBasicMotionCommand(motion.NewUpMotion())(e), nil
	case "i":
		e.SetMode(state.InsertMode)
	case "v":
		e.SetMode(state.VisualMode)
	case ":":
		e.SetMode(state.CommandMode)
	case "q":
		return e, tea.Quit
	case "G":
		return motion.CreateBasicMotionCommand(motion.NewEndOfBufferMotion())(e), nil
	case "g":
		// Handle 'gg' sequence through keytree
		if handled, model := h.keytree.Handle(msg.String(), e); handled {
			return model, nil
		}
	case "$":
		return motion.CreateBasicMotionCommand(motion.NewEndOfLineMotion())(e), nil
	case "^", "0":
		return motion.CreateBasicMotionCommand(motion.NewStartOfLineMotion())(e), nil
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
