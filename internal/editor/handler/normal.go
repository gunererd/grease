package handler

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/editor/command/motion"
	"github.com/gunererd/grease/internal/editor/keytree"
	"github.com/gunererd/grease/internal/editor/register"
	"github.com/gunererd/grease/internal/editor/state"
	"github.com/gunererd/grease/internal/editor/types"
)

type NormalMode struct {
	keytree  *keytree.KeyTree
	register *register.Register
	history  types.HistoryManager
	executor *CommandExecutor
	logger   types.Logger
}

func NewNormalMode(
	kt *keytree.KeyTree,
	register *register.Register,
	history types.HistoryManager,
	executor *CommandExecutor,
	logger types.Logger,
) *NormalMode {
	nm := &NormalMode{
		register: register,
		history:  history,
		executor: executor,
		logger:   logger,
	}

	// Vim style Jump to beginning of buffer
	kt.Add(state.NormalMode, []string{"g", "g"}, keytree.KeyAction{
		Before: func(e types.Editor) types.Editor {
			e.Buffer().ClearCursors()
			return e
		},
		Execute: func(e types.Editor) types.Editor {
			cursor, err := e.Buffer().GetPrimaryCursor()
			if err != nil {
				nm.logger.Println("Failed to get primary cursor:", err)
				return e
			}
			cmd := CreateGoToStartOfBufferCommand(cursor)
			return nm.executor.Execute(cmd, e)
		},
	})

	kt.Add(state.NormalMode, []string{"d", "d"}, keytree.KeyAction{
		Before: func(e types.Editor) types.Editor {
			e.Buffer().ClearCursors()
			return e
		},
		Execute: func(e types.Editor) types.Editor {
			cursor, err := e.Buffer().GetPrimaryCursor()
			if err != nil {
				nm.logger.Println("Failed to get primary cursor:", err)
				return e
			}
			cmd := CreateDeleteLineCommand(cursor)
			return nm.executor.Execute(cmd, e)
		},
	})

	kt.Add(state.NormalMode, []string{"c", "c"}, keytree.KeyAction{
		Before: func(e types.Editor) types.Editor {
			e.Buffer().ClearCursors()
			return e
		},
		Execute: func(e types.Editor) types.Editor {
			cursor, err := e.Buffer().GetPrimaryCursor()
			if err != nil {
				nm.logger.Println("Failed to get primary cursor:", err)
				return e
			}
			cmd := CreateChangeLineCommand(cursor)
			return nm.executor.Execute(cmd, e)
		},
	})

	kt.Add(state.NormalMode, []string{"d", "w"}, keytree.KeyAction{
		Execute: func(e types.Editor) types.Editor {
			cmd := CreateDeleteCommand(motion.NewWordMotion(false))
			return nm.executor.Execute(cmd, e)
		},
	})

	kt.Add(state.NormalMode, []string{"d", "W"}, keytree.KeyAction{
		Execute: func(e types.Editor) types.Editor {
			cmd := CreateDeleteCommand(motion.NewWordMotion(true))
			return nm.executor.Execute(cmd, e)
		},
	})

	kt.Add(state.NormalMode, []string{"d", "e"}, keytree.KeyAction{
		Execute: func(e types.Editor) types.Editor {
			cmd := CreateDeleteCommand(motion.NewWordEndMotion(false))
			return nm.executor.Execute(cmd, e)
		},
	})

	kt.Add(state.NormalMode, []string{"d", "E"}, keytree.KeyAction{
		Execute: func(e types.Editor) types.Editor {
			cmd := CreateDeleteCommand(motion.NewWordEndMotion(true))
			return nm.executor.Execute(cmd, e)
		},
	})

	kt.Add(state.NormalMode, []string{"d", "b"}, keytree.KeyAction{
		Execute: func(e types.Editor) types.Editor {
			cmd := CreateDeleteCommand(motion.NewWordBackMotion(false))
			return nm.executor.Execute(cmd, e)
		},
	})

	kt.Add(state.NormalMode, []string{"d", "B"}, keytree.KeyAction{
		Execute: func(e types.Editor) types.Editor {
			cmd := CreateDeleteCommand(motion.NewWordBackMotion(true))
			return nm.executor.Execute(cmd, e)
		},
	})

	kt.Add(state.NormalMode, []string{"c", "w"}, keytree.KeyAction{
		Execute: func(e types.Editor) types.Editor {
			cmd := CreateChangeCommand(motion.NewWordMotion(false))
			return nm.executor.Execute(cmd, e)
		},
	})

	kt.Add(state.NormalMode, []string{"c", "W"}, keytree.KeyAction{
		Execute: func(e types.Editor) types.Editor {
			cmd := CreateChangeCommand(motion.NewWordMotion(true))
			return nm.executor.Execute(cmd, e)
		},
	})

	kt.Add(state.NormalMode, []string{"c", "e"}, keytree.KeyAction{
		Execute: func(e types.Editor) types.Editor {
			cmd := CreateChangeCommand(motion.NewWordEndMotion(false))
			return nm.executor.Execute(cmd, e)
		},
	})

	kt.Add(state.NormalMode, []string{"c", "E"}, keytree.KeyAction{
		Execute: func(e types.Editor) types.Editor {
			cmd := CreateChangeCommand(motion.NewWordEndMotion(true))
			return nm.executor.Execute(cmd, e)
		},
	})

	kt.Add(state.NormalMode, []string{"c", "b"}, keytree.KeyAction{
		Execute: func(e types.Editor) types.Editor {
			cmd := CreateChangeCommand(motion.NewWordBackMotion(false))
			return nm.executor.Execute(cmd, e)
		},
	})

	kt.Add(state.NormalMode, []string{"c", "B"}, keytree.KeyAction{
		Execute: func(e types.Editor) types.Editor {
			cmd := CreateChangeCommand(motion.NewWordBackMotion(true))
			return nm.executor.Execute(cmd, e)
		},
	})

	// Word motion commands - yank
	kt.Add(state.NormalMode, []string{"y", "w"}, keytree.KeyAction{
		Execute: func(e types.Editor) types.Editor {
			cmd := CreateYankCommand(motion.NewWordMotion(false), nm.register)
			return nm.executor.Execute(cmd, e)
		},
	})
	kt.Add(state.NormalMode, []string{"y", "W"}, keytree.KeyAction{
		Execute: func(e types.Editor) types.Editor {
			cmd := CreateYankCommand(motion.NewWordMotion(true), nm.register)
			return nm.executor.Execute(cmd, e)
		},
	})

	kt.Add(state.NormalMode, []string{"y", "e"}, keytree.KeyAction{
		Execute: func(e types.Editor) types.Editor {
			cmd := CreateYankCommand(motion.NewWordEndMotion(false), nm.register)
			return nm.executor.Execute(cmd, e)
		},
	})
	kt.Add(state.NormalMode, []string{"y", "E"}, keytree.KeyAction{
		Execute: func(e types.Editor) types.Editor {
			cmd := CreateYankCommand(motion.NewWordEndMotion(true), nm.register)
			return nm.executor.Execute(cmd, e)
		},
	})

	kt.Add(state.NormalMode, []string{"y", "b"}, keytree.KeyAction{
		Execute: func(e types.Editor) types.Editor {
			cmd := CreateYankCommand(motion.NewWordBackMotion(false), nm.register)
			return nm.executor.Execute(cmd, e)
		},
	})
	kt.Add(state.NormalMode, []string{"y", "B"}, keytree.KeyAction{
		Execute: func(e types.Editor) types.Editor {
			cmd := CreateYankCommand(motion.NewWordBackMotion(true), nm.register)
			return nm.executor.Execute(cmd, e)
		},
	})

	nm.keytree = kt
	return nm
}

func (nm *NormalMode) Handle(msg tea.KeyMsg, e types.Editor) (types.Editor, tea.Cmd) {

	if handled, e := nm.keytree.Handle(msg.String(), e); handled {
		e.HandleCursorMovement()
		return e, nil
	}

	switch msg.String() {
	case "ctrl+c":
		return e, tea.Quit
	case "v":
		e.SetMode(state.VisualMode)
	case "q":
		return e, tea.Quit
	case "h":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := CreateMotionCommand(motion.NewLeftMotion(), cursor)
			e = nm.executor.Execute(cmd, e)
		}
	case "l":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := CreateMotionCommand(motion.NewRightMotion(), cursor)
			e = nm.executor.Execute(cmd, e)
		}
	case "j":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := CreateMotionCommand(motion.NewDownMotion(), cursor)
			e = nm.executor.Execute(cmd, e)
		}
	case "k":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := CreateMotionCommand(motion.NewUpMotion(), cursor)
			e = nm.executor.Execute(cmd, e)
		}
	case "G":
		e.Buffer().ClearCursors()
		cursor, err := e.Buffer().GetPrimaryCursor()
		if err != nil {
			nm.logger.Println("Failed to get primary cursor:", err)
			return e, nil
		}
		cmd := CreateGoToEndOfBufferCommand(cursor)
		e = nm.executor.Execute(cmd, e)
	case "$":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			e = CreateGoToEndOfLineCommand(cursor).Execute(e)
		}
	case "^", "0":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			e = CreateGoToStartOfLineCommand(cursor).Execute(e)
		}
	case "w":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := CreateWordMotionCommand(false, cursor)
			e = nm.executor.Execute(cmd, e)
		}
	case "W":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := CreateWordMotionCommand(true, cursor)
			e = nm.executor.Execute(cmd, e)
		}
	case "e":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := CreateWordEndMotionCommand(false, cursor)
			e = nm.executor.Execute(cmd, e)
		}
	case "E":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := CreateMotionCommand(motion.NewWordEndMotion(true), cursor)
			e = nm.executor.Execute(cmd, e)
		}
	case "b":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := CreateMotionCommand(motion.NewWordBackMotion(false), cursor)
			e = nm.executor.Execute(cmd, e)
		}
	case "B":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := CreateMotionCommand(motion.NewWordBackMotion(true), cursor)
			e = nm.executor.Execute(cmd, e)
		}
	case "o":
		cmd := CreateNewLineCommand(false)
		e = nm.executor.Execute(cmd, e)
	case "O":
		cmd := CreateNewLineCommand(true)
		e = nm.executor.Execute(cmd, e)
	case "D":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := CreateDeleteToEndOfLineCommand(cursor)
			e = nm.executor.Execute(cmd, e)
		}
	case "C":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := CreateChangeToEndOfLineCommand(cursor)
			e = nm.executor.Execute(cmd, e)
		}

	case "a":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := CreateAppendCommand(false, cursor)
			e = nm.executor.Execute(cmd, e)
		}
	case "A":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := CreateAppendCommand(true, cursor)
			e = nm.executor.Execute(cmd, e)
		}

	case "i":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := CreateInsertCommand(false, cursor)
			e = nm.executor.Execute(cmd, e)
		}
	case "I":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			cmd := CreateInsertCommand(true, cursor)
			e = nm.executor.Execute(cmd, e)
		}
	case "p":
		e.Buffer().ClearCursors()
		cursor, err := e.Buffer().GetPrimaryCursor()
		if err != nil {
			nm.logger.Println("Failed to get primary cursor:", err)
			return e, nil
		}
		cmd := CreatePasteCommand(cursor, nm.register, false)
		e = nm.executor.Execute(cmd, e)
	case "P":
		e.Buffer().ClearCursors()
		cursor, err := e.Buffer().GetPrimaryCursor()
		if err != nil {
			nm.logger.Println("Failed to get primary cursor:", err)
			return e, nil
		}
		cmd := CreatePasteCommand(cursor, nm.register, true)
		e = nm.executor.Execute(cmd, e)
	case "u":
		e = nm.history.Undo(e)
	case "ctrl+r":
		e = nm.history.Redo(e)
	case "ctrl+d":
		e.Buffer().ClearCursors()
		cursor, err := e.Buffer().GetPrimaryCursor()
		if err != nil {
			return e, nil
		}
		cmd := CreateHalfPageDownCommand(cursor, e.Viewport())
		e = nm.executor.Execute(cmd, e)
	case "ctrl+u":
		e.Buffer().ClearCursors()
		cursor, err := e.Buffer().GetPrimaryCursor()
		if err != nil {
			return e, nil
		}
		cmd := CreateHalfPageUpCommand(cursor, e.Viewport())
		e = nm.executor.Execute(cmd, e)

	}

	e.HandleCursorMovement()

	return e, nil
}
