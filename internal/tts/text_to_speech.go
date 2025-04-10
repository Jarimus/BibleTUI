package tts

import (
	"os"

	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
)

const AudioFolderPath = "BibleTUIaudio"

// Uses htgo-tts to play audio. Input language should be ISO 63
func SpeakText(text, lan string, audioStop chan bool) error {

	// Delete any previous audio files.
	if err := os.RemoveAll(AudioFolderPath); err != nil {
		return err
	}

	// Convert the
	lan = ISOtoTTScode(lan)
	speech := htgotts.Speech{Folder: AudioFolderPath, Language: lan, Handler: &handlers.Native{}}

	// Splits the text (a chapter) into verses
	verses, err := splitVerses(text)
	if err != nil {
		return err
	}

	// Listen to stop signal from audioStop
	var done bool
	listenForDone := func(bool, chan bool) {
		done = <-audioStop
	}
	go listenForDone(done, audioStop)

	// Play each verse after checking whether it is short enough for the tts to handle.
	for _, verse := range verses {
		parts := splitText(verse)
		for _, part := range parts {
			if done {
				break
			}
			if err := speech.Speak(part); err != nil {
				return err
			}

		}
	}

	// Delete the temp audio files.
	if err := os.RemoveAll(AudioFolderPath); err != nil {
		return err
	}

	return nil
}
