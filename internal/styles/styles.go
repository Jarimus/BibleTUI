package styles

import "github.com/charmbracelet/lipgloss"

// Styles for the TUI

const GreenColor = lipgloss.Color("#00ff00")
const RedColor = lipgloss.Color("#ff0000")
const BlueColor = lipgloss.Color("#0000ff")
const YellowColor = lipgloss.Color("#ffff00")
const GreyColor = lipgloss.Color("#999999")

var GreenText = lipgloss.NewStyle().Foreground(lipgloss.Color(GreenColor))
var RedText = lipgloss.NewStyle().Foreground(lipgloss.Color(RedColor))
var BlueText = lipgloss.NewStyle().Foreground(lipgloss.Color(BlueColor))
var YellowText = lipgloss.NewStyle().Foreground(lipgloss.Color(YellowColor))
var InfoText = lipgloss.NewStyle().Foreground(lipgloss.Color(GreyColor))

func GetStyleDimensions(style lipgloss.Style) (height int, width int) {
	// Get the combined size of borders, paddings and margins of a style. Used for the viewport style.
	// Why does style.GetFrameSize() not work directly? Why do I have to do this manually?

	// Border dimensions
	border := style.GetBorderStyle()
	borderV := border.GetTopSize() + border.GetBottomSize()
	borderH := border.GetLeftSize() + border.GetRightSize()

	// Padding dimensions
	paddingTop, paddingRight, paddingBottom, paddingLeft := style.GetPadding()
	paddingV := paddingTop + paddingBottom
	paddingH := paddingRight + paddingLeft

	// Margin dimensions
	marginV := style.GetMarginTop() + style.GetMarginBottom()
	marginH := style.GetMarginLeft() + style.GetMarginRight()

	// Return all the dimensions combined
	return borderV + paddingV + marginV, borderH + paddingH + marginH
}
