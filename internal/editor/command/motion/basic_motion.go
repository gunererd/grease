package motion

import (
	"log"

	"github.com/gunererd/grease/internal/editor/buffer"
	"github.com/gunererd/grease/internal/editor/types"
)

type LeftMotion struct{}

func NewLeftMotion() *LeftMotion {
	return &LeftMotion{}
}

func (m *LeftMotion) Calculate(lines []string, pos types.Position) types.Position {
	if pos.Column() > 0 {
		return buffer.NewPosition(pos.Line(), pos.Column()-1)
	}
	return pos
}

func (m *LeftMotion) Name() string {
	return "move_char_left"
}

type RightMotion struct{}

func NewRightMotion() *RightMotion {
	return &RightMotion{}
}

func (m *RightMotion) Calculate(lines []string, pos types.Position) types.Position {
	if pos.Line() < len(lines) {
		lineLen := len([]rune(lines[pos.Line()]))
		if pos.Column() < lineLen-1 {
			return buffer.NewPosition(pos.Line(), pos.Column()+1)
		}
	}
	return pos
}

func (m *RightMotion) Name() string {
	return "move_char_right"
}

type UpMotion struct{}

func NewUpMotion() *UpMotion {
	return &UpMotion{}
}

func (m *UpMotion) Calculate(lines []string, pos types.Position) types.Position {
	if pos.Line() <= 0 {
		return pos
	}

	// Get target line length
	targetLineLen := len([]rune(lines[pos.Line()-1]))

	// Keep same column if possible, otherwise go to end of line
	newCol := pos.Column()
	if newCol >= targetLineLen {
		newCol = targetLineLen - 1
		if newCol < 0 {
			newCol = 0
		}
	}

	return buffer.NewPosition(pos.Line()-1, newCol)
}

func (m *UpMotion) Name() string {
	return "move_line_up"
}

type DownMotion struct{}

func NewDownMotion() *DownMotion {
	return &DownMotion{}
}

func (m *DownMotion) Calculate(lines []string, pos types.Position) types.Position {
	if pos.Line() >= len(lines)-1 {
		return pos
	}

	// Get target line length
	targetLineLen := len([]rune(lines[pos.Line()+1]))

	// Keep same column if possible, otherwise go to end of line
	newCol := pos.Column()
	if newCol >= targetLineLen {
		newCol = targetLineLen - 1
		if newCol < 0 {
			newCol = 0
		}
	}
	return buffer.NewPosition(pos.Line()+1, newCol)
}

func (m *DownMotion) Name() string {
	return "move_line_down"
}

func CreateMotionCommand(motion Motion, cursorID int) func(types.Editor) types.Editor {
	return func(e types.Editor) types.Editor {
		buf := e.Buffer()
		cursor, err := buf.GetCursor(cursorID)
		if err != nil {
			return e
		}
		pos := cursor.GetPosition()
		log.Printf("type:<CreateMotionCommand>, name:<%s>, cursorID:<%d>, pos:<%v>\n", motion.Name(), cursorID, pos)
		cmd := NewMotionCommand(motion, cursorID)
		return cmd.Execute(e)
	}
}
