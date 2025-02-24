package api_query

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

func NewChapterQuery(url string) tea.Cmd {
	return func() tea.Msg {

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
		var query BookData
		err = json.Unmarshal(body, &query)
		if err != nil {
			fmt.Println(err)
		}

		return query
	}
}
func BibleChapterQuery() tea.Msg {

	url := "https://bible-api.com/john%201"

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
	var query BookData
	err = json.Unmarshal(body, &query)
	if err != nil {
		fmt.Println(err)
	}

	return query
}
