package main

import (
	"fmt"
	"strings"

	"github.com/Jarimus/BibleTUI/internal/api_query"
	"github.com/Jarimus/BibleTUI/internal/styles"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	wordwrap "github.com/mitchellh/go-wordwrap"
)

var (
	// Style for the title at the top of the viewport
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	// Styling the viewport
	vpStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 2)
	}()

	// Styling the bottom info panel
	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

// tea.Model for the reading screen
type bibleScreenModel struct {
	title     string
	bibleText string
	viewport  viewport.Model
}

// Initiates a new screen with a viewport to read the current Bible translation.
func newBibleScreen() bibleScreenModel {

	// Apply the chapter content and title to the model
	newBibleScreen := bibleScreenModel{
		title:     apiCfg.CurrentlyReading.ChapterData.Data.Reference,
		bibleText: apiCfg.CurrentlyReading.ChapterData.Data.Content,
	}

	// Generate a viewport from the dimensions of the global variables
	// Take into account the header, the footer and the styling for the viewport
	headerHeight := lipgloss.Height(newBibleScreen.headerView())
	footerHeight := lipgloss.Height(newBibleScreen.footerView())
	vpStyleHeight, vpStyleWidth := styles.GetStyleDimensions(vpStyle)
	verticalMargin := headerHeight + footerHeight + vpStyleHeight
	horizontalMargin := vpStyleWidth
	newBibleScreen.viewport = viewport.New(window_width-horizontalMargin, window_height-verticalMargin)

	// Move the viewport into position below the header
	newBibleScreen.viewport.YPosition = headerHeight

	// Apply the formatted Bible text to the viewport
	newBibleScreen.viewport.SetContent(formatBibleText(newBibleScreen.bibleText, newBibleScreen.viewport.Width))

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

	case error:
		m.viewport.SetContent(msg.Error())

	case audioDone:
		if apiCfg.CurrentlyReading.audioStop != nil {
			apiCfg.CurrentlyReading.audioStop <- true
			close(apiCfg.CurrentlyReading.audioStop)
			apiCfg.CurrentlyReading.audioStop = nil
		}

	// Window resize affects the viewport dimensions
	// The text needs to be reformatted for the new dimensions to get the linebreaks right
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		vpStyleHeight, vpStyleWidth := styles.GetStyleDimensions(vpStyle)
		verticalMargin := headerHeight + footerHeight + vpStyleHeight
		horizontalMargin := vpStyleWidth
		m.viewport.Width = msg.Width - horizontalMargin
		m.viewport.Height = msg.Height - verticalMargin

		m.viewport.SetContent(formatBibleText(m.bibleText, m.viewport.Width))

	// Retrieving the Bible chapter updates the viewport text
	case api_query.ChapterData:
		m.title = apiCfg.CurrentlyReading.ChapterData.Data.Reference
		m.bibleText = apiCfg.CurrentlyReading.ChapterData.Data.Content
		m.viewport.SetContent(formatBibleText(m.bibleText, m.viewport.Width))
		m.viewport.YOffset = 0

	// Pressing left and right moves between the chapters. Esc exits the reader
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEsc.String():
			cmds = append(cmds, stopAudio, backCmd)
		case tea.KeyLeft.String():
			cmds = append(cmds, stopAudio, toPreviousChapter)
		case tea.KeyRight.String():
			cmds = append(cmds, stopAudio, toNextChapter)
			// case "p":
			// 	if apiCfg.CurrentlyReading.audioStop != nil {
			// 		cmds = append(cmds, stopAudio)
			// 	} else {
			// 		cmds = append(cmds, playAudio(m.bibleText, apiCfg.CurrentlyReading.LanguageID))
			// 	}
		}

	}

	// update the viewport view, get commands
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// Style and render the viewport with the header and the footer
func (m bibleScreenModel) View() string {
	text := vpStyle.Render(m.viewport.View())
	return lipgloss.JoinVertical(lipgloss.Top, m.headerView(), text, m.footerView())
}

// Style and render the header
func (m bibleScreenModel) headerView() string {
	title := titleStyle.Render(m.title)
	line := strings.Repeat("-", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

// Style and render the footer
func (m bibleScreenModel) footerView() string {
	languageCode := "notImplemented"
	var help string
	var audioHelpText string
	if languageCode == "" {
		help = "↑↓: scroll | ← →: previous/next chapter | esc: quit"
	} else {
		if apiCfg.CurrentlyReading.audioStop != nil {
			audioHelpText = styles.GreenText.Render("p: stop audio")
		} else {
			audioHelpText = "p: play audio"
		}
		help = " | ↑↓: scroll | ← →: previous/next chapter | esc: quit"
	}

	info := infoStyle.Render(fmt.Sprintf("%s%s | %3.f%%", audioHelpText, help, m.viewport.ScrollPercent()*100))
	line := strings.Repeat("-", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)

}

// Query the previous chapter and return it as a tea.Msg for the model to process
func toPreviousChapter() tea.Msg {
	if apiCfg.CurrentlyReading.ChapterData.Data.Previous.ID == "" {
		return nil
	}
	chapterData := api_query.ChapterQuery(apiCfg.CurrentlyReading.TranslationID, apiCfg.CurrentlyReading.ChapterData.Data.Previous.ID, apiCfg.ApiKey)
	return chapterData
}

// Query the next chapter and return it as a tea.Msg for the model to process
func toNextChapter() tea.Msg {
	if apiCfg.CurrentlyReading.ChapterData.Data.Next.ID == "" {
		return nil
	}
	chapterData := api_query.ChapterQuery(apiCfg.CurrentlyReading.TranslationID, apiCfg.CurrentlyReading.ChapterData.Data.Next.ID, apiCfg.ApiKey)
	return chapterData
}

// formats the Bible text for the viewport. Linebreaks are needed for the viewport to handle scrolling properly.
func formatBibleText(text string, width int) string {

	if width < 0 {
		return text // Ensure width is non-negative to avoid overflow
	}
	formattedText := strings.ReplaceAll(text, "[", "\n[")
	formattedText = wordwrap.WrapString(formattedText, uint(width))

	return formattedText
}

// Plays the input text in the input language, if supported.
// Also listens for a stop signal in the audioStop channel.
// func playAudio(text, language string) tea.Cmd {
// 	return func() tea.Msg {
// 		// Create a channel to listen for when to stop playback.
// 		audioStop := make(chan bool, 1)
// 		apiCfg.CurrentlyReading.audioStop = audioStop
// 		err := tts.SpeakText(text, language, audioStop)
// 		if err != nil {
// 			return err
// 		}
// 		return audioDone{}
// 	}

// }

// Signals for the audio to stop in the audioStop channel.
// Also closes the channel and sets it to nil.
func stopAudio() tea.Msg {
	if apiCfg.CurrentlyReading.audioStop == nil {
		return nil
	}
	apiCfg.CurrentlyReading.audioStop <- true
	close(apiCfg.CurrentlyReading.audioStop)
	apiCfg.CurrentlyReading.audioStop = nil
	return nil
}
