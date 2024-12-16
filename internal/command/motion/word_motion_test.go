package motion

import (
	"strings"
	"testing"

	"github.com/gunererd/grease/internal/buffer"
	"github.com/gunererd/grease/internal/types"
	"github.com/stretchr/testify/suite"
)

type WordMotionTestSuite struct {
	suite.Suite
}

func (s *WordMotionTestSuite) TestWordMotion() {
	tests := []struct {
		name     string
		input    string
		pos      types.Position
		expected types.Position
		bigWord  bool
	}{
		{
			name:     "basic word movement",
			input:    "hello world",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 6),
			bigWord:  false,
		},
		{
			name:     "movement across punctuation",
			input:    "hello, world",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 5),
			bigWord:  false,
		},
		{
			name:     "big word movement",
			input:    "hello, world!",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 7),
			bigWord:  true,
		},
		{
			name:     "leading spaces",
			input:    "    hello",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 4),
			bigWord:  false,
		},
		{
			name:     "end of line to next line",
			input:    "hello\nworld",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(1, 0),
			bigWord:  false,
		},
		{
			name:     "move multiple lines",
			input:    "hello\n\n\nworld",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(1, 0),
			bigWord:  false,
		},
		{
			name:     "jump to next word from punctuation",
			input:    "hello, world",
			pos:      buffer.NewPosition(0, 5),
			expected: buffer.NewPosition(0, 7),
			bigWord:  false,
		},
		{
			name:     "jump to next punctuation from punctuation",
			input:    "hello,, world",
			pos:      buffer.NewPosition(0, 5),
			expected: buffer.NewPosition(0, 6),
			bigWord:  false,
		},
		{
			name:     "jump to next word from punctuation no space",
			input:    "hello,beautiful world",
			pos:      buffer.NewPosition(0, 5),
			expected: buffer.NewPosition(0, 6),
			bigWord:  false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			lines := strings.Split(tt.input, "\n")
			motion := NewWordMotion(tt.bigWord)
			result := motion.Calculate(lines, tt.pos)
			s.Equal(tt.expected, result, "positions should match")
		})
	}
}

func TestWordMotionSuite(t *testing.T) {
	suite.Run(t, new(WordMotionTestSuite))
}
