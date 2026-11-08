package ui

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"sync"

	"github.com/egoist/mygo/ui"
)

var (
	imageCacheLock sync.RWMutex
	imageCache     = make(map[string]ui.ImageSource)
)

type subImager interface {
	SubImage(r image.Rectangle) image.Image
}

// cropToSquare memotong bagian tengah gambar agar memiliki aspect ratio 1:1 tanpa distorsi
func cropToSquare(src image.Image) image.Image {
	bounds := src.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	if w <= 0 || h <= 0 || w == h {
		return src
	}

	size := w
	if h < size {
		size = h
	}

	x0 := bounds.Min.X + (w-size)/2
	y0 := bounds.Min.Y + (h-size)/2
	cropRect := image.Rect(x0, y0, x0+size, y0+size)

	if si, ok := src.(subImager); ok {
		return si.SubImage(cropRect)
	}

	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dst.Set(x, y, src.At(x0+x, y0+y))
		}
	}
	return dst
}

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

	// Coba decode image standar untuk meng-crop ke 1:1 square
	if rawImg, _, err := image.Decode(bytes.NewReader(data)); err == nil {
		sq := cropToSquare(rawImg)
		bmp := ui.NewBitmap(sq)
		imageCacheLock.Lock()
		imageCache[hash] = bmp
		imageCacheLock.Unlock()
		return bmp, true
	}

	// Fallback ke ui.DecodeImage bawaan jika format lain (seperti SVG)
	decoded, err := ui.DecodeImage(data)
	if err != nil {
		return nil, false
	}

	imageCacheLock.Lock()
	imageCache[hash] = decoded
	imageCacheLock.Unlock()

	return decoded, true
}
