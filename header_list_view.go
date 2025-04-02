package main

import (
	"fmt"
	"unicode/utf8"

	styles "github.com/Jarimus/BibleTUI/internal/styles"
	"github.com/charmbracelet/lipgloss"
)

type headerListModel interface {
	headerView() string
	getName(int) string
	getListLength() int
	getChoiceIndex() int
}

// Puts together the model's header and list items and returns them as a string.
// The string is centered to the window and header is placed above the list.
func getHeaderWithList(m headerListModel) string {
	var elements []string

	// header above the list
	elements = append(elements, m.headerView())

	// Space between the header and the list
	elements = append(elements, lipgloss.PlaceHorizontal(window_width, 0.5, " "))

	// When not all items fit the screen, we need to limit them.:
	listLength := m.getListLength()
	itemsShown := min(listLength, window_height-lipgloss.Height(m.headerView())-2) // -2 to leave empty space at the bottom and the top.
	// n: index for the topmost item shown.
	n := max(0, min(m.getChoiceIndex()-itemsShown/2, listLength-itemsShown))

	// show i items from the list, starting from n
	for i := range itemsShown {
		currentIndex := n + i
		currentListItem := m.getName(currentIndex)
		if utf8.RuneCountInString(currentListItem) > window_width {
			currentListItem = currentListItem[:window_width]
		}
		if m.getChoiceIndex() == currentIndex {
			choiceText := fmt.Sprint(styles.GreenText.Render(currentListItem))
			elements = append(elements, choiceText)
		} else {
			elements = append(elements, currentListItem)
		}
	}

	// join all the elements vertically
	return lipgloss.JoinVertical(0.5, elements...)
}
