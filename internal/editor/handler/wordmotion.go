package handler

import (
	"regexp"

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

	// Return current position if buffer or line is empty or line has only one character
	if len(runes) == 0 || len(runes) == 1 {
		return pos
	}

	// If we're at the end of the current line, try to move to the next line
	if col >= len(runes)-1 {
		return moveToNextLine(buf, pos)
	}

	if wm.bigWord {
		col = moveToNextBigWord(runes, col)
	} else {
		col = moveToNextWord(runes, col)
	}

	if col >= len(runes) {
		pos := moveToNextLine(buf, pos)
		// If first character of next line is not a word character, return current position
		line, _ := buf.GetLine(pos.Line())
		runes = []rune(line)
		if len(runes) > 0 && !isWordChar(runes[0]) {
			return wm.Calculate(buf, pos)
		}

		return pos
	}

	return buffer.NewPosition(pos.Line(), col)
}

func moveToNextBigWord(runes []rune, col int) int {
	col = skipWhile(runes, col, isWordChar)
	col = skipWhile(runes, col, func(r rune) bool { return !isWordChar(r) })
	return col
}

func moveToNextWord(runes []rune, col int) int {
	col = skipWhile(runes, col, func(r rune) bool { return !isWhitespace(r) && !isWordChar(r) })
	col = skipWhile(runes, col, isWordChar)
	col = skipWhile(runes, col, isWhitespace)
	return col
}

func skipWhile(runes []rune, start int, predicate func(r rune) bool) int {
	for start < len(runes) && predicate(runes[start]) {
		start++
	}
	return start
}

// Function to move to the next line or move to end of line if at the end
func moveToNextLine(buf types.Buffer, pos types.Position) types.Position {
	nextLine := pos.Line() + 1
	if nextLine < buf.LineCount() {
		return buffer.NewPosition(nextLine, 0)
	}
	line, err := buf.GetLine(pos.Line())
	if err != nil {
		return pos
	}

	return buffer.NewPosition(pos.Line(), len(line)-1)
}

// MotionCommand combines a motion with an optional operation
type MotionCommand struct {
	motion    Motion
	operation types.Operation
}

func NewMotionCommand(motion Motion, operation types.Operation) *MotionCommand {
	return &MotionCommand{
		motion:    motion,
		operation: operation,
	}
}

func (mc *MotionCommand) Execute(e types.Editor) types.Editor {
	buf := e.Buffer()
	cursor, _ := buf.GetPrimaryCursor()
	curPos := cursor.GetPosition()
	targetPos := mc.motion.Calculate(buf, curPos)

	if mc.operation != nil {
		return mc.operation.Execute(e, curPos, targetPos)
	}

	// If no operation, just move cursor
	cursor.SetPosition(targetPos)
	return e
}

// Factory function for word motion commands
func CreateWordMotionCommand(bigWord bool, operation types.Operation) func(e types.Editor) types.Editor {
	motion := NewWordMotion(bigWord)
	cmd := NewMotionCommand(motion, operation)
	return cmd.Execute
}

// Factory function for word end motion commands
func CreateWordEndMotionCommand(bigWord bool, operation types.Operation) func(e types.Editor) types.Editor {
	motion := NewWordEndMotion(bigWord)
	cmd := NewMotionCommand(motion, operation)
	return cmd.Execute
}

// Factory function for word back motion commands
func CreateWordBackMotionCommand(bigWord bool, operation types.Operation) func(e types.Editor) types.Editor {
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

	// Return current position if buffer or line is empty or line has only one character
	if len(runes) == 0 || len(runes) == 1 {
		return pos
	}

	if col >= len(runes)-1 {
		pos := moveToNextLine(buf, pos)
		line, _ := buf.GetLine(pos.Line())
		runes = []rune(line)
		if len(runes) > 0 {
			return wm.Calculate(buf, pos)
		}
		return pos
	}

	if wm.bigWord {
		col = moveToNextBigWordEnd(runes, col)
	} else {
		col = moveToNextWordEnd(runes, col)
	}

	if col >= len(runes) {
		return moveToNextLine(buf, pos)
	}

	return buffer.NewPosition(pos.Line(), col)
}

func moveToNextWordEnd(runes []rune, col int) int {
	if col >= len(runes) {
		return col
	}

	startType := getCharType(runes[col])
	col = skipWhile(runes, col, func(r rune) bool { return getCharType(r) == startType })
	col = skipWhile(runes, col, isWordChar)

	// Adjust position if moving into non-word character after whitespace
	if col > 0 && col < len(runes) && isWhitespace(runes[col-1]) && !isWhitespace(runes[col]) && !isWordChar(runes[col]) {
		return col
	}

	return max(col-1, 0)
}

func moveToNextBigWordEnd(runes []rune, col int) int {
	if col >= len(runes) {
		return col
	}

	startType := getCharType(runes[col])
	col = skipWhile(runes, col, func(r rune) bool { return getCharType(r) == startType })
	col = skipWhile(runes, col, func(r rune) bool { return !isWhitespace(r) })

	// Adjust position if moving into non-big-word character after whitespace
	if col > 0 && col < len(runes) && isWhitespace(runes[col-1]) && !isWhitespace(runes[col]) && !isWordChar(runes[col]) {
		return col
	}

	return max(col-1, 0)
}

// WordBackMotion implements Motion for word backward movements
type WordBackMotion struct {
	bigWord bool
}

func NewWordBackMotion(bigWord bool) *WordBackMotion {
	return &WordBackMotion{bigWord: bigWord}
}

func (wbm *WordBackMotion) Calculate(buf types.Buffer, pos types.Position) types.Position {
	line, _ := buf.GetLine(pos.Line())
	runes := []rune(line)
	col := pos.Column()

	// If we're at the start of the line, move to previous line
	if col <= 0 {
		prevPos := moveToPrevLine(buf, pos)
		if prevPos.Line() == pos.Line() {
			return pos // We couldn't move to previous line
		}

		line, _ := buf.GetLine(prevPos.Line())
		runes = []rune(line)
		if !isWhitespace(runes[prevPos.Column()]) && !isWordChar(runes[prevPos.Column()]) {
			return prevPos
		}

		return wbm.Calculate(buf, prevPos)
	}

	if wbm.bigWord {
		col = moveToPrevBigWord(runes, col)
	} else {
		col = moveToPrevWord(runes, col)
	}

	if col == 0 && !isWhitespace(runes[col]) {
		return buffer.NewPosition(pos.Line(), col)

	}

	if col == 0 && pos.Line() > 0 {
		return wbm.Calculate(buf, buffer.NewPosition(pos.Line(), col))
	}

	return buffer.NewPosition(pos.Line(), col)
}

func moveToPrevLine(buf types.Buffer, pos types.Position) types.Position {
	prevLine := pos.Line() - 1
	if prevLine >= 0 {
		line, err := buf.GetLine(prevLine)
		if err != nil {
			return pos
		}
		return buffer.NewPosition(prevLine, len([]rune(line))-1)
	}
	return buffer.NewPosition(pos.Line(), 0)
}

func moveToPrevWord(runes []rune, col int) int {
	if col <= 0 {
		return col
	}

	// If we're in whitespace, skip it
	col = skipBackWhile(runes, col, isWhitespace)
	if col == 0 {
		return col
	}

	// Get the type of character we're currently on
	startType := getCharType(runes[col-1])

	// Skip characters of the same type
	col = skipBackWhile(runes, col, func(r rune) bool { return getCharType(r) == startType })

	return col
}

func moveToPrevBigWord(runes []rune, col int) int {
	if col <= 0 {
		return col
	}

	// Skip trailing whitespace
	col = skipBackWhile(runes, col, isWhitespace)
	if col == 0 {
		return col
	}

	// Skip non-whitespace characters
	col = skipBackWhile(runes, col, func(r rune) bool { return !isWhitespace(r) })

	return col
}

func skipBackWhile(runes []rune, start int, predicate func(r rune) bool) int {
	for start > 0 && predicate(runes[start-1]) {
		start--
	}
	return start
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
