package filemanager

import (
	"fmt"
	"log"
	"os"

	eTypes "github.com/gunererd/grease/internal/editor/types"
	"github.com/gunererd/grease/internal/filemanager/directory"
	"github.com/gunererd/grease/internal/filemanager/operation"
	"github.com/gunererd/grease/internal/filemanager/types"
	"github.com/gunererd/grease/internal/filemanager/view"
)

type options struct {
	LogFile string
}

type Option func(*options)

func WithLog(filename string) Option {
	return func(o *options) {
		o.LogFile = filename
	}
}

func Initialize(editor eTypes.Editor, opts ...Option) (types.FileManager, error) {
	options := options{}

	for _, opt := range opts {
		opt(&options)
	}

	var logger types.Logger
	if options.LogFile != "" {
		fileLogger, err := NewFileLogger(options.LogFile, "FILEMANAGER")
		if err != nil {
			return nil, fmt.Errorf("failed to initialize logger: %w", err)
		}
		logger = fileLogger
	} else {
		logger = log.New(os.Stderr, "FILEMANAGER: ", log.Ldate|log.Ltime|log.Lmicroseconds)
	}

	dirManager := directory.NewDirectoryManager("", logger)
	opManager := operation.NewOperationManager(dirManager, logger)
	view := view.New(editor)

	fm := New(
		dirManager,
		opManager,
		view,
		editor,
		logger,
	)

	return fm, nil
}
