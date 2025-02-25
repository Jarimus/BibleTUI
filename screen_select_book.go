package main

import (
	"fmt"
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

func newBookSelectionScreen() bookSelectionModel {

	// Build the menu items from the Bible translation's data
	// The data was queried at start up and when the current translation is changed
	var menuItems = []bookMenuItem{}

	for _, book := range current.translationData.Books {
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

func (m bookSelectionModel) headerView() string {
	// Styling for the header
	topMsg := "* Choose a book *"
	topBottomBar := styles.YellowText.Render(strings.Repeat("*", len(topMsg)))
	topMsg = styles.YellowText.Render(topMsg)
	return lipgloss.JoinVertical(lipgloss.Left, topBottomBar, topMsg, topBottomBar)
}

func selectBook(bookID string) func() tea.Msg {
	// selects a book
	// Queries data for the book and opens a new model where a chapter can be chosen
	return func() tea.Msg {

		current.bookData = api_query.BookQuery(current.translationID, bookID)

		return newChapterSelectionScreen()
	}
}
