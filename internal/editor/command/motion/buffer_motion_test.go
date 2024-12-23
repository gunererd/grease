package motion

import (
	"strings"
	"testing"

	"github.com/gunererd/grease/internal/editor/buffer"
	"github.com/gunererd/grease/internal/editor/types"
	"github.com/stretchr/testify/suite"
)

type BufferMotionTestSuite struct {
	suite.Suite
}

func (s *BufferMotionTestSuite) TestStartOfBufferMotion() {
	tests := []struct {
		name     string
		input    string
		pos      types.Position
		expected types.Position
	}{
		{
			name:     "basic start of buffer",
			input:    "hello\nworld\ntest",
			pos:      buffer.NewPosition(1, 2),
			expected: buffer.NewPosition(0, 0),
		},
		{
			name:     "already at start",
			input:    "hello\nworld",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 0),
		},
		{
			name:     "empty buffer",
			input:    "",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 0),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			lines := strings.Split(tt.input, "\n")
			motion := NewStartOfBufferMotion()
			result := motion.Calculate(lines, tt.pos)
			s.Equal(tt.expected, result, "positions should match")
		})
	}
}

func (s *BufferMotionTestSuite) TestEndOfBufferMotion() {
	tests := []struct {
		name     string
		input    string
		pos      types.Position
		expected types.Position
	}{
		{
			name:     "basic end of buffer",
			input:    "hello\nworld\ntest",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(2, 0),
		},
		{
			name:     "already at end",
			input:    "hello\nworld",
			pos:      buffer.NewPosition(1, 0),
			expected: buffer.NewPosition(1, 0),
		},
		{
			name:     "empty buffer",
			input:    "",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 0),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			lines := strings.Split(tt.input, "\n")
			motion := NewEndOfBufferMotion()
			result := motion.Calculate(lines, tt.pos)
			s.Equal(tt.expected, result, "positions should match")
		})
	}
}

func TestBufferMotionSuite(t *testing.T) {
	suite.Run(t, new(BufferMotionTestSuite))
}
