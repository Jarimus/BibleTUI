package main

import (
	"fmt"

	styles "github.com/Jarimus/BibleTUI/internal/styles"
	"github.com/charmbracelet/lipgloss"
)

type headerListModel interface {
	headerView() string
	getName(int) string
	getListLength() int
	getChoiceIndex() int
}

func getHeaderWithList(m headerListModel) string {
	var options []string

	options = append(options, m.headerView())

	// When not all items fit the screen, we need to limit them:
	listLength := m.getListLength()
	itemsShown := min(listLength, window_height-lipgloss.Height(m.headerView()))
	// n: index for the topmost item shown.
	n := max(0, min(m.getChoiceIndex()-itemsShown/2, listLength-itemsShown))

	// show i items from the list, starting from n
	for i := range itemsShown {
		currentIndex := n + i
		if m.getChoiceIndex() == currentIndex {
			choiceText := fmt.Sprint(styles.GreenText.Render(m.getName(currentIndex)))
			options = append(options, choiceText)
		} else {
			options = append(options, m.getName(currentIndex))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, options...)
}
