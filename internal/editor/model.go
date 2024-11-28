package editor

import (
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/buffer"
	ioManager "github.com/gunererd/grease/internal/io"
	"github.com/gunererd/grease/internal/state"
	"github.com/gunererd/grease/internal/ui"
)

type Model struct {
	buffer          *buffer.Buffer
	viewport        *ui.Viewport
	mode            state.Mode
	width           int
	height          int
	io              *ioManager.Manager
	showLineNumbers bool
	cursorTimer     time.Time
	statusLine      *ui.StatusLine
}

func New(io *ioManager.Manager) *Model {
	return &Model{
		buffer:          buffer.New(),
		viewport:        ui.NewViewport(80, 24), // Default size
		mode:            state.NormalMode,
		io:              io,
		showLineNumbers: true,
		statusLine:      ui.NewStatusLine(),
	}
}

func (m *Model) Init() tea.Cmd {
	// Initialize cursor blink timer
	m.cursorTimer = time.Now()
	return tea.Tick(time.Millisecond*530, func(t time.Time) tea.Msg {
		return CursorBlinkMsg(t)
	})
}

type CursorBlinkMsg time.Time

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.UpdateViewport(msg.Width, msg.Height)
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case CursorBlinkMsg:
		m.viewport.ToggleCursor()
		return m, tea.Tick(time.Millisecond*530, func(t time.Time) tea.Msg {
			return CursorBlinkMsg(t)
		})
	}
	return m, nil
}

func (m *Model) View() string {
	// Get visible content from viewport
	content := m.viewport.View(m.buffer)

	// Add status line
	statusline := m.getStatusLine()

	// Combine content and status
	return strings.Join(content, "\n") + "\n" + statusline
}

func (m *Model) getStatusLine() string {
	cursor, _ := m.buffer.GetPrimaryCursor()
	mode := m.getModeString()
	x, y := m.viewport.GetRelativePosition(cursor.GetPosition())
	return m.statusLine.Render(mode, *cursor, m.buffer.LineCount(), x, y, m.width)
}

func (m *Model) getModeString() string {
	switch m.mode {
	case state.NormalMode:
		return "NORMAL"
	case state.InsertMode:
		return "INSERT"
	case state.CommandMode:
		return "COMMAND"
	default:
		return "UNKNOWN"
	}
}

func (m *Model) UpdateViewport(width, height int) {
	m.width = width
	m.height = height
	m.viewport.SetSize(width, height-1) // Reserve one line for status
}

func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case state.NormalMode:
		return m.handleNormalMode(msg)
	case state.InsertMode:
		return m.handleInsertMode(msg)
	case state.CommandMode:
		return m.handleCommandMode(msg)
	}
	return m, nil
}

func (m *Model) SetMode(mode state.Mode) {
	m.mode = mode
	m.viewport.SetMode(mode)
}

func (m *Model) handleNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	cursor, err := m.buffer.GetPrimaryCursor()
	if err != nil {
		return m, nil
	}

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "h":
		m.buffer.MoveCursor(cursor.GetID(), 0, -1)
	case "l":
		m.buffer.MoveCursor(cursor.GetID(), 0, 1)
	case "j":
		m.buffer.MoveCursor(cursor.GetID(), 1, 0)
	case "k":
		m.buffer.MoveCursor(cursor.GetID(), -1, 0)
	case "i":
		m.SetMode(state.InsertMode)
	case ":":
		m.SetMode(state.CommandMode)
	case "q":
		return m, tea.Quit
	case "z":
		// Center viewport on cursor
		cursor, _ := m.buffer.GetPrimaryCursor()
		m.viewport.CenterOn(cursor.GetPosition())
	case "zz":
		// Toggle line numbers
		m.showLineNumbers = !m.showLineNumbers
		m.UpdateViewport(m.width, m.height) // Refresh viewport size
	}
	m.handleCursorMovement()
	return m, nil
}

func (m *Model) handleInsertMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	cursor, err := m.buffer.GetPrimaryCursor()
	if err != nil {
		return m, nil
	}

	switch msg.Type {
	case tea.KeyEsc:
		m.SetMode(state.NormalMode)
	case tea.KeyRunes:
		m.buffer.Insert(string(msg.Runes))
	case tea.KeyEnter:
		m.buffer.Insert("\n")
	case tea.KeyBackspace:
		m.buffer.Delete(1)
	default:
		switch msg.String() {
		case "up":
			m.buffer.MoveCursor(cursor.GetID(), -1, 0)
		case "down":
			m.buffer.MoveCursor(cursor.GetID(), 1, 0)
		case "left":
			m.buffer.MoveCursor(cursor.GetID(), 0, -1)
		case "right":
			m.buffer.MoveCursor(cursor.GetID(), 0, 1)
		}
	}
	m.handleCursorMovement()
	return m, nil
}

func (m *Model) handleCommandMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.SetMode(state.NormalMode)
	case tea.KeyEnter:
		// TODO: Handle command execution
		m.SetMode(state.NormalMode)
	}
	return m, nil
}

func (m *Model) handleCursorMovement() {
	cursor, _ := m.buffer.GetPrimaryCursor()
	pos := cursor.GetPosition()
	m.viewport.SetCursor(pos) // This will also handle scrolling
}

// LoadFromStdin loads content from stdin into the buffer
func (m *Model) LoadFromStdin() error {
	return m.buffer.LoadFromReader(os.Stdin)
}

// AddCursor adds a new cursor at the specified position
func (m *Model) AddCursor(pos buffer.Position) error {
	_, err := m.buffer.AddCursor(pos, 50) // Regular cursors get normal priority
	return err
}

// RemoveCursor removes a cursor by its ID
func (m *Model) RemoveCursor(id int) {
	m.buffer.RemoveCursor(id)
}
