package main

import (
	"github.com/Jarimus/BibleTUI/internal/api_query"
	tea "github.com/charmbracelet/bubbletea"
)

// Global variables:
// Terminal width and height
var window_width int
var window_height int

type currentlyReading struct {
	translationName   string
	translationID     string
	translationData   api_query.TranslationData
	bookData          api_query.BookData
	chapterData       api_query.ChapterData
	currentBookStr    string
	currentChapterInt int
	currentChapterStr string
}

var current currentlyReading

func main() {

	current.translationName = "TEST"

	mainMenu := newMainMenu()
	root := newRootScreen([]tea.Model{mainMenu})
	//root := newRootScreen([]tea.Model{newBibleScreen()})
	p := tea.NewProgram(root)
	p.Run()
}
