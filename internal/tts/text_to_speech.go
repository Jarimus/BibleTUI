package tts

import (
	"os"
	"regexp"

	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
)

func SpeakText(text, lan string) error {
	lan = ISOtoTTScode(lan)
	speech := htgotts.Speech{Folder: "audio", Language: lan, Handler: &handlers.Native{}}

	reObject, err := regexp.Compile(`\[.*?\]`)
	if err != nil {
		return err
	}
	verses := reObject.Split(text, -1)
	for _, verse := range verses {
		if err := speech.Speak(verse); err != nil {
			return err
		}
		if err := os.RemoveAll("audio"); err != nil {
			return err
		}
	}

	return nil
}
