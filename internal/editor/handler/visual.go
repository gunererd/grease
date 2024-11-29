package handler

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/highlight"
	"github.com/gunererd/grease/internal/state"
	"github.com/gunererd/grease/internal/types"
)

type VisualMode struct {
	selectionStart types.Position
	highlightID    int
}

func NewVisualMode() *VisualMode {
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
	}

	switch msg.String() {
	case "esc":
		// Clear highlight when exiting visual mode
		e.HighlightManager().Remove(h.highlightID)
		h.highlightID = -1
		e.SetMode(state.NormalMode)
	case "h":
		e.Buffer().MoveCursor(cursor.GetID(), 0, -1)
	case "l":
		e.Buffer().MoveCursor(cursor.GetID(), 0, 1)
	case "j":
		e.Buffer().MoveCursor(cursor.GetID(), 1, 0)
	case "k":
		e.Buffer().MoveCursor(cursor.GetID(), -1, 0)
	case "i":
		// Clear highlight when entering insert mode
		e.HighlightManager().Remove(h.highlightID)
		h.highlightID = -1
		e.SetMode(state.InsertMode)
	case ":":
		// Clear highlight when entering command mode
		e.HighlightManager().Remove(h.highlightID)
		h.highlightID = -1
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
		e.HighlightManager().Update(
			h.highlightID,
			highlight.CreateVisualHighlight(h.selectionStart, currentPos),
		)
	}

	e.HandleCursorMovement()
	return e, nil
}
