package api_query

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetRandomVerse() string {

	url := "https://bible-api.com/data/web/random"

	// Request a new verse
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	// Read the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	// Unmarshal
	var query QueryFormat
	err = json.Unmarshal(body, &query)
	if err != nil {
		fmt.Println(err)
	}

	return query.RandomVerse.Text
}
