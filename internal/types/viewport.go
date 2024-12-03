package types

import (
	"github.com/gunererd/grease/internal/state"
)

type Viewport interface {
	SetSize(width, height int)
	Size() (width, height int)
	Offset() Position
	SetScrollOff(lines int)
	ScrollTo(pos Position)
	SetCursor(pos Position)
	Cursor() Position
	ToggleCursor()
	IsPositionVisible(pos Position) bool
	VisibleLines() (start, end int)
	VisibleColumns() (start, end int)
	SetMode(mode state.Mode)
	Render(buf Buffer) []string
	CenterOn(pos Position)
	BufferToViewportPosition(pos Position) (x, y int)
	ViewportToBufferPosition(x, y int) Position
	ScrollUp(lines int)
	ScrollDown(lines int)
	ScrollLeft(cols int)
	ScrollRight(cols int)
	SetHighlightManager(hm HighlightManager)
}
