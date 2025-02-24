package main

import (
	"fmt"
	"strings"

	styles "github.com/Jarimus/BibleTUI/internal/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type translationSelectionModel struct {
	menuItems   []item
	choiceIndex int
}

type item struct {
	name    string
	id      string
	command func(string, string) func() tea.Msg
}

func newTranslationScreen() translationSelectionModel {

	var menuItems = []item{
		{
			name:    "Simplified Chinese",
			id:      "7ea794434e9ea7ee-01",
			command: applyTranslation,
		},
		{
			name:    "Finnish New Testament",
			id:      "c739534f6a23acb2-01",
			command: applyTranslation,
		},
		{
			name:    "American Standard",
			id:      "685d1470fe4d5c3b-01",
			command: applyTranslation,
		},
		{
			name:    "King James",
			id:      "de4e12af7f28f599-01",
			command: applyTranslation,
		},
		{
			name:    "World English Bible",
			id:      "9879dbb7cfe39e4d-01",
			command: applyTranslation,
		},
		{
			name: "Back",
		},
	}

	return translationSelectionModel{
		menuItems: menuItems,
	}
}

func (m translationSelectionModel) Init() tea.Cmd {
	return nil
}

func (m translationSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.choiceIndex = (m.choiceIndex - 1 + len(m.menuItems)) % len(m.menuItems)
			return m, nil
		case "down":
			m.choiceIndex = (m.choiceIndex + 1) % len(m.menuItems)
			return m, nil
		case "enter":
			if m.menuItems[m.choiceIndex].name == "Back" {
				return m, func() tea.Msg {
					return goBackMsg{}
				}
			}
			name := m.menuItems[m.choiceIndex].name
			id := m.menuItems[m.choiceIndex].id
			return m, m.menuItems[m.choiceIndex].command(name, id)
		}
	}
	return m, nil
}

func (m translationSelectionModel) View() string {
	topMsg := "* Choose a translation *"
	topBottomBar := styles.YellowText.Render(strings.Repeat("*", len(topMsg)))
	topMsg = styles.YellowText.Render(topMsg)

	var options string
	for i, item := range m.menuItems {
		if m.choiceIndex == i {
			choiceText := fmt.Sprint(styles.GreenText.Render(item.name))
			options += choiceText + "\n"
		} else {
			options += item.name + "\n"
		}
	}
	return lipgloss.JoinVertical(lipgloss.Left, topBottomBar, topMsg, topBottomBar, options)
}

func applyTranslation(translationName, translationID string) func() tea.Msg {
	return func() tea.Msg {
		current.translationID = translationID
		current.translationName = translationName
		return goBackMsg{}
	}
}
