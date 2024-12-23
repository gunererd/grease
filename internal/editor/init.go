package editor

import (
	"flag"
	"log"
	"net/http"

	"github.com/gunererd/grease/internal/editor/buffer"
	"github.com/gunererd/grease/internal/editor/handler"
	"github.com/gunererd/grease/internal/editor/highlight"
	"github.com/gunererd/grease/internal/editor/history"
	"github.com/gunererd/grease/internal/editor/hook"
	ioManager "github.com/gunererd/grease/internal/editor/io"
	"github.com/gunererd/grease/internal/editor/keytree"
	"github.com/gunererd/grease/internal/editor/register"
	"github.com/gunererd/grease/internal/editor/ui"
)

type InitOptions struct {
	Filename string
	Profile  bool
}

func RegisterFlags() *InitOptions {
	opts := &InitOptions{}

	flag.BoolVar(&opts.Profile, "profile", false, "Enable pprof profiling on :6060")
	flag.StringVar(&opts.Filename, "f", "", "Input file path")
	flag.StringVar(&opts.Filename, "file", "", "Input file path")

	return opts
}

func Initialize(opts InitOptions) (*Editor, error) {
	// Setup logging

	if opts.Profile {
		go func() {
			log.Println("Starting pprof server on :6060")
			http.ListenAndServe(":6060", nil)
		}()
	}

	kt := keytree.NewKeyTree()
	manager := ioManager.New(ioManager.NewStdinSource(), ioManager.NewStdoutSink())
	highlightManager := highlight.New()
	buffer := buffer.New()
	statusLine := ui.NewStatusLine()
	viewport := ui.NewViewport(0, 0)
	viewport.SetHighlightManager(highlightManager)
	register := register.NewRegister()
	historyManager := history.New(100)
	hookManager := hook.NewManager()
	executor := handler.NewCommandExecutor(historyManager, hookManager)
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
	)

	if opts.Filename != "" {
		if err := editor.LoadFromFile(opts.Filename); err != nil {
			return nil, err
		}
	}

	return editor, nil
}
