package main

import (
	"log"
	"slices"
	"strings"

	"github.com/Jarimus/BibleTUI/internal/api_query"
	styles "github.com/Jarimus/BibleTUI/internal/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Structs for the menu model and the menu items
type translationSelectionModel struct {
	menuItems   []translationMenuItem
	choiceIndex int
}

type translationMenuItem struct {
	name    string
	id      string
	command func(string, string) func() tea.Msg
}

// return a new TeaModel for a menu to select a translation to use
func newTranslationScreen() translationSelectionModel {

	// Load list of translations from a file
	translations, err := loadTranslationsFromFile()
	if err != nil {
		log.Printf("error loading translations from a file: %s", err)
	}

	translations = append(translations, translationMenuItem{
		name:    "Add new translation",
		id:      "",
		command: openSelectLanguageScreen,
	})
	translations = append(translations, translationMenuItem{
		name: "Back",
	})

	return translationSelectionModel{
		menuItems: translations,
	}
}

func (m translationSelectionModel) Init() tea.Cmd {
	return nil
}

func (m translationSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Model handles navigating the menu
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "left":
			m.choiceIndex = (m.choiceIndex - 1 + len(m.menuItems)) % len(m.menuItems)
			return m, nil
		case "down", "right":
			m.choiceIndex = (m.choiceIndex + 1) % len(m.menuItems)
			return m, nil
		case "enter":
			if m.menuItems[m.choiceIndex].name == "Back" {
				return m, func() tea.Msg {
					return goBackMsg{}
				}
			}
			name := m.menuItems[m.choiceIndex].name
			id := m.menuItems[m.choiceIndex].id
			return m, m.menuItems[m.choiceIndex].command(name, id)
		case tea.KeyDelete.String():
			if m.menuItems[m.choiceIndex].name != "Add new translation" && m.menuItems[m.choiceIndex].name != "Back" {
				m.menuItems = slices.Delete(m.menuItems, m.choiceIndex, m.choiceIndex+1)
				translations := m.menuItems[:len(m.menuItems)-2]
				saveTranslationsToFile(translations)
			}
		}
	}
	return m, nil
}

func (m translationSelectionModel) View() string {

	helpText := "Press 'Del' to remove a translation from the list."
	helpText = lipgloss.Place(window_width, window_height-lipgloss.Height(getHeaderWithList(m))-1, 0.5, 1, helpText)
	helpText = styles.InfoText.Render(helpText)

	return lipgloss.JoinVertical(0, getHeaderWithList(m), helpText)
}

func (m translationSelectionModel) headerView() string {
	topMsg := "* Choose a translation *"
	topBottomBar := styles.YellowText.Render(strings.Repeat("*", len(topMsg)))
	topMsg = styles.YellowText.Render(topMsg)
	return lipgloss.JoinVertical(lipgloss.Left, topBottomBar, topMsg, topBottomBar)
}

func (m translationSelectionModel) getName(index int) string {
	return m.menuItems[index].name
}

func (m translationSelectionModel) getListLength() int {
	return len(m.menuItems)
}
func (m translationSelectionModel) getChoiceIndex() int {
	return m.choiceIndex
}

func selectTranslation(translationName, translationID string) func() tea.Msg {
	// Queries a translation and stores the data in memory
	// Translation data includes, among others, the names and IDs for the books in the translation
	return func() tea.Msg {
		current.translationID = translationID
		current.translationName = translationName

		current.translationData = api_query.TranslationQuery(current.translationID)

		return goBackMsg{}
	}
}

func openSelectLanguageScreen(string, string) func() tea.Msg {
	return func() tea.Msg {
		return newLanguagesScreen()
	}
}
