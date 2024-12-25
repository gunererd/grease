package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/editor"
	"github.com/gunererd/grease/internal/filemanager"
)

func main() {

	e, err := editor.Initialize(editor.WithLog("debug.log"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing editor: %v\n", err)
		os.Exit(1)
	}

	fm, err := filemanager.Initialize(e, filemanager.WithLog("debug.log"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing filemanager: %v\n", err)
		os.Exit(1)
	}

	// Get initial path (current directory if not specified)
	initialPath := "."
	if len(os.Args) > 1 {
		initialPath = os.Args[1]
	}

	// Load initial directory
	if err := fm.LoadDirectory(initialPath); err != nil {
		fmt.Printf("Error loading directory: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(fm,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
		tea.WithMouseAllMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
