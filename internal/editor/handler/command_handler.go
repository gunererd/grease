package handler

import (
	"github.com/gunererd/grease/internal/editor/command"
	"github.com/gunererd/grease/internal/editor/command/change"
	"github.com/gunererd/grease/internal/editor/command/clipboard"
	"github.com/gunererd/grease/internal/editor/command/delete"
	"github.com/gunererd/grease/internal/editor/command/insert"
	"github.com/gunererd/grease/internal/editor/command/motion"
	"github.com/gunererd/grease/internal/editor/register"
	"github.com/gunererd/grease/internal/editor/types"
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

func CreatePasteCommand(cursorID int, register *register.Register, before bool, history types.HistoryManager) Command {
	cmd := clipboard.NewPasteCommandAdapter(cursorID, register, before)
	return command.NewHistoryAwareCommand(cmd, history)
}

func CreateGoToStartOfLineCommand(cursorID int) Command {
	return CreateMotionCommand(motion.NewStartOfLineMotion(), cursorID)
}

func CreateGoToEndOfLineCommand(cursorID int) Command {
	return CreateMotionCommand(motion.NewEndOfLineMotion(), cursorID)
}

func CreateNewLineCommand(before bool, history types.HistoryManager) Command {
	cmd := insert.NewNewLineCommand(before)
	return command.NewHistoryAwareCommand(cmd, history)
}

func CreateDeleteToEndOfLineCommand(cursorID int, history types.HistoryManager) Command {
	cmd := delete.NewDeleteToEndOfLineCommand(cursorID)
	return command.NewHistoryAwareCommand(cmd, history)
}

func CreateDeleteLineCommand(cursorID int, history types.HistoryManager) Command {
	cmd := delete.NewDeleteLineCommand(cursorID)
	return command.NewHistoryAwareCommand(cmd, history)
}

func CreateGoToStartOfBufferCommand(cursorID int) Command {
	return CreateMotionCommand(motion.NewStartOfBufferMotion(), cursorID)
}

func CreateGoToEndOfBufferCommand(cursorID int) Command {
	return CreateMotionCommand(motion.NewEndOfBufferMotion(), cursorID)
}

func CreateDeleteCommand(motion motion.Motion, history types.HistoryManager) Command {
	cmd := delete.NewDeleteCommandAdapter(motion)
	return command.NewHistoryAwareCommand(cmd, history)
}

func CreateChangeCommand(motion motion.Motion, history types.HistoryManager) Command {
	cmd := change.NewChangeCommandAdapter(motion)
	return command.NewHistoryAwareCommand(cmd, history)
}

func CreateChangeLineCommand(cursorID int, history types.HistoryManager) Command {
	cmd := change.NewChangeLineCommand(cursorID)
	return command.NewHistoryAwareCommand(cmd, history)
}

func CreateChangeToEndOfLineCommand(cursorID int, history types.HistoryManager) Command {
	cmd := change.NewChangeToEndOfLineCommand(cursorID)
	return command.NewHistoryAwareCommand(cmd, history)
}

func CreateAppendCommand(endOfLine bool, cursorID int, history types.HistoryManager) Command {
	cmd := insert.NewAppendCommand(endOfLine, cursorID)
	return command.NewHistoryAwareCommand(cmd, history)
}

func CreateInsertCommand(startOfLine bool, cursorID int, history types.HistoryManager) Command {
	cmd := insert.NewInsertCommand(startOfLine, cursorID)
	return command.NewHistoryAwareCommand(cmd, history)
}

func CreateHalfPageDownCommand(cursorID int, viewport types.Viewport) Command {
	return CreateMotionCommand(motion.NewHalfPageDownMotion(viewport), cursorID)
}

func CreateHalfPageUpCommand(cursorID int, viewport types.Viewport) Command {
	return CreateMotionCommand(motion.NewHalfPageUpMotion(viewport), cursorID)
}

func bufferToLines(buf types.Buffer) []string {
	lines := make([]string, buf.LineCount())
	for i := 0; i < buf.LineCount(); i++ {
		line, _ := buf.GetLine(i)
		lines[i] = line
	}
	return lines
}
