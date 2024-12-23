package io

import (
	"io"
	"os"
)

type StdoutSink struct {
	writer io.Writer
}

// NewStdoutSink creates a new sink that writes to stdout
func NewStdoutSink() *StdoutSink {
	return &StdoutSink{
		writer: os.Stdout,
	}
}

func (s *StdoutSink) Write(data []byte) error {
	_, err := s.writer.Write(data)
	return err
}

func (s *StdoutSink) Flush() error {
	if f, ok := s.writer.(*os.File); ok {
		return f.Sync()
	}
	return nil
}

func (s *StdoutSink) Close() error {
	return nil // stdout doesn't need to be closed
}
