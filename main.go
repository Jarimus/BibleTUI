package main

import (
	"database/sql"
	"log"

	"github.com/Jarimus/BibleTUI/internal/api_query"
	"github.com/Jarimus/BibleTUI/internal/database"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
)

// Global variables:
// Terminal width and height. Necessary for the viewport in reading mode when it spawns.
var window_width int
var window_height int

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
	CurrentlyReading currentlyReading `json:"currently_reading"`
	CurrentUser      string
	dbQueries        *database.Queries
}

var apiCfg config

func main() {

	println("Loading...")

	var err error
	// Get settings
	apiCfg, err = loadSettings()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to database
	db, err := sql.Open("sqlite3", "bibletui.db")
	if err != nil {
		log.Fatalf("error opening connection to database: %s", err)
	}
	dbQueries := database.New(db)
	apiCfg.dbQueries = dbQueries

	// Initialize with the current translation
	apiCfg.CurrentlyReading.TranslationData = api_query.TranslationQuery(apiCfg.CurrentlyReading.TranslationID)

	mainMenu := newMainMenu()

	// Root screen holds the other models in a "stack" (slice) and displays the one at the top.
	root := newRootScreen([]tea.Model{mainMenu})
	p := tea.NewProgram(root)
	_, err = p.Run()
	if err != nil {
		log.Fatalf("error starting the program: %s", err)
	}
}
