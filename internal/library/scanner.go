package library

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
	"music-player/internal/model"
)

func ScanMusicDir() ([]model.Track, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	musicDir := filepath.Join(home, "Music")

	var tracks []model.Track
	err = filepath.Walk(musicDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".mp3", ".flac", ".m4a", ".aac", ".wav", ".ogg":
			// ekstensi didukung
		default:
			return nil
		}

		fileName := info.Name()
		cleanTitle := strings.TrimSuffix(fileName, filepath.Ext(fileName))

		t := model.Track{
			ID:     path,
			Title:  cleanTitle,
			Artist: "unknown artist",
			Album:  "unknown album",
			Path:   path,
		}

		if f, err := os.Open(path); err == nil {
			if m, err := tag.ReadFrom(f); err == nil {
				if strings.TrimSpace(m.Title()) != "" {
					t.Title = strings.TrimSpace(m.Title())
				}
				if strings.TrimSpace(m.Artist()) != "" {
					t.Artist = strings.TrimSpace(m.Artist())
				}
				if strings.TrimSpace(m.Album()) != "" {
					t.Album = strings.TrimSpace(m.Album())
				}
				if pic := m.Picture(); pic != nil && len(pic.Data) > 0 {
					t.CoverData = pic.Data
				}
			}
			f.Close()
		}

		tracks = append(tracks, t)
		return nil
	})

	return tracks, err
}
