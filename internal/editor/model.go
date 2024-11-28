package editor

import (
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/buffer"
	"github.com/gunererd/grease/internal/editor/handler"
	ioManager "github.com/gunererd/grease/internal/io"
	"github.com/gunererd/grease/internal/state"
	"github.com/gunererd/grease/internal/ui"
)

type Editor struct {
	buffer          *buffer.Buffer
	viewport        *ui.Viewport
	mode            state.Mode
	width           int
	height          int
	io              *ioManager.Manager
	showLineNumbers bool
	cursorTimer     time.Time
	statusLine      *ui.StatusLine
	handlers        map[state.Mode]handler.ModeHandler
}

func New(io *ioManager.Manager) *Editor {
	e := &Editor{
		buffer:          buffer.New(),
		viewport:        ui.NewViewport(80, 24), // Default size
		mode:            state.NormalMode,
		io:              io,
		showLineNumbers: true,
		statusLine:      ui.NewStatusLine(),
		handlers: map[state.Mode]handler.ModeHandler{
			state.NormalMode:  handler.NewNormalMode(),
			state.InsertMode:  handler.NewInsertMode(),
			state.VisualMode:  handler.NewVisualMode(),
			state.CommandMode: handler.NewCommandMode(),
		},
	}
	return e
}

func (e *Editor) Buffer() *buffer.Buffer {
	return e.buffer
}

func (e *Editor) Height() int {
	return e.height
}

func (e *Editor) Width() int {
	return e.width
}

func (e *Editor) Viewport() *ui.Viewport {
	return e.viewport
}

func (e *Editor) Init() tea.Cmd {
	// Initialize cursor blink timer
	e.cursorTimer = time.Now()
	return tea.Tick(time.Millisecond*530, func(t time.Time) tea.Msg {
		return CursorBlinkMsg(t)
	})
}

type CursorBlinkMsg time.Time

func (e *Editor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		e.UpdateViewport(msg.Width, msg.Height)
	case tea.KeyMsg:
		return e.handleKeyPress(msg)
	case CursorBlinkMsg:
		e.Viewport().ToggleCursor()
		return e, tea.Tick(time.Millisecond*530, func(t time.Time) tea.Msg {
			return CursorBlinkMsg(t)
		})
	}
	return e, nil
}

func (e *Editor) View() string {
	// Get visible content from viewport
	content := e.Viewport().View(e.Buffer())

	// Add status line
	statusline := e.getStatusLine()

	// Combine content and status
	return strings.Join(content, "\n") + "\n" + statusline
}

func (e *Editor) getStatusLine() string {
	cursor, _ := e.Buffer().GetPrimaryCursor()
	mode := e.getModeString()
	x, y := e.Viewport().GetRelativePosition(cursor.GetPosition())
	return e.statusLine.Render(mode, *cursor, e.Buffer().LineCount(), x, y, e.Width())
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

func (e *Editor) HandleCursorMovement() {
	cursor, _ := e.Buffer().GetPrimaryCursor()
	pos := cursor.GetPosition()
	e.Viewport().SetCursor(pos) // This will also handle scrolling
}

// LoadFromStdin loads content from stdin into the buffer
func (e *Editor) LoadFromStdin() error {
	return e.Buffer().LoadFromReader(os.Stdin)
}

// AddCursor adds a new cursor at the specified position
func (e *Editor) AddCursor(pos buffer.Position) error {
	_, err := e.Buffer().AddCursor(pos, 50) // Regular cursors get normal priority
	return err
}

// RemoveCursor removes a cursor by its ID
func (e *Editor) RemoveCursor(id int) {
	e.Buffer().RemoveCursor(id)
}
