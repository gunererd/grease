package handler

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
			expected: buffer.NewPosition(0, 5), // moves to the comma
			bigWord:  false,
		},
		{
			name:     "big word movement",
			input:    "hello, world!",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 7), // skips punctuation
			bigWord:  true,
		},
		{
			name:     "basic word movement skips whitespaces",
			input:    "    hello",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 4), // skips whitespace
			bigWord:  false,
		},
		{
			name:     "big word movement skips whitespaces",
			input:    "    hello",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 4), // skips whitespace
			bigWord:  true,
		},
		{
			name:     "end of line",
			input:    "hello",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 4),
			bigWord:  false,
		},
		{
			name:     "end of line for big word",
			input:    "hello",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 4),
			bigWord:  true,
		},
		{
			name:     "end of line to next line",
			input:    "hello\nworld",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(1, 0),
			bigWord:  false,
		},
		{
			name:     "next line with multiple spaces",
			input:    "hello\n    world",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(1, 4),
			bigWord:  false,
		},
		{
			name:     "multiple spaces",
			input:    "hello    world",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 9),
			bigWord:  false,
		},
		{
			name:     "empty buffer",
			input:    "",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 0),
			bigWord:  false,
		},
		{
			name:     "single character",
			input:    "a",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 0),
			bigWord:  false,
		},
		{
			name:     "start from middle of word",
			input:    "hello world",
			pos:      buffer.NewPosition(0, 2),
			expected: buffer.NewPosition(0, 6),
			bigWord:  false,
		},
		{
			name:     "mixed punctuation and spaces",
			input:    "hello...   world",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 5),
			bigWord:  false,
		},
		{
			name:     "big word over punctuation",
			input:    "hello...   world",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 11),
			bigWord:  true,
		},
		{
			name:     "at last position",
			input:    "hello",
			pos:      buffer.NewPosition(0, 4),
			expected: buffer.NewPosition(0, 4),
			bigWord:  false,
		},
		{
			name:     "punctuation",
			input:    "hello, world",
			pos:      buffer.NewPosition(0, 5),
			expected: buffer.NewPosition(0, 7),
			bigWord:  false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			buf := buffer.New()
			buf.LoadFromReader(strings.NewReader(tt.input))

			motion := NewWordMotion(tt.bigWord)

			result := motion.Calculate(buf, tt.pos)

			s.Equal(tt.expected, result, "positions should match")
		})
	}
}

func TestWordMotionSuite(t *testing.T) {
	suite.Run(t, new(WordMotionTestSuite))
}

type WordEndMotionTestSuite struct {
	suite.Suite
}

func (s *WordEndMotionTestSuite) TestWordEndMotion() {
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
			name:     "big word end movement",
			input:    "hello, world!",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 5),
			bigWord:  true,
		},
		{
			name:     "word end with leading spaces",
			input:    "    hello",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 8),
			bigWord:  false,
		},
		{
			name:     "big word end with leading spaces",
			input:    "    hello",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 8),
			bigWord:  true,
		},
		{
			name:     "empty buffer",
			input:    "",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 0),
			bigWord:  false,
		},
		{
			name:     "single character",
			input:    "a",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 0),
			bigWord:  false,
		},
		{
			name:     "start from middle of word",
			input:    "hello world",
			pos:      buffer.NewPosition(0, 2),
			expected: buffer.NewPosition(0, 4),
			bigWord:  false,
		},
		{
			name:     "mixed punctuation and spaces",
			input:    "hello...   world",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 4),
			bigWord:  false,
		},
		{
			name:     "big word over punctuation",
			input:    "hello...   world",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 7),
			bigWord:  true,
		},
		{
			name:     "end of line to next line",
			input:    "hello\nworld",
			pos:      buffer.NewPosition(0, 4),
			expected: buffer.NewPosition(1, 4),
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
			name:     "end of line to next line for big word",
			input:    "hello\nworld",
			pos:      buffer.NewPosition(0, 4),
			expected: buffer.NewPosition(1, 4),
			bigWord:  true,
		},
		{
			name:     "next line with multiple spaces for big word",
			input:    "hello\n    world",
			pos:      buffer.NewPosition(0, 4),
			expected: buffer.NewPosition(1, 8),
			bigWord:  true,
		},
		{
			name:     "jump to next line while on punctuation",
			input:    "hello,\n    world",
			pos:      buffer.NewPosition(0, 5),
			expected: buffer.NewPosition(1, 8),
			bigWord:  false,
		},
		{
			name:     "jump to next line while on punctuation for big word",
			input:    "hello,\n    world",
			pos:      buffer.NewPosition(0, 5),
			expected: buffer.NewPosition(1, 8),
			bigWord:  true,
		},
		{
			name:     "jump to next line punctuatuion",
			input:    "hello\n  ,   world",
			pos:      buffer.NewPosition(0, 4),
			expected: buffer.NewPosition(1, 2),
			bigWord:  false,
		},
		{
			name:     "jump to next line punctuatuion for big word",
			input:    "hello\n  ,   world",
			pos:      buffer.NewPosition(0, 4),
			expected: buffer.NewPosition(1, 2),
			bigWord:  true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			buf := buffer.New()
			buf.LoadFromReader(strings.NewReader(tt.input))

			motion := NewWordEndMotion(tt.bigWord)
			result := motion.Calculate(buf, tt.pos)

			s.Equal(tt.expected, result, "positions should match")
		})
	}
}

func TestWordEndMotionSuite(t *testing.T) {
	suite.Run(t, new(WordEndMotionTestSuite))
}

type WordBackMotionTestSuite struct {
	suite.Suite
}

func (s *WordBackMotionTestSuite) TestWordBackMotion() {
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
			name:     "big word back movement",
			input:    "hello, world!",
			pos:      buffer.NewPosition(0, 7),
			expected: buffer.NewPosition(0, 0),
			bigWord:  true,
		},
		{
			name:     "word back with trailing spaces",
			input:    "hello    world",
			pos:      buffer.NewPosition(0, 9),
			expected: buffer.NewPosition(0, 0),
			bigWord:  false,
		},
		{
			name:     "empty buffer",
			input:    "",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 0),
			bigWord:  false,
		},
		{
			name:     "single character",
			input:    "a",
			pos:      buffer.NewPosition(0, 0),
			expected: buffer.NewPosition(0, 0),
			bigWord:  false,
		},
		{
			name:     "from middle of word",
			input:    "hello world",
			pos:      buffer.NewPosition(0, 8),
			expected: buffer.NewPosition(0, 6),
			bigWord:  false,
		},
		{
			name:     "mixed punctuation and spaces",
			input:    "hello...   world",
			pos:      buffer.NewPosition(0, 11),
			expected: buffer.NewPosition(0, 5),
			bigWord:  false,
		},
		{
			name:     "big word over punctuation",
			input:    "hello...   world",
			pos:      buffer.NewPosition(0, 11),
			expected: buffer.NewPosition(0, 0),
			bigWord:  true,
		},
		{
			name:     "start of line to previous line",
			input:    "hello\nworld",
			pos:      buffer.NewPosition(1, 0),
			expected: buffer.NewPosition(0, 0),
			bigWord:  false,
		},
		{
			name:     "multiple spaces to previous line",
			input:    "hello\n    world",
			pos:      buffer.NewPosition(1, 4),
			expected: buffer.NewPosition(0, 0),
			bigWord:  false,
		},
		{
			name:     "jump to punctuation",
			input:    "hello\n, world",
			pos:      buffer.NewPosition(1, 2),
			expected: buffer.NewPosition(1, 0),
			bigWord:  false,
		},
		{
			name:     "start of first line",
			input:    "hello world",
			pos:      buffer.NewPosition(0, 0),
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
		{
			name:     "multiple punctuation",
			input:    "hello... world",
			pos:      buffer.NewPosition(0, 7),
			expected: buffer.NewPosition(0, 5),
			bigWord:  false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			buf := buffer.New()
			buf.LoadFromReader(strings.NewReader(tt.input))

			motion := NewWordBackMotion(tt.bigWord)
			result := motion.Calculate(buf, tt.pos)

			s.Equal(tt.expected, result, "positions should match")
		})
	}
}

func TestWordBackMotionSuite(t *testing.T) {
	suite.Run(t, new(WordBackMotionTestSuite))
}
