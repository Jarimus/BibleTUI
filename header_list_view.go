package main

import (
	"fmt"
	"strings"

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
	var elements []string

	// header above the list
	elements = append(elements, m.headerView())

	// Space between the header and the list
	elements = append(elements, strings.Repeat(" ", window_width))

	// When not all items fit the screen, we need to limit them.:
	listLength := m.getListLength()
	itemsShown := min(listLength, window_height-lipgloss.Height(m.headerView())-2) // -2 to leave empty space at the bottom and the top.
	// n: index for the topmost item shown.
	n := max(0, min(m.getChoiceIndex()-itemsShown/2, listLength-itemsShown))

	// show i items from the list, starting from n
	for i := range itemsShown {
		currentIndex := n + i
		if m.getChoiceIndex() == currentIndex {
			choiceText := fmt.Sprint(styles.GreenText.Render(m.getName(currentIndex)))
			elements = append(elements, choiceText)
		} else {
			elements = append(elements, m.getName(currentIndex))
		}
	}

	// join all the elements vertically
	list := lipgloss.JoinVertical(0.5, elements...)
	return lipgloss.PlaceHorizontal(window_width, lipgloss.Center, list)
}
