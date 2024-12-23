package clipboard

import (
	"github.com/gunererd/grease/internal/editor/register"
	"github.com/gunererd/grease/internal/editor/types"
)

type ClipboardCommand interface {
	Execute(lines []string, pos types.Position, register *register.Register) ([]string, types.Position)
}
