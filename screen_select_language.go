package main

import (
	"strings"

	"github.com/Jarimus/BibleTUI/internal/api_query"
	styles "github.com/Jarimus/BibleTUI/internal/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Structs for the menu model and the menu items
type languageSelectionModel struct {
	menuItemsList []menuItem
	choiceIndex   int
}

type menuItem struct {
	name    string
	id      string
	command func(string) func() tea.Msg
}

func newLanguagesScreen() languageSelectionModel {

	// Query for the languages
	biblesData := api_query.AllTranslationsQuery("")

	// Get all the different languages
	var menuItemsList = []menuItem{}
	for _, translation := range biblesData.Data {
		var found bool
		for _, menuItem := range menuItemsList {
			if menuItem.id == translation.Language.ID {
				found = true
			}
		}
		if !found {
			menuItemsList = append(menuItemsList, menuItem{
				name:    translation.Language.Name,
				id:      translation.Language.ID,
				command: selectLanguage,
			})
		}

	}

	menuItemsList = append(menuItemsList, menuItem{
		name: "Back",
	})

	return languageSelectionModel{
		menuItemsList: menuItemsList,
	}
}

func (m languageSelectionModel) Init() tea.Cmd {
	return nil
}

func (m languageSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Model handles navigating the menu
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.choiceIndex = (m.choiceIndex - 1 + len(m.menuItemsList)) % len(m.menuItemsList)
			return m, nil
		case tea.KeyPgUp.String(), "left":
			m.choiceIndex = max(0, m.choiceIndex-10)
		case "down":
			m.choiceIndex = (m.choiceIndex + 1) % len(m.menuItemsList)
			return m, nil
		case tea.KeyPgDown.String(), "right":
			m.choiceIndex = min(len(m.menuItemsList)-1, m.choiceIndex+10)
		case "enter":
			if m.menuItemsList[m.choiceIndex].command == nil {
				return m, func() tea.Msg {
					return goBackMsg{}
				}
			}
			id := m.menuItemsList[m.choiceIndex].id
			return m, m.menuItemsList[m.choiceIndex].command(id)
		}
	}
	return m, nil
}

func (m languageSelectionModel) View() string {
	return getHeaderWithList(m)
}

func (m languageSelectionModel) headerView() string {
	topMsg := "* Choose a language *"
	topBottomBar := styles.YellowText.Render(strings.Repeat("*", len(topMsg)))
	topMsg = styles.YellowText.Render(topMsg)
	return lipgloss.JoinVertical(lipgloss.Left, topBottomBar, topMsg, topBottomBar)
}

func (m languageSelectionModel) getName(index int) string {
	return m.menuItemsList[index].name
}

func (m languageSelectionModel) getListLength() int {
	return len(m.menuItemsList)
}
func (m languageSelectionModel) getChoiceIndex() int {
	return m.choiceIndex
}

// Queries for translations of a specific language
// Opens a new screen to choose a translation
func selectLanguage(languageID string) func() tea.Msg {
	biblesOfLanguage := api_query.AllTranslationsQuery(languageID)
	return func() tea.Msg {
		return biblesOfLanguage
	}
}
