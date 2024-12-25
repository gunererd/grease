package motion

import (
	"fmt"

	"github.com/gunererd/grease/internal/editor/types"
)

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
	cursor types.Cursor
}

func NewMotionCommand(motion Motion, cursor types.Cursor) *MotionCommand {
	return &MotionCommand{motion: motion, cursor: cursor}
}

func (c *MotionCommand) Execute(e types.Editor) types.Editor {

	// Calculate new position using the motion
	newPos := c.motion.Calculate(
		bufferToLines(e.Buffer()),
		c.cursor.GetPosition(),
	)

	c.cursor.SetPosition(newPos)
	return e
}

func (c *MotionCommand) Name() string {
	return c.motion.Name()
}

func (c *MotionCommand) Explain() string {
	return fmt.Sprintf("type:<MotionCommand>, cmd:<%s>, cursor:<%d>, pos:<%v>", c.Name(), c.cursor.ID(), c.cursor.GetPosition())
}
func bufferToLines(buf types.Buffer) []string {
	lines := make([]string, buf.LineCount())
	for i := 0; i < buf.LineCount(); i++ {
		line, _ := buf.GetLine(i)
		lines[i] = line
	}
	return lines
}
