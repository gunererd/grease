package clipboard

import (
	"github.com/gunererd/grease/internal/register"
	"github.com/gunererd/grease/internal/types"
)

type ClipboardCommand interface {
	Execute(lines []string, pos types.Position, register *register.Register) ([]string, types.Position)
}
