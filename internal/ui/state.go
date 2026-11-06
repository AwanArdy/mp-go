package ui

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"sync"

	"music-player/internal/audio"
	"music-player/internal/library"
	"music-player/internal/model"

	"github.com/dhowden/tag"
	"github.com/egoist/mygo"
)

type CardItem struct {
	Name      string
	Subtext   string
	CoverData []byte
}

type AppState struct {
	mu           sync.Mutex
	Tracks       []model.Track
	CurrentIndex int
	Engine       *audio.Engine
	Window       *mygo.Window

	DarkMode       bool
	ViewMode       string
	SelectedArtist string
	SelectedAlbum  string
	SearchQuery    string

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
		DarkMode:       false,
		ViewMode:       "tracks",
		SelectedArtist: "",
		SelectedAlbum:  "",
		SearchQuery:    "",
		Volume:         0.8,
	}

	app.Engine.SetHandlers(
		func() {
			if app.Repeat {
				if curr := app.CurrentTrack(); curr != nil {
					app.PlayTrack(*curr)
					return
				}
			}
			app.NextTrack()
		},
		func() { app.NextTrack() },
		func() { app.PrevTrack() },
	)

	return app
}

func (s *AppState) SetWindow(w *mygo.Window) {
	s.Window = w
}

func (s *AppState) ToggleDarkMode() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.DarkMode = !s.DarkMode
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

func (s *AppState) GetArtistCards() []CardItem {
	s.mu.Lock()
	defer s.mu.Unlock()

	type meta struct {
		count int
		cover []byte
	}
	group := make(map[string]*meta)
	var names []string

	for _, t := range s.Tracks {
		name := strings.TrimSpace(t.Artist)
		if name == "" {
			name = "Unknown Artist"
		}
		if _, exists := group[name]; !exists {
			group[name] = &meta{count: 0, cover: nil}
			names = append(names, name)
		}
		group[name].count++
		if group[name].cover == nil && len(t.CoverData) > 0 {
			group[name].cover = t.CoverData
		}
	}
	sort.Strings(names)

	var cards []CardItem
	for _, n := range names {
		cards = append(cards, CardItem{
			Name:      n,
			Subtext:   fmt.Sprintf("%d songs", group[n].count),
			CoverData: group[n].cover,
		})
	}
	return cards
}

func (s *AppState) GetAlbumCards() []CardItem {
	s.mu.Lock()
	defer s.mu.Unlock()

	type meta struct {
		artist string
		count  int
		cover  []byte
	}
	group := make(map[string]*meta)
	var names []string

	for _, t := range s.Tracks {
		name := strings.TrimSpace(t.Album)
		if name == "" {
			name = "Unknown Album"
		}
		if _, exists := group[name]; !exists {
			group[name] = &meta{artist: t.Artist, count: 0, cover: nil}
			names = append(names, name)
		}
		group[name].count++
		if group[name].cover == nil && len(t.CoverData) > 0 {
			group[name].cover = t.CoverData
		}
	}
	sort.Strings(names)

	var cards []CardItem
	for _, n := range names {
		cards = append(cards, CardItem{
			Name:      n,
			Subtext:   group[n].artist,
			CoverData: group[n].cover,
		})
	}
	return cards
}

func (s *AppState) ClearFilter() {
	s.SelectedArtist = ""
	s.SelectedAlbum = ""
	s.SearchQuery = ""
	s.ViewMode = "tracks"
}

func (s *AppState) PlayTrack(track model.Track) {
	s.mu.Lock()
	for i, t := range s.Tracks {
		if t.Path == track.Path {
			s.CurrentIndex = i
			if len(s.Tracks[i].CoverData) == 0 {
				if f, err := os.Open(t.Path); err == nil {
					if m, err := tag.ReadFrom(f); err == nil {
						if pic := m.Picture(); pic != nil {
							s.Tracks[i].CoverData = pic.Data
						}
					}
					f.Close()
				}
			}
			track = s.Tracks[i]
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
	tracks := s.FilteredTracks()
	if len(s.Tracks) == 0 {
		tracks = s.Tracks
	}
	if len(tracks) == 0 {
		return
	}

	curr := s.CurrentTrack()
	currIdx := -1
	if curr != nil {
		for i, t := range tracks {
			if t.Path == curr.Path {
				currIdx = i
				break
			}
		}
	}

	if s.Shuffle && len(tracks) > 1 {
		nextIdx := currIdx
		for nextIdx == currIdx {
			nextIdx = rand.Intn(len(tracks))
		}
		s.PlayTrack(tracks[nextIdx])
		return
	}

	nextIdx := (currIdx + 1) % len(tracks)
	s.PlayTrack(tracks[nextIdx])
}

func (s *AppState) PrevTrack() {
	tracks := s.FilteredTracks()
	if len(s.Tracks) == 0 {
		tracks = s.Tracks
	}
	if len(tracks) == 0 {
		return
	}

	curr := s.CurrentTrack()
	currIdx := -1
	if curr != nil {
		for i, t := range tracks {
			if t.Path == curr.Path {
				currIdx = i
				break
			}
		}
	}

	if s.Shuffle && len(tracks) > 1 {
		prevIdx := currIdx
		for prevIdx == currIdx {
			prevIdx = rand.Intn(len(tracks))
		}
		s.PlayTrack(tracks[prevIdx])
		return
	}

	prevIdx := (currIdx - 1 + len(tracks)) % len(tracks)
	s.PlayTrack(tracks[prevIdx])
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
