package motion

import (
	"regexp"

	"github.com/gunererd/grease/internal/editor/buffer"
	"github.com/gunererd/grease/internal/editor/types"
)

var (
	whitespaceRegex = regexp.MustCompile(`^[\s]$`)
	wordCharRegex   = regexp.MustCompile(`^[\p{L}\p{N}_]$`)
)

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

func (wm *WordMotion) Calculate(lines []string, pos types.Position) types.Position {
	if len(lines) == 0 || pos.Line() >= len(lines) {
		return pos
	}

	line := lines[pos.Line()]
	runes := []rune(line)
	col := pos.Column()

	// Return current position if line is empty or has only one character
	if len(runes) == 0 || len(runes) == 1 {
		return pos
	}

	// If we're at the end of the current line, try to move to the next line
	if col >= len(runes)-1 {
		moveToNextLine(lines, pos)
	}

	if wm.bigWord {
		col = moveToNextBigWord(runes, col)
	} else {
		col = moveToNextWord(runes, col)
	}

	if col >= len(runes) {
		pos := moveToNextLine(lines, pos)
		// If first character of next line is not a word character, return current position
		line := lines[pos.Line()]
		runes = []rune(line)
		if len(runes) > 0 && !isWordChar(runes[0]) {
			return wm.Calculate(lines, pos)
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
func moveToNextLine(lines []string, pos types.Position) types.Position {
	nextLine := pos.Line() + 1
	if nextLine < len(lines) {
		return buffer.NewPosition(nextLine, 0)
	}
	if pos.Line() < len(lines) {
		return buffer.NewPosition(pos.Line(), len([]rune(lines[pos.Line()]))-1)
	}
	return pos
}

// WordEndMotion implements Motion for word end movements
type WordEndMotion struct {
	bigWord bool
}

func NewWordEndMotion(bigWord bool) *WordEndMotion {
	return &WordEndMotion{bigWord: bigWord}
}

func (wm *WordEndMotion) Calculate(lines []string, pos types.Position) types.Position {
	if len(lines) == 0 || pos.Line() >= len(lines) {
		return pos
	}

	line := lines[pos.Line()]
	runes := []rune(line)
	col := pos.Column()

	// Return current position if line is empty or has only one character
	if len(runes) == 0 || len(runes) == 1 {
		return pos
	}

	if col >= len(runes)-1 {
		// Try to move to next line
		nextPos := moveToNextLine(lines, pos)
		if nextPos.Line() == pos.Line() {
			return pos // We couldn't move to next line
		}

		nextLine := lines[nextPos.Line()]
		nextRunes := []rune(nextLine)
		if len(nextRunes) > 0 {
			return wm.Calculate(lines, nextPos)
		}
		return nextPos
	}

	if wm.bigWord {
		col = moveToNextBigWordEnd(runes, col)
	} else {
		col = moveToNextWordEnd(runes, col)
	}

	if col >= len(runes) && pos.Line()+1 < len(lines) {
		return moveToNextLine(lines, pos)
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

func (wbm *WordBackMotion) Calculate(lines []string, pos types.Position) types.Position {
	if len(lines) == 0 || pos.Line() >= len(lines) {
		return pos
	}

	line := lines[pos.Line()]
	runes := []rune(line)
	col := pos.Column()

	// If we're at the start of the line, move to previous line
	if col <= 0 {
		prevPos := moveToPrevLine(lines, pos)
		if prevPos.Line() == pos.Line() {
			return pos // We couldn't move to previous line
		}

		line := lines[prevPos.Line()]
		runes = []rune(line)
		if !isWhitespace(runes[prevPos.Column()]) && !isWordChar(runes[prevPos.Column()]) {
			return prevPos
		}

		return wbm.Calculate(lines, prevPos)
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
		return wbm.Calculate(lines, buffer.NewPosition(pos.Line(), col))
	}

	return buffer.NewPosition(pos.Line(), col)
}

func moveToPrevLine(lines []string, pos types.Position) types.Position {
	prevLine := pos.Line() - 1
	if prevLine >= 0 {
		return buffer.NewPosition(prevLine, len([]rune(lines[prevLine]))-1)
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
