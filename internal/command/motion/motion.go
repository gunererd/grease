package motion

import "github.com/gunererd/grease/internal/types"

// Motion calculates a new position based on current position and buffer content
type Motion interface {
	// Calculate returns the target position for this motion
	Calculate(lines []string, pos types.Position) types.Position
	// Name returns the motion name (like "h", "w", "$")
	Name() string
}

// MotionCommand wraps a motion and handles the actual cursor movement
type MotionCommand struct {
	motion Motion
}

func NewMotionCommand(motion Motion) *MotionCommand {
	return &MotionCommand{motion: motion}
}

func (c *MotionCommand) Execute(e types.Editor) types.Editor {
	cursor, err := e.Buffer().GetPrimaryCursor()
	if err != nil {
		return e
	}

	// Calculate new position using the motion
	newPos := c.motion.Calculate(
		bufferToLines(e.Buffer()),
		cursor.GetPosition(),
	)

	// Move cursor to new position
	e.Buffer().MoveCursor(cursor.ID(), newPos.Line(), newPos.Column())
	e.HandleCursorMovement()

	return e
}

func (c *MotionCommand) Name() string {
	return c.motion.Name()
}

func bufferToLines(buf types.Buffer) []string {
	lines := make([]string, buf.LineCount())
	for i := 0; i < buf.LineCount(); i++ {
		line, _ := buf.GetLine(i)
		lines[i] = line
	}
	return lines
}
