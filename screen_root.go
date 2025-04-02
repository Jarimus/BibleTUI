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

// Return a tea.Model that serves as the root for other models (screens).
// It stores the other models in a slice and displays the one at the top of stack (last one).
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
	// New screens need to be returned as tea.Msg so that the root model can grab them and append them to the list of models.
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
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case api_query.ChapterData:
		apiCfg.CurrentlyReading.ChapterData = msg
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

// When the root receives a goBackMsg structm goBack function is called.
// The function removes the top model (screen) from the slice of models, effectively moving the interface to the next model in the stack.
// If there is only one model in the stack, the program quits.
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
