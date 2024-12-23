package handler

import (
	"log"

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
}

func NewNormalMode(kt *keytree.KeyTree, register *register.Register, history types.HistoryManager) *NormalMode {

	// Vim style Jump to beginning of buffer
	kt.Add(state.NormalMode, []string{"g", "g"}, keytree.KeyAction{
		Before: func(e types.Editor) types.Editor {
			e.Buffer().ClearCursors()
			return e
		},
		Execute: func(e types.Editor) types.Editor {
			cursor, err := e.Buffer().GetPrimaryCursor()
			if err != nil {
				log.Println("Failed to get primary cursor:", err)
				return e
			}
			return CreateGoToStartOfBufferCommand(cursor).Execute(e)
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
				log.Println("Failed to get primary cursor:", err)
				return e
			}
			return CreateDeleteLineCommand(cursor, history).Execute(e)
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
				log.Println("Failed to get primary cursor:", err)
				return e
			}
			return CreateChangeLineCommand(cursor, history).Execute(e)
		},
	})

	kt.Add(state.NormalMode, []string{"d", "w"}, keytree.KeyAction{
		Execute: CreateDeleteCommand(motion.NewWordMotion(false), history).Execute,
	})

	kt.Add(state.NormalMode, []string{"d", "W"}, keytree.KeyAction{
		Execute: CreateDeleteCommand(motion.NewWordMotion(true), history).Execute,
	})

	kt.Add(state.NormalMode, []string{"d", "e"}, keytree.KeyAction{
		Execute: CreateDeleteCommand(motion.NewWordEndMotion(false), history).Execute,
	})

	kt.Add(state.NormalMode, []string{"d", "E"}, keytree.KeyAction{
		Execute: CreateDeleteCommand(motion.NewWordEndMotion(true), history).Execute,
	})

	kt.Add(state.NormalMode, []string{"d", "b"}, keytree.KeyAction{
		Execute: CreateDeleteCommand(motion.NewWordBackMotion(false), history).Execute,
	})

	kt.Add(state.NormalMode, []string{"d", "B"}, keytree.KeyAction{
		Execute: CreateDeleteCommand(motion.NewWordBackMotion(true), history).Execute,
	})

	kt.Add(state.NormalMode, []string{"c", "w"}, keytree.KeyAction{
		Execute: CreateChangeCommand(motion.NewWordMotion(false), history).Execute,
	})

	kt.Add(state.NormalMode, []string{"c", "W"}, keytree.KeyAction{
		Execute: CreateChangeCommand(motion.NewWordMotion(true), history).Execute,
	})

	kt.Add(state.NormalMode, []string{"c", "e"}, keytree.KeyAction{
		Execute: CreateChangeCommand(motion.NewWordEndMotion(false), history).Execute,
	})

	kt.Add(state.NormalMode, []string{"c", "E"}, keytree.KeyAction{
		Execute: CreateChangeCommand(motion.NewWordEndMotion(true), history).Execute,
	})

	kt.Add(state.NormalMode, []string{"c", "b"}, keytree.KeyAction{
		Execute: CreateChangeCommand(motion.NewWordBackMotion(false), history).Execute,
	})

	kt.Add(state.NormalMode, []string{"c", "B"}, keytree.KeyAction{
		Execute: CreateChangeCommand(motion.NewWordBackMotion(true), history).Execute,
	})

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
		history:  history,
	}
}

func (nm *NormalMode) Handle(msg tea.KeyMsg, e types.Editor) (types.Editor, tea.Cmd) {

	if handled, e := nm.keytree.Handle(msg.String(), e); handled {
		e.HandleCursorMovement()
		return e, nil
	}

	switch msg.String() {
	case "ctrl+c":
		return e, tea.Quit
	// case "C":
	// 	log.Println("shift+c")
	// 	buf := e.Buffer()
	// 	cursor, err := buf.AddCursor()
	// 	if err != nil {
	// 		return e, nil
	// 	}
	// 	cmd := motion.CreateMotionCommand(motion.NewDownMotion(), cursor)
	// 	cmd(e)
	case "h":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			e = CreateMotionCommand(motion.NewLeftMotion(), cursor).Execute(e)
		}
	case "l":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			e = CreateMotionCommand(motion.NewRightMotion(), cursor).Execute(e)
		}
	case "j":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			e = CreateMotionCommand(motion.NewDownMotion(), cursor).Execute(e)
		}
	case "k":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			e = CreateMotionCommand(motion.NewUpMotion(), cursor).Execute(e)
		}
	case "v":
		e.SetMode(state.VisualMode)
	case "q":
		return e, tea.Quit
	case "G":
		e.Buffer().ClearCursors()
		cursor, err := e.Buffer().GetPrimaryCursor()
		if err != nil {
			log.Println("Failed to get primary cursor:", err)
			return e, nil
		}
		e = CreateGoToEndOfBufferCommand(cursor).Execute(e)
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
			e = CreateWordMotionCommand(false, cursor).Execute(e)
		}
	case "W":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			e = CreateWordMotionCommand(true, cursor).Execute(e)
		}
	case "e":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			e = CreateWordEndMotionCommand(false, cursor).Execute(e)
		}
	case "E":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			e = CreateMotionCommand(motion.NewWordEndMotion(true), cursor).Execute(e)
		}
	case "b":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			e = CreateMotionCommand(motion.NewWordBackMotion(false), cursor).Execute(e)
		}
	case "B":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			e = CreateMotionCommand(motion.NewWordBackMotion(true), cursor).Execute(e)
		}
	case "o":
		e = CreateNewLineCommand(false, nm.history).Execute(e)
	case "O":
		e = CreateNewLineCommand(true, nm.history).Execute(e)
	case "D":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			e = CreateDeleteToEndOfLineCommand(cursor, nm.history).Execute(e)
		}
	case "C":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			e = CreateChangeToEndOfLineCommand(cursor, nm.history).Execute(e)
		}

	case "a":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			e = CreateAppendCommand(false, cursor, nm.history).Execute(e)
		}
	case "A":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			e = CreateAppendCommand(true, cursor, nm.history).Execute(e)
		}

	case "i":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			e = CreateInsertCommand(false, cursor, nm.history).Execute(e)
		}
	case "I":
		cursors := e.Buffer().GetCursors()
		for _, cursor := range cursors {
			e = CreateInsertCommand(true, cursor, nm.history).Execute(e)
		}
	case "p":
		e.Buffer().ClearCursors()
		cursor, err := e.Buffer().GetPrimaryCursor()
		if err != nil {
			log.Println("Failed to get primary cursor:", err)
			return e, nil
		}
		e = CreatePasteCommand(cursor, nm.register, false, nm.history).Execute(e)
	case "P":
		e.Buffer().ClearCursors()
		cursor, err := e.Buffer().GetPrimaryCursor()
		if err != nil {
			log.Println("Failed to get primary cursor:", err)
			return e, nil
		}
		e = CreatePasteCommand(cursor, nm.register, true, nm.history).Execute(e)
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
		e = CreateHalfPageDownCommand(cursor, e.Viewport()).Execute(e)
	case "ctrl+u":
		e.Buffer().ClearCursors()
		cursor, err := e.Buffer().GetPrimaryCursor()
		if err != nil {
			return e, nil
		}
		e = CreateHalfPageUpCommand(cursor, e.Viewport()).Execute(e)

	}

	e.HandleCursorMovement()

	return e, nil
}
