package main

import tea "github.com/charmbracelet/bubbletea"

func main() {
	p := tea.NewProgram(newMainMenu())
	p.Run()
}
