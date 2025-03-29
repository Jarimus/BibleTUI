package main

import (
	"log"

	"github.com/Jarimus/BibleTUI/internal/api_query"
	"github.com/Jarimus/BibleTUI/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

// Global variables:
// Terminal width and height. Necessary for the viewport in reading mode when it spawns.
var window_width int
var window_height int

// struct for the data of the current Bible, its books and the current chapter being read.
var apiCfg config.Config

func main() {

	println("Loading...")

	// Get settings
	apiCfg, err := loadSettings()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize with the current  translation
	apiCfg.CurrentlyReading.TranslationData = api_query.TranslationQuery(apiCfg.CurrentlyReading.TranslationID)

	mainMenu := newMainMenu()

	// Root screen holds the other models in a "stack" (slice) and displays the one at the top.
	root := newRootScreen([]tea.Model{mainMenu})
	p := tea.NewProgram(root)
	_, err = p.Run()
	if err != nil {
		log.Fatalf("error starting the program: %s", err)
	}
}
