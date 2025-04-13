package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// A message that signals the root model to "go back" one screen
type goBackMsg struct{}

type newNotificationMsg struct {
	text  string
	style lipgloss.Style
}

type newErrorMsg struct {
	text string
}

type focusInput struct{}

type audioDone struct{}

var backCmd = func() tea.Msg { return goBackMsg{} }
