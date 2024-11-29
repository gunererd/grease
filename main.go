package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/buffer"
	"github.com/gunererd/grease/internal/editor"
	"github.com/gunererd/grease/internal/highlight"
	ioManager "github.com/gunererd/grease/internal/io"
	"github.com/gunererd/grease/internal/ui"
)

func main() {
	manager := ioManager.NewManager(ioManager.NewStdinSource(), ioManager.NewStdoutSink())
	highlightManager := highlight.NewManager()
	buffer := buffer.New()
	statusLine := ui.NewStatusLine()
	viewport := ui.NewViewport(0, 0)
	viewport.SetHighlightManager(highlightManager)
	m := editor.New(manager, buffer, statusLine, viewport, highlightManager)

	// Load content from stdin if it's not a terminal
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		if err := m.LoadFromStdin(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
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
