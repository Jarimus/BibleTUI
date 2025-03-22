package api_query

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Query the api for a translation.
// The response includes data for the books in the translation: id, name and other data
func TranslationQuery(translationID string) TranslationData {

	url := fmt.Sprintf("https://api.scripture.api.bible/v1/bibles/%s/books", translationID)

	// Set up the request with a corrent header (API-key needed for API-Bible)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request for translation data: %s", err)
		return TranslationData{}
	}

	req.Header.Set("api-key", getApiKey())

	// Perform the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error requesting data for translation: %s", err)
		return TranslationData{}
	}

	// Read the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading the body for translation query: %s", err)
		return TranslationData{}
	}
	defer resp.Body.Close()

	// Unmarshal the response
	var translationData TranslationData
	if err := json.Unmarshal(body, &translationData); err != nil {
		log.Printf("Error unmarshaling body for translation query: %s", err)
	}

	return translationData

}

func AllTranslationsQuery(languageID string) BiblesData {
	var url string
	if languageID == "" {
		url = "https://api.scripture.api.bible/v1/bibles"
	} else {
		url = fmt.Sprintf("https://api.scripture.api.bible/v1/bibles?language=%s", languageID)
	}

	// Set up the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request for all Bibles: %s", err)
		return BiblesData{}
	}
	req.Header.Set("api-key", getApiKey())

	// Perform the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error requesting data for translation: %s", err)
		return BiblesData{}
	}

	// Read the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading the body for Bibles query: %s", err)
		return BiblesData{}
	}
	defer resp.Body.Close()

	// Unmarshal the response body
	var biblesData BiblesData
	err = json.Unmarshal(body, &biblesData)
	if err != nil {
		log.Printf("error unmarshaling Bibles data response: %s", err)
	}

	return biblesData
}
