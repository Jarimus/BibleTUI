package api_query

// Old api
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

// Old api
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

// NEW API DOWN HERE

// query for the books of a translation. Use the Bible-id for the search:
// https://api.scripture.api.bible/v1/bibles/[BIBLE-ID]/books
type BooksQuery struct {
	Data []struct {
		ID           string `json:"id"`
		BibleID      string `json:"bibleId"`
		Abbreviation string `json:"abbreviation"`
		Name         string `json:"name"`
		NameLong     string `json:"nameLong"`
	} `json:"data"`
}

// Chapter structure. Query:
// https://api.scripture.api.bible/v1/bibles/[BIBLE-ID]/chapters/[CHAPTER-ID]?content-type=text
// &include-notes=false&include-titles=true&include-chapter-numbers=false&include-verse-numbers=true&include-verse-spans=false
type ChapterQuery struct {
	Data struct {
		ID         string `json:"id"`
		BibleID    string `json:"bibleId"`
		Number     string `json:"number"`
		BookID     string `json:"bookId"`
		Reference  string `json:"reference"`
		Copyright  string `json:"copyright"`
		VerseCount int    `json:"verseCount"`
		Content    string `json:"content"`
		Next       struct {
			ID     string `json:"id"`
			Number string `json:"number"`
			BookID string `json:"bookId"`
		} `json:"next"`
		Previous struct {
			ID     string `json:"id"`
			Number string `json:"number"`
			BookID string `json:"bookId"`
		} `json:"previous"`
	} `json:"data"`
	Meta struct {
		Fums          string `json:"fums"`
		FumsID        string `json:"fumsId"`
		FumsJsInclude string `json:"fumsJsInclude"`
		FumsJs        string `json:"fumsJs"`
		FumsNoScript  string `json:"fumsNoScript"`
	} `json:"meta"`
}
