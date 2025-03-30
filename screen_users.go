package main

import (
	"context"
	"log"
	"strings"

	styles "github.com/Jarimus/BibleTUI/internal/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// tea.Model for the users menu (to select and create new users).
type usersMenuModel struct {
	usersList   []userOption
	choiceIndex int
	textInput   textinput.Model
}

// Struct for users menu items.
type userOption struct {
	name    string
	command tea.Cmd
}

// Returns a tea.Model for the users menu.
func newUsersMenu() usersMenuModel {
	var m usersMenuModel

	// Get users from the database
	users, err := apiCfg.dbQueries.GetAllUsers(context.Background())
	if err != nil {
		log.Fatalf("error getting users from database: %s", err)
	}

	// Append users to the list
	var options []userOption
	for _, user := range users {
		options = append(options, userOption{name: user.Name})
	}

	// Default options for the users menu
	options = append(options, userOption{name: "default", command: func() tea.Msg { return goBackMsg{} }})
	options = append(options, userOption{name: "create new user", command: func() tea.Msg { return goBackMsg{} }})
	options = append(options, userOption{name: "back"})

	m.usersList = options

	// Settings for the text input
	m.textInput.CharLimit = 20
	m.textInput.EchoCharacter = '_'
	m.textInput.Placeholder = "enter user name"
	m.textInput.PlaceholderStyle = styles.InfoText
	m.textInput.TextStyle = styles.BlueText

	return m
}

func (m usersMenuModel) Init() tea.Cmd {
	return nil
}

func (m usersMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Similar to the main menu
	// The model handles navigating the menu
	// Calling a choice initiates a query to retrieve data for the chosen book
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.choiceIndex = (m.choiceIndex - 1 + len(m.usersList)) % len(m.usersList)
			return m, nil
		case tea.KeyPgUp.String(), "left":
			m.choiceIndex = max(0, m.choiceIndex-10)
		case "down":
			m.choiceIndex = (m.choiceIndex + 1) % len(m.usersList)
			return m, nil
		case tea.KeyPgDown.String(), "right":
			m.choiceIndex = min(len(m.usersList)-1, m.choiceIndex+10)
		case "enter":
			if m.usersList[m.choiceIndex].command == nil {
				return m, func() tea.Msg {
					return goBackMsg{}
				}
			}
			return m, m.usersList[m.choiceIndex].command
		}
	}
	return m, nil
}

func (m usersMenuModel) View() string {

	uiString := getHeaderWithList(m)

	// Help text
	helpText := styles.InfoText.Render("Press 'Del' to delete a user")
	helpText = lipgloss.PlaceHorizontal(window_width, 0.5, helpText)

	// Input field, if applicable
	var inputText string
	if m.textInput.Focused() {
		inputText = m.textInput.View()
	} else {
		inputText = ""
	}
	inputText = lipgloss.PlaceHorizontal(window_width, 0.5, inputText)

	// Vertical placement of elements
	inputText = lipgloss.PlaceVertical(window_height-lipgloss.Height(helpText)-lipgloss.Height(uiString), 0.5, inputText)
	helpText = lipgloss.PlaceVertical(window_height-lipgloss.Height(inputText)-lipgloss.Height(uiString), 1, helpText)

	return lipgloss.JoinVertical(0, uiString, inputText, helpText)
}

// Returns the header of the model as a string.
func (m usersMenuModel) headerView() string {
	// Styling for the header
	topMsg := "* Choose a user or create a new one *"
	topBottomBar := styles.YellowText.Render(strings.Repeat("*", len(topMsg)))
	topMsg = styles.YellowText.Render(topMsg)
	return lipgloss.JoinVertical(lipgloss.Center, topBottomBar, topMsg, topBottomBar)
}

// Returns the name of the menu item at index as a string.
func (m usersMenuModel) getName(index int) string {
	return m.usersList[index].name
}

// Return the length of the menu list as an integer.
func (m usersMenuModel) getListLength() int {
	return len(m.usersList)
}

// Return the current cursor placement as an integer.
func (m usersMenuModel) getChoiceIndex() int {
	return m.choiceIndex
}
