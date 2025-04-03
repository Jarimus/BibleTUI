package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Jarimus/BibleTUI/internal/api_query"
	styles "github.com/Jarimus/BibleTUI/internal/styles"
	"github.com/Jarimus/BibleTUI/internal/tts"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Structs for the menu model and the menu items
type addTranslationModel struct {
	menuItemsList []addTranslationMenuItem
	choiceIndex   int
	errorText     string
}

type addTranslationMenuItem struct {
	name          string
	translationID string
	languageID    string
	command       func(string, string, string) func() tea.Msg
}

// Creates a new tea.Model to display all the different translations from which a new translation can be picked.
func newAddTranslationScreen(biblesData api_query.BiblesData) addTranslationModel {

	// Get all the different translations for the language
	var menuItemsList = []addTranslationMenuItem{}
	for _, translation := range biblesData.Data {

		// Only add a single instance of each translation (not all different versions)
		var found bool
		for _, menuItem := range menuItemsList {
			if menuItem.name == translation.Name {
				found = true
			}
		}
		if !found {
			// Adds the description at the end of the translation's name if it's not empty or 'common'.
			var name string
			var audioAvailable string
			if _, ok := tts.LanguageMap[translation.Language.ID]; ok {
				audioAvailable = " [audio]"
			}
			if slices.Contains([]string{"", "common"}, strings.ToLower(translation.Description)) {
				name = fmt.Sprintf("%s%s", translation.Name, audioAvailable)
			} else {

				name = fmt.Sprintf("%s (%s)%s", translation.Name, translation.Description, audioAvailable)
			}
			menuItemsList = append(menuItemsList, addTranslationMenuItem{
				name:          name,
				translationID: translation.ID,
				languageID:    translation.Language.ID,
				command:       addNewTranslation,
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
	case error:
		m.errorText = msg.Error()
	case tea.KeyMsg:
		m.errorText = ""
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
			name := m.menuItemsList[m.choiceIndex].name
			translationID := m.menuItemsList[m.choiceIndex].translationID
			languageID := m.menuItemsList[m.choiceIndex].languageID
			return m, m.menuItemsList[m.choiceIndex].command(name, translationID, languageID)
		}
	}
	return m, nil
}

func (m addTranslationModel) View() string {
	return getHeaderWithList(m)
}

// Return the header as a string
func (m addTranslationModel) headerView() string {
	topMsg := "* Choose a translation to add *"
	topBottomBar := styles.YellowText.Render(strings.Repeat("*", len(topMsg)))
	topMsg = styles.YellowText.Render(topMsg)
	header := topBottomBar + "\n" + topMsg + "\n" + topBottomBar
	header = lipgloss.PlaceHorizontal(window_width, 0.5, header)
	if m.errorText != "" {
		errorText := lipgloss.PlaceHorizontal(window_width, 0.5, styles.RedText.Render(m.errorText))
		return lipgloss.JoinVertical(0, header, errorText)
	}
	return lipgloss.JoinVertical(0, header)
}

// Returns the name of the menu item at index as a string.
func (m addTranslationModel) getName(index int) string {
	return m.menuItemsList[index].name
}

// Return the length of the menu list as an integer.
func (m addTranslationModel) getListLength() int {
	return len(m.menuItemsList)
}

// Return the current cursor placement as an integer.
func (m addTranslationModel) getChoiceIndex() int {
	return m.choiceIndex
}

// Adds the chosen translation to the list of translations.
// Also sets the new translation as the current translation.
func addNewTranslation(translationName, translationID, languageID string) func() tea.Msg {

	// Add the translation to the database
	translation, err := addTranslationToDatabase(translationName, translationID, languageID)
	if err != nil {
		return func() tea.Msg {
			return err
		}
	}

	// Set the current translation to the newly added translation
	apiCfg.CurrentlyReading.TranslationName = translation.Name
	apiCfg.CurrentlyReading.TranslationID = translation.ApiID
	apiCfg.CurrentlyReading.LanguageID = translation.LanguageID
	apiCfg.CurrentlyReading.TranslationData, err = api_query.TranslationQuery(translation.ApiID, apiCfg.ApiKey)
	if err != nil {
		msg := styles.RedText.Render("Unable to access the api. Enter a valid API Key to access the Bible translations.")
		tea.Println(msg)
		// time.Sleep(1500 * time.Millisecond)
	}

	err = saveSettings()
	if err != nil {
		errorF := func() tea.Msg {
			return err
		}
		return errorF
	}

	goBack := func() tea.Msg {
		return goBackMsg{}
	}

	return tea.Batch(goBack, goBack, goBack)
}
