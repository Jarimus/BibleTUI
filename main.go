package main

import (
	_ "embed"
	"log"

	"github.com/Jarimus/BibleTUI/internal/api_query"
	"github.com/Jarimus/BibleTUI/internal/database"
	tea "github.com/charmbracelet/bubbletea"
	_ "modernc.org/sqlite"
)

// Global variables:
// Terminal width and height. Necessary for the viewport in reading mode when it spawns.
var window_width int
var window_height int

// Embed the schema files
//
//go:embed sql/schema/001_users.sql
var usersSchema string

//go:embed sql/schema/002_translations.sql
var translationsSchema string

// Database filepath
var dbFilePath string

// Struct for the data about the current translation being read.
type currentlyReading struct {
	TranslationName string                    `json:"translation_name"`
	TranslationID   string                    `json:"translation_id"`
	TranslationData api_query.TranslationData `json:"translation_data"`
	BookData        api_query.BookData        `json:"book_data"`
	ChapterData     api_query.ChapterData     `json:"chapter_data"`
}

// A config struct for current translation, user, database queries
type config struct {
	CurrentUser      string `json:"current_user"`
	CurrentUserID    int64  `json:"current_user_id"`
	ApiKey           string `json:"api_key"`
	dbQueries        *database.Queries
	CurrentlyReading currentlyReading `json:"currently_reading"`
}

// Api config struct to store the config file's data in memory.
var apiCfg config

func main() {

	println("Loading...")

	var err error

	// Load settings
	err = loadSettings()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to database
	dbFilePath = "BibleTUI.db"

	err = initializeDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize with the current translation
	apiCfg.CurrentlyReading.TranslationData, err = api_query.TranslationQuery(apiCfg.CurrentlyReading.TranslationID, apiCfg.ApiKey)
	if err != nil {

	}

	// Create a new main menu tea.Model
	mainMenu := newMainMenu()

	// Root screen holds the other models in a "stack" (slice) and displays the one at the top.
	root := newRootScreen([]tea.Model{mainMenu})
	p := tea.NewProgram(root)
	_, err = p.Run()
	if err != nil {
		log.Fatalf("error starting the program: %v", err)
	}
}
