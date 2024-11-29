package types

import (
	"github.com/gunererd/grease/internal/state"
)

type Viewport interface {
	SetSize(width, height int)
	GetSize() (width, height int)
	GetOffset() Position
	SetScrollOff(lines int)
	ScrollTo(pos Position)
	SetCursor(pos Position)
	GetCursor() Position
	ToggleCursor()
	IsPositionVisible(pos Position) bool
	GetVisibleLines() (start, end int)
	GetVisibleColumns() (start, end int)
	SetMode(mode state.Mode)
	View(buf Buffer) []string
	CenterOn(pos Position)
	GetRelativePosition(pos Position) (x, y int)
	GetAbsolutePosition(x, y int) Position
	ScrollUp(lines int)
	ScrollDown(lines int)
	ScrollLeft(cols int)
	ScrollRight(cols int)
}
