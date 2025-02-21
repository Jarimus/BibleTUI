package api_query

type QueryFormat struct {
	Translation struct {
		Identifier   string `json:"identifier"`
		Name         string `json:"name"`
		Language     string `json:"language"`
		LanguageCode string `json:"language_code"`
		License      string `json:"license"`
	} `json:"translation"`
	RandomVerse struct {
		BookID  string `json:"book_id"`
		Book    string `json:"book"`
		Chapter int    `json:"chapter"`
		Verse   int    `json:"verse"`
		Text    string `json:"text"`
	} `json:"random_verse"`
}
