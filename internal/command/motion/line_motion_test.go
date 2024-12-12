package motion

import (
	"strings"
	"testing"

	"github.com/gunererd/grease/internal/buffer"
	"github.com/gunererd/grease/internal/types"
	"github.com/stretchr/testify/suite"
)

type LineMotionTestSuite struct {
	suite.Suite
}

func (s *LineMotionTestSuite) TestStartOfLineMotion() {
	tests := []struct {
		name     string
		input    string
		pos      types.Position
		expected types.Position
	}{
		{
			name:     "basic start of line",
			input:    "hello world",
			pos:      buffer.NewPosition(0, 5),
			expected: buffer.NewPosition(0, 0),
		},
		{
			name:     "already at start of line",
			input:    "hello world",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 0),
		},
		{
			name:     "empty line",
			input:    "\nhello",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 0),
		},
		{
			name:     "multiple lines",
			input:    "hello\nworld",
			pos:      buffer.NewPosition(1, 3),
			expected: buffer.NewPosition(1, 0),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			lines := strings.Split(tt.input, "\n")
			motion := NewStartOfLineMotion()
			result := motion.Calculate(lines, tt.pos)
			s.Equal(tt.expected, result, "positions should match")
		})
	}
}

func (s *LineMotionTestSuite) TestEndOfLineMotion() {
	tests := []struct {
		name     string
		input    string
		pos      types.Position
		expected types.Position
	}{
		{
			name:     "basic end of line",
			input:    "hello world",
			pos:      buffer.NewPosition(0, 5),
			expected: buffer.NewPosition(0, 10),
		},
		{
			name:     "already at end of line",
			input:    "hello world",
			pos:      buffer.NewPosition(0, 10),
			expected: buffer.NewPosition(0, 10),
		},
		{
			name:     "empty line",
			input:    "\nhello",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 0),
		},
		{
			name:     "multiple lines",
			input:    "hello\nworld",
			pos:      buffer.NewPosition(1, 2),
			expected: buffer.NewPosition(1, 4),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			lines := strings.Split(tt.input, "\n")
			motion := NewEndOfLineMotion()
			result := motion.Calculate(lines, tt.pos)
			s.Equal(tt.expected, result, "positions should match")
		})
	}
}

func TestLineMotionSuite(t *testing.T) {
	suite.Run(t, new(LineMotionTestSuite))
}
