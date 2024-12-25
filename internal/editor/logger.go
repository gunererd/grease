package editor

import (
	"fmt"
	"log"
	"os"

	"github.com/gunererd/grease/internal/editor/types"
)

type FileLogger struct {
	logger *log.Logger
	file   *os.File
}

func NewFileLogger(filename, prefix string) (types.Logger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return &FileLogger{
		logger: log.New(file, prefix+": ", log.Ldate|log.Ltime|log.Lmicroseconds),
		file:   file,
	}, nil
}

func (l *FileLogger) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

func (l *FileLogger) Println(v ...interface{}) {
	l.logger.Println(v...)
}

func (l *FileLogger) Close() error {
	return l.file.Close()
}
