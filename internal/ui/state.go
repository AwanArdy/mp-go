package ui

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"

	"github.com/egoist/mygo"
	"music-player/internal/audio"
	"music-player/internal/library"
	"music-player/internal/model"
)

type AppState struct {
	mu           sync.Mutex
	Tracks       []model.Track
	CurrentIndex int
	Engine       *audio.Engine
	Window       *mygo.Window

	ViewMode       string
	SelectedArtist string
	SelectedAlbum  string
	SearchQuery    string

	CurrentPage int
	PageSize    int

	Volume   float64
	Progress float64
	IsLiked  bool
	Shuffle  bool
	Repeat   bool
}

func NewAppState() *AppState {
	tracks, _ := library.ScanMusicDir()
	app := &AppState{
		Tracks:         tracks,
		CurrentIndex:   -1,
		Engine:         audio.NewEngine(),
		ViewMode:       "tracks",
		SelectedArtist: "",
		SelectedAlbum:  "",
		SearchQuery:    "",
		CurrentPage:    0,
		PageSize:       9,
		Volume:         0.8,
	}

	app.Engine.SetOnEnd(func() {
		app.NextTrack()
	})

	return app
}

func (s *AppState) SetWindow(w *mygo.Window) {
	s.Window = w
}

func (s *AppState) FilteredTracks() []model.Track {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []model.Track
	query := strings.ToLower(strings.TrimSpace(s.SearchQuery))

	for _, t := range s.Tracks {
		if s.SelectedArtist != "" && !strings.EqualFold(t.Artist, s.SelectedArtist) {
			continue
		}
		if s.SelectedAlbum != "" && !strings.EqualFold(t.Album, s.SelectedAlbum) {
			continue
		}
		if query != "" {
			matchTitle := strings.Contains(strings.ToLower(t.Title), query)
			matchArtist := strings.Contains(strings.ToLower(t.Artist), query)
			matchAlbum := strings.Contains(strings.ToLower(t.Album), query)
			if !matchTitle && !matchArtist && !matchAlbum {
				continue
			}
		}
		result = append(result, t)
	}
	return result
}

func (s *AppState) GetArtists() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	seen := make(map[string]bool)
	var list []string
	for _, t := range s.Tracks {
		name := strings.TrimSpace(t.Artist)
		if name == "" {
			name = "Unknown Artist"
		}
		if !seen[name] {
			seen[name] = true
			list = append(list, name)
		}
	}
	sort.Strings(list)
	return list
}

func (s *AppState) GetAlbums() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	seen := make(map[string]bool)
	var list []string
	for _, t := range s.Tracks {
		name := strings.TrimSpace(t.Album)
		if name == "" {
			name = "Unknown Album"
		}
		if !seen[name] {
			seen[name] = true
			list = append(list, name)
		}
	}
	sort.Strings(list)
	return list
}

func (s *AppState) ClearFilter() {
	s.SelectedArtist = ""
	s.SelectedAlbum = ""
	s.SearchQuery = ""
	s.CurrentPage = 0
	s.ViewMode = "tracks"
}

func (s *AppState) PlayTrack(track model.Track) {
	s.mu.Lock()
	for i, t := range s.Tracks {
		if t.Path == track.Path {
			s.CurrentIndex = i
			break
		}
	}
	s.mu.Unlock()

	if err := s.Engine.Play(track); err != nil {
		log.Printf("[AUDIO ERROR] Gagal memutar %s: %v", track.Title, err)
	}
}

func (s *AppState) TogglePlay() {
	if s.CurrentIndex == -1 && len(s.Tracks) > 0 {
		s.PlayTrack(s.Tracks[0])
		return
	}
	s.Engine.TogglePlay()
}

func (s *AppState) NextTrack() {
	if len(s.Tracks) == 0 {
		return
	}
	next := (s.CurrentIndex + 1) % len(s.Tracks)
	s.PlayTrack(s.Tracks[next])
}

func (s *AppState) PrevTrack() {
	if len(s.Tracks) == 0 {
		return
	}
	prev := (s.CurrentIndex - 1 + len(s.Tracks)) % len(s.Tracks)
	s.PlayTrack(s.Tracks[prev])
}

func (s *AppState) Seek(prog float64) {
	s.Engine.Seek(prog)
}

func (s *AppState) CurrentTrack() *model.Track {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.CurrentIndex >= 0 && s.CurrentIndex < len(s.Tracks) {
		return &s.Tracks[s.CurrentIndex]
	}
	return nil
}

func (s *AppState) CloseApp() {
	if s.Window != nil {
		s.Window.Close()
	} else {
		mygo.App.Quit()
	}
}

func (s *AppState) MinimizeApp() {
	if s.Window != nil {
		s.Window.Minimize()
	}
}

func (s *AppState) CurrentTimeString() string {
	if s.CurrentTrack() == nil {
		return "00:00"
	}
	pos := s.Engine.Position()
	sec := int(pos.Seconds())
	return fmt.Sprintf("%02d:%02d", sec/60, sec%60)
}

func (s *AppState) TotalTimeString() string {
	if s.CurrentTrack() == nil {
		return "00:00"
	}
	dur := s.Engine.Duration()
	if dur <= 0 {
		return "00:00"
	}
	sec := int(dur.Seconds())
	return fmt.Sprintf("%02d:%02d", sec/60, sec%60)
}
