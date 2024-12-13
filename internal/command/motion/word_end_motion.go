package motion

import (
	"unicode"

	"github.com/gunererd/grease/internal/buffer"
	"github.com/gunererd/grease/internal/types"
)

type WordEndMotion struct {
	bigWord bool
}

func NewWordEndMotion(bigWord bool) *WordEndMotion {
	return &WordEndMotion{
		bigWord: bigWord,
	}
}

func (m *WordEndMotion) Calculate(lines []string, pos types.Position) types.Position {
	if len(lines) == 0 {
		return pos
	}

	line := lines[pos.Line()]
	col := pos.Column()

	// Find end of next word
	if col < len(line)-1 {
		startCol := col
		for col < len(line) && !isWordBoundary(rune(line[col]), m.bigWord) {
			col++
		}
		if col > startCol {
			return buffer.NewPosition(pos.Line(), col-1)
		}
	}

	// If we're at end of line, try next line
	if col >= len(line)-1 && pos.Line() < len(lines)-1 {
		nextLine := lines[pos.Line()+1]
		col = 0
		// Skip leading spaces in next line
		for col < len(nextLine) && unicode.IsSpace(rune(nextLine[col])) {
			col++
		}
		// Find end of first word in next line
		startCol := col
		for col < len(nextLine)-1 && !isWordBoundary(rune(nextLine[col]), m.bigWord) {
			col++
		}
		if col >= startCol {
			return buffer.NewPosition(pos.Line()+1, col)
		}
	}

	return pos
}

func (m *WordEndMotion) Name() string {
	if m.bigWord {
		return "move_next_word_end"
	}
	return "move_next_long_word_end"
}
