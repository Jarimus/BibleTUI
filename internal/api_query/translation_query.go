package api_query

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Query the api for a translation.
// The response includes data for the books in the translation: id, name and other data
func TranslationQuery(translationID, apiKey string) (TranslationData, error) {

	url := fmt.Sprintf("https://api.scripture.api.bible/v1/bibles/%s/books", translationID)

	// Set up the request with a corrent header (API-key needed for API-Bible)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return TranslationData{}, err
	}

	req.Header.Set("api-key", apiKey)

	// Perform the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return TranslationData{}, err
	}

	// Read the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TranslationData{}, err
	}
	defer resp.Body.Close()

	// Unmarshal the response
	var translationData TranslationData
	if err := json.Unmarshal(body, &translationData); err != nil {
		return TranslationData{}, err
	}

	return translationData, nil

}

func AllTranslationsQuery(languageID, apiKey string) (BiblesData, error) {
	var url string
	if languageID == "" {
		url = "https://api.scripture.api.bible/v1/bibles"
	} else {
		url = fmt.Sprintf("https://api.scripture.api.bible/v1/bibles?language=%s", languageID)
	}

	// Set up the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return BiblesData{}, err
	}
	req.Header.Set("api-key", apiKey)

	// Perform the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return BiblesData{}, err
	}

	// Read the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return BiblesData{}, err
	}
	defer resp.Body.Close()

	// Unmarshal the response body
	var biblesData BiblesData
	err = json.Unmarshal(body, &biblesData)
	if err != nil {
		return BiblesData{}, err
	}

	return biblesData, nil
}
