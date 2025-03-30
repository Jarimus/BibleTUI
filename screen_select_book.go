package main

import (
	"strings"

	"github.com/Jarimus/BibleTUI/internal/api_query"
	styles "github.com/Jarimus/BibleTUI/internal/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Structs to build the list of books
type bookSelectionModel struct {
	menuItems   []bookMenuItem
	choiceIndex int
}

type bookMenuItem struct {
	name    string
	id      string
	command func(string) func() tea.Msg
}

// Returns a new model that displays a list of books in the current translations as a list.
func newBookSelectionScreen() bookSelectionModel {

	// Build the menu items from the Bible translation's data
	// The data was queried at start up and when the current translation is changed
	var menuItems = []bookMenuItem{}

	for _, book := range apiCfg.CurrentlyReading.TranslationData.Books {
		item := bookMenuItem{
			name:    book.Name,
			id:      book.ID,
			command: selectBook,
		}
		menuItems = append(menuItems, item)
	}

	menuItems = append(menuItems, bookMenuItem{name: "Back"})

	return bookSelectionModel{
		menuItems: menuItems,
	}
}

func (m bookSelectionModel) Init() tea.Cmd {
	// No initialization needed
	return nil
}

func (m bookSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Similar to the main menu
	// The model handles navigating the menu
	// Calling a choice initiates a query to retrieve data for the chosen book
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.choiceIndex = (m.choiceIndex - 1 + len(m.menuItems)) % len(m.menuItems)
			return m, nil
		case tea.KeyPgUp.String(), "left":
			m.choiceIndex = max(0, m.choiceIndex-10)
		case "down":
			m.choiceIndex = (m.choiceIndex + 1) % len(m.menuItems)
			return m, nil
		case tea.KeyPgDown.String(), "right":
			m.choiceIndex = min(len(m.menuItems)-1, m.choiceIndex+10)
		case "enter":
			if m.menuItems[m.choiceIndex].command == nil {
				return m, func() tea.Msg {
					return goBackMsg{}
				}
			}
			bookID := m.menuItems[m.choiceIndex].id
			return m, m.menuItems[m.choiceIndex].command(bookID)
		}
	}
	return m, nil
}

func (m bookSelectionModel) View() string {

	return getHeaderWithList(m)
}

// Returns the header of the model as a string.
func (m bookSelectionModel) headerView() string {
	// Styling for the header
	topMsg := "* Choose a book *"
	topBottomBar := styles.YellowText.Render(strings.Repeat("*", len(topMsg)))
	topMsg = styles.YellowText.Render(topMsg)
	return lipgloss.JoinVertical(lipgloss.Center, topBottomBar, topMsg, topBottomBar)
}

// Returns the name of the menu item at index as a string.
func (m bookSelectionModel) getName(index int) string {
	return m.menuItems[index].name
}

// Return the length of the menu list as an integer.
func (m bookSelectionModel) getListLength() int {
	return len(m.menuItems)
}

// Return the current cursor placement as an integer.
func (m bookSelectionModel) getChoiceIndex() int {
	return m.choiceIndex
}

// selects a book
// Queries data for the book and opens a new model where a chapter can be chosen
func selectBook(bookID string) func() tea.Msg {
	return func() tea.Msg {

		apiCfg.CurrentlyReading.BookData = api_query.BookQuery(apiCfg.CurrentlyReading.TranslationID, bookID, apiCfg.apiKey)

		return newChapterSelectionScreen()
	}
}
