package io

import (
	"io"
	"os"
)

// Source represents a source of content for the editor
type Source interface {
	// Read returns the entire content from the source
	Read() ([]byte, error)
	// Name returns an identifier for the source
	Name() string
	// Close releases any resources associated with the source
	Close() error
}

// StdinSource implements Source for standard input
type StdinSource struct {
	reader io.Reader
	name   string
}

// NewStdinSource creates a new source that reads from stdin
func NewStdinSource() *StdinSource {
	return &StdinSource{
		reader: os.Stdin,
		name:   "stdin",
	}
}

func (s *StdinSource) Read() ([]byte, error) {
	return io.ReadAll(s.reader)
}

func (s *StdinSource) Name() string {
	return s.name
}

func (s *StdinSource) Close() error {
	return nil // stdin doesn't need to be closed
}
