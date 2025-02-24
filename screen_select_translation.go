package main

import (
	"fmt"
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

func newTranslationScreen() translationSelectionModel {
	// Translations are retrieved from the API: https://scripture.api.bible/
	// Translation IDs: https://docs.api.bible/guides/bibles
	var menuItems = []translationMenuItem{
		{
			name:    "Simplified Chinese",
			id:      "7ea794434e9ea7ee-01",
			command: selectTranslation,
		},
		{
			name:    "Finnish New Testament",
			id:      "c739534f6a23acb2-01",
			command: selectTranslation,
		},
		{
			name:    "American Standard",
			id:      "685d1470fe4d5c3b-01",
			command: selectTranslation,
		},
		{
			name:    "King James",
			id:      "de4e12af7f28f599-01",
			command: selectTranslation,
		},
		{
			name:    "World English Bible",
			id:      "9879dbb7cfe39e4d-01",
			command: selectTranslation,
		},
		{
			name:    "Open Hebrew Living New Testament",
			id:      "a8a97eebae3c98e4-01",
			command: selectTranslation,
		},
		{
			name:    "Brenton Greek Septuagint",
			id:      "c114c33098c4fef1-01",
			command: selectTranslation,
		},
		{
			name:    "Solid Rock Greek New Testament",
			id:      "47f396bad37936f0-01",
			command: selectTranslation,
		},
		{
			name: "Back",
		},
	}

	return translationSelectionModel{
		menuItems: menuItems,
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
		}
	}
	return m, nil
}

func (m translationSelectionModel) View() string {
	topMsg := "* Choose a translation *"
	topBottomBar := styles.YellowText.Render(strings.Repeat("*", len(topMsg)))
	topMsg = styles.YellowText.Render(topMsg)

	var options string
	for i, item := range m.menuItems {
		if m.choiceIndex == i {
			choiceText := fmt.Sprint(styles.GreenText.Render(item.name))
			options += choiceText + "\n"
		} else {
			options += item.name + "\n"
		}
	}
	return lipgloss.JoinVertical(lipgloss.Left, topBottomBar, topMsg, topBottomBar, options)
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
