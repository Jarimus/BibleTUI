package styles

import "github.com/charmbracelet/lipgloss"

// Styles for the TUI

const GreenColor = lipgloss.Color("#00ff00")
const RedColor = lipgloss.Color("#ff0000")
const BlueColor = lipgloss.Color("#0000ff")
const YellowColor = lipgloss.Color("#ffff00")
const GreyColor = lipgloss.Color("#7777777")

var GreenText = lipgloss.NewStyle().Foreground(lipgloss.Color(GreenColor))
var RedText = lipgloss.NewStyle().Foreground(lipgloss.Color(RedColor))
var BlueText = lipgloss.NewStyle().Foreground(lipgloss.Color(BlueColor))
var YellowText = lipgloss.NewStyle().Foreground(lipgloss.Color(YellowColor))
var InfoText = lipgloss.NewStyle().Foreground(lipgloss.Color(GreyColor))
