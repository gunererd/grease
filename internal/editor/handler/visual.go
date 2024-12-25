package handler

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/editor/command/motion"
	"github.com/gunererd/grease/internal/editor/highlight"
	"github.com/gunererd/grease/internal/editor/keytree"
	"github.com/gunererd/grease/internal/editor/register"
	"github.com/gunererd/grease/internal/editor/state"
	"github.com/gunererd/grease/internal/editor/types"
)

type VisualMode struct {
	selectionStart types.Position
	highlightID    int
	keytree        *keytree.KeyTree
	hlm            types.HighlightManager
	logger         types.Logger
}

func NewVisualMode(
	kt *keytree.KeyTree,
	register *register.Register,
	hlm types.HighlightManager,
	logger types.Logger,
) *VisualMode {
	vm := &VisualMode{
		hlm:         hlm,
		logger:      logger,
		highlightID: -1,
	}

	kt.Add(state.VisualMode, []string{"g", "g"}, keytree.KeyAction{
		Before: func(e types.Editor) types.Editor {
			e.Buffer().ClearCursors()
			return e
		},
		Execute: func(e types.Editor) types.Editor {
			cursor, err := e.Buffer().GetPrimaryCursor()
			if err != nil {
				vm.logger.Println("Failed to get primary cursor:", err)
				return e
			}
			return CreateGoToStartOfBufferCommand(cursor).Execute(e)
		},
	})

	vm.keytree = kt

	return vm
}

func (vm *VisualMode) Handle(msg tea.KeyMsg, e types.Editor) (types.Editor, tea.Cmd) {

	e.Buffer().ClearCursors()
	cursor, err := e.Buffer().GetPrimaryCursor()
	if err != nil {
		return e, nil
	}

	// Initialize visual selection if not already done
	if vm.highlightID == -1 {
		vm.selectionStart = cursor.GetPosition()
		vm.highlightID = e.HighlightManager().Add(
			highlight.CreateVisualHighlight(vm.selectionStart, vm.selectionStart),
		)
		if vm.highlightID == -1 {
			vm.logger.Printf("Failed to create visual highlight at position %v", vm.selectionStart)
		}
	}

	// Handle key sequences
	if handled, e := vm.keytree.Handle(msg.String(), e); handled {
		vm.updateHighlight(cursor, e.HighlightManager())
		e.HandleCursorMovement()
		return e, nil
	}

	var cmd tea.Cmd
	switch msg.String() {
	case "esc":
		vm.cleanup(e)
		e.SetMode(state.NormalMode)
	case "h":
		e = CreateMotionCommand(motion.NewLeftMotion(), cursor).Execute(e)
	case "l":
		e = CreateMotionCommand(motion.NewRightMotion(), cursor).Execute(e)
	case "j":
		e = CreateMotionCommand(motion.NewDownMotion(), cursor).Execute(e)
	case "k":
		e = CreateMotionCommand(motion.NewUpMotion(), cursor).Execute(e)
	case "G":
		e = CreateMotionCommand(motion.NewEndOfBufferMotion(), cursor).Execute(e)
	case "i":
		vm.cleanup(e)
		e.SetMode(state.InsertMode)
	case "q":
		return e, tea.Quit
	case "$":
		e = CreateMotionCommand(motion.NewEndOfLineMotion(), cursor).Execute(e)
	case "^", "0":
		e = CreateMotionCommand(motion.NewStartOfLineMotion(), cursor).Execute(e)
	case "w":

		CreateWordMotionCommand(false, cursor).Execute(e)
	case "W":
		CreateWordMotionCommand(true, cursor).Execute(e)
	case "e":
		CreateWordEndMotionCommand(false, cursor).Execute(e)
	case "E":
		CreateWordEndMotionCommand(true, cursor).Execute(e)
	case "b":
		CreateWordBackMotionCommand(false, cursor).Execute(e)
	case "B":
		CreateWordBackMotionCommand(true, cursor).Execute(e)
		// case "y":
		// 	e = vm.operationManager.Execute(types.OpYank, e, vm.selectionStart, cursor.GetPosition())
		// 	vm.cleanup(e)
		// case "d":
		// 	e = vm.operationManager.Execute(types.OpDelete, e, vm.selectionStart, cursor.GetPosition())
		// 	vm.cleanup(e)
		// 	e.Buffer().MoveCursor(cursor.ID(), vm.selectionStart.Line(), vm.selectionStart.Column())
		// 	e.HandleCursorMovement()
		// case "c":
		// 	e = vm.operationManager.Execute(types.OpChange, e, vm.selectionStart, cursor.GetPosition())
		// 	vm.cleanup(e)
		// 	e.Buffer().MoveCursor(cursor.ID(), vm.selectionStart.Line(), vm.selectionStart.Column())
		// 	e.HandleCursorMovement()
	}

	e.HandleCursorMovement()
	vm.updateHighlight(cursor, e.HighlightManager())

	return e, cmd
}

func (vm *VisualMode) updateHighlight(cursor types.Cursor, hm types.HighlightManager) {
	if vm.highlightID != -1 {
		currentPos := cursor.GetPosition()
		var iposition types.Position = currentPos
		if !hm.Update(
			vm.highlightID,
			highlight.CreateVisualHighlight(vm.selectionStart, iposition),
		) {
			vm.logger.Printf("Failed to update visual highlight %d from %v to %v",
				vm.highlightID, vm.selectionStart, currentPos)
			vm.highlightID = -1
		}
	}
}

func (vm *VisualMode) cleanup(e types.Editor) {
	if vm.highlightID != -1 {
		e.HighlightManager().Remove(vm.highlightID)
		vm.highlightID = -1
	}
	e.SetMode(state.NormalMode)
}
