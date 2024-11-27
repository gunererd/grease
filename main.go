package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunererd/grease/internal/editor"
)

func main() {
	f, err := tea.LogToFile("debug.log", "DEBUG")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	p := tea.NewProgram(
		editor.New(),
		tea.WithAltScreen(),
	)

	log.Println("Starting")
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
