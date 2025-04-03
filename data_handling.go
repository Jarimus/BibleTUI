package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Jarimus/BibleTUI/internal/database"
	"github.com/Jarimus/BibleTUI/internal/styles"
	tea "github.com/charmbracelet/bubbletea"
)

const settingsFilePath = "settings.json"
const logFilePath = "error_log.txt"

func GetTranslationForUserById(translationID string) (database.Translation, error) {
	params := database.GetTranslationForUserByIdParams{
		ApiID:  translationID,
		UserID: apiCfg.CurrentUserID,
	}
	translation, err := apiCfg.dbQueries.GetTranslationForUserById(context.Background(), params)
	if err != nil {
		return database.Translation{}, err
	}
	return translation, nil
}

// Adds the chosen translation to the database for the active user.
// Also sets the current translation as the newly added translation.
func addTranslationToDatabase(translationName, translationId, languageID string) (database.Translation, tea.Msg) {

	// First check if the translation already exists
	getParams := database.GetTranslationForUserByIdParams{
		ApiID:  translationId,
		UserID: apiCfg.CurrentUserID,
	}
	dbTranslation, err := apiCfg.dbQueries.GetTranslationForUserById(context.Background(), getParams)
	if err != nil && err != sql.ErrNoRows {
		return database.Translation{}, err
	}
	if dbTranslation.ApiID == translationId {
		return database.Translation{}, errors.New("translation already added for the current user")
	}

	// Create a translation entry in the database
	createParams := database.CreateTranslationParams{
		Name:       translationName,
		ApiID:      translationId,
		LanguageID: languageID,
		UserID:     apiCfg.CurrentUserID,
	}
	translation, err := apiCfg.dbQueries.CreateTranslation(context.Background(), createParams)
	if err != nil {
		return database.Translation{}, err
	}

	return translation, nil
}

// Return all translations for the current user
func loadTranslationsForUser() ([]database.Translation, error) {

	translations, err := apiCfg.dbQueries.GetTranslationsForUser(context.Background(), apiCfg.CurrentUserID)
	if err != nil {
		return nil, err
	}

	return translations, nil
}

// Loads and return the apiCfg from a json-file.
// If file is not found, returns an empty config file with a default Bible translation as the current translation.
func loadSettings() error {

	fileData, err := os.ReadFile(settingsFilePath)

	if err != nil {

		println("No settings file found. ")
		time.Sleep(1000 * time.Millisecond)

		apiCfg.CurrentlyReading.TranslationName = "No translation"
		apiCfg.CurrentlyReading.TranslationID = ""
		apiCfg.ApiKey = styles.RedText.Render("Enter your API Key to access the Bible translations!")
		apiCfg.CurrentUser = "Default"
		apiCfg.CurrentUserID = int64(1)

		println("Creating new settings.json...")

		err = saveSettings()
		if err != nil {
			return err
		}

		time.Sleep(1000 * time.Millisecond)

		return nil
	}

	err = json.Unmarshal(fileData, &apiCfg)
	if err != nil {
		return err
	}

	return nil
}

// Saves the apiCfg to a json-file.
func saveSettings() error {

	jsonData, err := json.MarshalIndent(apiCfg, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(settingsFilePath, jsonData, 0600)
	if err != nil {
		return err
	}

	return nil
}

// Initializes the database, creating a .db-file if necessary
func initializeDB() error {

	// If the database file does not exist, create a new one.
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		fmt.Printf("Database file %s does not exist. Creating...\n", dbFilePath)
		time.Sleep(1500 * time.Millisecond)

		// Create the database file
		file, err := os.Create(dbFilePath)
		if err != nil {
			return err
		}
		err = file.Close()
		if err != nil {
			return err
		}

		// Open the database
		db, err := sql.Open("sqlite", dbFilePath)
		if err != nil {
			return err
		}

		// Apply the schema
		if _, err := db.Exec(usersSchema); err != nil {
			return err
		}
		if _, err := db.Exec(translationsSchema); err != nil {
			return err
		}

		println("Database initialized successfully!")
		time.Sleep(1500 * time.Millisecond)

		dbQueries := database.New(db)
		apiCfg.dbQueries = dbQueries

		// Create a 'Default' user
		if _, err := apiCfg.dbQueries.CreateUser(context.Background(), "Default"); err != nil {
			return err
		}

		return nil

	}

	// If the database exists, open the existing database and store the connection in the database
	db, err := sql.Open("sqlite", dbFilePath)
	if err != nil {
		return err
	}

	dbQueries := database.New(db)
	apiCfg.dbQueries = dbQueries

	return nil
}

func logError(err error) {

	var file *os.File

	errorText := err.Error()

	if _, err = os.Stat(logFilePath); os.IsNotExist(err) {

		// Create the log file
		file, err := os.Create(logFilePath)
		if err != nil {
			logError(err)
			return
		}
		err = file.Close()
		if err != nil {
			logError(err)
			return
		}
	}

	file, err = os.Open(logFilePath)
	if err != nil {
		logError(err)
		return
	}

	_, errW := file.Write([]byte(errorText))
	if errW != nil {
		log.Fatalf("error writing errors to log file: %s", err)
	}
}
