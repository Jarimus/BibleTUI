package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Jarimus/BibleTUI/internal/api_query"
	styles "github.com/Jarimus/BibleTUI/internal/styles"
)

// type Model interface methods: Init() tea.Cmd, Update(tea.Msg) (tea.Model, tea.Cmd), View
type mainMenuModel struct {
	textField string
	options   []option
	cursor    int
}

type option struct {
	text    string
	command tea.Cmd
}

func newMainMenu() mainMenuModel {

	newOptions := []option{
		{text: "Random verse", command: func() tea.Msg { return updateRandomVerse{} }},
		{text: "Read the Bible", command: func() tea.Msg { return newBibleScreen() }},
		{text: "Exit BibleTUI", command: tea.Quit},
	}

	m := mainMenuModel{
		textField: "",
		options:   newOptions,
		cursor:    0,
	}

	m.newRandomVerse()

	return m
}

func (m mainMenuModel) Init() tea.Cmd {
	return nil
}

func (m mainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case updateRandomVerse:
		m.newRandomVerse()
		return m, nil
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

	view := "\n" + m.textField + "\n"

	welcomeMsg := "* Welcome to BibleTui! *"
	topBottomBar := strings.Repeat("*", len(welcomeMsg))

	view += fmt.Sprintf("%s\n%s\n%s\n", topBottomBar, welcomeMsg, topBottomBar)

	for i, option := range m.options {
		if m.cursor == i {
			choiceText := fmt.Sprint(styles.GreenText.Render(option.text))
			view = strings.Join([]string{view, choiceText, "\n"}, "")
		} else {
			view += option.text + "\n"
		}
	}

	return view
}

type updateRandomVerse struct{}

func (m *mainMenuModel) newRandomVerse() {

	query := api_query.GetRandomVerse()

	// Apply the new random verse
	m.textField = fmt.Sprintf(`%s
	- %s %d:%d (%s)
	`,
		query.RandomVerse.Text,
		query.RandomVerse.Book,
		query.RandomVerse.Chapter,
		query.RandomVerse.Verse,
		query.Translation.Name)

	// Tried out red text styling
	// m.textField = styles.RedText.Render(m.textField)

}
