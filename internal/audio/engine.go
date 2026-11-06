package audio

import (
	"crypto/md5"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/audio"
	_ "github.com/egoist/mygo/audio/mp3"
	_ "github.com/egoist/mygo/audio/vorbis"
	"music-player/internal/model"
)

type Engine struct {
	mu          sync.Mutex
	player      *audio.Player
	currentFile *os.File
	isPlaying   bool
	volume      float64
	onEnd       func()
	onNext      func()
	onPrev      func()
}

func NewEngine() *Engine {
	e := &Engine{
		volume: 0.8,
	}

	// Signature callback: func(cmd mygo.MediaCommand, position time.Duration)
	mygo.NowPlaying.OnCommand(func(cmd mygo.MediaCommand, position time.Duration) {
		switch cmd {
		case mygo.MediaPlay:
			if !e.IsPlaying() {
				e.TogglePlay()
			}
		case mygo.MediaPause:
			if e.IsPlaying() {
				e.TogglePlay()
			}
		case mygo.MediaNext:
			e.mu.Lock()
			nextFn := e.onNext
			e.mu.Unlock()
			if nextFn != nil {
				nextFn()
			}
		case mygo.MediaPrevious:
			e.mu.Lock()
			prevFn := e.onPrev
			e.mu.Unlock()
			if prevFn != nil {
				prevFn()
			}
		}
	})

	return e
}

func (e *Engine) SetHandlers(onEnd, onNext, onPrev func()) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.onEnd = onEnd
	e.onNext = onNext
	e.onPrev = onPrev
}

func (e *Engine) Play(t model.Track) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	ext := strings.ToLower(filepath.Ext(t.Path))
	var f *os.File
	var err error

	if ext == ".mp3" || ext == ".wav" || ext == ".ogg" {
		f, err = os.Open(t.Path)
		if err == nil {
			stream, errStream := audio.OpenStream(f)
			if errStream == nil {
				e.setupPlayer(stream, f, t)
				return nil
			}
			f.Close()
		}
	}

	f, err = transcodeToWav(t.Path)
	if err != nil {
		return fmt.Errorf("transcode audio gagal: %w", err)
	}

	stream, err := audio.OpenStream(f)
	if err != nil {
		f.Close()
		return fmt.Errorf("decode audio stream gagal: %w", err)
	}

	e.setupPlayer(stream, f, t)
	return nil
}

func (e *Engine) setupPlayer(stream *audio.Player, f *os.File, t model.Track) {
	if e.player != nil {
		e.player.Close()
	}
	if e.currentFile != nil {
		e.currentFile.Close()
	}

	e.currentFile = f
	e.player = stream
	e.player.SetVolume(e.volume)
	e.player.Play()
	e.isPlaying = true

	e.player.OnEnd(func() {
		e.mu.Lock()
		e.isPlaying = false
		handler := e.onNext
		if handler == nil {
			handler = e.onEnd
		}
		e.mu.Unlock()

		mygo.NowPlaying.Sync()
		if handler != nil {
			handler()
		}
	})

	mygo.NowPlaying.Control(e.player)

	// Artwork langsung menerima raw []byte
	mygo.NowPlaying.Set(mygo.NowPlayingInfo{
		Title:    t.Title,
		Artist:   t.Artist,
		Album:    t.Album,
		Artwork:  t.CoverData,
		Commands: mygo.MediaPlay | mygo.MediaPause | mygo.MediaNext | mygo.MediaPrevious,
	})
	mygo.NowPlaying.Sync()
}

func transcodeToWav(srcPath string) (*os.File, error) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		ext := filepath.Ext(srcPath)
		return nil, fmt.Errorf("format %s memerlukan ffmpeg. Pasang dengan: sudo apt install -y ffmpeg", ext)
	}

	hash := fmt.Sprintf("%x", md5.Sum([]byte(srcPath)))
	tempWav := filepath.Join(os.TempDir(), fmt.Sprintf("mygo_cache_%s.wav", hash))

	if info, err := os.Stat(tempWav); err == nil && info.Size() > 100 {
		return os.Open(tempWav)
	}

	cmd := exec.Command("ffmpeg", "-y", "-v", "error", "-i", srcPath, "-acodec", "pcm_s16le", "-ar", "44100", "-ac", "2", tempWav)
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("ffmpeg error: %w (%s)", err, string(out))
	}

	return os.Open(tempWav)
}

func (e *Engine) TogglePlay() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.player == nil {
		return false
	}

	if e.isPlaying {
		e.player.Pause()
		e.isPlaying = false
	} else {
		e.player.Play()
		e.isPlaying = true
	}
	mygo.NowPlaying.Sync()
	return e.isPlaying
}

func (e *Engine) IsPlaying() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.isPlaying
}

func (e *Engine) Position() time.Duration {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.player == nil {
		return 0
	}
	return e.player.Position()
}

func (e *Engine) Duration() time.Duration {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.player == nil {
		return 0
	}
	return e.player.Duration()
}

func (e *Engine) Progress() float64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.player == nil {
		return 0
	}
	dur := e.player.Duration()
	if dur <= 0 {
		return 0
	}
	p := float64(e.player.Position()) / float64(dur)
	if p > 1.0 {
		p = 1.0
	}
	return p
}

func (e *Engine) Seek(progress float64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.player == nil {
		return
	}
	dur := e.player.Duration()
	if dur <= 0 {
		return
	}
	target := time.Duration(progress * float64(dur))
	e.player.Seek(target)
}

func (e *Engine) SetVolume(val float64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.volume = val
	if e.player != nil {
		e.player.SetVolume(val)
	}
}
