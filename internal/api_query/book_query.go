package api_query

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func BookQuery(translationID, bookID string) BookData {
	// Key data: names and IDs for chapters of the chosen book

	url := fmt.Sprintf("https://api.scripture.api.bible/v1/bibles/%s/books/%s/chapters", translationID, bookID)

	// Set up the request with a corrent header (API-key needed for API-Bible)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request for book data: %v", err)
		return BookData{}
	}

	req.Header.Set("api-key", getApi())

	// Perform the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error requesting data for book: %v", err)
		return BookData{}
	}

	// Read the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading the body for book query: %v", err)
		return BookData{}
	}
	defer resp.Body.Close()

	// Unmarshal the response
	var bookData BookData
	if err := json.Unmarshal(body, &bookData); err != nil {
		fmt.Printf("Error unmarshaling body for book query: %v", err)
	}

	return bookData

}
