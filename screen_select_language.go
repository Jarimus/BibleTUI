package main

import (
	"sort"
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
	errorText     string
}

type menuItem struct {
	name    string
	id      string
	command func(string) func() tea.Msg
}

// Returns a model that displays the languages available for all the translations in the api.
func newLanguagesScreen() languageSelectionModel {

	// Query for the languages
	biblesData, err := api_query.AllTranslationsQuery("", apiCfg.ApiKey)
	if err != nil {

	}

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
			name := translation.Language.Name
			// var name string
			// if tts.ISOtoTTScode(translation.Language.ID) != "" {
			// 	name = translation.Language.Name + " [audio available]"
			// } else {
			// 	name = translation.Language.Name
			// }
			menuItemsList = append(menuItemsList, menuItem{
				name:    name,
				id:      translation.Language.ID,
				command: selectLanguage,
			})
		}

	}

	// Sort alphabetically
	sort.Slice(menuItemsList, func(i, j int) bool {
		return menuItemsList[i].name < menuItemsList[j].name
	})

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
	case error:
		logError(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEsc.String():
			return m, func() tea.Msg { return goBackMsg{} }
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

	errorText := styles.RedText.Render(m.errorText)
	errorText = lipgloss.PlaceHorizontal(window_width, 0.5, errorText)

	return lipgloss.JoinVertical(0, getHeaderWithList(m), errorText)
}

// Returns the header of the model as a string.
func (m languageSelectionModel) headerView() string {
	topMsg := "* Choose a language *"
	topBottomBar := styles.YellowText.Render(strings.Repeat("*", len(topMsg)))
	topMsg = styles.YellowText.Render(topMsg)
	return lipgloss.JoinVertical(lipgloss.Left, topBottomBar, topMsg, topBottomBar)
}

// Returns the name of the menu item at index as a string.
func (m languageSelectionModel) getName(index int) string {
	return m.menuItemsList[index].name
}

// Return the length of the menu list as an integer.
func (m languageSelectionModel) getListLength() int {
	return len(m.menuItemsList)
}

// Return the current cursor placement as an integer.
func (m languageSelectionModel) getChoiceIndex() int {
	return m.choiceIndex
}

// Queries for translations of a specific language
// Opens a new screen to choose a translation
func selectLanguage(languageID string) func() tea.Msg {
	biblesOfLanguage, err := api_query.AllTranslationsQuery(languageID, apiCfg.ApiKey)
	if err != nil {
		return func() tea.Msg {
			return err
		}
	}
	return func() tea.Msg {
		return newAddTranslationScreen(biblesOfLanguage)
	}
}
