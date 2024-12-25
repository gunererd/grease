package editor

import (
	"fmt"
	"log"
	"os"

	"github.com/gunererd/grease/internal/editor/buffer"
	"github.com/gunererd/grease/internal/editor/handler"
	"github.com/gunererd/grease/internal/editor/highlight"
	"github.com/gunererd/grease/internal/editor/history"
	"github.com/gunererd/grease/internal/editor/hook"
	ioManager "github.com/gunererd/grease/internal/editor/io"
	"github.com/gunererd/grease/internal/editor/keytree"
	"github.com/gunererd/grease/internal/editor/register"
	"github.com/gunererd/grease/internal/editor/types"
	"github.com/gunererd/grease/internal/editor/ui"
)

type options struct {
	filename string
	profile  bool
	logFile  string
}

type Option func(*options)

func WithFilename(filename string) Option {
	return func(o *options) {
		o.filename = filename
	}
}

func WithProfiling(enabled bool) Option {
	return func(o *options) {
		o.profile = enabled
	}
}

func WithLog(logFile string) Option {
	return func(o *options) {
		o.logFile = logFile
	}
}

func Initialize(opts ...Option) (*Editor, error) {
	options := &options{
		filename: "",
		profile:  false,
		logFile:  "",
	}

	for _, opt := range opts {
		opt(options)
	}

	var logger types.Logger
	if options.logFile != "" {
		fileLogger, err := NewFileLogger(options.logFile, "EDITOR")
		if err != nil {
			return nil, fmt.Errorf("failed to initialize logger: %w", err)
		}
		logger = fileLogger
	} else {
		logger = log.New(os.Stderr, "EDITOR: ", log.Ldate|log.Ltime|log.Lmicroseconds)
	}

	kt := keytree.NewKeyTree()
	manager := ioManager.New(ioManager.NewStdinSource(), ioManager.NewStdoutSink())
	highlightManager := highlight.New()
	buffer := buffer.New(logger)
	statusLine := ui.NewStatusLine()
	viewport := ui.NewViewport(0, 0)
	viewport.SetHighlightManager(highlightManager)
	register := register.NewRegister()
	historyManager := history.New(100)
	hookManager := hook.NewManager()
	executor := handler.NewCommandExecutor(historyManager, hookManager, logger)

	editor := New(
		manager,
		buffer,
		statusLine,
		viewport,
		highlightManager,
		kt,
		historyManager,
		register,
		executor,
		hookManager,
		logger,
	)

	if options.filename != "" {
		if err := editor.LoadFromFile(options.filename); err != nil {
			return nil, err
		}
	}

	return editor, nil
}
