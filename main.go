package main

import (
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/editor"
	"github.com/gunererd/grease/internal/filemanager"
)

func main() {
	f, err := tea.LogToFile("debug.log", "DEBUG")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up logging: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	opts := editor.RegisterFlags()
	flag.Parse()

	e, err := editor.Initialize(*opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing editor: %v\n", err)
		os.Exit(1)
	}

	path := "."
	if opts.Filename != "" {
		path = opts.Filename
	}

	fm := filemanager.New(path, e)
	if err := fm.LoadDirectory(); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading directory: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(e,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
		tea.WithMouseAllMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
