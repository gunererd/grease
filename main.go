package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/buffer"
	"github.com/gunererd/grease/internal/editor"
	"github.com/gunererd/grease/internal/editor/handler"
	"github.com/gunererd/grease/internal/highlight"
	ioManager "github.com/gunererd/grease/internal/io"
	"github.com/gunererd/grease/internal/keytree"
	"github.com/gunererd/grease/internal/ui"
)

func main() {

	f, err := tea.LogToFile("debug.log", "DEBUG")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	profile := flag.Bool("profile", false, "Enable pprof profiling on :6060")
	filename := flag.String("f", "", "Input file path")
	flag.StringVar(filename, "file", "", "Input file path")
	flag.Parse()
	if *profile {
		go func() {
			log.Println("Starting pprof server on :6060")
			http.ListenAndServe(":6060", nil)
		}()
	}

	kt := keytree.NewKeyTree()
	manager := ioManager.NewManager(ioManager.NewStdinSource(), ioManager.NewStdoutSink())
	highlightManager := highlight.NewManager()
	buffer := buffer.New()
	statusLine := ui.NewStatusLine()
	viewport := ui.NewViewport(0, 0)
	viewport.SetHighlightManager(highlightManager)
	historyManager := handler.NewHistoryManager(100)
	operationManager := handler.NewOperationManager(historyManager)
	m := editor.New(manager, buffer, statusLine, viewport, highlightManager, kt, historyManager, operationManager)

	// Check for file input first
	if *filename != "" {
		if err := m.LoadFromFile(*filename); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Load content from stdin if it's not a terminal
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			if err := m.LoadFromStdin(); err != nil {
				fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
				os.Exit(1)
			}
		}
	}

	p := tea.NewProgram(m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
		tea.WithMouseAllMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
