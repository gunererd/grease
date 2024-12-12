package motion

import "github.com/gunererd/grease/internal/types"

type Motion interface {
	// Calculate returns the target position without side effects
	Calculate(lines []string, pos types.Position) types.Position

	// GetRange returns the range of text affected by this motion
	// GetRange(text [][]rune, pos types.Position) (from types.Position, to types.Position)
}
