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

func NewVisualMode(kt *keytree.KeyTree) *VisualMode {
	return &VisualMode{
		highlightID: -1, // Invalid highlight ID
	}
}

func (h *VisualMode) Handle(msg tea.KeyMsg, e types.Editor) (tea.Model, tea.Cmd) {
	cursor, err := e.Buffer().GetPrimaryCursor()
	if err != nil {
		return e, nil
	}

	// Initialize visual selection if not already done
	if h.highlightID == -1 {
		h.selectionStart = cursor.GetPosition()
		h.highlightID = e.HighlightManager().Add(
			highlight.CreateVisualHighlight(h.selectionStart, h.selectionStart),
		)
		if h.highlightID == -1 {
			log.Printf("Failed to create visual highlight at position %v", h.selectionStart)
		}
	}

	switch msg.String() {
	case "esc":
		// Clear highlight when exiting visual mode
		if h.highlightID != -1 {
			e.HighlightManager().Remove(h.highlightID)
			h.highlightID = -1
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
		if h.highlightID != -1 {
			e.HighlightManager().Remove(h.highlightID)
			h.highlightID = -1
		}
		e.SetMode(state.InsertMode)
	case ":":
		// Clear highlight when entering command mode
		if h.highlightID != -1 {
			e.HighlightManager().Remove(h.highlightID)
			h.highlightID = -1
		}
		e.SetMode(state.CommandMode)
	case "q":
		return e, tea.Quit
	case "z":
		// Center viewport on cursor
		cursor, _ := e.Buffer().GetPrimaryCursor()
		e.Viewport().CenterOn(cursor.GetPosition())
	}

	// Update highlight to match current cursor position
	if h.highlightID != -1 {
		currentPos := cursor.GetPosition()
		var iposition types.Position = currentPos
		if !e.HighlightManager().Update(
			h.highlightID,
			highlight.CreateVisualHighlight(h.selectionStart, iposition),
		) {
			log.Printf("Failed to update visual highlight %d from %v to %v",
				h.highlightID, h.selectionStart, currentPos)
			// Reset highlight state since update failed
			h.highlightID = -1
		}
	}

	e.HandleCursorMovement()
	return e, nil
}
