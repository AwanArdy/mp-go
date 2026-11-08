package ui

import (
	"crypto/md5"
	"fmt"
	"sync"

	"github.com/egoist/mygo/ui"
)

var (
	imageCacheLock sync.RWMutex
	imageCache     = make(map[string]ui.ImageSource)
)

func GetOrDecodeImage(data []byte) (ui.ImageSource, bool) {
	if len(data) == 0 {
		return nil, false
	}

	hash := fmt.Sprintf("%x", md5.Sum(data))

	imageCacheLock.RLock()
	img, found := imageCache[hash]
	imageCacheLock.RUnlock()
	if found {
		return img, true
	}

	decoded, err := ui.DecodeImage(data)
	if err != nil {
		return nil, false
	}

	imageCacheLock.Lock()
	imageCache[hash] = decoded
	imageCacheLock.Unlock()

	return decoded, true
}
