package main

import (
	"encoding/json"
	"log"
	"os"
	"sort"
)

// Look for the file "data/translations.json".
// If not found, create a set of basic set of items for translations menu
func LoadTranslationsFromFile() []translationMenuItem {

	type translationJsonItems struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	}

	data, err := os.ReadFile("data/translations.json")
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
			log.Printf("error marshaling translations to write to a file: %s", err)
		}

		err = os.WriteFile("data/translations.json", marshalledTranslation, 0644)
		if err != nil {
			log.Printf("error writing to file 'data/translations.json': %s", err)
		}

		var resultVals []translationMenuItem

		for _, item := range translations {
			resultVals = append(resultVals, translationMenuItem{
				name:    item.Name,
				id:      item.Id,
				command: selectTranslation,
			})
		}

		return resultVals
	}

	var translations []translationJsonItems
	err = json.Unmarshal(data, &translations)
	if err != nil {
		log.Printf("error unmarshaling translations data from file: %s", err)
	}

	var resultVals []translationMenuItem

	for _, item := range translations {
		resultVals = append(resultVals, translationMenuItem{
			name:    item.Name,
			id:      item.Id,
			command: selectTranslation,
		})
	}

	return resultVals
}

func SaveTranslationsToFile(translations []translationMenuItem) error {
	return nil
}
