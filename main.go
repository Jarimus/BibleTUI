package main

import (
	"log"

	"github.com/Jarimus/BibleTUI/internal/api_query"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

// Global variables:
// Terminal width and height. Necessary for the viewport in reading mode when it spawns.
var window_width int
var window_height int

// struct for the data of the current Bible, its books and the current chapter being read.
type currentlyReading struct {
	translationName string
	translationID   string
	translationData api_query.TranslationData
	bookData        api_query.BookData
	chapterData     api_query.ChapterData
}

var current currentlyReading

func main() {

	// Get environmental variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize with the Finnish translation
	current.translationName = "Finnish New Testament"
	current.translationID = "c739534f6a23acb2-01"
	current.translationData = api_query.TranslationQuery(current.translationID)

	mainMenu := newMainMenu()

	// Root screen holds the other moddels in a "stack" (slice) and displays the one at the top.
	root := newRootScreen([]tea.Model{mainMenu})
	p := tea.NewProgram(root)
	p.Run()
}
