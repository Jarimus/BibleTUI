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

// NEW API DOWN HERE

// Data for the books of a translation. Use the Bible-id for the search:
// https://api.scripture.api.bible/v1/bibles/[BIBLE-ID]/books
type TranslationData struct {
	Books []struct {
		ID           string `json:"id"`
		BibleID      string `json:"bibleId"`
		Abbreviation string `json:"abbreviation"`
		Name         string `json:"name"`
		NameLong     string `json:"nameLong"`
	} `json:"data"`
}

// Data for the chapters of a book. Need Bible-id and book-id (TranslationData.Books[i].ID):
// https://api.scripture.api.bible/v1/bibles/[BIBLE-ID]/books/[BOOK-ID]/chapters
type BookData struct {
	Chapters []struct {
		ID        string `json:"id"`
		BibleID   string `json:"bibleId"`
		BookID    string `json:"bookId"`
		Number    string `json:"number"`
		Reference string `json:"reference"`
	} `json:"data"`
}

// Chapter structure. Query:
// https://api.scripture.api.bible/v1/bibles/[BIBLE-ID]/chapters/[CHAPTER-ID]?content-type=text
// &include-notes=false&include-titles=true&include-chapter-numbers=false&include-verse-numbers=true&include-verse-spans=false
type ChapterData struct {
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
