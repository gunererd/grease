package highlight

import (
	"sort"
	"sync"

	"github.com/gunererd/grease/internal/types"
)

// Manager implements types.HighlightManager interface
type Manager struct {
	highlights map[int]types.Highlight
	nextID     int
	mu         sync.RWMutex
}

// NewManager creates a new highlight manager
func NewManager() types.HighlightManager {
	return &Manager{
		highlights: make(map[int]types.Highlight),
		nextID:     1,
	}
}

// Add adds a new highlight and returns its ID
func (m *Manager) Add(h types.Highlight) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Create a new highlight with the next ID
	newHighlight := NewHighlight(
		h.GetStartPosition(),
		h.GetEndPosition(),
		h.GetType(),
		h.GetPriority(),
	)
	
	// Set the ID
	id := m.nextID
	m.highlights[id] = newHighlight
	m.nextID++
	return id
}

// Remove removes a highlight by its ID
func (m *Manager) Remove(id int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.highlights, id)
}

// Clear removes all highlights of a given type
func (m *Manager) Clear(highlightType types.HighlightType) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, h := range m.highlights {
		if h.GetType() == highlightType {
			delete(m.highlights, id)
		}
	}
}

// Get returns a highlight by its ID
func (m *Manager) Get(id int) (types.Highlight, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	h, ok := m.highlights[id]
	return h, ok
}

// GetForLine returns all highlights that intersect with the given line
func (m *Manager) GetForLine(line int) []types.Highlight {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []types.Highlight
	for _, h := range m.highlights {
		if line >= h.GetStartPosition().Line() && line <= h.GetEndPosition().Line() {
			result = append(result, h)
		}
	}

	// Sort by priority (highest first)
	sort.Slice(result, func(i, j int) bool {
		return result[i].GetPriority() > result[j].GetPriority()
	})

	return result
}

// GetForPosition returns all highlights that contain the given position
func (m *Manager) GetForPosition(pos types.Position) []types.Highlight {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []types.Highlight
	for _, h := range m.highlights {
		if h.Contains(pos) {
			result = append(result, h)
		}
	}

	// Sort by priority (highest first)
	sort.Slice(result, func(i, j int) bool {
		return result[i].GetPriority() > result[j].GetPriority()
	})

	return result
}

// Update updates an existing highlight
func (m *Manager) Update(id int, h types.Highlight) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.highlights[id]; !exists {
		return false
	}

	// Create new highlight with same ID but updated positions
	newHighlight := NewHighlight(
		h.GetStartPosition(),
		h.GetEndPosition(),
		h.GetType(),
		h.GetPriority(),
	)
	
	m.highlights[id] = newHighlight
	return true
}
