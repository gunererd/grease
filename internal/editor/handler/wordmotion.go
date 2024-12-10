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

	// Return current position if buffer or line is empty
	if len(runes) == 0 {
		return pos
	}

	// If we're at the end of the current line, try to move to the next line
	if col >= len(runes) {
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
		if !isWordChar(runes[0]) {
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
func CreateWordMotionCommand(bigWord bool, operation types.Operation) func(e types.Editor) (types.Editor, tea.Cmd) {
	motion := NewWordMotion(bigWord)
	cmd := NewMotionCommand(motion, operation)
	return cmd.Execute
}

// Factory function for word end motion commands
func CreateWordEndMotionCommand(bigWord bool, operation types.Operation) func(e types.Editor) (types.Editor, tea.Cmd) {
	motion := NewWordEndMotion(bigWord)
	cmd := NewMotionCommand(motion, operation)
	return cmd.Execute
}

// Factory function for word back motion commands
func CreateWordBackMotionCommand(bigWord bool, operation types.Operation) func(e types.Editor) (types.Editor, tea.Cmd) {
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
