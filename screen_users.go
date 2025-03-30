package main

import (
	"context"
	"fmt"
	"log"
	"slices"
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
	helpText    string
	errorText   string
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
	options = append(options, userOption{name: "create new user", command: func() tea.Msg { return "focus" }})
	options = append(options, userOption{name: "back"})

	m.usersList = options

	// Settings for the text input
	m.textInput = textinput.New()
	m.textInput.Cursor.Style = styles.BlueText
	m.textInput.CharLimit = 20
	m.textInput.EchoCharacter = '_'
	m.textInput.Placeholder = "enter user name"
	m.textInput.PlaceholderStyle = styles.InfoText
	m.textInput.TextStyle = styles.BlueText

	// Initial help text
	m.helpText = "Press 'Del' to delete a user"

	// Initial error text
	m.errorText = "test"

	return m
}

func (m usersMenuModel) Init() tea.Cmd {
	return m.textInput.Cursor.BlinkCmd()
}

func (m usersMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// If text input is not in focus, keyMsg navigate the menu
	if !m.textInput.Focused() {
		switch msg := msg.(type) {
		case string:
			if msg == "focus" {
				m.textInput.Focus()
			}
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
			case "focus":
				m.textInput.Focus()
			}
		}
	} else { // If text input is in focus, user input controls the text input.
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up":
				m.textInput.CursorStart()
				return m, nil
			case "down":
				m.textInput.CursorEnd()
				return m, nil
			case "enter":
				userInput := m.textInput.Value()

				// Check whether the user already exists in the database
				dbUser, _ := apiCfg.dbQueries.GetUser(context.Background(), userInput)
				if dbUser.Name == userInput {
					m.errorText = "User already exists"
					return m, nil
				}

				// If not, create a new user with the given name
				userInfo, err := apiCfg.dbQueries.CreateUser(context.Background(), userInput)
				if err != nil {
					m.errorText = err.Error()
					return m, nil
				}
				m.usersList = slices.Insert(m.usersList, 0, userOption{name: userInfo.Name})
				m.helpText = fmt.Sprintf("User created:%s\n", userInfo.Name)
				m.textInput.Reset()
				m.textInput.Blur()
			}
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m usersMenuModel) View() string {

	uiString := getHeaderWithList(m)

	// Help text
	helpText := styles.InfoText.Render(m.helpText)
	helpText = lipgloss.PlaceHorizontal(window_width, 0.5, helpText)

	// Error text
	errorText := styles.RedText.Render(m.errorText)
	errorText = lipgloss.PlaceHorizontal(window_width, 0.5, errorText)

	// Input field, if applicable
	var inputText string
	if m.textInput.Focused() {
		inputText = m.textInput.View()
	} else {
		inputText = "New user: "
	}
	inputText = lipgloss.PlaceHorizontal(window_width, 0.5, inputText)

	// Vertical placement of elements
	inputText = lipgloss.PlaceVertical(window_height-lipgloss.Height(helpText)-lipgloss.Height(uiString)-lipgloss.Height(errorText), 0.5, inputText)
	helpText = lipgloss.PlaceVertical(window_height-lipgloss.Height(inputText)-lipgloss.Height(uiString)-lipgloss.Height(errorText), 1, helpText)
	errorText = lipgloss.PlaceVertical(window_height-lipgloss.Height(inputText)-lipgloss.Height(uiString)-lipgloss.Height(helpText), 1, errorText)

	return lipgloss.JoinVertical(0, uiString, inputText, errorText, helpText)
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
