package ui

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

var (
	colorBlack     = ui.RGB(0, 0, 0)
	colorWhite     = ui.RGB(255, 255, 255)
	colorLightGray = ui.RGB(220, 220, 220)
	colorDarkGray  = ui.RGB(45, 45, 45)
)

func customButton(c *ui.Context, label string, active bool, bg, fg ui.Color) ui.Element {
	btn := ui.Button(c, label).Border(1, fg).Radius(0)
	btnBg := bg
	btnFg := fg
	if btn.Hovered() || active {
		btnBg = fg
		btnFg = bg
	}
	btn.Background(btnBg).TextColor(btnFg)
	return btn
}

func customIconButton(c *ui.Context, w, h float32, active bool, bg, fg ui.Color, drawIcon func(p *ui.Painter, r ui.Rect, fg ui.Color)) ui.Element {
	btn := ui.Button(c, "").Width(w).Height(h).Border(1, fg).Radius(0)
	btnBg := bg
	btnFg := fg
	if btn.Hovered() || active {
		btnBg = fg
		btnFg = bg
	}
	btn.Background(btnBg)
	btn.Draw(func(p *ui.Painter, r ui.Rect) {
		if drawIcon != nil {
			drawIcon(p, r, btnFg)
		}
	})
	return btn
}

func renderCoverBox(c *ui.Context, coverData []byte, size float32, radius float32, fg ui.Color) {
	box := ui.Row(c).Width(size).Height(size).Border(1, fg).Radius(radius).Clip().Center()
	if len(coverData) > 0 {
		img, err := ui.DecodeImage(coverData)
		if err == nil {
			box.Children(func() {
				ui.Image(c, img).Width(size).Height(size)
			})
			return
		}
	}
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		drawNoteIcon(p, r, fg)
	})
}

func customSlider(c *ui.Context, val *float64, min, max float64, fg, muted ui.Color) ui.Element {
	s := ui.SliderBase(c, val, min, max).Height(20).Radius(0)
	ratio := float32((*val - min) / (max - min))
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	s.Draw(func(p *ui.Painter, r ui.Rect) {
		trackH := float32(4)
		trackY := r.Y + (r.H-trackH)/2

		p.Fill(ui.Rect{X: r.X, Y: trackY, W: r.W, H: trackH}, muted, 0)
		p.Fill(ui.Rect{X: r.X, Y: trackY, W: r.W * ratio, H: trackH}, fg, 0)

		thumbW := float32(8)
		thumbH := float32(14)
		thumbX := r.X + (r.W-thumbW)*ratio
		thumbY := r.Y + (r.H-thumbH)/2
		p.Fill(ui.Rect{X: thumbX, Y: thumbY, W: thumbW, H: thumbH}, fg, 0)
	})

	return s
}

// -------------------------------------------------------------
// Vektor Ikon Monokrom
// -------------------------------------------------------------
func drawMic(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2 - 2

	// 1. Kapsul kepala mikrofon
	p.Fill(ui.Rect{X: cx - 3, Y: cy - 7, W: 6, H: 9}, fg, 2)

	// 2. Kerangka dudukan U di sekeliling mikrofon
	p.Fill(ui.Rect{X: cx - 5, Y: cy - 3, W: 1.5, H: 6}, fg, 0)
	p.Fill(ui.Rect{X: cx + 3.5, Y: cy - 3, W: 1.5, H: 6}, fg, 0)
	p.Fill(ui.Rect{X: cx - 5, Y: cy + 3, W: 10, H: 1.5}, fg, 0)

	// 3. Tangkai tiang
	p.Fill(ui.Rect{X: cx - 0.75, Y: cy + 4.5, W: 1.5, H: 4}, fg, 0)

	// 4. Garis alas bawah
	p.Fill(ui.Rect{X: cx - 4, Y: cy + 8.5, W: 8, H: 1.5}, fg, 0)
}

func drawMoon(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	radius := float32(6)
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			d1 := x*x + y*y
			dx2 := x - 3
			dy2 := y + 1
			d2 := dx2*dx2 + dy2*dy2
			if d1 <= radius*radius && d2 >= (radius-1)*(radius-1) {
				p.Fill(ui.Rect{X: cx + x, Y: cy + y, W: 1.2, H: 1.2}, fg, 0)
			}
		}
	}
}

func drawSun(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	p.Fill(ui.Rect{X: cx - 3, Y: cy - 3, W: 6, H: 6}, fg, 0)
	p.Fill(ui.Rect{X: cx - 1, Y: cy - 7, W: 2, H: 3}, fg, 0)
	p.Fill(ui.Rect{X: cx - 1, Y: cy + 4, W: 2, H: 3}, fg, 0)
	p.Fill(ui.Rect{X: cx - 7, Y: cy - 1, W: 3, H: 2}, fg, 0)
	p.Fill(ui.Rect{X: cx + 4, Y: cy - 1, W: 3, H: 2}, fg, 0)
	p.Fill(ui.Rect{X: cx - 5, Y: cy - 5, W: 2, H: 2}, fg, 0)
	p.Fill(ui.Rect{X: cx + 3, Y: cy - 5, W: 2, H: 2}, fg, 0)
	p.Fill(ui.Rect{X: cx - 5, Y: cy + 3, W: 2, H: 2}, fg, 0)
	p.Fill(ui.Rect{X: cx + 3, Y: cy + 3, W: 2, H: 2}, fg, 0)
}

func drawSearchIcon(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2 - 1
	cy := r.Y + r.H/2 - 1
	p.Fill(ui.Rect{X: cx - 4, Y: cy - 4, W: 8, H: 1.5}, fg, 0)
	p.Fill(ui.Rect{X: cx - 4, Y: cy + 3, W: 8, H: 1.5}, fg, 0)
	p.Fill(ui.Rect{X: cx - 4, Y: cy - 3, W: 1.5, H: 6}, fg, 0)
	p.Fill(ui.Rect{X: cx + 3, Y: cy - 3, W: 1.5, H: 6}, fg, 0)
	for i := float32(0); i < 4; i++ {
		p.Fill(ui.Rect{X: cx + 3 + i, Y: cy + 3 + i, W: 1.5, H: 1.5}, fg, 0)
	}
}

func drawMinimize(p *ui.Painter, r ui.Rect, fg ui.Color) {
	w := float32(12)
	h := float32(2)
	x := r.X + (r.W-w)/2
	y := r.Y + (r.H-h)/2 + 2
	p.Fill(ui.Rect{X: x, Y: y, W: w, H: h}, fg, 0)
}

func drawClose(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	size := float32(8)
	half := size / 2
	for i := float32(0); i <= size; i++ {
		p.Fill(ui.Rect{X: cx - half + i, Y: cy - half + i, W: 1.5, H: 1.5}, fg, 0)
		p.Fill(ui.Rect{X: cx + half - i, Y: cy - half + i, W: 1.5, H: 1.5}, fg, 0)
	}
}

func drawPrev(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	p.Fill(ui.Rect{X: cx - 6, Y: cy - 6, W: 2, H: 12}, fg, 0)
	for col := float32(0); col < 8; col++ {
		h := 12 - col*1.5
		y := cy - h/2
		p.Fill(ui.Rect{X: cx + 4 - col, Y: y, W: 1.5, H: h}, fg, 0)
	}
}

func drawNext(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	for col := float32(0); col < 8; col++ {
		h := 12 - col*1.5
		y := cy - h/2
		p.Fill(ui.Rect{X: cx - 4 + col, Y: y, W: 1.5, H: h}, fg, 0)
	}
	p.Fill(ui.Rect{X: cx + 5, Y: cy - 6, W: 2, H: 12}, fg, 0)
}

func drawPlay(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2 - 1
	cy := r.Y + r.H/2
	for col := float32(0); col < 10; col++ {
		h := 14 - col*1.4
		y := cy - h/2
		p.Fill(ui.Rect{X: cx - 4 + col, Y: y, W: 1.5, H: h}, fg, 0)
	}
}

func drawPause(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	p.Fill(ui.Rect{X: cx - 5, Y: cy - 7, W: 3, H: 14}, fg, 0)
	p.Fill(ui.Rect{X: cx + 2, Y: cy - 7, W: 3, H: 14}, fg, 0)
}

func drawShuffle(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	p.Fill(ui.Rect{X: cx - 7, Y: cy - 5, W: 4, H: 2}, fg, 0)
	p.Fill(ui.Rect{X: cx - 7, Y: cy + 3, W: 4, H: 2}, fg, 0)
	for i := float32(0); i < 8; i++ {
		p.Fill(ui.Rect{X: cx - 3 + i, Y: cy - 4 + i, W: 1.5, H: 1.5}, fg, 0)
		p.Fill(ui.Rect{X: cx - 3 + i, Y: cy + 4 - i, W: 1.5, H: 1.5}, fg, 0)
	}
	p.Fill(ui.Rect{X: cx + 4, Y: cy + 3, W: 4, H: 2}, fg, 0)
	p.Fill(ui.Rect{X: cx + 4, Y: cy - 5, W: 4, H: 2}, fg, 0)
}

func drawRepeat(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	p.Fill(ui.Rect{X: cx - 6, Y: cy - 5, W: 12, H: 2}, fg, 0)
	p.Fill(ui.Rect{X: cx - 6, Y: cy + 4, W: 12, H: 2}, fg, 0)
	p.Fill(ui.Rect{X: cx - 6, Y: cy - 4, W: 2, H: 9}, fg, 0)
	p.Fill(ui.Rect{X: cx + 4, Y: cy - 4, W: 2, H: 9}, fg, 0)
	p.Fill(ui.Rect{X: cx + 5, Y: cy + 2, W: 2, H: 2}, fg, 0)
	p.Fill(ui.Rect{X: cx + 5, Y: cy + 6, W: 2, H: 2}, fg, 0)
}

func drawNoteIcon(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	p.Fill(ui.Rect{X: cx - 4, Y: cy - 7, W: 2, H: 10}, fg, 0)
	p.Fill(ui.Rect{X: cx - 7, Y: cy + 1, W: 5, H: 4}, fg, 0)
	p.Fill(ui.Rect{X: cx + 3, Y: cy - 5, W: 2, H: 8}, fg, 0)
	p.Fill(ui.Rect{X: cx, Y: cy + 1, W: 5, H: 4}, fg, 0)
	p.Fill(ui.Rect{X: cx - 4, Y: cy - 7, W: 9, H: 3}, fg, 0)
}

// -------------------------------------------------------------
// Root View
// -------------------------------------------------------------
func (s *AppState) View(c *ui.Context) {
	themeBg := colorWhite
	themeFg := colorBlack
	themeMuted := colorLightGray
	if s.DarkMode {
		themeBg = colorBlack
		themeFg = colorWhite
		themeMuted = colorDarkGray
	}

	if s.Engine.IsPlaying() {
		c.After(250 * time.Millisecond)
		s.Progress = s.Engine.Progress()
	}

	ui.Column(c).Fill().Background(themeBg).Radius(0).Children(func() {
		renderTopBar(c, s, themeBg, themeFg)

		ui.Row(c).Fill().MinHeight(0).Padding(10).Gap(10).Children(func() {
			renderSidebar(c, s, themeBg, themeFg)
			renderMainContent(c, s, themeBg, themeFg)
		})

		renderBottomPlaybackBar(c, s, themeBg, themeFg, themeMuted)
	})
}

func renderTopBar(c *ui.Context, s *AppState, bg, fg ui.Color) {
	ui.Row(c).FillWidth().Height(46).Shrink(0).Padding(6, 12).Border(1, fg).Radius(0).Background(bg).Center().Gap(12).Children(func() {
		ui.Row(c).Center().Gap(8).Children(func() {
			ui.Text(c, "MUSIC PLAYER").FontSize(13).Bold().TextColor(fg)

			themeIconDraw := drawMoon
			if s.DarkMode {
				themeIconDraw = drawSun
			}
			if customIconButton(c, 28, 28, false, bg, fg, themeIconDraw).Clicked() {
				s.ToggleDarkMode()
			}
		})

		ui.Row(c).FillWidth().Height(30).Border(1, fg).Radius(0).Background(bg).Center().Children(func() {
			iconBox := ui.Row(c).Width(28).Height(28).Center()
			iconBox.Draw(func(p *ui.Painter, r ui.Rect) {
				drawSearchIcon(p, r, fg)
			})

			ui.TextInput(c, &s.SearchQuery).
				Border(0, bg).
				Radius(0).
				FillWidth().
				Height(28).
				Background(bg).
				TextColor(fg)
		})

		ui.Row(c).Gap(6).Children(func() {
			if customIconButton(c, 36, 28, false, bg, fg, drawMinimize).Clicked() {
				s.MinimizeApp()
			}
			if customIconButton(c, 36, 28, false, bg, fg, drawClose).Clicked() {
				s.CloseApp()
			}
		})
	})
}

func renderSidebar(c *ui.Context, s *AppState, bg, fg ui.Color) {
	ui.Column(c).Width(220).FillHeight().MinHeight(0).Border(1, fg).Radius(0).Padding(10).Gap(8).Children(func() {
		ui.Text(c, "NAVIGATION").FontSize(11).Bold().TextColor(fg)

		if customButton(c, "ALL TRACKS", s.ViewMode == "tracks" && s.SelectedArtist == "" && s.SelectedAlbum == "", bg, fg).FillWidth().Height(30).Clicked() {
			s.ClearFilter()
		}
		if customButton(c, "BY ARTISTS", s.ViewMode == "artists", bg, fg).FillWidth().Height(30).Clicked() {
			s.ViewMode = "artists"
		}
		if customButton(c, "BY ALBUMS", s.ViewMode == "albums", bg, fg).FillWidth().Height(30).Clicked() {
			s.ViewMode = "albums"
		}

		ui.Row(c).FillWidth().Height(1).Border(1, fg).Radius(0)

		ui.Text(c, "ACTIVE FILTER").FontSize(11).Bold().TextColor(fg)
		filterText := "None (All Songs)"
		if s.SelectedArtist != "" {
			filterText = "Artist: " + s.SelectedArtist
		} else if s.SelectedAlbum != "" {
			filterText = "Album: " + s.SelectedAlbum
		}
		ui.Text(c, filterText).FontSize(10).TextColor(fg)

		if s.SelectedArtist != "" || s.SelectedAlbum != "" || s.SearchQuery != "" {
			if customButton(c, "[RESET FILTER]", false, bg, fg).FillWidth().Height(26).Clicked() {
				s.ClearFilter()
			}
		}

		ui.Column(c).Fill().MinHeight(0)
	})
}

// -------------------------------------------------------------
// 3. Konten Utama (Daftar Lagu / Artis / Album / Lirik)
// -------------------------------------------------------------
func renderMainContent(c *ui.Context, s *AppState, bg, fg ui.Color) {
	ui.Column(c).Fill().FillHeight().MinHeight(0).Border(1, fg).Radius(0).Padding(10).Gap(8).Children(func() {
		switch s.ViewMode {
		case "artists":
			renderArtistsList(c, s, bg, fg)
		case "albums":
			renderAlbumsList(c, s, bg, fg)
		case "lyrics":
			renderLyricsView(c, s, bg, fg)
		default:
			renderTracksList(c, s, bg, fg)
		}
	})
}

// -------------------------------------------------------------
// Tampilan Lirik Lagu Tersinkronisasi
// -------------------------------------------------------------
func renderLyricsView(c *ui.Context, s *AppState, bg, fg ui.Color) {
	current := s.CurrentTrack()

	// 1. Header Bar Lirik
	ui.Row(c).FillWidth().Height(26).Shrink(0).Center().Children(func() {
		title := "NO TRACK PLAYING"
		if current != nil {
			title = fmt.Sprintf("LYRICS: %s — %s", current.Title, current.Artist)
		}
		ui.Text(c, title).FontSize(11).Bold().TextColor(fg)

		ui.Row(c).FillWidth().Right(0).Children(func() {
			if customButton(c, "BACK TO TRACKS", false, bg, fg).Height(24).Clicked() {
				s.ViewMode = "tracks"
			}
		})
	})

	// 2. Status Penanganan
	if current == nil {
		ui.Column(c).Fill().Center().Children(func() {
			ui.Text(c, "PUTAR LAGU UNTUK MELIHAT LIRIK").FontSize(13).Bold().TextColor(fg)
		})
		return
	}

	if s.LyricsStatus == "loading" {
		ui.Column(c).Fill().Center().Gap(6).Children(func() {
			ui.Text(c, "MENCARI LIRIK TERSINKRONISASI...").FontSize(13).Bold().TextColor(fg)
			ui.Text(c, "(YTMusic API > LRCLIB)").FontSize(10).TextColor(fg)
		})
		return
	}

	if s.LyricsStatus == "not_found" || len(s.LyricsLines) == 0 {
		ui.Column(c).Fill().Center().Gap(6).Children(func() {
			ui.Text(c, "LIRIK TIDAK DITEMUKAN").FontSize(14).Bold().TextColor(fg)
			ui.Text(c, "Hanya lirik tersinkronisasi (synced lyrics) yang didukung").FontSize(10).TextColor(fg)
		})
		return
	}

	// 3. Tentukan Indeks Lirik Berdasarkan Waktu Pemutaran
	currentPos := s.Engine.Position()
	activeIdx := 0
	for i, line := range s.LyricsLines {
		if line.Time <= currentPos {
			activeIdx = i
		} else {
			break
		}
	}

	// 4. Logika Easing Lerp untuk Animasi Scroll
	targetPos := float64(activeIdx)
	diff := targetPos - s.LyricVisualPos

	if math.Abs(diff) > 4.0 {
		s.LyricVisualPos = targetPos
	} else if math.Abs(diff) > 0.001 {
		step := diff * 0.22
		minStep := 0.035

		if diff > 0 {
			if step < minStep {
				step = minStep
			}
			if s.LyricVisualPos >= targetPos {
				s.LyricVisualPos = targetPos
			} else {
				s.LyricVisualPos += step
			}
		} else {
			if step > -minStep {
				step = -minStep
			}
			if s.LyricVisualPos+step <= targetPos {
				s.LyricVisualPos = targetPos
			} else {
				s.LyricVisualPos += step
			}
		}

		if s.LyricVisualPos != targetPos {
			c.After(16 * time.Millisecond)
		}
	} else {
		s.LyricVisualPos = targetPos
	}

	// 5. Warna Abu-Abu Adaptif Sesuai Mode
	grayMuted := ui.RGB(145, 145, 145)
	grayDim := ui.RGB(90, 90, 90)
	if s.DarkMode {
		grayMuted = ui.RGB(160, 160, 160)
		grayDim = ui.RGB(85, 85, 85)
	}

	// 6. Parameter Viewport 7 Baris
	const lineH float32 = 46
	const viewportH float32 = 7 * lineH // 322px

	baseIdx := int(math.Floor(s.LyricVisualPos))
	frac := float32(s.LyricVisualPos - float64(baseIdx))

	// Kontainer Terluar: Berada tepat di tengah layar secara vertikal & horizontal
	ui.Column(c).Fill().Center().Children(func() {
		// Viewport Frameless setinggi 7 baris dengan .Clip()
		ui.Column(c).FillWidth().Height(viewportH).Clip().Children(func() {
			// Render 8 baris (-3 sampai +4): Baris teratas menyusut seiring frac, mendorong baris bawah naik
			for k := -3; k <= 4; k++ {
				lineIdx := baseIdx + k

				// Baris teratas menyusut tingginya dari 46px ke 0px
				rowH := lineH
				if k == -3 {
					rowH = (1.0 - frac) * lineH
				}

				// Baris kosong di luar jangkauan lagu (misal pada intro atau outro)
				if lineIdx < 0 || lineIdx >= len(s.LyricsLines) {
					ui.Row(c).FillWidth().Height(rowH)
					continue
				}

				line := s.LyricsLines[lineIdx]
				displayText := line.Text
				if displayText == "" {
					displayText = "• • •" // Indikator jeda musik / instrumental
				}

				// Hitung jarak visual ke pusat untuk menentukan ukuran & kontras warna
				dist := float64(lineIdx) - s.LyricVisualPos
				absDist := math.Abs(dist)

				fontSize := float32(14)
				textColor := grayDim
				isBold := false

				if absDist < 0.6 {
					fontSize = 24
					textColor = fg // Hitam pekat di Light Mode, Putih murni di Dark Mode
					isBold = true
				} else if absDist < 1.6 {
					fontSize = 18
					textColor = grayMuted
				} else if absDist < 2.6 {
					fontSize = 16
					textColor = grayMuted
				}

				// Render Baris Teks Lirik
				row := ui.Row(c).FillWidth().Height(rowH).Clip().Center()
				if rowH > 14 {
					row.Children(func() {
						txt := ui.Text(c, displayText).FontSize(fontSize).TextColor(textColor)
						if isBold {
							txt.Bold()
						}
					})
				}

				// Klik baris lirik untuk melompat (seek) ke waktu baris tersebut
				if row.Clicked() {
					dur := s.Engine.Duration()
					if dur > 0 {
						s.Seek(float64(line.Time) / float64(dur))
					}
				}
			}
		})
	})
}

func renderArtistsList(c *ui.Context, s *AppState, bg, fg ui.Color) {
	cards := s.GetArtistCards()

	ui.Row(c).FillWidth().Height(24).Shrink(0).Children(func() {
		ui.Text(c, fmt.Sprintf("ARTISTS DIRECTORY (%d)", len(cards))).FontSize(12).Bold().TextColor(fg)
	})

	ui.Scroll(c).Fill().MinHeight(0).Children(func() {
		ui.Column(c).FillWidth().Gap(10).Children(func() {
			total := len(cards)
			for i := 0; i < total; i += 4 {
				end := i + 4
				if end > total {
					end = total
				}
				rowCards := cards[i:end]

				ui.Row(c).FillWidth().Gap(10).Children(func() {
					for _, item := range rowCards {
						card := item
						isSelected := s.SelectedArtist == card.Name

						cardBtn := ui.Button(c, "").
							Grow(1).
							Basis(0).
							MinWidth(0).
							Height(185).
							Border(1, fg).
							Radius(8)

						cardBg := bg
						cardFg := fg
						if cardBtn.Hovered() || isSelected {
							cardBg = fg
							cardFg = bg
						}
						cardBtn.Background(cardBg)

						displayName := card.Name
						if len([]rune(displayName)) > 22 {
							displayName = string([]rune(displayName)[:20]) + "..."
						}

						cardBtn.Children(func() {
							ui.Column(c).Fill().Padding(8).Gap(6).Center().Children(func() {
								renderCoverBox(c, card.CoverData, 100, 6, cardFg)
								ui.Text(c, displayName).FontSize(11).Bold().TextColor(cardFg)
								ui.Text(c, card.Subtext).FontSize(9).TextColor(cardFg)
							})
						})

						if cardBtn.Clicked() {
							s.SelectedArtist = card.Name
							s.SelectedAlbum = ""
							s.ViewMode = "tracks"
						}
					}

					for k := 0; k < 4-len(rowCards); k++ {
						ui.Row(c).Grow(1).Basis(0).MinWidth(0)
					}
				})
			}
		})
	})
}

func renderAlbumsList(c *ui.Context, s *AppState, bg, fg ui.Color) {
	cards := s.GetAlbumCards()

	ui.Row(c).FillWidth().Height(24).Shrink(0).Children(func() {
		ui.Text(c, fmt.Sprintf("ALBUMS DIRECTORY (%d)", len(cards))).FontSize(12).Bold().TextColor(fg)
	})

	ui.Scroll(c).Fill().MinHeight(0).Children(func() {
		ui.Column(c).FillWidth().Gap(10).Children(func() {
			total := len(cards)
			for i := 0; i < total; i += 4 {
				end := i + 4
				if end > total {
					end = total
				}
				rowCards := cards[i:end]

				ui.Row(c).FillWidth().Gap(10).Children(func() {
					for _, item := range rowCards {
						card := item
						isSelected := s.SelectedAlbum == card.Name

						cardBtn := ui.Button(c, "").
							Grow(1).
							Basis(0).
							MinWidth(0).
							Height(185).
							Border(1, fg).
							Radius(8)

						cardBg := bg
						cardFg := fg
						if cardBtn.Hovered() || isSelected {
							cardBg = fg
							cardFg = bg
						}
						cardBtn.Background(cardBg)

						displayName := card.Name
						if len([]rune(displayName)) > 22 {
							displayName = string([]rune(displayName)[:20]) + "..."
						}

						cardBtn.Children(func() {
							ui.Column(c).Fill().Padding(8).Gap(6).Center().Children(func() {
								renderCoverBox(c, card.CoverData, 100, 6, cardFg)
								ui.Text(c, displayName).FontSize(11).Bold().TextColor(cardFg)
								ui.Text(c, card.Subtext).FontSize(9).TextColor(cardFg)
							})
						})

						if cardBtn.Clicked() {
							s.SelectedAlbum = card.Name
							s.SelectedArtist = ""
							s.ViewMode = "tracks"
						}
					}

					for k := 0; k < 4-len(rowCards); k++ {
						ui.Row(c).Grow(1).Basis(0).MinWidth(0)
					}
				})
			}
		})
	})
}

func renderTracksList(c *ui.Context, s *AppState, bg, fg ui.Color) {
	tracks := s.FilteredTracks()

	ui.Row(c).FillWidth().Height(24).Shrink(0).Center().Children(func() {
		ui.Text(c, fmt.Sprintf("TRACKS (%d items)", len(tracks))).FontSize(12).Bold().TextColor(fg)
	})

	ui.Scroll(c).Fill().MinHeight(0).Children(func() {
		ui.Column(c).FillWidth().Gap(4).Children(func() {
			if len(tracks) == 0 {
				ui.Row(c).Fill().Center().Children(func() {
					ui.Text(c, "NO TRACKS FOUND").FontSize(12).Bold().TextColor(fg)
				})
				return
			}

			for idx, t := range tracks {
				track := t
				isSelected := s.CurrentTrack() != nil && s.CurrentTrack().Path == track.Path
				rowNum := idx + 1

				rowBg := bg
				textColor := fg
				btnBorderColor := fg
				if isSelected {
					rowBg = fg
					textColor = bg
					btnBorderColor = bg
				}

				ui.Row(c).FillWidth().Height(36).Padding(4, 10).Border(1, fg).Radius(0).Background(rowBg).Center().Children(func() {
					prefix := fmt.Sprintf("[%02d]", rowNum)
					if isSelected {
						prefix = "[*]"
					}
					ui.Text(c, prefix).Width(38).FontSize(11).Bold().TextColor(textColor)

					ui.Column(c).FillWidth().Children(func() {
						ui.Text(c, track.Title).FontSize(12).Bold().TextColor(textColor)
						ui.Text(c, track.Artist).FontSize(10).TextColor(textColor)
					})

					ui.Text(c, track.Album).Width(180).FontSize(10).TextColor(textColor)

					if isSelected {
						playIconDraw := drawPause
						if !s.Engine.IsPlaying() {
							playIconDraw = drawPlay
						}

						btn := ui.Button(c, "").Width(36).Height(26).Border(1, btnBorderColor).Radius(0)
						bBg := bg
						bFg := fg
						if btn.Hovered() {
							bBg = fg
							bFg = bg
						}
						btn.Background(bBg)
						btn.Draw(func(p *ui.Painter, r ui.Rect) {
							playIconDraw(p, r, bFg)
						})
						if btn.Clicked() {
							s.TogglePlay()
						}
					} else {
						btn := customIconButton(c, 36, 26, false, bg, fg, drawPlay)
						if btn.Clicked() {
							s.PlayTrack(track)
						}
					}
				})
			}
		})
	})
}

// -------------------------------------------------------------
// 4. Playback Bar di Bawah: Ikon Mikrofon di Sisi Kanan Info Lagu
// -------------------------------------------------------------
func renderBottomPlaybackBar(c *ui.Context, s *AppState, bg, fg, muted ui.Color) {
	current := s.CurrentTrack()

	ui.Row(c).FillWidth().Height(78).Shrink(0).Padding(10, 14).Border(1, fg).Radius(0).Background(bg).Center().Children(func() {
		// --- KIRI: Cover Art + Info Lagu + Tombol Mikrofon ---
		ui.Row(c).Width(310).Center().Gap(10).Children(func() {
			var cover []byte
			if current != nil {
				cover = current.CoverData
			}
			renderCoverBox(c, cover, 52, 6, fg)

			ui.Column(c).FillWidth().Children(func() {
				title := "IDLE"
				artist := "READY"
				if current != nil {
					title = current.Title
					artist = current.Artist
				}
				ui.Text(c, title).FontSize(12).Bold().TextColor(fg)
				ui.Text(c, artist).FontSize(10).TextColor(fg)
			})

			// Tombol Mikrofon Vektor (Klik untuk berganti ke/dari mode lirik)
			isLyricsMode := s.ViewMode == "lyrics"
			if customIconButton(c, 30, 30, isLyricsMode, bg, fg, drawMic).Clicked() {
				if s.ViewMode == "lyrics" {
					s.ViewMode = "tracks"
				} else {
					s.ViewMode = "lyrics"
					if current != nil {
						s.LoadLyrics(*current)
					}
				}
			}
		})

		// --- TENGAH: Kontrol Navigasi & Timeline Slider ---
		ui.Column(c).FillWidth().Center().Gap(4).Children(func() {
			ui.Row(c).Center().Gap(6).Children(func() {
				if customIconButton(c, 36, 26, s.Shuffle, bg, fg, drawShuffle).Clicked() {
					s.Shuffle = !s.Shuffle
				}
				if customIconButton(c, 36, 26, false, bg, fg, drawPrev).Clicked() {
					s.PrevTrack()
				}

				playDraw := drawPlay
				if s.Engine.IsPlaying() {
					playDraw = drawPause
				}
				if customIconButton(c, 44, 26, s.Engine.IsPlaying(), bg, fg, playDraw).Clicked() {
					s.TogglePlay()
				}

				if customIconButton(c, 36, 26, false, bg, fg, drawNext).Clicked() {
					s.NextTrack()
				}
				if customIconButton(c, 36, 26, s.Repeat, bg, fg, drawRepeat).Clicked() {
					s.Repeat = !s.Repeat
				}
			})

			ui.Row(c).FillWidth().Center().Padding(0, 16).Gap(8).Children(func() {
				ui.Text(c, s.CurrentTimeString()).FontSize(10).Bold().TextColor(fg)

				progSlider := customSlider(c, &s.Progress, 0.0, 1.0, fg, muted).FillWidth()
				if progSlider.Changed() {
					s.Seek(s.Progress)
				}

				ui.Text(c, s.TotalTimeString()).FontSize(10).Bold().TextColor(fg)
			})
		})

		// --- KANAN: Volume Slider ---
		ui.Row(c).Width(160).Right(0).Center().Gap(6).Children(func() {
			ui.Text(c, "VOL:").FontSize(10).Bold().TextColor(fg)

			volSlider := customSlider(c, &s.Volume, 0.0, 1.0, fg, muted).Width(70)
			if volSlider.Changed() {
				s.Engine.SetVolume(s.Volume)
			}

			ui.Text(c, fmt.Sprintf("%d%%", int(s.Volume*100))).FontSize(10).TextColor(fg)
		})
	})
}
