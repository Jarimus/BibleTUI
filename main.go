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
	translationName string
	translationID   string
	translationData api_query.TranslationData
	bookData        api_query.BookData
	chapterData     api_query.ChapterData
}

var current currentlyReading

func main() {

	current.translationName = "Finnish New Testament"
	current.translationID = "c739534f6a23acb2-01"
	current.translationData = api_query.TranslationQuery(current.translationID)

	mainMenu := newMainMenu()
	root := newRootScreen([]tea.Model{mainMenu})
	p := tea.NewProgram(root)
	p.Run()
}
