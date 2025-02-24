package main

import tea "github.com/charmbracelet/bubbletea"

// Global variables:
// Terminal width and height
var window_width int
var window_height int

type currentlyReading struct {
	translationName string
	translationID   string
	bookData        string
	chapterData     int
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
