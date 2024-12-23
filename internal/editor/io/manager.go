package io

import (
	"fmt"
	"sync"
)

// Manager handles all IO operations for the editor
type Manager struct {
	source Source
	sink   Sink
	mu     sync.RWMutex
}

// New creates a new IO manager with the given source and sink
func New(source Source, sink Sink) *Manager {
	return &Manager{
		source: source,
		sink:   sink,
	}
}

// LoadContent reads content from the current source
func (m *Manager) LoadContent() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.source == nil {
		return nil, fmt.Errorf("no source configured")
	}

	return m.source.Read()
}

// SaveContent writes content to the current sink
func (m *Manager) SaveContent(content []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.sink == nil {
		return fmt.Errorf("no sink configured")
	}

	if err := m.sink.Write(content); err != nil {
		return fmt.Errorf("failed to write content: %w", err)
	}

	if err := m.sink.Flush(); err != nil {
		return fmt.Errorf("failed to flush content: %w", err)
	}

	return nil
}

// SetSource changes the current source
func (m *Manager) SetSource(source Source) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.source != nil {
		if err := m.source.Close(); err != nil {
			return fmt.Errorf("failed to close previous source: %w", err)
		}
	}

	m.source = source
	return nil
}

// SetSink changes the current sink
func (m *Manager) SetSink(sink Sink) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.sink != nil {
		if err := m.sink.Close(); err != nil {
			return fmt.Errorf("failed to close previous sink: %w", err)
		}
	}

	m.sink = sink
	return nil
}

// Close releases all resources
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error

	if m.source != nil {
		if err := m.source.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close source: %w", err))
		}
	}

	if m.sink != nil {
		if err := m.sink.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close sink: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing manager: %v", errs)
	}

	return nil
}
