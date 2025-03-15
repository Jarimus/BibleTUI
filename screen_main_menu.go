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

// tea.Model for the mainMenu
type mainMenuModel struct {
	randomVerseVP viewport.Model
	options       []option
	choiceIndex   int
}

// struct for the different menu items
type option struct {
	text    string
	command tea.Cmd
}

func newMainMenu() mainMenuModel {

	var m mainMenuModel

	// Options for the main menu. Command are tea.Cmd's for the model's Update function to process.
	newOptions := []option{
		{text: "Random verse", command: tea.Batch(m.newRandomVerse, waitingForVerse)},
		{text: "Read the Bible", command: func() tea.Msg { return newBookSelectionScreen() }},   // Select book -> select chapter -> read
		{text: "Change translation", command: func() tea.Msg { return newTranslationScreen() }}, // Open a screen to choose the translation
		{text: "Exit BibleTUI", command: tea.Quit},
	}

	m = mainMenuModel{
		randomVerseVP: viewport.New(window_width, window_height-lipgloss.Height(getHeaderWithList(m))),
		options:       newOptions,
		choiceIndex:   0,
	}

	m.randomVerseVP.SetContent("Loading a random verse...")

	return m
}

func (m mainMenuModel) Init() tea.Cmd {
	// Initiate the main menu with a random verse
	return m.newRandomVerse
}

func (m mainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case api_query.RandomQuery:
		// when a random query is returned, apply it to the verse viewport
		m.applyRandomVerse(msg)
	case tea.KeyMsg:
		// keys navigate the menu. Enter calls the tea.Cmd for the selected option
		switch msg.String() {
		case "up", "left":
			m.choiceIndex = (m.choiceIndex - 1 + len(m.options)) % len(m.options)
			return m, nil
		case "down", "right":
			m.choiceIndex = (m.choiceIndex + 1) % len(m.options)
			return m, nil
		case "enter":
			return m, m.options[m.choiceIndex].command
		}
	case tea.WindowSizeMsg:
		// Resizing the terminal window affects the viewport dimensions
		m.randomVerseVP.Width = msg.Width
		m.randomVerseVP.Height = msg.Height - lipgloss.Height(getHeaderWithList(m))
	case string:
		// The only string tea.Msg affects the viewport only
		m.randomVerseVP.SetContent(msg)
	}

	return m, nil
}

func (m mainMenuModel) View() string {

	// Join the menu and the verse viewport vertically
	return lipgloss.JoinVertical(lipgloss.Center, getHeaderWithList(m), m.randomVerseVP.View())
}

func (m mainMenuModel) headerView() string {
	// Styling for the header
	topMsg := "* Welcome to BibleTUI! *"
	topBottomBar := styles.YellowText.Render(strings.Repeat("*", len(topMsg)))
	topMsg = styles.YellowText.Render(topMsg)

	translation := fmt.Sprintf("Current translation: %s", apiCfg.CurrentlyReading.TranslationName)
	translation = styles.InfoText.Render(translation)

	return lipgloss.JoinVertical(lipgloss.Center, topBottomBar, topMsg, topBottomBar, translation)
}

func (m mainMenuModel) getName(index int) string {
	return m.options[index].text
}

func (m mainMenuModel) getListLength() int {
	return len(m.options)
}
func (m mainMenuModel) getChoiceIndex() int {
	return m.choiceIndex
}

func (m mainMenuModel) newRandomVerse() tea.Msg {
	// Queries a new random verse and return it as a tea.Msg. Model processes it to the viewport.

	query := api_query.GetRandomVerse()

	return query

}

func waitingForVerse() tea.Msg {
	// A message while the random verse is being fetched.
	return "Loading a random verse..."
}

func (m *mainMenuModel) applyRandomVerse(query api_query.RandomQuery) {
	// Apply the new random verse, formatted
	bibleText, _ := strings.CutSuffix(query.RandomVerse.Text, "\n")

	line := strings.Repeat("-", window_width)

	s := fmt.Sprintf("%s\n%s\n- %s %d:%d (%s)\n%s",
		line,
		bibleText,
		query.RandomVerse.Book,
		query.RandomVerse.Chapter,
		query.RandomVerse.Verse,
		query.Translation.Name,
		line,
	)
	m.randomVerseVP.SetContent(s)
}
