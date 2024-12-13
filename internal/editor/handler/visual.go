package handler

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/command/motion"
	"github.com/gunererd/grease/internal/highlight"
	"github.com/gunererd/grease/internal/keytree"
	"github.com/gunererd/grease/internal/register"
	"github.com/gunererd/grease/internal/state"
	"github.com/gunererd/grease/internal/types"
)

type VisualMode struct {
	selectionStart   types.Position
	highlightID      int
	operationManager types.OperationManager
	keytree          *keytree.KeyTree
}

func NewVisualMode(kt *keytree.KeyTree, register *register.Register) *VisualMode {

	kt.Add([]string{"g", "g"}, keytree.KeyAction{
		Before: func(e types.Editor) types.Editor {
			e.Buffer().ClearCursors()
			return e
		},
		Execute: motion.CreateBasicMotionCommand(motion.NewStartOfBufferMotion(), -1),
	})

	return &VisualMode{
		highlightID: -1, // Invalid highlight ID
		// operationManager: om,
		keytree: kt,
	}
}

func (vm *VisualMode) Handle(msg tea.KeyMsg, e types.Editor) (types.Editor, tea.Cmd) {

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
			log.Printf("Failed to create visual highlight at position %v", vm.selectionStart)
		}
	}

	// Handle key sequences
	if handled, model := vm.keytree.Handle(msg.String(), e); handled {
		vm.updateHighlight(cursor, e.HighlightManager())
		e.HandleCursorMovement()
		return model, nil
	}

	var model types.Editor = e
	var cmd tea.Cmd
	switch msg.String() {
	case "esc":
		vm.cleanup(e)
		e.SetMode(state.NormalMode)
	case "h":
		model = motion.CreateBasicMotionCommand(motion.NewLeftMotion(), cursor.ID())(e)
	case "l":
		model = motion.CreateBasicMotionCommand(motion.NewRightMotion(), cursor.ID())(e)
	case "j":
		model = motion.CreateBasicMotionCommand(motion.NewDownMotion(), cursor.ID())(e)
	case "k":
		model = motion.CreateBasicMotionCommand(motion.NewUpMotion(), cursor.ID())(e)
	case "G":
		model = motion.CreateBasicMotionCommand(motion.NewEndOfBufferMotion(), cursor.ID())(e)
	case "i":
		vm.cleanup(e)
		e.SetMode(state.InsertMode)
	case "q":
		return e, tea.Quit
	case "$":
		motion.CreateBasicMotionCommand(motion.NewEndOfLineMotion(), cursor.ID())
	case "^", "0":
		motion.CreateBasicMotionCommand(motion.NewStartOfLineMotion(), cursor.ID())
	case "w":

		CreateWordMotionCommand(false, cursor.ID()).Execute(e)
	case "W":
		CreateWordMotionCommand(true, cursor.ID()).Execute(e)
	case "e":
		CreateWordEndMotionCommand(false, cursor.ID()).Execute(e)
	case "E":
		CreateWordEndMotionCommand(true, cursor.ID()).Execute(e)
	case "b":
		CreateWordBackMotionCommand(false, cursor.ID()).Execute(e)
	case "B":
		CreateWordBackMotionCommand(true, cursor.ID()).Execute(e)
	case "y":
		model = vm.operationManager.Execute(types.OpYank, e, vm.selectionStart, cursor.GetPosition())
		vm.cleanup(model)
	case "d":
		model = vm.operationManager.Execute(types.OpDelete, e, vm.selectionStart, cursor.GetPosition())
		vm.cleanup(model)
		model.Buffer().MoveCursor(cursor.ID(), vm.selectionStart.Line(), vm.selectionStart.Column())
		model.HandleCursorMovement()
	case "c":
		model = vm.operationManager.Execute(types.OpChange, e, vm.selectionStart, cursor.GetPosition())
		vm.cleanup(model)
		model.Buffer().MoveCursor(cursor.ID(), vm.selectionStart.Line(), vm.selectionStart.Column())
		model.HandleCursorMovement()
	}

	vm.updateHighlight(cursor, e.HighlightManager())

	return model, cmd
}

func (vm *VisualMode) updateHighlight(cursor types.Cursor, hm types.HighlightManager) {
	if vm.highlightID != -1 {
		currentPos := cursor.GetPosition()
		var iposition types.Position = currentPos
		if !hm.Update(
			vm.highlightID,
			highlight.CreateVisualHighlight(vm.selectionStart, iposition),
		) {
			log.Printf("Failed to update visual highlight %d from %v to %v",
				vm.highlightID, vm.selectionStart, currentPos)
			vm.highlightID = -1
		}
	}
}

func (vm *VisualMode) cleanup(model types.Editor) {
	if vm.highlightID != -1 {
		model.HighlightManager().Remove(vm.highlightID)
		vm.highlightID = -1
	}
	model.SetMode(state.NormalMode)
}
