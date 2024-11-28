package editor

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/buffer"
	ioManager "github.com/gunererd/grease/internal/io"
	"github.com/gunererd/grease/internal/viewport"
)

type Mode int

const (
	NormalMode Mode = iota
	InsertMode
	CommandMode
)

type Model struct {
	buffer          *buffer.Buffer
	viewport        *viewport.Viewport
	mode            Mode
	width           int
	height          int
	io              *ioManager.Manager
	showLineNumbers bool
}

func New(io *ioManager.Manager) *Model {
	return &Model{
		buffer:          buffer.New(),
		viewport:        viewport.New(80, 24), // Default size
		mode:            NormalMode,
		io:              io,
		showLineNumbers: true,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.UpdateViewport(msg.Width, msg.Height)
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}
	return m, nil
}

func (m *Model) View() string {
	// Get visible content from viewport with line numbers
	content := m.viewport.View(m.buffer, m.showLineNumbers)

	// Add status line
	status := m.getStatusLine()

	// Combine content and status
	return strings.Join(content, "\n") + "\n" + status
}

func (m *Model) getStatusLine() string {
	cursor, _ := m.buffer.GetPrimaryCursor()
	mode := m.getModeString()
	pos := cursor.GetPosition()

	// Get relative cursor position for display
	x, y := m.viewport.GetRelativePosition(pos)

	return fmt.Sprintf("%s - Buffer[%d,%d] View[%d,%d] - %d%%",
		mode, pos.Line+1, pos.Column+1, y+1, x+1,
		int(float64(pos.Line+1)/float64(m.buffer.LineCount())*100))
}

func (m *Model) getModeString() string {
	switch m.mode {
	case NormalMode:
		return "NORMAL"
	case InsertMode:
		return "INSERT"
	case CommandMode:
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
	case NormalMode:
		return m.handleNormalMode(msg)
	case InsertMode:
		return m.handleInsertMode(msg)
	case CommandMode:
		return m.handleCommandMode(msg)
	}
	return m, nil
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
		m.mode = InsertMode
	case ":":
		m.mode = CommandMode
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
	switch msg.Type {
	case tea.KeyEsc:
		m.mode = NormalMode
	case tea.KeyRunes:
		m.buffer.Insert(string(msg.Runes))
	case tea.KeyEnter:
		m.buffer.Insert("\n")
	case tea.KeyBackspace:
		m.buffer.Delete(1)
	}
	m.handleCursorMovement()
	return m, nil
}

func (m *Model) handleCommandMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.mode = NormalMode
	case tea.KeyEnter:
		// TODO: Handle command execution
		m.mode = NormalMode
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
