package motion

import (
	"unicode"

	"github.com/gunererd/grease/internal/editor/buffer"
	"github.com/gunererd/grease/internal/editor/types"
)

type WordMotion struct {
	bigWord bool
}

func NewWordMotion(bigWord bool) *WordMotion {
	return &WordMotion{
		bigWord: bigWord,
	}
}

func (m *WordMotion) Calculate(lines []string, pos types.Position) types.Position {
	if len(lines) == 0 {
		return pos
	}

	line := lines[pos.Line()]

	// Handle empty or single character line
	if len(line) <= 1 {
		return m.handleEmptyOrSingleCharLine(lines, pos)
	}

	// Get next position based on current position
	var nextPos types.Position
	if m.bigWord {
		nextPos = m.handleBigWordMotion(line, pos)
	} else {
		nextPos = m.handleWordMotion(line, pos)
	}

	// If we didn't move or reached end of line, try next line
	if (nextPos.Line() == pos.Line() && nextPos.Column() >= len(line)) ||
		(nextPos.Equal(pos) && pos.Line() < len(lines)-1) {
		return buffer.NewPosition(pos.Line()+1, 0)
	}

	return nextPos
}

func (m *WordMotion) handleBigWordMotion(line string, pos types.Position) types.Position {
	col := pos.Column()
	// Skip until we find a space
	for col < len(line) && !unicode.IsSpace(rune(line[col])) {
		col++
	}
	// Skip spaces
	for col < len(line) && unicode.IsSpace(rune(line[col])) {
		col++
	}
	return buffer.NewPosition(pos.Line(), col)
}

func (m *WordMotion) handleWordMotion(line string, pos types.Position) types.Position {
	col := pos.Column()

	// If we're on punctuation
	if col < len(line) && isPunctuation(rune(line[col])) {
		return m.handlePunctuationMotion(line, pos)
	}

	// Skip current word
	for col < len(line) && !unicode.IsSpace(rune(line[col])) && !isPunctuation(rune(line[col])) {
		col++
	}

	// If we hit punctuation, stop there
	if col < len(line) && isPunctuation(rune(line[col])) {
		return buffer.NewPosition(pos.Line(), col)
	}

	// Skip spaces
	for col < len(line) && unicode.IsSpace(rune(line[col])) {
		col++
	}

	return buffer.NewPosition(pos.Line(), col)
}

func (m *WordMotion) handlePunctuationMotion(line string, pos types.Position) types.Position {
	col := pos.Column()

	// If next character is also punctuation
	if col+1 < len(line) && isPunctuation(rune(line[col+1])) {
		return buffer.NewPosition(pos.Line(), col+1)
	}

	// Skip current punctuation and spaces
	col++
	for col < len(line) && (unicode.IsSpace(rune(line[col])) || isPunctuation(rune(line[col]))) {
		col++
	}

	return buffer.NewPosition(pos.Line(), col)
}

func (m *WordMotion) handleEmptyOrSingleCharLine(lines []string, pos types.Position) types.Position {
	if pos.Line() < len(lines)-1 {
		return buffer.NewPosition(pos.Line()+1, 0)
	}
	return pos
}

func (m *WordMotion) Name() string {
	if m.bigWord {
		return "move_next_long_word_start"
	}
	return "move_next_word_start"
}

func isPunctuation(r rune) bool {
	return !unicode.IsSpace(r) && !unicode.IsLetter(r) && !unicode.IsNumber(r)
}
