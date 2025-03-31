package main

import "github.com/charmbracelet/lipgloss"

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