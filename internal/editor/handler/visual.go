package handler

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/highlight"
	"github.com/gunererd/grease/internal/keytree"
	"github.com/gunererd/grease/internal/state"
	"github.com/gunererd/grease/internal/types"
)

type VisualMode struct {
	selectionStart types.Position
	highlightID    int
}

func NewVisualMode(kt *keytree.KeyTree, hm types.HistoryManager) *VisualMode {
	return &VisualMode{
		highlightID: -1, // Invalid highlight ID
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
		model, cmd = CreateWordMotionCommand(false, nil)(e)
		model.HandleCursorMovement()
	case "W":
		model, cmd = CreateWordMotionCommand(true, nil)(e)
		model.HandleCursorMovement()
	case "e":
		model, cmd = CreateWordEndMotionCommand(false, nil)(e)
		model.HandleCursorMovement()
	case "E":
		model, cmd = CreateWordEndMotionCommand(true, nil)(e)
		model.HandleCursorMovement()
	case "b":
		model, cmd = CreateWordBackMotionCommand(false, nil)(e)
		model.HandleCursorMovement()
	case "B":
		model, cmd = CreateWordBackMotionCommand(true, nil)(e)
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
		// Yank the selected text
		yankOp := NewYankOperation()
		model, cmd = yankOp.Execute(e, vm.selectionStart, cursor.GetPosition())

		// Clear highlight and exit visual mode
		if vm.highlightID != -1 {
			e.HighlightManager().Remove(vm.highlightID)
			vm.highlightID = -1
		}
		model.SetMode(state.NormalMode)
	case "d":
		// Delete the selected text
		deleteOp := NewHistoryAwareOperation(NewDeleteOperation(), e.HistoryManager())
		model, cmd = deleteOp.Execute(e, vm.selectionStart, cursor.GetPosition())

		// Clear highlight and exit visual mode
		if vm.highlightID != -1 {
			e.HighlightManager().Remove(vm.highlightID)
			vm.highlightID = -1
		}
		model.SetMode(state.NormalMode)
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
