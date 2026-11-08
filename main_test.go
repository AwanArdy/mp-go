package main

import (
	"testing"

	"github.com/egoist/mygo/ui"
	appUI "music-player/internal/ui"
)

func TestView(t *testing.T) {
	state := appUI.NewAppState()
	tt := ui.NewTester(state.View, 1080, 720)
	if !tt.HasText("Music") {
		t.Errorf("expected Music in header, got texts %q", tt.Texts())
	}
}
