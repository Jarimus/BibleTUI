package main

import (
	"slices"

	"github.com/Jarimus/BibleTUI/internal/api_query"
	tea "github.com/charmbracelet/bubbletea"
)

type rootScreenModel struct {
	models []tea.Model
}

func newRootScreen(models []tea.Model) rootScreenModel {
	return rootScreenModel{
		models: models,
	}
}

func (m rootScreenModel) Init() tea.Cmd {
	return m.models[len(m.models)-1].Init()
}

func (m rootScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.Model:
		m.models = append(m.models, msg)
		return m, nil
	case goBackMsg:
		if len(m.models) == 1 {
			return m, tea.Quit
		}
		m = m.goBack()
	case tea.WindowSizeMsg:
		window_height = msg.Height
		window_width = msg.Width
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if len(m.models) == 1 {
				return m, tea.Quit
			}
			m = m.goBack()
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case api_query.TranslationData:
		current.translationData = msg
	}

	currentModel := m.models[len(m.models)-1]
	updatedModel, cmd := currentModel.Update(msg)
	m.models[len(m.models)-1] = updatedModel

	return m, cmd
}

func (m rootScreenModel) View() string {
	return m.models[len(m.models)-1].View() // + fmt.Sprintf("\n\nModel count: %d", len(m.models))
}

type goBackMsg struct{}

func (m rootScreenModel) goBack() rootScreenModel {
	if len(m.models) > 1 {
		m.models = slices.Delete(m.models, len(m.models)-1, len(m.models))
		return m
	} else {
		return m
	}
}
