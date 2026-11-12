package ui

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

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
	ViewMode       string // "tracks", "artists", "albums", "lyrics"
	SelectedArtist string
	SelectedAlbum  string
	SearchQuery    string

	SelectedTrackPath string
	lastClickTime     time.Time
	lastClickPath     string

	ArtistCards []CardItem
	AlbumCards  []CardItem

	// State Lirik Lagu
	LyricsLines    []library.LyricLine
	LyricsStatus   string // "idle", "loading", "ready", "not_found"
	LyricVisualPos float64
	lyricsTrack    string

	Volume   float64
	Progress float64
	IsLiked  bool
	Shuffle  bool
	Repeat   bool

	Playlists          []model.Playlist
	SelectedPlaylistID string

	IsAddingPlaylist    bool
	EditingPlaylistID   string
	InlinePlaylistInput string
}

func NewAppState() *AppState {
	tracks, _ := library.ScanMusicDir()
	app := &AppState{
		Tracks:         tracks,
		CurrentIndex:   -1,
		Engine:         audio.NewEngine(),
		DarkMode:       true,
		ViewMode:       "tracks",
		SelectedArtist: "",
		SelectedAlbum:  "",
		SearchQuery:    "",
		LyricsStatus:   "idle",
		Volume:         0.8,
	}

	app.BuildLibraryCaches()

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

func (s *AppState) StartAddPlaylist() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.IsAddingPlaylist = true
	s.EditingPlaylistID = ""
	s.InlinePlaylistInput = fmt.Sprintf("Playlist #%d", len(s.Playlists)+1)
}

func (s *AppState) CommitAddPlaylist() {
	s.mu.Lock()
	defer s.mu.Unlock()

	trimmed := strings.TrimSpace(s.InlinePlaylistInput)
	if trimmed == "" {
		trimmed = fmt.Sprintf("Playlist #%d", len(s.Playlists)+1)
	}

	id := fmt.Sprintf("pl_%d", time.Now().UnixNano())
	pl := model.Playlist{
		ID:         id,
		Name:       trimmed,
		TrackPaths: make([]string, 0),
	}
	s.Playlists = append(s.Playlists, pl)
	s.SelectedPlaylistID = id
	s.SelectedArtist = ""
	s.SelectedAlbum = ""
	s.ViewMode = "playlist"

	s.IsAddingPlaylist = false
	s.InlinePlaylistInput = ""
}

func (s *AppState) CancelInlinePlaylist() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.IsAddingPlaylist = false
	s.EditingPlaylistID = ""
	s.InlinePlaylistInput = ""
}

func (s *AppState) StartRenamePlaylist(id, currentName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.EditingPlaylistID = id
	s.IsAddingPlaylist = false
	s.InlinePlaylistInput = currentName
}

func (s *AppState) CommitRenamePlaylist() {
	s.mu.Lock()
	defer s.mu.Unlock()

	trimmed := strings.TrimSpace(s.InlinePlaylistInput)
	if trimmed != "" && s.EditingPlaylistID != "" {
		for i := range s.Playlists {
			if s.Playlists[i].ID == s.EditingPlaylistID {
				s.Playlists[i].Name = trimmed
				break
			}
		}
	}
	s.EditingPlaylistID = ""
	s.InlinePlaylistInput = ""
}

func (s *AppState) DeletePlaylist(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := -1
	for i, pl := range s.Playlists {
		if pl.ID == id {
			idx = i
			break
		}
	}

	if idx != -1 {
		s.Playlists = append(s.Playlists[:idx], s.Playlists[idx+1:]...)
	}

	if s.SelectedPlaylistID == id {
		s.SelectedPlaylistID = ""
		s.ViewMode = "tracks"
	}
	if s.EditingPlaylistID == id {
		s.EditingPlaylistID = ""
	}
}

func (s *AppState) GetSelectedPlaylist() *model.Playlist {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.Playlists {
		if s.Playlists[i].ID == s.SelectedPlaylistID {
			return &s.Playlists[i]
		}
	}
	return nil
}

func (s *AppState) PlaylistTracks(pl *model.Playlist) []model.Track {
	s.mu.Lock()
	defer s.mu.Unlock()

	if pl == nil {
		return nil
	}

	trackMap := make(map[string]model.Track, len(s.Tracks))
	for _, t := range s.Tracks {
		trackMap[t.Path] = t
	}

	var result []model.Track
	for _, p := range pl.TrackPaths {
		if t, ok := trackMap[p]; ok {
			result = append(result, t)
		}
	}
	return result
}

func (s *AppState) GetSelectedTrack() *model.Track {
	s.mu.Lock()
	defer s.mu.Unlock()

	targetPath := s.SelectedTrackPath
	if targetPath == "" {
		if s.CurrentIndex >= 0 && s.CurrentIndex < len(s.Tracks) {
			return &s.Tracks[s.CurrentIndex]
		}
		return nil
	}

	for i := range s.Tracks {
		if s.Tracks[i].Path == targetPath {
			return &s.Tracks[i]
		}
	}
	return nil
}

func (s *AppState) BuildLibraryCaches() {
	s.mu.Lock()
	defer s.mu.Unlock()

	type metaArt struct {
		count int
		cover []byte
	}
	artGroup := make(map[string]*metaArt)
	var artNames []string

	for _, t := range s.Tracks {
		name := strings.TrimSpace(t.Artist)
		if name == "" {
			name = "Unknown Artist"
		}
		if _, exists := artGroup[name]; !exists {
			artGroup[name] = &metaArt{count: 0, cover: nil}
			artNames = append(artNames, name)
		}
		artGroup[name].count++
		if artGroup[name].cover == nil && len(t.CoverData) > 0 {
			artGroup[name].cover = t.CoverData
		}
	}
	sort.Strings(artNames)

	s.ArtistCards = make([]CardItem, 0, len(artNames))
	for _, n := range artNames {
		s.ArtistCards = append(s.ArtistCards, CardItem{
			Name:      n,
			Subtext:   fmt.Sprintf("%d songs", artGroup[n].count),
			CoverData: artGroup[n].cover,
		})
	}

	type metaAlb struct {
		artist string
		count  int
		cover  []byte
	}
	albGroup := make(map[string]*metaAlb)
	var albNames []string

	for _, t := range s.Tracks {
		name := strings.TrimSpace(t.Album)
		if name == "" {
			name = "Unknown Album"
		}
		if _, exists := albGroup[name]; !exists {
			albGroup[name] = &metaAlb{artist: t.Artist, count: 0, cover: nil}
			albNames = append(albNames, name)
		}
		albGroup[name].count++
		if albGroup[name].cover == nil && len(t.CoverData) > 0 {
			albGroup[name].cover = t.CoverData
		}
	}
	sort.Strings(albNames)

	s.AlbumCards = make([]CardItem, 0, len(albNames))
	for _, n := range albNames {
		s.AlbumCards = append(s.AlbumCards, CardItem{
			Name:      n,
			Subtext:   albGroup[n].artist,
			CoverData: albGroup[n].cover,
		})
	}
}

func (s *AppState) SetWindow(w *mygo.Window) {
	s.Window = w
}

func (s *AppState) ToggleDarkMode() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.DarkMode = true
}

func (s *AppState) LoadLyrics(track model.Track) {
	s.mu.Lock()
	if s.lyricsTrack == track.Path && (s.LyricsStatus == "ready" || s.LyricsStatus == "loading") {
		s.mu.Unlock()
		return
	}
	s.lyricsTrack = track.Path
	s.LyricsStatus = "loading"
	s.LyricsLines = nil
	s.LyricVisualPos = 0
	s.mu.Unlock()

	go func(t model.Track, dur time.Duration) {
		lines, err := library.FetchSyncedLyrics(t.Title, t.Artist, t.Album, dur)
		s.mu.Lock()
		defer s.mu.Unlock()
		if s.lyricsTrack != t.Path {
			return
		}
		if err != nil || len(lines) == 0 {
			s.LyricsStatus = "not_found"
			s.LyricsLines = nil
		} else {
			s.LyricsStatus = "ready"
			s.LyricsLines = lines
			s.LyricVisualPos = 0
		}
	}(track, s.Engine.Duration())
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
	return s.ArtistCards
}

func (s *AppState) GetAlbumCards() []CardItem {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.AlbumCards
}

func (s *AppState) ClearFilter() {
	s.SelectedArtist = ""
	s.SelectedAlbum = ""
	s.SelectedPlaylistID = ""
	s.SearchQuery = ""
	s.ViewMode = "tracks"
}

func (s *AppState) PlayTrack(track model.Track) {
	s.mu.Lock()
	s.SelectedTrackPath = track.Path
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

	// Jika mode lirik sedang terbuka, langsung muat lirik untuk lagu baru
	if s.ViewMode == "lyrics" {
		s.LoadLyrics(track)
	} else {
		s.lyricsTrack = ""
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
	if len(tracks) == 0 {
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
	if len(tracks) == 0 {
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
