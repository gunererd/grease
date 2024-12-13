package handler

import (
	"github.com/gunererd/grease/internal/command/clipboard"
	"github.com/gunererd/grease/internal/command/motion"
	"github.com/gunererd/grease/internal/register"
	"github.com/gunererd/grease/internal/types"
)

// Command is the interface that all commands must implement
type Command interface {
	Execute(e types.Editor) types.Editor
	Name() string
}

// CreateMotionCommand creates a command from a motion
func CreateMotionCommand(motion motion.Motion, cursorID int) Command {
	return &MotionCommand{
		motion:   motion,
		cursorID: cursorID,
	}
}

// MotionCommand wraps a motion into a command
type MotionCommand struct {
	motion   motion.Motion
	cursorID int
}

func (mc *MotionCommand) Execute(e types.Editor) types.Editor {
	buf := e.Buffer()
	cursor, err := buf.GetCursor(mc.cursorID)
	if err != nil {
		return e
	}

	curPos := cursor.GetPosition()
	targetPos := mc.motion.Calculate(
		bufferToLines(buf),
		curPos,
	)
	cursor.SetPosition(targetPos)
	return e
}

func (mc *MotionCommand) Name() string {
	return mc.motion.Name()
}

// Factory functions for different command types
func CreateWordMotionCommand(bigWord bool, cursorID int) Command {
	return CreateMotionCommand(motion.NewWordMotion(bigWord), cursorID)
}

func CreateWordEndMotionCommand(bigWord bool, cursorID int) Command {
	return CreateMotionCommand(motion.NewWordEndMotion(bigWord), cursorID)
}

func CreateWordBackMotionCommand(bigWord bool, cursorID int) Command {
	return CreateMotionCommand(motion.NewWordBackMotion(bigWord), cursorID)
}

func CreateYankCommand(motion motion.Motion, register *register.Register) Command {
	return clipboard.NewYankCommandAdapter(motion, register)
}

func CreatePasteCommand(motion motion.Motion, register *register.Register, before bool) Command {
	return clipboard.NewPasteCommandAdapter(motion, register, before)
}

func bufferToLines(buf types.Buffer) []string {
	lines := make([]string, buf.LineCount())
	for i := 0; i < buf.LineCount(); i++ {
		line, _ := buf.GetLine(i)
		lines[i] = line
	}
	return lines
}
