package highlight

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/gunererd/grease/internal/editor/types"
)

// Style holds the styles for different types of highlights
type Style struct {
	Visual lipgloss.Style
	Search lipgloss.Style
}

// NewStyle creates a new Style with default values
func NewStyle() *Style {
	return &Style{
		Visual: lipgloss.NewStyle().
			Background(lipgloss.Color("#444444")).
			Foreground(lipgloss.Color("#ffffff")),

		Search: lipgloss.NewStyle().
			Background(lipgloss.Color("#755800")).
			Foreground(lipgloss.Color("#ffffff")),
	}
}

// GetStyle returns the appropriate style for a given highlight type
func (s *Style) GetStyle(t types.HighlightType) lipgloss.Style {
	switch t {
	case types.VisualHighlight:
		return s.Visual
	case types.SearchHighlight:
		return s.Search
	default:
		return lipgloss.NewStyle()
	}
}
