package clipboard

import (
	"strings"
	"testing"

	"github.com/gunererd/grease/internal/buffer"
	"github.com/gunererd/grease/internal/register"
	"github.com/gunererd/grease/internal/types"
	"github.com/stretchr/testify/suite"
)

type PasteTestSuite struct {
	suite.Suite
	register *register.Register
}

func (s *PasteTestSuite) SetupTest() {
	s.register = register.NewRegister()
}

func (s *PasteTestSuite) TestPasteCommand() {
	tests := []struct {
		name          string
		input         string
		pasteText     string
		pos           types.Position
		before        bool
		expectedLines []string
		expectedPos   types.Position
	}{
		{
			name:          "paste single word after cursor",
			input:         "hello world",
			pasteText:     "test",
			pos:           buffer.NewPosition(0, 5),
			before:        false,
			expectedLines: []string{"hello testworld"},
			expectedPos:   buffer.NewPosition(0, 9),
		},
		{
			name:          "paste single word before cursor",
			input:         "hello world",
			pasteText:     "test",
			pos:           buffer.NewPosition(0, 5),
			before:        true,
			expectedLines: []string{"hellotest world"},
			expectedPos:   buffer.NewPosition(0, 8),
		},
		{
			name:          "paste multiline text",
			input:         "hello\nworld",
			pasteText:     "test\ntext",
			pos:           buffer.NewPosition(0, 4),
			before:        false,
			expectedLines: []string{"hellotest", "text", "world"},
			expectedPos:   buffer.NewPosition(1, 3),
		},
		{
			name:          "paste at empty line",
			input:         "\n",
			pasteText:     "test",
			pos:           buffer.NewPosition(0, 0),
			before:        false,
			expectedLines: []string{"test", ""},
			expectedPos:   buffer.NewPosition(0, 3),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			lines := strings.Split(tt.input, "\n")
			s.register.Set(tt.pasteText)
			cmd := NewPasteCommand(tt.before)

			resultLines, resultPos := cmd.Execute(lines, tt.pos, s.register)

			s.Equal(tt.expectedLines, resultLines, "lines should match after paste")
			s.Equal(tt.expectedPos, resultPos, "cursor position should match")
		})
	}
}

func TestPasteSuite(t *testing.T) {
	suite.Run(t, new(PasteTestSuite))
}
