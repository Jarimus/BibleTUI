package api_query

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func TranslationQuery(translationID string) TranslationData {

	url := fmt.Sprintf("https://api.scripture.api.bible/v1/bibles/%s/books", translationID)

	// Set up the request with a corrent header (API-key needed for API-Bible)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request for translation data: %v", err)
		return TranslationData{}
	}

	req.Header.Set("api-key", getApi())

	// Perform the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error requesting data for translation: %v", err)
		return TranslationData{}
	}

	// Read the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading the body for translation query: %v", err)
		return TranslationData{}
	}
	defer resp.Body.Close()

	// Unmarshal the response
	var translationData TranslationData
	if err := json.Unmarshal(body, &translationData); err != nil {
		fmt.Printf("Error unmarshaling body for translation query: %v", err)
	}

	return translationData

}
