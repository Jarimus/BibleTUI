package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"sort"
)

type translationJsonItems struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

const translationsFilePath = "translations.json"
const settingsFilePath = "settings.json"

// Look for the translation json-file.
// If not found, create a set of basic set of items for the translations menu.
func loadTranslationsFromFile() ([]translationMenuItem, error) {

	data, err := os.ReadFile(translationsFilePath)
	if err != nil {

		translations := []translationJsonItems{
			{
				Name: "Simplified Chinese",
				Id:   "7ea794434e9ea7ee-01",
			},
			{
				Name: "Finnish New Testament",
				Id:   "c739534f6a23acb2-01",
			},
			{
				Name: "American Standard",
				Id:   "685d1470fe4d5c3b-01",
			},
			{
				Name: "King James",
				Id:   "de4e12af7f28f599-01",
			},
			{
				Name: "World English Bible",
				Id:   "9879dbb7cfe39e4d-01",
			},
			{
				Name: "Open Hebrew Living New Testament",
				Id:   "a8a97eebae3c98e4-01",
			},
			{
				Name: "Brenton Greek Septuagint",
				Id:   "c114c33098c4fef1-01",
			},
			{
				Name: "Solid Rock Greek New Testament",
				Id:   "47f396bad37936f0-01",
			},
		}
		sort.Slice(translations, func(i, j int) bool {
			return translations[i].Name < translations[j].Name
		})

		marshalledTranslation, err := json.MarshalIndent(translations, "", "  ")
		if err != nil {
			return nil, err
		}

		err = os.WriteFile(translationsFilePath, marshalledTranslation, 0600)
		if err != nil {
			return nil, err
		}

		var resultVals []translationMenuItem

		for _, item := range translations {
			resultVals = append(resultVals, translationMenuItem{
				name:    item.Name,
				id:      item.Id,
				command: selectTranslation,
			})
		}

		return resultVals, nil
	}

	var translations []translationJsonItems
	err = json.Unmarshal(data, &translations)
	if err != nil {
		log.Printf("error unmarshaling translations data from file: %s", err)
	}

	sort.Slice(translations, func(i, j int) bool {
		return translations[i].Name < translations[j].Name
	})

	var resultVals []translationMenuItem

	for _, item := range translations {
		resultVals = append(resultVals, translationMenuItem{
			name:    item.Name,
			id:      item.Id,
			command: selectTranslation,
		})
	}

	return resultVals, nil
}

// Loads the current list of translation from a file, adds the new translation, and writes the new list to the file.
func addTranslationToFile(translationName, translationId string) error {

	translationMenuItems, err := loadTranslationsFromFile()
	if err != nil {
		return err
	}

	var translations []translationJsonItems
	for _, translation := range translationMenuItems {
		translations = append(translations, translationJsonItems{
			Name: translation.name,
			Id:   translation.id,
		})
	}
	translations = append(translations, translationJsonItems{
		Name: translationName,
		Id:   translationId,
	})

	dataMarshalled, err := json.MarshalIndent(translations, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(translationsFilePath, dataMarshalled, 0600)
	if err != nil {
		return err
	}

	return nil
}

// Saves the list of translations (open in the current screen) to a file.
func saveTranslationsToFile(translationsMenuItems []translationMenuItem) error {
	var translations []translationJsonItems
	for _, translation := range translationsMenuItems {
		translations = append(translations, translationJsonItems{
			Name: translation.name,
			Id:   translation.id,
		})
	}
	dataMarshalled, err := json.MarshalIndent(translations, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(translationsFilePath, dataMarshalled, 0600)
	if err != nil {
		return err
	}
	return nil
}

// Loads and return the apiCfg from a json-file.
// If file is not found, returns an empty config file with a default Bible translation as the current translation.
func loadSettings() (config, error) {

	var cfg config

	fileData, err := os.ReadFile(settingsFilePath)
	if err != nil {
		cfg.CurrentlyReading.TranslationName = "Finnish New Testament"
		cfg.CurrentlyReading.TranslationID = "c739534f6a23acb2-01"
		cfg.apiKey = getApiKey()
		return cfg, nil
	}

	err = json.Unmarshal(fileData, &cfg)
	if err != nil {
		var emptyConfig config
		return emptyConfig, err
	}

	return cfg, nil
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
func initializeDB() (*sql.DB, error) {
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		log.Printf("Database file %s does not exist. Creating...", dbFilePath)

		// Create the database file
		file, err := os.Create(dbFilePath)
		if err != nil {
			return nil, err
		}
		file.Close()

		// Open the database
		db, err := sql.Open("sqlite3", dbFilePath)
		if err != nil {
			return nil, err
		}

		// Apply the schema
		if _, err := db.Exec(usersSchema); err != nil {
			return nil, err
		}
		if _, err := db.Exec(translationsSchema); err != nil {
			return nil, err
		}

		log.Println("Database initialized successfully.")
		return db, nil
	}

	// If the database exists, open the existing database
	return sql.Open("sqlite3", dbFilePath)
}
