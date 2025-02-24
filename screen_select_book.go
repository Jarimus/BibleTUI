package main

import (
	"fmt"
	"strings"

	"github.com/Jarimus/BibleTUI/internal/api_query"
	styles "github.com/Jarimus/BibleTUI/internal/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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

	var menuItems = []bookMenuItem{}

	for _, book := range current.translationData.Books {
		item := bookMenuItem{
			name:    book.Name,
			id:      book.ID,
			command: applyBook,
		}
		menuItems = append(menuItems, item)
	}

	menuItems = append(menuItems, bookMenuItem{name: "Back"})

	return bookSelectionModel{
		menuItems: menuItems,
	}
}

func (m bookSelectionModel) Init() tea.Cmd {
	return nil
}

func (m bookSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	topMsg := "* Choose a book *"
	topBottomBar := styles.YellowText.Render(strings.Repeat("*", len(topMsg)))
	topMsg = styles.YellowText.Render(topMsg)

	var options string

	// Show choices:
	listLength := len(m.menuItems)
	itemsShown := min(30, listLength)
	// n: top index to show list items from.
	n := max(0, min(m.choiceIndex-itemsShown/2, listLength-itemsShown))

	for i := range itemsShown {
		currentIndex := n + i
		if m.choiceIndex == currentIndex {
			choiceText := fmt.Sprint(styles.GreenText.Render(m.menuItems[currentIndex].name))
			options += choiceText + "\n"
		} else {
			options += m.menuItems[currentIndex].name + "\n"
		}
	}
	return lipgloss.JoinVertical(lipgloss.Left, topBottomBar, topMsg, topBottomBar, options)
}

func applyBook(bookID string) func() tea.Msg {
	return func() tea.Msg {

		current.bookData = api_query.BookQuery(current.translationID, bookID)

		return newChapterSelectionScreen()
	}
}
