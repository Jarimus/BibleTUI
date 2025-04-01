package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mitchellh/go-wordwrap"

	"github.com/Jarimus/BibleTUI/internal/api_query"
	styles "github.com/Jarimus/BibleTUI/internal/styles"
)

// tea.Model for the mainMenu
type mainMenuModel struct {
	randomVerse string
	options     []option
	choiceIndex int
	textInput   textinput.Model
}

// struct for the different menu items
type option struct {
	text    string
	command tea.Cmd
}

// Return a tea.model that displays the main menu.
func newMainMenu() mainMenuModel {

	var m mainMenuModel

	// Options for the main menu. Command are tea.Cmd's for the model's Update function to process.
	newOptions := []option{
		{text: "Random verse", command: tea.Batch(m.newRandomVerse, waitingForVerse)},
		{text: "Read the Bible", command: func() tea.Msg { return newBookSelectionScreen() }},   // Select book -> select chapter -> read
		{text: "Change translation", command: func() tea.Msg { return newTranslationScreen() }}, // Open a screen to choose the translation
		{text: "Change user", command: func() tea.Msg { return newUsersMenu() }},                // Open users menu to switch current user and create new users
		{text: "Change API Key", command: func() tea.Msg { return focusInput{} }},
		{text: "Exit BibleTUI", command: tea.Quit},
	}

	inputField := textinput.New()

	m = mainMenuModel{
		randomVerse: "",
		options:     newOptions,
		choiceIndex: 0,
		textInput:   inputField,
	}

	// Input field settings
	m.textInput.TextStyle = styles.BlueText
	m.textInput.PromptStyle = styles.BlueText
	m.textInput.PlaceholderStyle = styles.BlueText
	m.textInput.Cursor.Style = styles.BlueText
	m.textInput.Prompt = "> "
	m.textInput.Placeholder = "Enter your API key"

	m.randomVerse = "Loading a random verse..."

	return m
}

func (m mainMenuModel) Init() tea.Cmd {
	// Initiate the main menu with a random verse
	return tea.Batch(m.newRandomVerse, m.textInput.Cursor.BlinkCmd())
}

func (m mainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	// A notification replaces the random vers
	case newNotificationMsg:
		m.randomVerse = msg.text

	// Put the input field into focus.
	case focusInput:
		m.textInput.Focus()
		return m, nil

	// when a random query is returned, apply it to the verse viewport
	case api_query.RandomQuery:
		m.applyRandomVerse(msg)

	// keys navigate the menu. Enter calls the tea.Cmd for the selected option
	case tea.KeyMsg:
		if !m.textInput.Focused() {
			switch msg.String() {
			case "up", "left":
				m.choiceIndex = (m.choiceIndex - 1 + len(m.options)) % len(m.options)
				return m, nil
			case "down", "right":
				m.choiceIndex = (m.choiceIndex + 1) % len(m.options)
				return m, nil
			case "enter":
				return m, m.options[m.choiceIndex].command
			case tea.KeyEsc.String():
				return m, func() tea.Msg { return goBackMsg{} }
			}
		} else {
			switch msg.String() {
			case tea.KeyEsc.String():
				m.textInput.Blur()
				m.textInput.Reset()
			case "up":
				m.textInput.CursorStart()
				return m, nil
			case "down":
				m.textInput.CursorEnd()
				return m, nil

			// First test if the API Key works. If not, display an error message
			case "enter":
				newApiKey := m.textInput.Value()
				if err := api_query.TestConnection(newApiKey); err != nil {
					if window_width < 0 {
						m.randomVerse = styles.RedText.Render(err.Error())
					} else {
						m.randomVerse = styles.RedText.Render(wordwrap.WrapString(err.Error(), uint(window_width)))
					}
				} else {
					apiCfg.ApiKey = m.textInput.Value()
					err = saveSettings()
					if err != nil {
						m.randomVerse = err.Error()
						return m, nil
					}
					m.textInput.Blur()
					m.textInput.Reset()
				}

			}
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m mainMenuModel) View() string {

	// Join the elements to form a ui
	var inputField string
	if m.textInput.Focused() {
		inputField = m.textInput.View()
	} else {
		inputField = ""
	}

	inputField = lipgloss.PlaceHorizontal(window_width, 0.5, inputField)
	inputField = lipgloss.PlaceVertical(window_height-lipgloss.Height(getHeaderWithList(m))-lipgloss.Height(m.randomVerse), 0.5, inputField)

	return lipgloss.JoinVertical(0, getHeaderWithList(m), m.randomVerse, inputField)
}

// Returns the header as a string.
func (m mainMenuModel) headerView() string {
	// Styling for the header
	topMsg := "* Welcome to BibleTUI! *"
	topBottomBar := styles.YellowText.Render(strings.Repeat("*", len(topMsg)))
	topMsg = styles.YellowText.Render(topMsg)

	// Display the active user above the menu
	user := fmt.Sprintf("Active user: %s", apiCfg.CurrentUser)
	user = styles.InfoText.Render(user)

	// Display API Key
	apiKey := fmt.Sprintf("API Key: %s", apiCfg.ApiKey)
	apiKey = styles.InfoText.Render(apiKey)

	// Display the current translation above the menu
	translation := fmt.Sprintf("Current translation: %s", apiCfg.CurrentlyReading.TranslationName)
	translation = styles.InfoText.Render(translation)

	return lipgloss.JoinVertical(lipgloss.Center, topBottomBar, topMsg, topBottomBar, user, apiKey, translation)
}

// Returns the name of the menu item at index as a string.
func (m mainMenuModel) getName(index int) string {
	return m.options[index].text
}

// Return the length of the menu list as an integer.
func (m mainMenuModel) getListLength() int {
	return len(m.options)
}

// Return the current cursor placement as an integer.
func (m mainMenuModel) getChoiceIndex() int {
	return m.choiceIndex
}

// Uses the internal api_query package to query a random verse from the api and returns the query as a tea.Msg (struct).
func (m mainMenuModel) newRandomVerse() tea.Msg {
	// Queries a new random verse and return it as a tea.Msg. Model processes it to the viewport.

	query := api_query.GetRandomVerse()

	return query

}

// Returns a tea.Msg (string) to be displayed in the main menu.
// This is a placeholder text while the api is being queried for a random verse.
func waitingForVerse() tea.Msg {
	// A message while the random verse is being fetched.
	return newNotificationMsg{
		text: "Loading a random verse...",
	}
}

// When the tea.Model receives a random verse as a tea.Msg, it uses this function to apply the verse to the model to be displayed.
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
	s = strings.ReplaceAll(s, "\n\n", "\n")
	if window_width < 0 {
		m.randomVerse = s
	} else {
		m.randomVerse = wordwrap.WrapString(s, uint(window_width))
	}
}
