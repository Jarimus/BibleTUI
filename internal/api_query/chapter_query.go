package api_query

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func ChapterQuery(translationID, chapterID string) ChapterData {
	// Key data in this query: content of the current chapter, IDs for the previous and next chapter

	url := fmt.Sprintf("https://api.scripture.api.bible/v1/bibles/%s/chapters/%s?content-type=text&include-notes=false&include-titles=true&include-chapter-numbers=false&include-verse-numbers=true&include-verse-spans=false",
		translationID,
		chapterID)

	// Set up the request with a corrent header (API-key needed for API-Bible)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request for book data: %v", err)
		return ChapterData{}
	}

	// If you do not have an api-key, you have to make an account at https://scripture.api.bible/ to get one.
	req.Header.Set("api-key", os.Getenv("API_KEY"))

	// Perform the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error requesting data for book: %v", err)
		return ChapterData{}
	}

	// Read the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading the body for book query: %v", err)
		return ChapterData{}
	}
	defer resp.Body.Close()

	// Unmarshal the response
	var chapterData ChapterData
	if err := json.Unmarshal(body, &chapterData); err != nil {
		fmt.Printf("Error unmarshaling body for book query: %v", err)
	}

	return chapterData

}
