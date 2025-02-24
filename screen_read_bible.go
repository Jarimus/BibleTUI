package main

import (
	"fmt"
	"strings"

	"github.com/Jarimus/BibleTUI/internal/api_query"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	wordwrap "github.com/mitchellh/go-wordwrap"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	vpStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 2)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type bibleScreenModel struct {
	title     string
	bibleText string
	viewport  viewport.Model
}

func newBibleScreen() bibleScreenModel {
	// Get the chapter for reading
	bibleText := current.chapterData.Data.Content
	bibleText = strings.ReplaceAll(bibleText, "[", "\n[")
	bibleText = wordwrap.WrapString(bibleText, uint(window_width-2))
	title := current.chapterData.Data.Reference

	newBibleScreen := bibleScreenModel{
		title:     title,
		bibleText: bibleText,
	}
	// Generate a viewport from the dimensions of the global variables
	headerHeight := lipgloss.Height(newBibleScreen.headerView())
	footerHeight := lipgloss.Height(newBibleScreen.footerView())
	vpStyleHeight, vpStyleWidth := getStyleDimensions(vpStyle)
	verticalMargin := headerHeight + footerHeight + vpStyleHeight
	horizontalMargin := vpStyleWidth
	newBibleScreen.viewport = viewport.New(window_width-horizontalMargin, window_height-verticalMargin)
	newBibleScreen.viewport.YPosition = headerHeight
	newBibleScreen.viewport.SetContent(newBibleScreen.bibleText)

	// Return the model
	return newBibleScreen
}

func (m bibleScreenModel) Init() tea.Cmd {
	return nil
}

func (m bibleScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	// Window resize affects the viewport
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		vpStyleHeight, vpStyleWidth := getStyleDimensions(vpStyle)
		verticalMargin := headerHeight + footerHeight + vpStyleHeight
		horizontalMargin := vpStyleWidth
		m.viewport.Width = msg.Width - horizontalMargin
		m.viewport.Height = msg.Height - verticalMargin

	// Retrieving the Bible chapter updates the viewport text
	case api_query.ChapterData:
		m.title = current.chapterData.Data.Reference
		m.bibleText = current.chapterData.Data.Content
		m.bibleText = strings.ReplaceAll(m.bibleText, "[", "\n[")
		m.bibleText = wordwrap.WrapString(m.bibleText, uint(window_width-2))
		m.viewport.YPosition = lipgloss.Height(m.headerView())
		m.viewport.SetContent(m.bibleText)
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyLeft.String():
			cmds = append(cmds, toPreviousChapter)
		case tea.KeyRight.String():
			cmds = append(cmds, toNextChapter)
		}

	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m bibleScreenModel) View() string {
	text := vpStyle.Render(m.viewport.View())
	return lipgloss.JoinVertical(lipgloss.Top, m.headerView(), text, m.footerView())
}

func (m bibleScreenModel) headerView() string {
	title := titleStyle.Render(m.title)
	line := strings.Repeat("-", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m bibleScreenModel) footerView() string {
	help := "↑↓: scroll | esc: quit | ← →: previous/next chapter"
	info := infoStyle.Render(fmt.Sprintf("%s | %3.f%%", help, m.viewport.ScrollPercent()*100))
	line := strings.Repeat("-", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)

}

func toPreviousChapter() tea.Msg {
	if current.chapterData.Data.Previous.ID == "" {
		return nil
	}
	chapterData := api_query.ChapterQuery(current.translationID, current.chapterData.Data.Previous.ID)
	return chapterData
}

func toNextChapter() tea.Msg {
	if current.chapterData.Data.Next.ID == "" {
		return nil
	}
	chapterData := api_query.ChapterQuery(current.translationID, current.chapterData.Data.Next.ID)
	return chapterData
}

func getStyleDimensions(style lipgloss.Style) (height int, width int) {
	border := style.GetBorderStyle()
	paddingTop, _, _, paddingLeft := style.GetPadding()
	marginTop, marginLeft := style.GetMarginTop(), style.GetMarginLeft()
	return (border.GetTopSize() + paddingTop + marginTop) * 2, (border.GetLeftSize() + paddingLeft + marginLeft) * 2
}
