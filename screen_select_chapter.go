package main

import (
	"fmt"
	"strings"

	"github.com/Jarimus/BibleTUI/internal/api_query"
	styles "github.com/Jarimus/BibleTUI/internal/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type chapterSelectionModel struct {
	menuItems   []chapterMenuItem
	choiceIndex int
}

type chapterMenuItem struct {
	name    string
	id      string
	command func(string) func() tea.Msg
}

func newChapterSelectionScreen() chapterSelectionModel {

	var menuItems = []chapterMenuItem{}

	for _, chapter := range current.bookData.Chapters {
		item := chapterMenuItem{
			name:    chapter.Reference,
			id:      chapter.ID,
			command: applyChapter,
		}
		menuItems = append(menuItems, item)
	}

	menuItems = append(menuItems, chapterMenuItem{name: "Back"})

	return chapterSelectionModel{
		menuItems: menuItems,
	}
}

func (m chapterSelectionModel) Init() tea.Cmd {
	return nil
}

func (m chapterSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			chapterID := m.menuItems[m.choiceIndex].id
			return m, m.menuItems[m.choiceIndex].command(chapterID)
		}
	}
	return m, nil
}

func (m chapterSelectionModel) View() string {
	topMsg := "* Choose a chapter *"
	topBottomBar := styles.YellowText.Render(strings.Repeat("*", len(topMsg)))
	topMsg = styles.YellowText.Render(topMsg)

	var options string

	// Show choices:
	listLength := len(m.menuItems)
	itemsShown := min(listLength, window_height-5)
	// n: starting index, which is 4 less than choiceIndex, min 0, max listLength - itemsShown
	n := max(0, min(m.choiceIndex-itemsShown/2, listLength-itemsShown))

	for i := 0; i < itemsShown; i++ {
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

func applyChapter(chapterID string) func() tea.Msg {
	return func() tea.Msg {

		current.chapterData = api_query.ChapterQuery(current.translationID, chapterID)

		return newBibleScreen()
	}
}
