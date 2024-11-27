package editor

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/navigator"
)

// readDirectoryMsg is sent when a directory has been read
type readDirectoryMsg struct {
	err error
}

// Editor represents a modal file system editor
type Editor struct {
	view      *View
	state     *State
	navigator *navigator.Navigator
	handlers  map[Mode]ModeHandler
}

// New creates a new Editor instance with default settings
func New() Editor {
	nav := navigator.New()
	state := NewState()
	view := NewView(nav, state)
	return Editor{
		view:      view,
		state:     state,
		navigator: nav,
		handlers: map[Mode]ModeHandler{
			NormalMode: NewNormalModeHandler(),
			InsertMode: NewInsertModeHandler(),
			VisualMode: NewVisualModeHandler(),
		},
	}
}

// Init initializes the editor
func (e Editor) Init() tea.Cmd {
	return e.navigator.ReadDirectory(".")
}

// Update handles incoming messages and updates editor state
func (e Editor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return e.handleKeypress(msg)

	case tea.WindowSizeMsg:
		e.view.UpdateSize(msg.Width, msg.Height)
		e.view.EnsureVisible(e.navigator.Cursor.Row, e.navigator.NumEntries())
		return e, nil

	case readDirectoryMsg:
		if msg.err != nil {
			return e, tea.Quit
		}
		e.navigator.SetCursor(0)
		e.view.ScrollOffset = 0
		e.state.ClearSelection()
		return e, nil
	}

	return e, nil
}

// View renders the editor state
func (e Editor) View() string {
	var b strings.Builder

	// Render entries
	numEntries := e.navigator.NumEntries()
	viewHeight := e.view.Height - 2 // Reserve space for status line

	for i := e.view.ScrollOffset; i < numEntries && i < e.view.ScrollOffset+viewHeight; i++ {
		entry, ok := e.navigator.GetEntry(i)
		if !ok {
			continue
		}
		b.WriteString(e.view.RenderEntry(entry, i))
		b.WriteString("\n")
	}

	// Add empty lines to push status line to bottom
	for i := numEntries; i < e.view.ScrollOffset+viewHeight; i++ {
		b.WriteString("\n")
	}

	// Render status line
	mode := e.state.GetMode()
	b.WriteString(e.view.RenderStatusLine(mode.String(), e.navigator.Buffer.CurrentDir, e.navigator.Cursor.Row))

	return b.String()
}

func (e Editor) handleKeypress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if handler, ok := e.handlers[e.state.GetMode()]; ok {
		return handler.Handle(msg, &e)
	}
	return e, nil
}
