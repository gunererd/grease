package motion

import (
	"unicode"

	"github.com/gunererd/grease/internal/buffer"
	"github.com/gunererd/grease/internal/types"
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
	col := pos.Column()

	// Skip current word/WORD
	for col < len(line) && !m.isWordBoundary(line[col:], m.bigWord) {
		col++
	}

	// Skip spaces
	for col < len(line) && unicode.IsSpace(rune(line[col])) {
		col++
	}

	// If we reached end of line, try next line
	if col >= len(line) && pos.Line() < len(lines)-1 {
		return buffer.NewPosition(pos.Line()+1, 0)
	}

	// If we're still in bounds, return new position
	if col < len(line) {
		return buffer.NewPosition(pos.Line(), col)
	}

	return pos
}

func (m *WordMotion) isWordBoundary(text string, bigWord bool) bool {
	if len(text) == 0 {
		return false
	}

	if bigWord {
		return unicode.IsSpace(rune(text[0]))
	}

	// For normal words, consider punctuation as boundary
	return unicode.IsSpace(rune(text[0])) || unicode.IsPunct(rune(text[0]))
}

func (m *WordMotion) Name() string {
	if m.bigWord {
		return "move_next_long_word_start"
	}
	return "move_next_word_start"
}
