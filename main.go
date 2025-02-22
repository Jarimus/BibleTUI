package main

import tea "github.com/charmbracelet/bubbletea"

func main() {
	mainMenu := newMainMenu()
	root := newRootScreen([]tea.Model{mainMenu})
	//root := newRootScreen([]tea.Model{newBibleScreen()})
	p := tea.NewProgram(root)
	p.Run()
}
