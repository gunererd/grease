package model

import (
	"strings"

	"github.com/gunererd/grease/internal/buffer"
	"github.com/gunererd/grease/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

// readDirectoryMsg is sent when a directory has been read
type readDirectoryMsg struct {
	err error
}

type Model struct {
	Mode Mode

	Buffer *buffer.Buffer

	Cursor Cursor

	// Window dimensions
	Width  int
	Height int

	// Viewport for scrolling
	ScrollOffset int
	ViewHeight   int // Available height for content

}

func New() Model {
	return Model{
		Mode:         NormalMode,
		Buffer:       buffer.NewBuffer(),
		Cursor:       NewCursor(),
		ViewHeight:   10, // Default height, will be updated when window size is received
		ScrollOffset: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return m.readDirectory(".")
}

func (m Model) readDirectory(path string) tea.Cmd {
	return func() tea.Msg {
		err := m.Buffer.ReadDirectory(path)
		return readDirectoryMsg{err: err}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeypress(msg)

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		// Reserve 1 line for statusline
		m.ViewHeight = m.Height - 1
		m.ensureCursorVisible()
		return m, nil

	case readDirectoryMsg:
		if msg.err != nil {
			// Handle error (we'll improve this later)
			return m, tea.Quit
		}

		// Reset scroll when changing directories
		m.ScrollOffset = 0
		m.Cursor.Row = 0
		return m, nil
	}

	return m, nil
}

// ensureCursorVisible adjusts scroll offset to keep cursor in view
func (m *Model) ensureCursorVisible() {
	// If cursor is above viewport, scroll up
	if m.Cursor.Row < m.ScrollOffset {
		m.ScrollOffset = m.Cursor.Row
	}

	// If cursor is below viewport, scroll down
	if m.Cursor.Row >= m.ScrollOffset+m.ViewHeight {
		m.ScrollOffset = m.Cursor.Row - m.ViewHeight + 1
	}
}

func (m Model) handleKeypress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	}

	switch m.Mode {
	case NormalMode:
		return m.handleNormalMode(msg)
	case InsertMode:
		return m.handleInsertMode(msg)
	case VisualMode:
		return m.handleVisualMode(msg)
	}

	return m, nil
}

func (m Model) handleNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "i":
		m.Mode = InsertMode

	case "v":
		m.Mode = VisualMode

	case "h":
		m.Cursor.MoveLeft()

	case "l":
		m.Cursor.MoveRight(m.Buffer.GetLineLength(m.Cursor.Row))

	case "j":
		if m.Cursor.Row < m.Buffer.NumLines()-1 {
			m.Cursor.MoveDown(m.Buffer.GetLineLength(m.Cursor.Row + 1))
			m.ensureCursorVisible()
		}

	case "k":
		if m.Cursor.Row > 0 {
			m.Cursor.MoveUp(m.Buffer.GetLineLength(m.Cursor.Row - 1))
			m.ensureCursorVisible()
		}

	case "-":
		if parentDir, err := m.Buffer.GetParentDir(); err == nil {
			return m, m.readDirectory(parentDir)
		}

	case "enter":
		if entry, ok := m.Buffer.GetEntry(m.Cursor.Row); ok && entry.IsDir {
			return m, m.readDirectory(entry.Path)
		}
	}

	return m, nil
}

func (m Model) handleInsertMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.Mode = NormalMode
	}
	return m, nil
}

func (m Model) handleVisualMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.Mode = NormalMode

	case "j":
		if m.Cursor.Row < m.Buffer.NumLines()-1 {
			m.Cursor.MoveLeft()
		}

	case "k":
		if m.Cursor.Row > 0 {
			m.Cursor.MoveRight(m.Buffer.GetLineLength(m.Cursor.Row))
		}
	}

	return m, nil
}

func (m Model) View() string {
	var s strings.Builder

	// Calculate visible range
	endIdx := min(m.ScrollOffset+m.ViewHeight, m.Buffer.NumLines())

	// Calculate available height for content
	// Reserve 1 line for statusline
	contentHeight := m.Height - 1

	// Render visible buffer content
	for i := m.ScrollOffset; i < endIdx && i-m.ScrollOffset < contentHeight; i++ {
		line := m.Buffer.GetLine(i)

		// Handle cursor rendering
		if i == m.Cursor.Row {
			if m.Cursor.Col >= len(line) {
				line = line + " " // Add space for cursor at end
			}
			s.WriteString(line[:m.Cursor.Col])
			if m.Cursor.Col < len(line) {
				s.WriteString(ui.CursorStyle.Render(string(line[m.Cursor.Col])))
				s.WriteString(line[m.Cursor.Col+1:])
			} else {
				s.WriteString(ui.CursorStyle.Render(" "))
			}
		} else {
			s.WriteString(line)
		}
		s.WriteString("\n")
	}

	// Fill remaining space with empty lines to push statusline to bottom
	currentLines := strings.Count(s.String(), "\n")
	for i := 0; i < contentHeight-currentLines-1; i++ {
		s.WriteString("\n")
	}

	// Render status line at the bottom
	s.WriteString(ui.RenderStatusLine(
		m.Mode.String(),
		m.Buffer.GetCurrentDir(),
		m.Width,
		m.Cursor.Row,
		m.Cursor.Col,
	))

	return s.String()
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
