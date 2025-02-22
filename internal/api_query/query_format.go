package api_query

type RandomQuery struct {
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


type BookQuery struct {
	Reference string `json:"reference"`
	Verses    []struct {
		BookID   string `json:"book_id"`
		BookName string `json:"book_name"`
		Chapter  int    `json:"chapter"`
		Verse    int    `json:"verse"`
		Text     string `json:"text"`
	} `json:"verses"`
	Text            string `json:"text"`
	TranslationID   string `json:"translation_id"`
	TranslationName string `json:"translation_name"`
	TranslationNote string `json:"translation_note"`
}