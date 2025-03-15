package config

import (
	"github.com/Jarimus/BibleTUI/internal/api_query"
)

type currentlyReading struct {
	TranslationName string
	TranslationID   string
	TranslationData api_query.TranslationData
	BookData        api_query.BookData
	ChapterData     api_query.ChapterData
}

type Config struct {
	CurrentlyReading currentlyReading
}
