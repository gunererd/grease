package motion

import (
	"strings"
	"testing"

	"github.com/gunererd/grease/internal/buffer"
	"github.com/gunererd/grease/internal/types"
	"github.com/stretchr/testify/suite"
)

type WordBackMotionTestSuite struct {
	suite.Suite
}

func (s *WordMotionTestSuite) TestWordBackMotion() {
	tests := []struct {
		name     string
		input    string
		pos      types.Position
		expected types.Position
		bigWord  bool
	}{
		{
			name:     "basic word back movement",
			input:    "hello world",
			pos:      buffer.NewPosition(0, 6),
			expected: buffer.NewPosition(0, 0),
			bigWord:  false,
		},
		{
			name:     "movement across punctuation",
			input:    "hello, world",
			pos:      buffer.NewPosition(0, 7),
			expected: buffer.NewPosition(0, 5),
			bigWord:  false,
		},
		{
			name:     "start of line to previous line",
			input:    "hello\nworld",
			pos:      buffer.NewPosition(1, 0),
			expected: buffer.NewPosition(0, 0),
			bigWord:  false,
		},
		{
			name:     "from word to previous line with punctuation",
			input:    "hello,\nworld",
			pos:      buffer.NewPosition(1, 0),
			expected: buffer.NewPosition(0, 5),
			bigWord:  false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			lines := strings.Split(tt.input, "\n")
			motion := NewWordBackMotion(tt.bigWord)
			result := motion.Calculate(lines, tt.pos)
			s.Equal(tt.expected, result, "positions should match")
		})
	}
}

func TestWordBackMotionSuite(t *testing.T) {
	suite.Run(t, new(WordBackMotionTestSuite))
}
