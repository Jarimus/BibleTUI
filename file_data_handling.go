package main

import (
	"encoding/json"
	"log"
	"os"
	"sort"
)

type translationJsonItems struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type settings struct {
	TranslationName string `json:"name"`
	TranslationID   string `json:"id"`
}

const translationsFilePath = "translations.json"
const settingsFilePath = "settings.json"

// Look for the file "translations.json".
// If not found, create a set of basic set of items for translations menu
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

		marshalledTranslation, err := json.Marshal(translations)
		if err != nil {
			return nil, err
		}

		err = os.WriteFile(translationsFilePath, marshalledTranslation, 0644)
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

	dataMarshalled, err := json.Marshal(translations)
	if err != nil {
		return err
	}

	os.WriteFile(translationsFilePath, dataMarshalled, 0644)
	return nil
}

func saveTranslationsToFile(translationsMenuItems []translationMenuItem) error {
	var translations []translationJsonItems
	for _, translation := range translationsMenuItems {
		translations = append(translations, translationJsonItems{
			Name: translation.name,
			Id:   translation.id,
		})
	}
	dataMarshalled, err := json.Marshal(translations)
	if err != nil {
		return err
	}

	os.WriteFile(translationsFilePath, dataMarshalled, 0644)
	return nil
}

func loadSettings() (settings, error) {

	var settingsData settings

	fileData, err := os.ReadFile(settingsFilePath)
	if err != nil {
		settingsData.TranslationName = "Finnish New Testament"
		settingsData.TranslationID = "c739534f6a23acb2-01"
		return settingsData, nil
	}

	err = json.Unmarshal(fileData, &settingsData)
	if err != nil {
		return settings{}, err
	}

	return settingsData, nil
}

func saveSettings() error {

	settingsData := settings{
		TranslationName: apiCfg.CurrentlyReading.TranslationName,
		TranslationID:   apiCfg.CurrentlyReading.TranslationID,
	}

	jsonData, err := json.Marshal(settingsData)
	if err != nil {
		return err
	}

	os.WriteFile(settingsFilePath, jsonData, 0644)

	return nil
}
