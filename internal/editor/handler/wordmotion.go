package handler

import (
	"regexp"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/buffer"
	"github.com/gunererd/grease/internal/types"
)

var (
	whitespaceRegex = regexp.MustCompile(`^[\s]$`)
	wordCharRegex   = regexp.MustCompile(`^[\p{L}\p{N}_]$`)
)

// Motion defines how to calculate target position from current position
type Motion interface {
	Calculate(buf types.Buffer, pos types.Position) types.Position
}

type WordMotion struct {
	bigWord bool
}

func NewWordMotion(bigWord bool) *WordMotion {
	return &WordMotion{bigWord: bigWord}
}

func isWordChar(r rune) bool {
	return wordCharRegex.MatchString(string(r))
}

func isWhitespace(r rune) bool {
	return whitespaceRegex.MatchString(string(r))
}

func getCharType(r rune) int {
	if isWhitespace(r) {
		return 0 // whitespace
	}
	if isWordChar(r) {
		return 1 // word character
	}
	return 2 // punctuation/symbol
}

func (wm *WordMotion) Calculate(buf types.Buffer, pos types.Position) types.Position {
	line, _ := buf.GetLine(pos.Line())
	runes := []rune(line)
	col := pos.Column()

	// If we're at the end of the current line, try to move to the next line
	if col >= len(runes) {
		nextLine := pos.Line() + 1
		if nextLine < buf.LineCount() {
			return buffer.NewPosition(nextLine, 0)
		}
		return pos
	}

	if wm.bigWord {
		// For 'W', move to next non-whitespace after whitespace
		// Skip non-whitespace
		for col < len(runes) && !isWhitespace(runes[col]) {
			col++
		}
		// Skip whitespace
		for col < len(runes) && isWhitespace(runes[col]) {
			col++
		}
	} else {
		// For 'w', handle word characters and punctuation separately
		if col < len(runes) {
			startType := getCharType(runes[col])
			col++
			// Skip characters of the same type
			for col < len(runes) && getCharType(runes[col]) == startType {
				col++
			}
			// Skip any whitespace
			for col < len(runes) && isWhitespace(runes[col]) {
				col++
			}
		}
	}

	// If we reached the end of line, move to the start of next line
	if col >= len(runes) {
		nextLine := pos.Line() + 1
		if nextLine < buf.LineCount() {
			return buffer.NewPosition(nextLine, 0)
		}
	}

	return buffer.NewPosition(pos.Line(), col)
}

// MotionCommand combines a motion with an optional operation
type MotionCommand struct {
	motion    Motion
	operation Operation
}

func NewMotionCommand(motion Motion, operation Operation) *MotionCommand {
	return &MotionCommand{
		motion:    motion,
		operation: operation,
	}
}

func (mc *MotionCommand) Execute(e types.Editor) (types.Editor, tea.Cmd) {
	buf := e.Buffer()
	cursor, _ := buf.GetPrimaryCursor()
	curPos := cursor.GetPosition()
	targetPos := mc.motion.Calculate(buf, curPos)

	if mc.operation != nil {
		return mc.operation.Execute(e, curPos, targetPos)
	}

	// If no operation, just move cursor
	cursor.SetPosition(targetPos)
	return e, nil
}

// Factory function for word motion commands
func CreateWordMotionCommand(bigWord bool, operation Operation) func(e types.Editor) (types.Editor, tea.Cmd) {
	motion := NewWordMotion(bigWord)
	cmd := NewMotionCommand(motion, operation)
	return cmd.Execute
}

// Factory function for word end motion commands
func CreateWordEndMotionCommand(bigWord bool, operation Operation) func(e types.Editor) (types.Editor, tea.Cmd) {
	motion := NewWordEndMotion(bigWord)
	cmd := NewMotionCommand(motion, operation)
	return cmd.Execute
}

// Factory function for word back motion commands
func CreateWordBackMotionCommand(bigWord bool, operation Operation) func(e types.Editor) (types.Editor, tea.Cmd) {
	motion := NewWordBackMotion(bigWord)
	cmd := NewMotionCommand(motion, operation)
	return cmd.Execute
}

// WordEndMotion implements Motion for word end movements
type WordEndMotion struct {
	bigWord bool
}

func NewWordEndMotion(bigWord bool) *WordEndMotion {
	return &WordEndMotion{bigWord: bigWord}
}

func (wm *WordEndMotion) Calculate(buf types.Buffer, pos types.Position) types.Position {
	line, _ := buf.GetLine(pos.Line())
	runes := []rune(line)
	col := pos.Column()

	// If we're at the end of the current line, try to move to the next line
	if col >= len(runes) {
		nextLine := pos.Line() + 1
		if nextLine < buf.LineCount() {
			return buffer.NewPosition(nextLine, 0)
		}
		return pos
	}

	if wm.bigWord {
		// For 'E', move to end of current WORD
		for col < len(runes)-1 && !isWhitespace(runes[col+1]) {
			col++
		}
	} else {
		// For 'e', handle word characters and punctuation separately
		if col < len(runes)-1 {
			startType := getCharType(runes[col+1])
			col++
			// Move to the last character of the current word
			for col < len(runes)-1 && getCharType(runes[col+1]) == startType {
				col++
			}
		}
	}

	return buffer.NewPosition(pos.Line(), col)
}

// WordBackMotion implements Motion for word backward movements
type WordBackMotion struct {
	bigWord bool
}

func NewWordBackMotion(bigWord bool) *WordBackMotion {
	return &WordBackMotion{bigWord: bigWord}
}

func (wm *WordBackMotion) Calculate(buf types.Buffer, pos types.Position) types.Position {
	line, _ := buf.GetLine(pos.Line())
	runes := []rune(line)
	col := pos.Column()

	// If we're at the start of the current line, try to move to the previous line
	if col <= 0 {
		prevLine := pos.Line() - 1
		if prevLine >= 0 {
			prevLineLen, _ := buf.LineLen(prevLine)
			return buffer.NewPosition(prevLine, prevLineLen)
		}
		return pos
	}

	if wm.bigWord {
		// For 'B', move backward to start of current WORD
		// Skip whitespace
		for col > 0 && isWhitespace(runes[col-1]) {
			col--
		}
		// Move to start of current WORD
		for col > 0 && !isWhitespace(runes[col-1]) {
			col--
		}
	} else {
		// For 'b', handle word characters and punctuation separately
		if col > 0 {
			// Skip whitespace
			for col > 0 && isWhitespace(runes[col-1]) {
				col--
			}
			if col > 0 {
				startType := getCharType(runes[col-1])
				// Move to start of current word
				for col > 0 && getCharType(runes[col-1]) == startType {
					col--
				}
			}
		}
	}

	return buffer.NewPosition(pos.Line(), col)
}
