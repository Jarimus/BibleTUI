package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/Jarimus/BibleTUI/internal/database"
	"github.com/Jarimus/BibleTUI/internal/styles"
	tea "github.com/charmbracelet/bubbletea"
)

const settingsFilePath = "settings.json"
const logFilePath = "error_log.txt"

// Adds the chosen translation to the database for the active user.
// Also sets the current translation as the newly added translation.
func addTranslationToDatabase(translationName, TranslationId string) (database.Translation, tea.Msg) {

	// Create a translation entry in the database
	params := database.CreateTranslationParams{
		Name:   translationName,
		ApiID:  TranslationId,
		UserID: apiCfg.CurrentUserID,
	}
	translation, err := apiCfg.dbQueries.CreateTranslation(context.Background(), params)
	if err != nil {
		return database.Translation{}, err
	}

	return translation, nil
}

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

		log.Print("No settings file found. ")
		time.Sleep(1000 * time.Millisecond)

		apiCfg.CurrentlyReading.TranslationName = "No translation"
		apiCfg.CurrentlyReading.TranslationID = ""
		apiCfg.ApiKey = styles.RedText.Render("Enter your API Key to access the Bible translations!")
		apiCfg.CurrentUser = "Default"
		apiCfg.CurrentUserID = int64(1)

		log.Print("Creating new settings.json...\n")
		time.Sleep(1000 * time.Millisecond)

		saveSettings()

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
		log.Printf("Database file %s does not exist. Creating...", dbFilePath)
		time.Sleep(1500 * time.Millisecond)

		// Create the database file
		file, err := os.Create(dbFilePath)
		if err != nil {
			return err
		}
		file.Close()

		// Open the database
		db, err := sql.Open("sqlite3", dbFilePath)
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

		log.Println("Database initialized successfully!")
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
	db, err := sql.Open("sqlite3", dbFilePath)
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
			return
		}
		file.Close()
	}

	file, err = os.Open(logFilePath)
	if err != nil {
		return
	}

	_, errW := file.Write([]byte(errorText))
	if errW != nil {
		log.Fatalf("error writing errors to log file: %s", err)
	}

}
