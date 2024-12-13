package handler

// import (
// 	"github.com/gunererd/grease/internal/motion"
// 	"github.com/gunererd/grease/internal/types"
// )

// // Motion defines how to calculate target position from current position

// // MotionCommand combines a motion with an optional operation
// type MotionCommand struct {
// 	motion    motion.Motion
// 	operation types.Operation
// }

// func NewMotionCommand(motion motion.Motion, operation types.Operation) *MotionCommand {
// 	return &MotionCommand{
// 		motion:    motion,
// 		operation: operation,
// 	}
// }

// func (mc *MotionCommand) Execute(e types.Editor) types.Editor {
// 	buf := e.Buffer()
// 	cursor, _ := buf.GetPrimaryCursor()
// 	curPos := cursor.GetPosition()
// 	targetPos := mc.motion.Calculate(
// 		bufferToLines(buf),
// 		curPos,
// 	)

// 	if mc.operation != nil {
// 		return mc.operation.Execute(e, curPos, targetPos)
// 	}

// 	// If no operation, just move cursor
// 	cursor.SetPosition(targetPos)
// 	return e
// }

// // Factory function for word motion commands
// func CreateWordMotionCommand(bigWord bool, operation types.Operation) func(e types.Editor) types.Editor {
// 	motion := motion.NewWordMotion(bigWord)
// 	cmd := NewMotionCommand(motion, operation)
// 	return cmd.Execute
// }

// // Factory function for word end motion commands
// func CreateWordEndMotionCommand(bigWord bool, operation types.Operation) func(e types.Editor) types.Editor {
// 	motion := motion.NewWordEndMotion(bigWord)
// 	cmd := NewMotionCommand(motion, operation)
// 	return cmd.Execute
// }

// // Factory function for word back motion commands
// func CreateWordBackMotionCommand(bigWord bool, operation types.Operation) func(e types.Editor) types.Editor {
// 	motion := motion.NewWordBackMotion(bigWord)
// 	cmd := NewMotionCommand(motion, operation)
// 	return cmd.Execute
// }

// func bufferToLines(buf types.Buffer) []string {
// 	lines := make([]string, buf.LineCount())
// 	for i := 0; i < buf.LineCount(); i++ {
// 		line, _ := buf.GetLine(i)
// 		lines[i] = line
// 	}
// 	return lines
// }
