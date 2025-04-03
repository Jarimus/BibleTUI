package main

import (
	"context"
	"log"
	"slices"
	"strings"

	"github.com/Jarimus/BibleTUI/internal/api_query"
	"github.com/Jarimus/BibleTUI/internal/database"
	styles "github.com/Jarimus/BibleTUI/internal/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Structs for the menu model and the menu items
type translationSelectionModel struct {
	menuItems   []translationMenuItem
	choiceIndex int
	errorText   string
}

type translationMenuItem struct {
	name          string
	translationID string
	languageID    string
	command       func(string, string, string) func() tea.Msg
}

// return a new tea.Model for a menu to select a translation to use.
func newTranslationScreen() translationSelectionModel {

	// Load list of translations from database for the active user.
	translationsDB, err := loadTranslationsForUser()
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
		} else {
			log.Fatalf("error loading translations from database: %s", err)
		}
	}

	var translations []translationMenuItem

	for _, item := range translationsDB {
		translations = append(translations, translationMenuItem{
			name:          item.Name,
			translationID: item.ApiID,
			languageID:    item.LanguageID,
			command:       selectTranslation,
		})
	}

	translations = append(translations, translationMenuItem{
		name:    "Add new translation",
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
	case error:
		m.errorText = msg.Error()
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEsc.String():
			return m, func() tea.Msg { return goBackMsg{} }
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
			translationID := m.menuItems[m.choiceIndex].translationID
			languageID := m.menuItems[m.choiceIndex].languageID
			return m, m.menuItems[m.choiceIndex].command(name, translationID, languageID)
		case tea.KeyDelete.String():
			if m.menuItems[m.choiceIndex].name != "Add new translation" && m.menuItems[m.choiceIndex].name != "Back" {

				// Delete the translation for the user from the database.
				params := database.DeleteTranslationForUserParams{
					UserID: apiCfg.CurrentUserID,
					ApiID:  m.menuItems[m.choiceIndex].translationID,
				}
				err := apiCfg.dbQueries.DeleteTranslationForUser(context.Background(), params)
				if err != nil {
					m.errorText = err.Error()
					return m, nil
				}

				// Delete the menu item from the slice.
				m.menuItems = slices.Delete(m.menuItems, m.choiceIndex, m.choiceIndex+1)
			}
		}
	}
	return m, nil
}

func (m translationSelectionModel) View() string {

	errorText := m.errorText
	errorText = styles.RedText.Render(errorText)
	errorText = lipgloss.PlaceHorizontal(window_width, 0.5, errorText)

	helpText := "Press 'Del' to remove a translation from the list."
	helpText = lipgloss.PlaceHorizontal(window_width, 0.5, helpText)
	helpText = styles.InfoText.Render(helpText)

	return lipgloss.JoinVertical(0, getHeaderWithList(m), errorText, helpText)
}

// Returns the header of the model as a string.
func (m translationSelectionModel) headerView() string {
	topMsg := "* Choose a translation *"
	topBottomBar := styles.YellowText.Render(strings.Repeat("*", len(topMsg)))
	topMsg = styles.YellowText.Render(topMsg)
	return lipgloss.JoinVertical(lipgloss.Left, topBottomBar, topMsg, topBottomBar)
}

// Returns the name of the menu item at index as a string.
func (m translationSelectionModel) getName(index int) string {
	return m.menuItems[index].name
}

// Return the length of the menu list as an integer.
func (m translationSelectionModel) getListLength() int {
	return len(m.menuItems)
}

// Return the current cursor placement as an integer.
func (m translationSelectionModel) getChoiceIndex() int {
	return m.choiceIndex
}

// Queries a translation and stores the data in memory
// Translation data includes, among others, the names and IDs for the books in the translation
func selectTranslation(translationName, translationID, languageID string) func() tea.Msg {
	return func() tea.Msg {
		var err error
		apiCfg.CurrentlyReading.TranslationID = translationID
		apiCfg.CurrentlyReading.TranslationName = translationName
		apiCfg.CurrentlyReading.LanguageID = languageID
		apiCfg.CurrentlyReading.TranslationData, err = api_query.TranslationQuery(apiCfg.CurrentlyReading.TranslationID, apiCfg.ApiKey)
		if err != nil {
			return err
		}

		err = saveSettings()
		if err != nil {
			return err
		}

		return goBackMsg{}
	}
}

// Return a function the returns a tea.Model screen to select a language in.
func openSelectLanguageScreen(string, string, string) func() tea.Msg {
	return func() tea.Msg {
		return newLanguagesScreen()
	}
}
