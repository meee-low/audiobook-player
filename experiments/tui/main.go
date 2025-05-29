package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/meee-low/audiobook-player/internal/tui"
	"log"
)

func main() {
	mainModel := tui.NewModel()
	program := tea.NewProgram(&mainModel, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		log.Fatalf("Error running bubbletea program %v", err)
	}
}
