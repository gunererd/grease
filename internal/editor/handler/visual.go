package handler

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/command/motion"
	"github.com/gunererd/grease/internal/highlight"
	"github.com/gunererd/grease/internal/keytree"
	"github.com/gunererd/grease/internal/state"
	"github.com/gunererd/grease/internal/types"
)

type VisualMode struct {
	selectionStart   types.Position
	highlightID      int
	operationManager types.OperationManager
	keytree          *keytree.KeyTree
}

func NewVisualMode(kt *keytree.KeyTree, hm types.HistoryManager, om types.OperationManager) *VisualMode {

	kt.Add([]string{"g", "g"}, keytree.KeyAction{
		Execute: motion.CreateBasicMotionCommand(motion.NewStartOfBufferMotion()),
	})

	return &VisualMode{
		highlightID:      -1, // Invalid highlight ID
		operationManager: om,
		keytree:          kt,
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
		// Clear highlight when exiting visual mode
		if vm.highlightID != -1 {
			e.HighlightManager().Remove(vm.highlightID)
			vm.highlightID = -1
		}
		e.SetMode(state.NormalMode)
	case "h":
		model = motion.CreateBasicMotionCommand(motion.NewLeftMotion())(e)
	case "l":
		model = motion.CreateBasicMotionCommand(motion.NewRightMotion())(e)
	case "j":
		model = motion.CreateBasicMotionCommand(motion.NewDownMotion())(e)
	case "k":
		model = motion.CreateBasicMotionCommand(motion.NewUpMotion())(e)
	case "gg":
		model = motion.CreateBasicMotionCommand(motion.NewStartOfBufferMotion())(e)
	case "G":
		model = motion.CreateBasicMotionCommand(motion.NewEndOfBufferMotion())(e)
	case "i":
		// Clear highlight when entering insert mode
		if vm.highlightID != -1 {
			e.HighlightManager().Remove(vm.highlightID)
			vm.highlightID = -1
		}
		e.SetMode(state.InsertMode)
	case ":":
		// Clear highlight when entering command mode
		if vm.highlightID != -1 {
			e.HighlightManager().Remove(vm.highlightID)
			vm.highlightID = -1
		}
		e.SetMode(state.CommandMode)
	case "q":
		return e, tea.Quit
	case "w":
		model = CreateWordMotionCommand(false, nil)(e)
		model.HandleCursorMovement()
	case "W":
		model = CreateWordMotionCommand(true, nil)(e)
		model.HandleCursorMovement()
	case "e":
		model = CreateWordEndMotionCommand(false, nil)(e)
		model.HandleCursorMovement()
	case "E":
		model = CreateWordEndMotionCommand(true, nil)(e)
		model.HandleCursorMovement()
	case "b":
		model = CreateWordBackMotionCommand(false, nil)(e)
		model.HandleCursorMovement()
	case "B":
		model = CreateWordBackMotionCommand(true, nil)(e)
		model.HandleCursorMovement()
	case "$":
		model = motion.CreateBasicMotionCommand(motion.NewEndOfLineMotion())(e)
	case "^", "0":
		model = motion.CreateBasicMotionCommand(motion.NewStartOfLineMotion())(e)
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
