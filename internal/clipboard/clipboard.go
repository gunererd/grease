package clipboard

import "sync"

// Manager handles storing and retrieving clipboard content
type Manager struct {
	content string
	mu      sync.RWMutex
}

// New creates a new clipboard manager
func New() *Manager {
	return &Manager{}
}

// Set stores text in the clipboard
func (m *Manager) Set(text string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.content = text
}

// Get retrieves text from the clipboard
func (m *Manager) Get() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.content
}
