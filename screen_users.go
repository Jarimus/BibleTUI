package main

import (
	"context"
	"database/sql"
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
	usersList         []userOption
	choiceIndex       int
	textInput         textinput.Model
	infoText          string
	notificationText  string
	notificationStyle lipgloss.Style
	errorText         string
	confirmDelete     bool
}

// Struct for users menu items.
type userOption struct {
	name    string
	command func(string) tea.Cmd
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
		options = append(options, userOption{
			name:    user.Name,
			command: selectUser,
		})
	}

	// Default options for the users menu
	options = append(options, userOption{name: "Create new user", command: focusInputField})
	options = append(options, userOption{name: "Back"})

	m.usersList = options

	// Settings for the text input
	m.textInput = textinput.New()
	m.textInput.TextStyle = styles.BlueText
	m.textInput.PromptStyle = styles.BlueText
	m.textInput.PlaceholderStyle = styles.BlueText
	m.textInput.CharLimit = 20
	m.textInput.Prompt = "> "
	m.textInput.Placeholder = "Enter user name"

	// Initial info text
	m.infoText = "Press 'Del' to delete a user"

	// Initial error text
	m.errorText = ""

	// Initial notification text
	m.notificationText = ""
	m.notificationStyle = styles.PurpleText

	return m
}

func (m usersMenuModel) Init() tea.Cmd {
	return m.textInput.Cursor.BlinkCmd()
}

func (m usersMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case focusInput:
		m.textInput.Focus()
	case newNotificationMsg:
		m.notificationText = msg.text
		m.notificationStyle = msg.style
	case newErrorMsg:
		m.errorText = msg.text
		return m, nil
	}

	// If text input is not in focus, keyMsg navigate the menu
	if !m.textInput.Focused() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			// Hide notification and error texts when user navigates the menu.
			// Also reset the user deletion confirmation.
			m.notificationText = ""
			m.errorText = ""
			if msg.String() != tea.KeyDelete.String() {
				m.confirmDelete = false
			}

			// Check keystrokes
			switch msg.String() {
			case tea.KeyEsc.String():
				return m, func() tea.Msg { return goBackMsg{} }
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

			// Calls the command for the active list item, goes back a screen if no command is specified.
			case "enter":
				if m.usersList[m.choiceIndex].command == nil {
					return m, func() tea.Msg {
						return goBackMsg{}
					}
				}
				cmd = m.usersList[m.choiceIndex].command(m.usersList[m.choiceIndex].name)
				return m, cmd

			// Pressing delete attempts to delete the highlighted user.
			case tea.KeyDelete.String():

				activeListItem := m.usersList[m.choiceIndex]

				// If the cursor is not on the default menu items, attempt to delete the user from the database and the list.
				if activeListItem.name != "Back" && activeListItem.name != "Create new user" && activeListItem.name != "Default" {

					// First press shows a notification about confirming the deletion.
					if !m.confirmDelete {
						m.errorText = fmt.Sprintf("Are you sure you want to delete user '%s'?", activeListItem.name)
						m.confirmDelete = true
						return m, nil
					}

					// Second press deletes the user.
					if activeListItem.name == apiCfg.CurrentUser {
						apiCfg.CurrentUser = "Default"
						apiCfg.CurrentUserID = 1
						apiCfg.CurrentlyReading = currentlyReading{}
						apiCfg.CurrentlyReading.TranslationName = "No translation"
					}

					err := apiCfg.dbQueries.DeleteUser(context.Background(), activeListItem.name)
					if err != nil {
						m.errorText = err.Error()
						return m, nil
					}

					for i, item := range m.usersList {
						if item.name == activeListItem.name {
							m.usersList = slices.Delete(m.usersList, i, i+1)
						}
					}
					m.errorText = fmt.Sprintf("User deleted: '%s'", activeListItem.name)
					m.confirmDelete = false
				}
			}
		}
	} else { // If text input is in focus, user input controls the text input.
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case tea.KeyEsc.String():
				m.textInput.Blur()
				m.textInput.Reset()
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
					m.errorText = "User already exists."
					return m, nil
				}

				// If not, create a new user with the given name
				userInfo, err := apiCfg.dbQueries.CreateUser(context.Background(), userInput)
				if err != nil {
					m.errorText = err.Error()
					return m, nil
				}
				m.usersList = slices.Insert(m.usersList, 0, userOption{name: userInfo.Name, command: selectUser})
				m.notificationText = fmt.Sprintf("User created: %s", userInfo.Name)
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

	// Notification text
	notificationText := m.notificationText
	notificationText = m.notificationStyle.Render(notificationText)
	notificationText = lipgloss.PlaceHorizontal(window_width, 0.5, notificationText)

	// Help text
	infoText := styles.InfoText.Render(m.infoText)
	infoText = lipgloss.PlaceHorizontal(window_width, 0.5, infoText)

	// Error text
	errorText := styles.RedText.Render(m.errorText)
	errorText = lipgloss.PlaceHorizontal(window_width, 0.5, errorText)

	// Input field, if applicable
	var inputText string
	if m.textInput.Focused() {
		inputText = m.textInput.View()
	} else {
		inputText = ""
	}
	inputText = lipgloss.PlaceHorizontal(window_width, 0.5, inputText)

	// Vertical placement of elements
	inputText = lipgloss.PlaceVertical(2, 1, inputText)
	notificationText = lipgloss.PlaceVertical(2, 1, notificationText)
	errorText = lipgloss.PlaceVertical(2, 1, errorText)
	infoText = lipgloss.PlaceVertical(window_height-lipgloss.Height(uiString)-lipgloss.Height(inputText)-lipgloss.Height(notificationText)-lipgloss.Height(errorText), 1, infoText)

	return lipgloss.JoinVertical(0, uiString, inputText, notificationText, errorText, infoText)
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

func focusInputField(string) tea.Cmd {
	return func() tea.Msg {
		return focusInput{}
	}
}

func selectUser(user string) tea.Cmd {
	return func() tea.Msg {
		// Get user from database
		user, err := apiCfg.dbQueries.GetUser(context.Background(), user)
		if err != nil {
			return func() tea.Msg {
				return newErrorMsg{
					text: err.Error(),
				}
			}
		}

		apiCfg.CurrentUser = user.Name
		apiCfg.CurrentUserID = user.ID

		// Set the current translation to the first translation the user has on their list.
		userTranslations, err := apiCfg.dbQueries.GetTranslationsForUser(context.Background(), user.ID)
		if err == sql.ErrNoRows || len(userTranslations) == 0 {
			apiCfg.CurrentlyReading.TranslationID = ""
			apiCfg.CurrentlyReading.TranslationName = "No translation"
		} else if err != nil {
			return newErrorMsg{
				text: err.Error(),
			}
		} else {
			firstTranslation := userTranslations[0]
			apiCfg.CurrentlyReading.TranslationID = firstTranslation.ApiID
			apiCfg.CurrentlyReading.TranslationName = firstTranslation.Name
		}

		err = saveSettings()
		if err != nil {
			return err
		}
		return newNotificationMsg{
			text:  fmt.Sprintf("Current user set to %s.", apiCfg.CurrentUser),
			style: styles.PurpleText}
	}
}
