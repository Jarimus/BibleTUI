package main

import (
	"fmt"
	"strings"

	"github.com/Jarimus/BibleTUI/internal/api_query"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
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
	fullQuery := api_query.BibleChapterQuery()
	switch fullQuery := fullQuery.(type) {
	case api_query.BookQuery:
		bibleText := fullQuery.Text
		title := fmt.Sprintf("%s: %d", fullQuery.Verses[0].BookName, fullQuery.Verses[0].Chapter)
		newBibleScreen := bibleScreenModel{
			title:     title,
			bibleText: bibleText,
		}

		// Generate a viewport from the dimensions of the global variables
		headerHeight := lipgloss.Height(newBibleScreen.headerView())
		footerHeight := lipgloss.Height(newBibleScreen.footerView())
		verticalMarginHeight := headerHeight + footerHeight
		newBibleScreen.viewport = viewport.New(window_width, window_height-verticalMarginHeight)
		newBibleScreen.viewport.YPosition = headerHeight
		newBibleScreen.viewport.SetContent(newBibleScreen.bibleText)

		// Return the model
		return newBibleScreen
	}
	return bibleScreenModel{}
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
		verticalMarginHeight := headerHeight + footerHeight
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - verticalMarginHeight

	// Retrieving the Bible chapter updates the viewport text
	case api_query.BookQuery:
		m.title = fmt.Sprintf("%s: %d", msg.Verses[0].BookName, msg.Verses[0].Chapter)
		m.bibleText = msg.Text
		m.viewport.SetContent(m.bibleText)
		return m, tea.Quit
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m bibleScreenModel) View() string {
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m bibleScreenModel) headerView() string {
	title := titleStyle.Render(m.title)
	line := strings.Repeat("-", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m bibleScreenModel) footerView() string {
	help := "↑↓: scroll | ctrl+c, esc: quit | ← →: next chapter"
	info := infoStyle.Render(fmt.Sprintf("%s | %3.f%%", help, m.viewport.ScrollPercent()*100))
	line := strings.Repeat("-", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)

}
