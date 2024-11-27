package model

type Mode int

const (
	NormalMode Mode = iota
	InsertMode
	VisualMode
	VisualBlockMode
)

func (m Mode) String() string {
	switch m {
	case NormalMode:
		return "NORMAL"
	case InsertMode:
		return "INSERT"
	case VisualMode:
		return "VISUAL"
	case VisualBlockMode:
		return "VISUAL-BLOCK"
	default:
		return "UNKNOWN"
	}
}
