package io

import (
	"io"
	"os"
)

type StdinSource struct {
	reader io.Reader
	name   string
}

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

type FileSource struct {
	reader io.Reader
	name   string
}

func NewFileSource(reader io.Reader, name string) *FileSource {
	return &FileSource{
		reader: reader,
		name:   name,
	}
}

func (s *FileSource) Read() ([]byte, error) {
	return io.ReadAll(s.reader)
}

func (s *FileSource) Name() string {
	return s.name
}

func (s *FileSource) Close() error {
	if closer, ok := s.reader.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
