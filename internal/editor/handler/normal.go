package handler

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/command/motion"
	"github.com/gunererd/grease/internal/keytree"
	"github.com/gunererd/grease/internal/register"
	"github.com/gunererd/grease/internal/state"
	"github.com/gunererd/grease/internal/types"
)

type NormalMode struct {
	keytree  *keytree.KeyTree
	register *register.Register
	// history types.HistoryManager
}

func NewNormalMode(kt *keytree.KeyTree, register *register.Register) *NormalMode {

	// Vim style Jump to beginning of buffer
	kt.Add(state.NormalMode, []string{"g", "g"}, keytree.KeyAction{
		Before: func(e types.Editor) types.Editor {
			e.Buffer().ClearCursors()
			return e
		},
		Execute: motion.CreateMotionCommand(motion.NewStartOfBufferMotion(), 0),
	})

	kt.Add(state.NormalMode, []string{"g", "e"}, keytree.KeyAction{
		Before: func(e types.Editor) types.Editor {
			e.Buffer().ClearCursors()
			return e
		},
		Execute: motion.CreateMotionCommand(motion.NewEndOfBufferMotion(), 0),
	})

	// // Undo command
	// kt.Add([]string{"u"}, keytree.KeyAction{
	// 	Execute: func(e types.Editor) types.Editor {
	// 		return history.Undo(e)
	// 	},
	// })

	// // Redo command
	// kt.Add([]string{"ctrl+r"}, keytree.KeyAction{
	// 	Execute: func(e types.Editor) types.Editor {
	// 		return history.Redo(e)
	// 	},
	// })

	// Word motion commands - change
	// kt.Add([]string{"c", "w"}, keytree.KeyAction{
	// 	Execute: CreateWordMotionCommand(false, NewHistoryAwareOperation(NewChangeOperation(), history)),
	// })
	// kt.Add([]string{"c", "W"}, keytree.KeyAction{
	// 	Execute: CreateWordMotionCommand(true, NewHistoryAwareOperation(NewChangeOperation(), history)),
	// })

	// kt.Add([]string{"c", "e"}, keytree.KeyAction{
	// 	Execute: CreateWordEndMotionCommand(false, NewHistoryAwareOperation(NewChangeOperation(), history)),
	// })

	// kt.Add([]string{"c", "E"}, keytree.KeyAction{
	// 	Execute: CreateWordEndMotionCommand(true, NewHistoryAwareOperation(NewChangeOperation(), history)),
	// })

	// kt.Add([]string{"c", "b"}, keytree.KeyAction{
	// 	Execute: CreateWordBackMotionCommand(false, NewHistoryAwareOperation(NewChangeOperation(), history)),
	// })
	// kt.Add([]string{"c", "B"}, keytree.KeyAction{
	// 	Execute: CreateWordBackMotionCommand(true, NewHistoryAwareOperation(NewChangeOperation(), history)),
	// })

	// // Word motion commands - delete
	// kt.Add([]string{"d", "w"}, keytree.KeyAction{
	// 	Execute: CreateWordMotionCommand(false, NewHistoryAwareOperation(NewDeleteOperation(), history)),
	// })
	// kt.Add([]string{"d", "W"}, keytree.KeyAction{
	// 	Execute: CreateWordMotionCommand(true, NewHistoryAwareOperation(NewDeleteOperation(), history)),
	// })

	// kt.Add([]string{"d", "e"}, keytree.KeyAction{
	// 	Execute: CreateWordEndMotionCommand(false, NewHistoryAwareOperation(NewDeleteOperation(), history)),
	// })

	// kt.Add([]string{"d", "E"}, keytree.KeyAction{
	// 	Execute: CreateWordEndMotionCommand(true, NewHistoryAwareOperation(NewDeleteOperation(), history)),
	// })

	// kt.Add([]string{"d", "b"}, keytree.KeyAction{
	// 	Execute: CreateWordBackMotionCommand(false, NewHistoryAwareOperation(NewDeleteOperation(), history)),
	// })
	// kt.Add([]string{"d", "B"}, keytree.KeyAction{
	// 	Execute: CreateWordBackMotionCommand(true, NewHistoryAwareOperation(NewDeleteOperation(), history)),
	// })

	// Word motion commands - yank
	kt.Add(state.NormalMode, []string{"y", "w"}, keytree.KeyAction{
		Execute: CreateYankCommand(motion.NewWordMotion(false), register).Execute,
	})
	kt.Add(state.NormalMode, []string{"y", "W"}, keytree.KeyAction{
		Execute: CreateYankCommand(motion.NewWordMotion(true), register).Execute,
	})

	kt.Add(state.NormalMode, []string{"y", "e"}, keytree.KeyAction{
		Execute: CreateYankCommand(motion.NewWordEndMotion(false), register).Execute,
	})
	kt.Add(state.NormalMode, []string{"y", "E"}, keytree.KeyAction{
		Execute: CreateYankCommand(motion.NewWordEndMotion(true), register).Execute,
	})

	kt.Add(state.NormalMode, []string{"y", "b"}, keytree.KeyAction{
		Execute: CreateYankCommand(motion.NewWordBackMotion(false), register).Execute,
	})
	kt.Add(state.NormalMode, []string{"y", "B"}, keytree.KeyAction{
		Execute: CreateYankCommand(motion.NewWordBackMotion(true), register).Execute,
	})

	return &NormalMode{
		keytree:  kt,
		register: register,
		// history: history,
	}
}

func (h *NormalMode) Handle(msg tea.KeyMsg, e types.Editor) (types.Editor, tea.Cmd) {

	model := e

	// Handle key sequences
	if handled, model := h.keytree.Handle(msg.String(), e); handled {
		e.HandleCursorMovement()
		return model, nil
	}

	// cursor, err := e.Buffer().GetPrimaryCursor()
	// if err != nil {
	// 	return e, nil
	// }

	switch msg.String() {
	case "ctrl+c":
		return e, tea.Quit
	case "C":
		log.Println("shift+c")
		buf := e.Buffer()
		cursor, err := buf.AddCursor()
		if err != nil {
			return e, nil
		}
		cmd := motion.CreateMotionCommand(motion.NewDownMotion(), cursor.ID())
		cmd(e)
	case "h":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := motion.CreateMotionCommand(motion.NewLeftMotion(), cursor.ID())
			cmd(e)
		}
	case "l":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := motion.CreateMotionCommand(motion.NewRightMotion(), cursor.ID())
			cmd(e)
		}
	case "j":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := motion.CreateMotionCommand(motion.NewDownMotion(), cursor.ID())
			cmd(e)
		}
	case "k":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := motion.CreateMotionCommand(motion.NewUpMotion(), cursor.ID())
			cmd(e)
		}
	case "i":
		e.SetMode(state.InsertMode)
	case "v":
		e.SetMode(state.VisualMode)
	case ":":
		e.SetMode(state.CommandMode)
	case "q":
		return e, tea.Quit
	case "$":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := motion.CreateMotionCommand(motion.NewEndOfLineMotion(), cursor.ID())
			cmd(e)
		}
	case "^", "0":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := motion.CreateMotionCommand(motion.NewStartOfLineMotion(), cursor.ID())
			cmd(e)
		}

		_ = 5
	case "w":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := motion.CreateMotionCommand(motion.NewWordMotion(false), cursor.ID())
			cmd(e)
		}
	case "W":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := motion.CreateMotionCommand(motion.NewWordMotion(true), cursor.ID())
			cmd(e)
		}
	case "e":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := motion.CreateMotionCommand(motion.NewWordEndMotion(false), cursor.ID())
			cmd(e)
		}
	case "E":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := motion.CreateMotionCommand(motion.NewWordEndMotion(true), cursor.ID())
			cmd(e)
		}
	case "b":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := motion.CreateMotionCommand(motion.NewWordBackMotion(false), cursor.ID())
			cmd(e)
		}
	case "B":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := motion.CreateMotionCommand(motion.NewWordBackMotion(true), cursor.ID())
			cmd(e)
		}
		// case "p":
		// 	cursors := e.Buffer().GetCursors()
		// 	for _, cursor := range cursors {
		// 		CreatePasteCommand(motion.NewWordMotion(false), h.register, false).Execute(e)
		// 	}
		// case "P":
		// 	model = CreatePasteCommand(motion.NewWordMotion(true), h.register, true).Execute(e)
	}

	e.HandleCursorMovement()
	return model, nil
}
