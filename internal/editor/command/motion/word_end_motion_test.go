package motion

import (
	"strings"
	"testing"

	"github.com/gunererd/grease/internal/editor/buffer"
	"github.com/gunererd/grease/internal/editor/types"
	"github.com/stretchr/testify/suite"
)

type WordEndMotionTestSuite struct {
	suite.Suite
}

func (s *WordMotionTestSuite) TestWordEndMotion() {
	tests := []struct {
		name     string
		input    string
		pos      types.Position
		expected types.Position
		bigWord  bool
	}{
		{
			name:     "basic word end movement",
			input:    "hello world",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 4),
			bigWord:  false,
		},
		{
			name:     "movement across punctuation",
			input:    "hello, world",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 4),
			bigWord:  false,
		},
		{
			name:     "next line with multiple spaces",
			input:    "hello\n    world",
			pos:      buffer.NewPosition(0, 4),
			expected: buffer.NewPosition(1, 8),
			bigWord:  false,
		},
		{
			name:     "jump to next line punctuation",
			input:    "hello\n  ,   world",
			pos:      buffer.NewPosition(0, 4),
			expected: buffer.NewPosition(1, 2),
			bigWord:  false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			lines := strings.Split(tt.input, "\n")
			motion := NewWordEndMotion(tt.bigWord)
			result := motion.Calculate(lines, tt.pos)
			s.Equal(tt.expected, result, "positions should match")
		})
	}
}

func TestWordEndMotionSuite(t *testing.T) {
	suite.Run(t, new(WordMotionTestSuite))
}
