package main

import (
	"slices"

	"github.com/Jarimus/BibleTUI/internal/api_query"
	tea "github.com/charmbracelet/bubbletea"
)

// The root model holds the other models and display the latest one
type rootScreenModel struct {
	models []tea.Model
}

func newRootScreen(models []tea.Model) rootScreenModel {
	// Initialize with a slice of models
	return rootScreenModel{
		models: models,
	}
}

func (m rootScreenModel) Init() tea.Cmd {
	// Root initializes the current model
	return m.models[len(m.models)-1].Init()
}

func (m rootScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Root model handles tea.Msg's that affect which model is being displayed:
	// new models are added to the stack and some tea.Msg's remove the top model from the stack
	switch msg := msg.(type) {
	case tea.Model:
		m.models = append(m.models, msg)
		return m, nil
	case goBackMsg:
		return m.goBack()
	case tea.WindowSizeMsg:
		window_height = msg.Height
		window_width = msg.Width
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m.goBack()
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case api_query.ChapterData:
		current.chapterData = msg
	}

	// Get the current model from the top of the model stack and call its Update method, and update it
	currentModel := m.models[len(m.models)-1]
	updatedModel, cmd := currentModel.Update(msg)
	m.models[len(m.models)-1] = updatedModel

	return m, cmd
}

func (m rootScreenModel) View() string {
	// Root calls the model at the top of the stack to view
	return m.models[len(m.models)-1].View() // + fmt.Sprintf("\n\nModel count: %d", len(m.models))
}

type goBackMsg struct{}

func (m rootScreenModel) goBack() (rootScreenModel, tea.Cmd) {
	// a goBackMsg from other models removes the top model from the stack
	// If there is only one model on the stack, quit the program
	if len(m.models) > 1 {
		m.models = slices.Delete(m.models, len(m.models)-1, len(m.models))
		return m, nil
	} else {
		return m, tea.Quit
	}
}
