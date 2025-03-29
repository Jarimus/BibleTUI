package config

import (
	"github.com/Jarimus/BibleTUI/internal/api_query"
)

type currentlyReading struct {
	TranslationName string `json:"translation_name"`
	TranslationID   string `json:"translation_id"`
	TranslationData api_query.TranslationData `json:"translation_data"`
	BookData        api_query.BookData `json:"book_data"`
	ChapterData     api_query.ChapterData `json:"chapter_data"`
}

type Config struct {
	CurrentlyReading currentlyReading `json:"currently_reading"`
}
