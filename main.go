package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/hotkey"
)

const (
	dictBaseURL string = "https://api.dictionaryapi.dev/api/v2/entries/en/"
)

type dictRes struct {
	Phonetic string `json:"phonetic"`
	Origin   string `json:"origin"`
	Meanings []struct {
		PartOfSpeech string `json:"partOfSpeech"`
		Definitions  []struct {
			Definition string `json:"definition"`
			Example    string `json:"example"`
		} `json:"definitions"`
	} `json:"meanings"`
}

func main() {
	a := app.New()
	log.Println("cliptionary: app started")
	dummy := a.NewWindow("cliptionary dummy")
	dummy.Hide()

	clip := a.Clipboard()

	go func() {
		hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModAlt}, hotkey.KeyD)
		err := hk.Register()
		if err != nil {
			panic(fmt.Sprintf("cliptionary: hotkey registraion failed: %v\n", err))
		}
		log.Printf("cliptionary: hotkey %v registered\n", hk)
		for range hk.Keydown() {
			copy()
			time.Sleep(200 * time.Millisecond)
			word := clip.Content()
			log.Println("cliptonary: new window")
			fyne.Do(func() {
				res, err := http.Get(dictBaseURL + word)
				if err != nil {
					log.Printf("cliptionary: err fetching dictionary entry for %s: %v\n", word, err)
				} else if res.StatusCode != http.StatusOK {
					log.Println("cliptionary: entr res", res.StatusCode)
				} else {
					defer res.Body.Close()
					w := a.NewWindow(word)
					d := json.NewDecoder(res.Body)
					var mean []dictRes
					err = d.Decode(&mean)
					if err != nil {
						log.Printf("cliptionary: error decodin: %v", err)
					} else {
						w.SetContent(widget.NewLabel(fmt.Sprintf("Phonetic: %s\nOrigin:%s\nMeaning: %v", mean[0].Phonetic, mean[0].Phonetic, mean[0].Meanings[0].Definitions[0].Definition)))
						w.Show()
					}
				}
			})
		}
	}()

	a.Run()
}

func copy() {
	// TODO
	// Copy text by simulating Ctrl+C
}
