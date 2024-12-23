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

func CreateMotionCommand(m motion.Motion, cursor types.Cursor) command.Command {
	return motion.NewMotionCommand(m, cursor)
}

// Factory functions for different command types
func CreateWordMotionCommand(bigWord bool, cursor types.Cursor) command.Command {
	return CreateMotionCommand(motion.NewWordMotion(bigWord), cursor)
}

func CreateWordEndMotionCommand(bigWord bool, cursor types.Cursor) command.Command {
	return CreateMotionCommand(motion.NewWordEndMotion(bigWord), cursor)
}

func CreateWordBackMotionCommand(bigWord bool, cursor types.Cursor) command.Command {
	return CreateMotionCommand(motion.NewWordBackMotion(bigWord), cursor)
}

func CreateYankCommand(motion motion.Motion, register *register.Register) command.Command {
	return clipboard.NewYankCommandAdapter(motion, register)
}

func CreatePasteCommand(cursor types.Cursor, register *register.Register, before bool, history types.HistoryManager) command.Command {
	cmd := clipboard.NewPasteCommandAdapter(cursor, register, before)
	return command.NewHistoryAwareCommand(cmd, history)
}

func CreateGoToStartOfLineCommand(cursor types.Cursor) command.Command {
	return CreateMotionCommand(motion.NewStartOfLineMotion(), cursor)
}

func CreateGoToEndOfLineCommand(cursor types.Cursor) command.Command {
	return CreateMotionCommand(motion.NewEndOfLineMotion(), cursor)
}

func CreateNewLineCommand(before bool, history types.HistoryManager) command.Command {
	cmd := insert.NewNewLineCommand(before)
	return command.NewHistoryAwareCommand(cmd, history)
}

func CreateDeleteToEndOfLineCommand(cursor types.Cursor, history types.HistoryManager) command.Command {
	cmd := delete.NewDeleteToEndOfLineCommand(cursor)
	return command.NewHistoryAwareCommand(cmd, history)
}

func CreateDeleteLineCommand(cursor types.Cursor, history types.HistoryManager) command.Command {
	cmd := delete.NewDeleteLineCommand(cursor)
	return command.NewHistoryAwareCommand(cmd, history)
}

func CreateGoToStartOfBufferCommand(cursor types.Cursor) command.Command {
	return CreateMotionCommand(motion.NewStartOfBufferMotion(), cursor)
}

func CreateGoToEndOfBufferCommand(cursor types.Cursor) command.Command {
	return CreateMotionCommand(motion.NewEndOfBufferMotion(), cursor)
}

func CreateDeleteCommand(motion motion.Motion, history types.HistoryManager) command.Command {
	cmd := delete.NewDeleteCommandAdapter(motion)
	return command.NewHistoryAwareCommand(cmd, history)
}

func CreateChangeCommand(motion motion.Motion, history types.HistoryManager) command.Command {
	cmd := change.NewChangeCommandAdapter(motion)
	return command.NewHistoryAwareCommand(cmd, history)
}

func CreateChangeLineCommand(cursor types.Cursor, history types.HistoryManager) command.Command {
	cmd := change.NewChangeLineCommand(cursor)
	return command.NewHistoryAwareCommand(cmd, history)
}

func CreateChangeToEndOfLineCommand(cursor types.Cursor, history types.HistoryManager) command.Command {
	cmd := change.NewChangeToEndOfLineCommand(cursor)
	return command.NewHistoryAwareCommand(cmd, history)
}

func CreateAppendCommand(endOfLine bool, cursor types.Cursor, history types.HistoryManager) command.Command {
	cmd := insert.NewAppendCommand(endOfLine, cursor)
	return command.NewHistoryAwareCommand(cmd, history)
}

func CreateInsertCommand(startOfLine bool, cursor types.Cursor, history types.HistoryManager) command.Command {
	cmd := insert.NewInsertCommand(startOfLine, cursor)
	return command.NewHistoryAwareCommand(cmd, history)
}

func CreateHalfPageDownCommand(cursor types.Cursor, viewport types.Viewport) command.Command {
	return CreateMotionCommand(motion.NewHalfPageDownMotion(viewport), cursor)
}

func CreateHalfPageUpCommand(cursor types.Cursor, viewport types.Viewport) command.Command {
	return CreateMotionCommand(motion.NewHalfPageUpMotion(viewport), cursor)
}
