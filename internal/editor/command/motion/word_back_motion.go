package motion

import (
	"unicode"

	"github.com/gunererd/grease/internal/editor/buffer"
	"github.com/gunererd/grease/internal/editor/types"
)

type WordBackMotion struct {
	bigWord bool
}

func NewWordBackMotion(bigWord bool) *WordBackMotion {
	return &WordBackMotion{
		bigWord: bigWord,
	}
}

func (m *WordBackMotion) Calculate(lines []string, pos types.Position) types.Position {
	if len(lines) == 0 || (pos.Line() == 0 && pos.Column() == 0) {
		return pos
	}

	// here line is "hello"
	line := lines[pos.Line()]
	col := pos.Column()

	// Skip spaces backwards
	for col > 0 && unicode.IsSpace(rune(line[col-1])) {
		col--
	}

	// If we're at start of line and not first line, go to previous line
	if col == 0 && pos.Line() > 0 {
		pos = buffer.NewPosition(pos.Line()-1, len(lines[pos.Line()-1]))
		line = lines[pos.Line()]
		col = pos.Column()
	}

	// Find start of current/previous word
	if col > 0 {
		isCurrentWordChar := !isWordBoundary(rune(line[col-1]), m.bigWord)
		col--

		for col > 0 {
			isPrevWordChar := !isWordBoundary(rune(line[col-1]), m.bigWord)
			if isPrevWordChar != isCurrentWordChar {
				break
			}
			col--
		}
	}

	return buffer.NewPosition(pos.Line(), col)
}

func (m *WordBackMotion) Name() string {
	if m.bigWord {
		return "move_prev_long_word_start"
	}
	return "move_prev_word_start"
}

func isWordBoundary(r rune, bigWord bool) bool {
	if bigWord {
		return unicode.IsSpace(r)
	}
	return unicode.IsSpace(r) || unicode.IsPunct(r)
}
