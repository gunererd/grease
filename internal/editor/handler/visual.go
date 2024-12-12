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
}

func NewVisualMode(kt *keytree.KeyTree, hm types.HistoryManager, om types.OperationManager) *VisualMode {
	return &VisualMode{
		highlightID:      -1, // Invalid highlight ID
		operationManager: om,
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
		return motion.CreateBasicMotionCommand(motion.NewLeftMotion())(e), nil
	case "l":
		return motion.CreateBasicMotionCommand(motion.NewRightMotion())(e), nil
	case "j":
		return motion.CreateBasicMotionCommand(motion.NewDownMotion())(e), nil
	case "k":
		return motion.CreateBasicMotionCommand(motion.NewUpMotion())(e), nil
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
		// Vim style jump to end of line
		cursor, _ := e.Buffer().GetPrimaryCursor()
		line := cursor.GetPosition().Line()
		lineLength, _ := e.Buffer().LineLen(line)
		e.Buffer().MoveCursor(cursor.ID(), line, lineLength-1)
		e.HandleCursorMovement()
	case "^", "0":
		// Vim style jump to beginning of line
		cursor, _ := e.Buffer().GetPrimaryCursor()
		line := cursor.GetPosition().Line()
		e.Buffer().MoveCursor(cursor.ID(), line, 0)
		e.HandleCursorMovement()
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

	// Update highlight to match current cursor position
	if vm.highlightID != -1 {
		currentPos := cursor.GetPosition()
		var iposition types.Position = currentPos
		if !e.HighlightManager().Update(
			vm.highlightID,
			highlight.CreateVisualHighlight(vm.selectionStart, iposition),
		) {
			log.Printf("Failed to update visual highlight %d from %v to %v",
				vm.highlightID, vm.selectionStart, currentPos)
			// Reset highlight state since update failed
			vm.highlightID = -1
		}
	}

	return model, cmd
}

func (vm *VisualMode) cleanup(model types.Editor) {
	if vm.highlightID != -1 {
		model.HighlightManager().Remove(vm.highlightID)
		vm.highlightID = -1
	}
	model.SetMode(state.NormalMode)
}
