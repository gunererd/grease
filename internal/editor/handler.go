package editor

import (
	tea "github.com/charmbracelet/bubbletea"
)

// ModeHandler handles keypresses for a specific mode
type ModeHandler interface {
	Handle(msg tea.KeyMsg, editor *Editor) (tea.Model, tea.Cmd)
}

// NormalModeHandler handles keypresses in normal mode
type NormalModeHandler struct{}

func NewNormalModeHandler() *NormalModeHandler {
	return &NormalModeHandler{}
}

func (h *NormalModeHandler) Handle(msg tea.KeyMsg, e *Editor) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return e, tea.Quit
	case "h":
		e.navigator.MoveCursorLeft()
		return e, nil

	case "l":
		e.navigator.MoveCursorRight()
		return e, nil

	case "-":
		if parent, err := e.navigator.GetParentDir(); err == nil {
			return e, e.navigator.ReadDirectory(parent)
		}
	case "enter":
		if entry, ok := e.navigator.GetCurrentEntry(); ok && entry.IsDir {
			return e, e.navigator.ReadDirectory(entry.Path)
		}
	case "j":
		if e.navigator.MoveCursor(1) {
			e.view.EnsureVisible(e.navigator.Cursor.Row, e.navigator.NumEntries())
		}
	case "k":
		if e.navigator.MoveCursor(-1) {
			e.view.EnsureVisible(e.navigator.Cursor.Row, e.navigator.NumEntries())
		}
	case "i":
		e.state.SetMode(InsertMode)
	case "v":
		e.state.SetMode(VisualMode)
		e.state.StartSelection(e.navigator.Cursor.Row, e.navigator.Cursor.Col)
	}

	return e, nil
}

// InsertModeHandler handles keypresses in insert mode
type InsertModeHandler struct{}

func NewInsertModeHandler() *InsertModeHandler {
	return &InsertModeHandler{}
}

func (h *InsertModeHandler) Handle(msg tea.KeyMsg, e *Editor) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		e.state.SetMode(NormalMode)
		return e, nil
	case "backspace":
		row := e.navigator.Cursor.Row
		col := e.navigator.Cursor.Col
		if col > 0 {
			e.navigator.Buffer.DeleteCharAtCursor(row, col)
			e.navigator.Cursor.MoveLeft()
		}
	case "up":
		e.navigator.MoveCursorUp()
	case "down":
		e.navigator.MoveCursorDown()
	case "left":
		e.navigator.MoveCursorLeft()
	case "right":
		e.navigator.MoveCursorRight()
	default:
		if len(msg.String()) == 1 {
			row := e.navigator.Cursor.Row
			col := e.navigator.Cursor.Col
			e.navigator.Buffer.InsertCharAtCursor(msg.String(), row, col)
			e.navigator.MoveCursorRight()
		}
	}
	return e, nil
}

// VisualModeHandler handles keypresses in visual mode
type VisualModeHandler struct{}

func NewVisualModeHandler() *VisualModeHandler {
	return &VisualModeHandler{}
}

func (h *VisualModeHandler) Handle(msg tea.KeyMsg, e *Editor) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		e.state.SetMode(NormalMode)
		e.state.ClearSelection()
		return e, nil
	case "j":
		if e.navigator.MoveCursor(1) {
			e.state.UpdateSelection(e.navigator.Cursor.Row, e.navigator.Cursor.Col)
			e.view.EnsureVisible(e.navigator.Cursor.Row, e.navigator.NumEntries())
		}
	case "k":
		if e.navigator.MoveCursor(-1) {
			e.state.UpdateSelection(e.navigator.Cursor.Row, e.navigator.Cursor.Col)
			e.view.EnsureVisible(e.navigator.Cursor.Row, e.navigator.NumEntries())
		}
	case "h":
		e.navigator.MoveCursorLeft()
		e.state.UpdateSelection(e.navigator.Cursor.Row, e.navigator.Cursor.Col)
	case "l":
		e.navigator.MoveCursorRight()
		e.state.UpdateSelection(e.navigator.Cursor.Row, e.navigator.Cursor.Col)
	}
	return e, nil
}
