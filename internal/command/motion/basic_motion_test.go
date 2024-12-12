package motion

import (
	"strings"
	"testing"

	"github.com/gunererd/grease/internal/buffer"
	"github.com/gunererd/grease/internal/types"
	"github.com/stretchr/testify/suite"
)

type BasicMotionTestSuite struct {
	suite.Suite
}

func (s *BasicMotionTestSuite) TestLeftMotion() {
	tests := []struct {
		name     string
		input    string
		pos      types.Position
		expected types.Position
	}{
		{
			name:     "basic left movement",
			input:    "hello",
			pos:      buffer.NewPosition(0, 1),
			expected: buffer.NewPosition(0, 0),
		},
		{
			name:     "left at start of line",
			input:    "hello",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 0),
		},
		{
			name:     "left with multiple lines",
			input:    "hello\nworld",
			pos:      buffer.NewPosition(1, 1),
			expected: buffer.NewPosition(1, 0),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			lines := strings.Split(tt.input, "\n")
			motion := NewLeftMotion()
			result := motion.Calculate(lines, tt.pos)
			s.Equal(tt.expected, result, "positions should match")
		})
	}
}

func (s *BasicMotionTestSuite) TestRightMotion() {
	tests := []struct {
		name     string
		input    string
		pos      types.Position
		expected types.Position
	}{
		{
			name:     "basic right movement",
			input:    "hello",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 1),
		},
		{
			name:     "right at end of line",
			input:    "hello",
			pos:      buffer.NewPosition(0, 4),
			expected: buffer.NewPosition(0, 4),
		},
		{
			name:     "right with multiple lines",
			input:    "hello\nworld",
			pos:      buffer.NewPosition(1, 0),
			expected: buffer.NewPosition(1, 1),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			lines := strings.Split(tt.input, "\n")
			motion := NewRightMotion()
			result := motion.Calculate(lines, tt.pos)
			s.Equal(tt.expected, result, "positions should match")
		})
	}
}

func (s *BasicMotionTestSuite) TestUpMotion() {
	tests := []struct {
		name     string
		input    string
		pos      types.Position
		expected types.Position
	}{
		{
			name:     "basic up movement",
			input:    "hello\nworld",
			pos:      buffer.NewPosition(1, 0),
			expected: buffer.NewPosition(0, 0),
		},
		{
			name:     "up at first line",
			input:    "hello\nworld",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 0),
		},
		{
			name:     "up to shorter line",
			input:    "hi\nworld",
			pos:      buffer.NewPosition(1, 4),
			expected: buffer.NewPosition(0, 1),
		},
		{
			name:     "up to empty line",
			input:    "\nworld",
			pos:      buffer.NewPosition(1, 4),
			expected: buffer.NewPosition(0, 0),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			lines := strings.Split(tt.input, "\n")
			motion := NewUpMotion()
			result := motion.Calculate(lines, tt.pos)
			s.Equal(tt.expected, result, "positions should match")
		})
	}
}

func (s *BasicMotionTestSuite) TestDownMotion() {
	tests := []struct {
		name     string
		input    string
		pos      types.Position
		expected types.Position
	}{
		{
			name:     "basic down movement",
			input:    "hello\nworld",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(1, 0),
		},
		{
			name:     "down at last line",
			input:    "hello\nworld",
			pos:      buffer.NewPosition(1, 0),
			expected: buffer.NewPosition(1, 0),
		},
		{
			name:     "down to shorter line",
			input:    "world\nhi",
			pos:      buffer.NewPosition(0, 4),
			expected: buffer.NewPosition(1, 1),
		},
		{
			name:     "down to empty line",
			input:    "world\n",
			pos:      buffer.NewPosition(0, 4),
			expected: buffer.NewPosition(1, 0),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			lines := strings.Split(tt.input, "\n")
			motion := NewDownMotion()
			result := motion.Calculate(lines, tt.pos)
			s.Equal(tt.expected, result, "positions should match")
		})
	}
}

func TestBasicMotionSuite(t *testing.T) {
	suite.Run(t, new(BasicMotionTestSuite))
}
