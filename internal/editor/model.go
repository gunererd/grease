package editor

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/editor/handler"
	ioManager "github.com/gunererd/grease/internal/io"
	"github.com/gunererd/grease/internal/keytree"
	"github.com/gunererd/grease/internal/register"
	"github.com/gunererd/grease/internal/state"
	"github.com/gunererd/grease/internal/types"
)

type Editor struct {
	buffer           types.Buffer
	viewport         types.Viewport
	mode             state.Mode
	width            int
	height           int
	io               *ioManager.Manager
	showLineNumbers  bool
	statusLine       types.StatusLine
	handlers         map[state.Mode]types.ModeHandler
	highlightManager types.HighlightManager
	historyManager   types.HistoryManager
}

func New(
	io *ioManager.Manager,
	b types.Buffer,
	sl types.StatusLine,
	wp types.Viewport,
	hlm types.HighlightManager,
	kt *keytree.KeyTree,
	// hm types.HistoryManager,
	// om types.OperationManager,
	register *register.Register,
) *Editor {
	e := &Editor{
		buffer:          b,
		viewport:        wp, // Default size
		mode:            state.NormalMode,
		io:              io,
		showLineNumbers: true,
		statusLine:      sl,
		handlers: map[state.Mode]types.ModeHandler{
			state.NormalMode: handler.NewNormalMode(kt, register),
			state.InsertMode: handler.NewInsertMode(),
			state.VisualMode: handler.NewVisualMode(kt, register),
		},
		// highlightManager: hlm,
		// historyManager:   hm,
	}
	return e
}

func (e *Editor) Buffer() types.Buffer {
	return e.buffer
}

func (e *Editor) Height() int {
	return e.height
}

func (e *Editor) Width() int {
	return e.width
}

func (e *Editor) Viewport() types.Viewport {
	return e.viewport
}

func (e *Editor) HighlightManager() types.HighlightManager {
	return e.highlightManager
}

func (e *Editor) HistoryManager() types.HistoryManager {
	return e.historyManager
}

func (e *Editor) Init() tea.Cmd {
	return nil
}

type CursorBlinkMsg time.Time

func (e *Editor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		e.UpdateViewport(msg.Width, msg.Height)
	case tea.KeyMsg:
		return e.handleKeyPress(msg)
	}
	return e, nil
}

func (e *Editor) View() string {
	// Get visible content from viewport
	content := e.Viewport().Render(e.Buffer())

	// Add status line
	statusline := e.getStatusLine()

	// Combine content and status
	return strings.Join(content, "\n") + "\n" + statusline
}

func (e *Editor) getStatusLine() string {
	cursor, _ := e.Buffer().GetPrimaryCursor()
	mode := e.getModeString()
	x, y := e.Viewport().BufferToViewportPosition(cursor.GetPosition())
	return e.statusLine.Render(mode, cursor, e.Buffer().LineCount(), x, y, e.Width())
}

func (e *Editor) getModeString() string {
	switch e.mode {
	case state.NormalMode:
		return "NORMAL"
	case state.InsertMode:
		return "INSERT"
	case state.VisualMode:
		return "VISUAL"
	case state.CommandMode:
		return "COMMAND"
	default:
		return "UNKNOWN"
	}
}

func (e *Editor) UpdateViewport(width, height int) {
	e.width = width
	e.height = height
	e.Viewport().SetSize(width, height-1) // Reserve one line for status
}

func (e *Editor) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if handler, ok := e.handlers[e.mode]; ok {
		return handler.Handle(msg, e)
	}
	return e, nil
}

func (e *Editor) SetMode(mode state.Mode) {
	e.mode = mode
	e.Viewport().SetMode(mode)
}

// func (e *Editor) HandleCursorMovement() {
// 	cursor, _ := e.Buffer().GetPrimaryCursor()
// 	pos := cursor.GetPosition()
// 	e.Viewport().SetCursor(pos) // This will also handle scrolling
// }

func (e *Editor) HandleCursorMovement() {
	for _, cursor := range e.Buffer().GetCursors() {
		pos := cursor.GetPosition()
		e.Viewport().SetCursor(pos) // This will also handle scrolling
	}
}

// LoadFromStdin loads content from stdin into the buffer
func (e *Editor) LoadFromStdin() error {
	content, err := e.io.LoadContent()
	if err != nil {
		return fmt.Errorf("error reading from stdin: %w", err)
	}
	return e.buffer.LoadFromReader(strings.NewReader(string(content)))
}

// LoadFromFile loads content from a file into the editor's buffer
func (e *Editor) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	if err := e.io.SetSource(ioManager.NewFileSource(file, filename)); err != nil {
		return fmt.Errorf("error setting source: %w", err)
	}

	content, err := e.io.LoadContent()
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return e.buffer.LoadFromReader(strings.NewReader(string(content)))
}

// AddCursor adds a new cursor at the specified position
func (e *Editor) AddCursor(pos types.Position) error {
	_, err := e.Buffer().AddCursor()
	return err
}

// RemoveCursor removes a cursor by its ID
func (e *Editor) RemoveCursor(id int) {
	e.Buffer().RemoveCursor(id)
}
