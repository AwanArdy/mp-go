package library

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type LyricLine struct {
	Time time.Duration
	Text string
}

var lrcRegex = regexp.MustCompile(`\[(\d{1,2}):(\d{2})(?:\.(\d{2,3}))?\](.*)`)

// ParseLRC mengekstrak baris bertimestamp dari format LRC
func ParseLRC(lrc string) []LyricLine {
	var lines []LyricLine
	for _, rawLine := range strings.Split(lrc, "\n") {
		rawLine = strings.TrimSpace(rawLine)
		matches := lrcRegex.FindStringSubmatch(rawLine)
		if len(matches) >= 5 {
			min, _ := strconv.Atoi(matches[1])
			sec, _ := strconv.Atoi(matches[2])
			msec := 0
			if matches[3] != "" {
				msec, _ = strconv.Atoi(matches[3])
				if len(matches[3]) == 2 {
					msec *= 10
				}
			}
			text := strings.TrimSpace(matches[4])
			if text == "" && len(lines) == 0 {
				continue
			}
			t := time.Duration(min)*time.Minute + time.Duration(sec)*time.Second + time.Duration(msec)*time.Millisecond
			lines = append(lines, LyricLine{
				Time: t,
				Text: text,
			})
		}
	}
	return lines
}

// getPrimaryArtist memotong nama kolaborasi agar pencarian database lirik akurat
func getPrimaryArtist(raw string) string {
	raw = strings.TrimSpace(raw)
	separators := []string{",", "/", ";", " feat.", " ft.", " feat ", " ft ", " with ", " x ", " & "}
	for _, sep := range separators {
		if idx := strings.Index(strings.ToLower(raw), sep); idx != -1 {
			raw = raw[:idx]
		}
	}
	return strings.TrimSpace(raw)
}

// cleanSongTitle membersihkan penomoran file dan keterangan album/edisi
func cleanSongTitle(raw string) string {
	re := regexp.MustCompile(`(?i)(#\d+|special edition|\(.*\)|\[.*\])`)
	cleaned := re.ReplaceAllString(raw, "")
	return strings.TrimSpace(cleaned)
}

// FetchSyncedLyrics mengambil lirik berurutan: YTMusic -> LRCLIB
func FetchSyncedLyrics(title, artist, album string, dur time.Duration) ([]LyricLine, error) {
	cleanTitle := cleanSongTitle(title)
	primaryArtist := getPrimaryArtist(artist)

	log.Printf("[LYRICS] Mencari lirik untuk: '%s' oleh '%s' (artis asli: '%s')", cleanTitle, primaryArtist, artist)

	// 1. Prioritas Utama: YTMusic API Wrapper Lokal
	if lines, err := fetchFromYTMusic(cleanTitle, primaryArtist); err == nil && len(lines) > 0 {
		log.Printf("[LYRICS] Berhasil didapatkan dari YTMusic API (%d baris)", len(lines))
		return lines, nil
	}

	// 2. Prioritas Kedua: LRCLIB API
	if lines, err := fetchFromLRCLIB(cleanTitle, primaryArtist, album, int(dur.Seconds())); err == nil && len(lines) > 0 {
		log.Printf("[LYRICS] Berhasil didapatkan dari LRCLIB (%d baris)", len(lines))
		return lines, nil
	}

	log.Printf("[LYRICS] Lirik tersinkronisasi tidak ditemukan di semua provider")
	return nil, fmt.Errorf("lirik tidak ditemukan")
}

// 1. Handler YTMusic API
func fetchFromYTMusic(title, artist string) ([]LyricLine, error) {
	query := url.QueryEscape(fmt.Sprintf("%s %s", title, artist))
	client := &http.Client{Timeout: 3 * time.Second}

	endpoints := []string{
		fmt.Sprintf("http://127.0.0.1:8989/lyrics?query=%s", query),
		fmt.Sprintf("http://127.0.0.1:5000/lyrics?query=%s", query),
	}

	for _, ep := range endpoints {
		resp, err := client.Get(ep)
		if err != nil || resp.StatusCode != http.StatusOK {
			continue
		}
		defer resp.Body.Close()

		var res struct {
			Synced bool   `json:"synced"`
			LRC    string `json:"lrc"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&res); err == nil && res.Synced && res.LRC != "" {
			lines := ParseLRC(res.LRC)
			if len(lines) > 0 {
				return lines, nil
			}
		}
	}
	return nil, fmt.Errorf("ytmusicapi synced lyrics unavailable")
}

// 2. Handler LRCLIB API
type lrclibRecord struct {
	TrackName    string `json:"trackName"`
	ArtistName   string `json:"artistName"`
	SyncedLyrics string `json:"syncedLyrics"`
}

func fetchFromLRCLIB(title, artist, album string, durSec int) ([]LyricLine, error) {
	client := &http.Client{Timeout: 6 * time.Second}

	// Tahap 1: Gunakan search query yang persis seperti uji coba curl (title + primary artist)
	searchQuery := strings.TrimSpace(title + " " + artist)
	searchURL := fmt.Sprintf("https://lrclib.net/api/search?q=%s", url.QueryEscape(searchQuery))

	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("User-Agent", "GnomeMusicPlayer/1.0 (Linux; x86_64)")

	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		var list []lrclibRecord
		if err := json.Unmarshal(bodyBytes, &list); err == nil {
			for _, item := range list {
				// Pastikan hanya mengambil lirik bertimestamp (synced)
				if strings.TrimSpace(item.SyncedLyrics) != "" {
					lines := ParseLRC(item.SyncedLyrics)
					if len(lines) > 0 {
						return lines, nil
					}
				}
			}
		}
	}

	// Tahap 2: Fallback ke direct get bila search tahap 1 tidak membuahkan hasil
	baseURL := "https://lrclib.net/api/get?"
	params := url.Values{}
	params.Set("track_name", title)
	if artist != "" {
		params.Set("artist_name", artist)
	}
	if durSec > 0 {
		params.Set("duration", strconv.Itoa(durSec))
	}

	getReq, _ := http.NewRequest("GET", baseURL+params.Encode(), nil)
	getReq.Header.Set("User-Agent", "GnomeMusicPlayer/1.0 (Linux; x86_64)")

	if getResp, err := client.Do(getReq); err == nil && getResp.StatusCode == http.StatusOK {
		defer getResp.Body.Close()
		var item lrclibRecord
		if err := json.NewDecoder(getResp.Body).Decode(&item); err == nil && strings.TrimSpace(item.SyncedLyrics) != "" {
			lines := ParseLRC(item.SyncedLyrics)
			if len(lines) > 0 {
				return lines, nil
			}
		}
	}

	return nil, fmt.Errorf("lirik tidak ditemukan di LRCLIB")
}
