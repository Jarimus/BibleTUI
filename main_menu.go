package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// type Model interface methods: Init() tea.Cmd, Update(tea.Msg) (tea.Model, tea.Cmd), View
type mainMenuModel struct {
	options []option
	cursor  int
}

type option struct {
	text    string
	command tea.Cmd
}

func newMainMenu() mainMenuModel {

	newOptions := []option{
		{text: "Random verse", command: func() tea.Msg { return nil }},
		{text: "Read the Bible", command: func() tea.Msg { return nil }},
		{text: "Exit the Bible", command: tea.Quit},
	}

	return mainMenuModel{
		options: newOptions,
		cursor:  0,
	}
}

func (m mainMenuModel) Init() tea.Cmd {
	return nil
}

func (m mainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.cursor = (m.cursor - 1) % len(m.options)
			if m.cursor < 0 {
				m.cursor += len(m.options)
			}
			return m, nil
		case "down":
			m.cursor = (m.cursor + 1) % len(m.options)
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, m.options[m.cursor].command
		}

	}
	return m, nil
}

func (m mainMenuModel) View() string {

	var view string

	welcomeMsg := "* Welcome to BibleTui! *"
	topBottomBar := strings.Repeat("*", len(welcomeMsg))

	view = fmt.Sprintf("%s\n%s\n%s\n", topBottomBar, welcomeMsg, topBottomBar)

	for i, option := range m.options {
		if m.cursor == i {
			view = view + "-> " + option.text + "\n"
		} else {
			view = view + "   " + option.text + "\n"
		}
	}

	return view
}
