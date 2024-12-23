package handler

import (
	"github.com/gunererd/grease/internal/editor/command/change"
	"github.com/gunererd/grease/internal/editor/command/clipboard"
	"github.com/gunererd/grease/internal/editor/command/delete"
	"github.com/gunererd/grease/internal/editor/command/insert"
	"github.com/gunererd/grease/internal/editor/command/motion"
	"github.com/gunererd/grease/internal/editor/register"
	"github.com/gunererd/grease/internal/editor/types"
)

func CreateMotionCommand(m motion.Motion, cursor types.Cursor) types.Command {
	return motion.NewMotionCommand(m, cursor)
}

// Factory functions for different command types
func CreateWordMotionCommand(bigWord bool, cursor types.Cursor) types.Command {
	return CreateMotionCommand(motion.NewWordMotion(bigWord), cursor)
}

func CreateWordEndMotionCommand(bigWord bool, cursor types.Cursor) types.Command {
	return CreateMotionCommand(motion.NewWordEndMotion(bigWord), cursor)
}

func CreateWordBackMotionCommand(bigWord bool, cursor types.Cursor) types.Command {
	return CreateMotionCommand(motion.NewWordBackMotion(bigWord), cursor)
}

func CreateYankCommand(motion motion.Motion, register *register.Register) types.Command {
	return clipboard.NewYankCommandAdapter(motion, register)
}

func CreatePasteCommand(cursor types.Cursor, register *register.Register, before bool) types.Command {
	return clipboard.NewPasteCommandAdapter(cursor, register, before)
}

func CreateGoToStartOfLineCommand(cursor types.Cursor) types.Command {
	return CreateMotionCommand(motion.NewStartOfLineMotion(), cursor)
}

func CreateGoToEndOfLineCommand(cursor types.Cursor) types.Command {
	return CreateMotionCommand(motion.NewEndOfLineMotion(), cursor)
}

func CreateNewLineCommand(before bool) types.Command {
	return insert.NewNewLineCommand(before)
}

func CreateDeleteToEndOfLineCommand(cursor types.Cursor) types.Command {
	return delete.NewDeleteToEndOfLineCommand(cursor)
}

func CreateDeleteLineCommand(cursor types.Cursor) types.Command {
	return delete.NewDeleteLineCommand(cursor)
}

func CreateGoToStartOfBufferCommand(cursor types.Cursor) types.Command {
	return CreateMotionCommand(motion.NewStartOfBufferMotion(), cursor)
}

func CreateGoToEndOfBufferCommand(cursor types.Cursor) types.Command {
	return CreateMotionCommand(motion.NewEndOfBufferMotion(), cursor)
}

func CreateDeleteCommand(motion motion.Motion) types.Command {
	return delete.NewDeleteCommandAdapter(motion)
}

func CreateChangeCommand(motion motion.Motion) types.Command {
	return change.NewChangeCommandAdapter(motion)
}

func CreateChangeLineCommand(cursor types.Cursor) types.Command {
	return change.NewChangeLineCommand(cursor)
}

func CreateChangeToEndOfLineCommand(cursor types.Cursor) types.Command {
	return change.NewChangeToEndOfLineCommand(cursor)
}

func CreateAppendCommand(endOfLine bool, cursor types.Cursor) types.Command {
	return insert.NewAppendCommand(endOfLine, cursor)
}

func CreateInsertCommand(startOfLine bool, cursor types.Cursor) types.Command {
	return insert.NewInsertCommand(startOfLine, cursor)
}

func CreateHalfPageDownCommand(cursor types.Cursor, viewport types.Viewport) types.Command {
	return CreateMotionCommand(motion.NewHalfPageDownMotion(viewport), cursor)
}

func CreateHalfPageUpCommand(cursor types.Cursor, viewport types.Viewport) types.Command {
	return CreateMotionCommand(motion.NewHalfPageUpMotion(viewport), cursor)
}
