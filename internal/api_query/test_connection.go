package api_query

import (
	"net/http"

	"errors"
)

func TestConnection(apiKey string) error {
	url := "https://api.scripture.api.bible/v1/bibles"

	// Set up the request with a corrent header (API-key needed for API-Bible)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("api-key", apiKey)

	// Perform the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == 200 {
		return nil
	} else {
		return errors.New("invalid API Key\nGet a personal API Key by registering at https://scripture.api.bible/")
	}
}
