package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Jarimus/BibleTUI/internal/api_query"
	styles "github.com/Jarimus/BibleTUI/internal/styles"
)

// type Model interface methods: Init() tea.Cmd, Update(tea.Msg) (tea.Model, tea.Cmd), View
type mainMenuModel struct {
	randomVerseVP viewport.Model
	options       []option
	cursor        int
}

type option struct {
	text    string
	command tea.Cmd
}

func newMainMenu() mainMenuModel {

	var m mainMenuModel

	newOptions := []option{
		{text: "Random verse", command: m.newRandomVerse},
		{text: "Read the Bible", command: func() tea.Msg { return newBibleScreen() }}, // Pass initial width and height
		{text: "Exit BibleTUI", command: tea.Quit},
	}

	m = mainMenuModel{
		randomVerseVP: viewport.New(window_width, window_height-lipgloss.Height(m.menuView())),
		options:       newOptions,
		cursor:        0,
	}

	m.randomVerseVP.SetContent("Loading...")

	return m
}

func (m mainMenuModel) Init() tea.Cmd {
	return m.newRandomVerse
}

func (m mainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case api_query.RandomQuery:
		m.applyRandomVerse(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.cursor = (m.cursor - 1 + len(m.options)) % len(m.options)
			return m, nil
		case "down":
			m.cursor = (m.cursor + 1) % len(m.options)
			return m, nil
		case "enter":
			return m, m.options[m.cursor].command
		}
	case tea.WindowSizeMsg:
		m.randomVerseVP.Width = msg.Width
		m.randomVerseVP.Height = msg.Height - lipgloss.Height(m.menuView())
	}
	return m, nil
}

func (m mainMenuModel) View() string {
	emptyFullLine := strings.Repeat(" ", window_width)

	return lipgloss.JoinVertical(lipgloss.Center, emptyFullLine, m.menuView(), m.randomVerseVP.View())
}

func (m mainMenuModel) menuView() string {
	welcomeMsg := "* Welcome to BibleTUI! *"
	topBottomBar := strings.Repeat("*", len(welcomeMsg))

	var options string
	for i, option := range m.options {
		if m.cursor == i {
			choiceText := fmt.Sprint(styles.GreenText.Render(option.text))
			options += choiceText + "\n"
		} else {
			options += option.text + "\n"
		}
	}
	return lipgloss.JoinVertical(lipgloss.Center, topBottomBar, welcomeMsg, topBottomBar, options)
}

func (m mainMenuModel) newRandomVerse() tea.Msg {

	query := api_query.GetRandomVerse()

	return query

}

func (m *mainMenuModel) applyRandomVerse(query api_query.RandomQuery) {
	// Apply the new random verse
	bibleText, _ := strings.CutSuffix(query.RandomVerse.Text, "\n")

	s := fmt.Sprintf("%s\n- %s %d:%d (%s)",
		bibleText,
		query.RandomVerse.Book,
		query.RandomVerse.Chapter,
		query.RandomVerse.Verse,
		query.Translation.Name)
	m.randomVerseVP.SetContent(s)
}
