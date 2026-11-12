package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	appUI "music-player/internal/ui"
)

func main() {
	state := appUI.NewAppState()

	mygo.App.WhenReady(func() {
		win := mygo.NewWindow(mygo.WindowOptions{
			Title:     "Music",
			Width:     1080,
			Height:    720,
			Frameless: true,
			Content:   ui.View(state.View),
		})

		state.SetWindow(win)
	})

	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
