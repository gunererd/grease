package types

type StatusLine interface {
	Render(mode string, cursor Cursor, bufferLineCount int, viewX, viewY int, width int) string
}
