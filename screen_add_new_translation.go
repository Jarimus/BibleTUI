package main

import (
	"strings"

	"github.com/Jarimus/BibleTUI/internal/api_query"
	styles "github.com/Jarimus/BibleTUI/internal/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Structs for the menu model and the menu items
type addTranslationModel struct {
	menuItemsList []addTranslationMenuItem
	choiceIndex   int
}

type addTranslationMenuItem struct {
	name    string
	id      string
	command func(string, string) func() tea.Msg
}

func newAddTranslationScreen(biblesData api_query.BiblesData) addTranslationModel {

	// Get all the different languages
	var menuItemsList = []addTranslationMenuItem{}
	for _, translation := range biblesData.Data {
		var found bool
		for _, menuItem := range menuItemsList {
			if menuItem.id == translation.Language.ID {
				found = true
			}
		}
		if !found {
			menuItemsList = append(menuItemsList, addTranslationMenuItem{
				name:    translation.Name,
				id:      translation.ID,
				command: addNewTranslation,
			})
		}

	}

	menuItemsList = append(menuItemsList, addTranslationMenuItem{
		name: "Back",
	})

	return addTranslationModel{
		menuItemsList: menuItemsList,
	}
}

func (m addTranslationModel) Init() tea.Cmd {
	return nil
}

func (m addTranslationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			name := m.menuItemsList[m.choiceIndex].name
			id := m.menuItemsList[m.choiceIndex].id
			return m, m.menuItemsList[m.choiceIndex].command(name, id)
		}
	}
	return m, nil
}

func (m addTranslationModel) View() string {
	return getHeaderWithList(m)
}

func (m addTranslationModel) headerView() string {
	topMsg := "* Choose a language *"
	topBottomBar := styles.YellowText.Render(strings.Repeat("*", len(topMsg)))
	topMsg = styles.YellowText.Render(topMsg)
	return lipgloss.JoinVertical(lipgloss.Left, topBottomBar, topMsg, topBottomBar)
}

func (m addTranslationModel) getName(index int) string {
	return m.menuItemsList[index].name
}

func (m addTranslationModel) getListLength() int {
	return len(m.menuItemsList)
}
func (m addTranslationModel) getChoiceIndex() int {
	return m.choiceIndex
}

// Queries for translations of a specific language
// Opens a new screen to choose a translation
func addNewTranslation(translationName, translationID string) func() tea.Msg {
	biblesOfLanguage := api_query.AllTranslationsQuery(translationID)
	return func() tea.Msg {
		return biblesOfLanguage
	}
}
