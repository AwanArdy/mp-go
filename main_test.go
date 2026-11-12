package main

import (
	"testing"

	"github.com/egoist/mygo/ui"
	appUI "music-player/internal/ui"
)

func TestView(t *testing.T) {
	state := appUI.NewAppState()
	tt := ui.NewTester(state.View, 1080, 720)

	// Verify Header Branding & M3 badge
	if !tt.HasText("Music") {
		t.Errorf("expected Music in header, got texts %q", tt.Texts())
	}
	if !tt.HasText("M3") {
		t.Errorf("expected M3 badge in header, got texts %q", tt.Texts())
	}

	// Verify Sidebar Navigation Drawer items
	if !tt.HasText("LIBRARY") {
		t.Errorf("expected LIBRARY section, got texts %q", tt.Texts())
	}
	if !tt.HasText("All Tracks") {
		t.Errorf("expected All Tracks nav item, got texts %q", tt.Texts())
	}
	if !tt.HasText("Artists") {
		t.Errorf("expected Artists nav item, got texts %q", tt.Texts())
	}
	if !tt.HasText("Albums") {
		t.Errorf("expected Albums nav item, got texts %q", tt.Texts())
	}
	if !tt.HasText("New Playlist") {
		t.Errorf("expected New Playlist FAB button, got texts %q", tt.Texts())
	}

	// Verify Dark Mode is enabled
	if !state.DarkMode {
		t.Errorf("expected DarkMode to be true")
	}

	// Test switching to Artists mode
	state.ViewMode = "artists"
	tt = ui.NewTester(state.View, 1080, 720)
	if !tt.HasText("Artists") {
		t.Errorf("expected Artists view, got texts %q", tt.Texts())
	}

	// Test switching to Albums mode
	state.ViewMode = "albums"
	tt = ui.NewTester(state.View, 1080, 720)
	if !tt.HasText("Albums") {
		t.Errorf("expected Albums view, got texts %q", tt.Texts())
	}

	// Test switching to Tracks mode
	state.ViewMode = "tracks"
	tt = ui.NewTester(state.View, 1080, 720)
	if !tt.HasText("Tracks") {
		t.Errorf("expected Tracks view, got texts %q", tt.Texts())
	}
}
