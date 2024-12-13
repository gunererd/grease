package clipboard

import (
	"strings"
	"testing"

	"github.com/gunererd/grease/internal/buffer"
	"github.com/gunererd/grease/internal/command/motion"
	"github.com/gunererd/grease/internal/register"
	"github.com/gunererd/grease/internal/types"
	"github.com/stretchr/testify/suite"
)

type YankTestSuite struct {
	suite.Suite
	register *register.Register
}

func (s *YankTestSuite) SetupTest() {
	s.register = register.NewRegister()
}

func (s *YankTestSuite) TestYankCommand() {
	tests := []struct {
		name          string
		input         string
		pos           types.Position
		motion        motion.Motion
		expectedText  string
		expectedLines []string
		expectedPos   types.Position
	}{
		{
			name:          "yank single word",
			input:         "hello world",
			pos:           buffer.NewPosition(0, 0),
			motion:        motion.NewWordMotion(false),
			expectedText:  "hello ",
			expectedLines: []string{"hello world"},
			expectedPos:   buffer.NewPosition(0, 0),
		},
		{
			name:          "yank to end of line",
			input:         "hello world",
			pos:           buffer.NewPosition(0, 6),
			motion:        motion.NewEndOfLineMotion(),
			expectedText:  "world",
			expectedLines: []string{"hello world"},
			expectedPos:   buffer.NewPosition(0, 6),
		},
		{
			name:          "yank multiple lines",
			input:         "hello\nworld\ntest",
			pos:           buffer.NewPosition(0, 0),
			motion:        motion.NewEndOfBufferMotion(),
			expectedText:  "hello\nworld\nt",
			expectedLines: []string{"hello", "world", "test"},
			expectedPos:   buffer.NewPosition(0, 0),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			lines := strings.Split(tt.input, "\n")
			cmd := NewYankCommand(tt.motion)

			resultLines, resultPos := cmd.Execute(lines, tt.pos, s.register)

			s.Equal(tt.expectedText, s.register.Get(), "yanked text should match")
			s.Equal(tt.expectedLines, resultLines, "lines should remain unchanged")
			s.Equal(tt.expectedPos, resultPos, "position should remain unchanged")
		})
	}
}

func TestYankSuite(t *testing.T) {
	suite.Run(t, new(YankTestSuite))
}
