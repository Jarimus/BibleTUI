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

	var options []string

	options = append(options, m.headerView())

	// When not all items fit the screen, we need to limit them:
	listLength := len(m.menuItems)
	itemsShown := min(listLength, window_height-lipgloss.Height(m.headerView()))
	// n: index for the topmost item shown.
	n := max(0, min(m.choiceIndex-itemsShown/2, listLength-itemsShown))

	// show i items from the list, starting from n
	for i := range itemsShown {
		currentIndex := n + i
		if m.choiceIndex == currentIndex {
			choiceText := fmt.Sprint(styles.GreenText.Render(m.menuItems[currentIndex].name))
			options = append(options, choiceText)
		} else {
			options = append(options, m.menuItems[currentIndex].name)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, options...)
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
