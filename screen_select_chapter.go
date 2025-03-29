package main

import (
	"strings"

	"github.com/Jarimus/BibleTUI/internal/api_query"
	styles "github.com/Jarimus/BibleTUI/internal/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Chapter selection works similarly to book selection
type chapterSelectionModel struct {
	menuItems   []chapterMenuItem
	choiceIndex int
}

type chapterMenuItem struct {
	name    string
	id      string
	command func(string) func() tea.Msg
}

// Returns a model that displays the chapters of the current book as a list.
func newChapterSelectionScreen() chapterSelectionModel {
	var menuItems = []chapterMenuItem{}

	for _, chapter := range apiCfg.CurrentlyReading.BookData.Chapters {
		item := chapterMenuItem{
			name:    chapter.Reference,
			id:      chapter.ID,
			command: selectChapter,
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

	return getHeaderWithList(m)
}

// Returns the header of the model as a string.
func (m chapterSelectionModel) headerView() string {
	topMsg := "* Choose a chapter *"
	topBottomBar := styles.YellowText.Render(strings.Repeat("*", len(topMsg)))
	topMsg = styles.YellowText.Render(topMsg)
	return lipgloss.JoinVertical(lipgloss.Center, topBottomBar, topMsg, topBottomBar)
}

// Returns the name of the menu item at index as a string.
func (m chapterSelectionModel) getName(index int) string {
	return m.menuItems[index].name
}

// Return the length of the menu list as an integer.
func (m chapterSelectionModel) getListLength() int {
	return len(m.menuItems)
}

// Return the current cursor placement as an integer.
func (m chapterSelectionModel) getChoiceIndex() int {
	return m.choiceIndex
}

// selects a chapter
// Queries data for the chapter and opens a new model where the chapter can be read
func selectChapter(chapterID string) func() tea.Msg {
	return func() tea.Msg {

		apiCfg.CurrentlyReading.ChapterData = api_query.ChapterQuery(apiCfg.CurrentlyReading.TranslationID, chapterID)

		return newBibleScreen()
	}
}
