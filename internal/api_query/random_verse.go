package api_query

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func GetRandomVerse() RandomQuery {

	url := "https://bible-api.com/data/web/random"

	// Request a new verse
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}

	// Read the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	// Unmarshal
	var query RandomQuery
	err = json.Unmarshal(body, &query)
	if err != nil {
		log.Println(err)
	}

	return query
}
