package ui

import (
	"fmt"

	"github.com/egoist/mygo/ui"
)

var (
	colorBlack     = ui.RGB(0, 0, 0)
	colorWhite     = ui.RGB(255, 255, 255)
	colorLightGray = ui.RGB(220, 220, 220)
)

// Helper tombol teks biasa monokrom (0px radius)
func customButton(c *ui.Context, label string, active bool) ui.Element {
	btn := ui.Button(c, label).Border(1, colorBlack).Radius(0)
	if active {
		btn.Background(colorBlack).TextColor(colorWhite)
	} else {
		btn.Background(colorWhite).TextColor(colorBlack)
	}
	return btn
}

// Helper tombol ikon vektor monokrom
func customIconButton(c *ui.Context, w, h float32, active bool, drawIcon func(p *ui.Painter, r ui.Rect, fg ui.Color)) ui.Element {
	btn := ui.Button(c, "").Width(w).Height(h).Border(1, colorBlack).Radius(0)
	bg := colorWhite
	fg := colorBlack
	if active {
		bg = colorBlack
		fg = colorWhite
	}
	btn.Background(bg)
	btn.Draw(func(p *ui.Painter, r ui.Rect) {
		if drawIcon != nil {
			drawIcon(p, r, fg)
		}
	})
	return btn
}

// Helper gambar cover art: tampilkan gambar jika ada, fallback ke vektor not jika tidak ada
func renderCoverBox(c *ui.Context, coverData []byte, size float32, radius float32) {
	box := ui.Row(c).Width(size).Height(size).Border(1, colorBlack).Radius(radius).Clip().Center()
	if len(coverData) > 0 {
		img, err := ui.DecodeImage(coverData)
		if err == nil {
			box.Children(func() {
				ui.Image(c, img).Width(size).Height(size)
			})
			return
		}
	}
	// Fallback ke vektor not lagu
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		drawNoteIcon(p, r, colorBlack)
	})
}

// -------------------------------------------------------------
// Vektor Ikon Presisi
// -------------------------------------------------------------
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

func customSlider(c *ui.Context, val *float64, min, max float64) ui.Element {
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

		p.Fill(ui.Rect{X: r.X, Y: trackY, W: r.W, H: trackH}, colorLightGray, 0)
		p.Fill(ui.Rect{X: r.X, Y: trackY, W: r.W * ratio, H: trackH}, colorBlack, 0)

		thumbW := float32(8)
		thumbH := float32(14)
		thumbX := r.X + (r.W-thumbW)*ratio
		thumbY := r.Y + (r.H-thumbH)/2
		p.Fill(ui.Rect{X: thumbX, Y: thumbY, W: thumbW, H: thumbH}, colorBlack, 0)
	})

	return s
}

// -------------------------------------------------------------
// Root View
// -------------------------------------------------------------
func (s *AppState) View(c *ui.Context) {
	ui.Column(c).Fill().Background(colorWhite).Border(2, colorBlack).Radius(0).Children(func() {
		renderTopBar(c, s)

		ui.Row(c).Fill().MinHeight(0).Padding(10).Gap(10).Children(func() {
			renderSidebar(c, s)
			renderMainContent(c, s)
		})

		renderBottomPlaybackBar(c, s)
	})
}

func renderTopBar(c *ui.Context, s *AppState) {
	ui.Row(c).FillWidth().Height(46).Shrink(0).Padding(6, 12).Border(1, colorBlack).Radius(0).Background(colorWhite).Center().Children(func() {
		ui.Text(c, "MUSIC PLAYER").FontSize(13).Bold().TextColor(colorBlack).Width(130)

		ui.Row(c).FillWidth().Height(30).Border(1, colorBlack).Radius(0).Background(colorWhite).Center().Children(func() {
			iconBox := ui.Row(c).Width(28).Height(28).Center()
			iconBox.Draw(func(p *ui.Painter, r ui.Rect) {
				drawSearchIcon(p, r, colorBlack)
			})

			ui.TextInput(c, &s.SearchQuery).
				Border(0, colorWhite).
				Radius(0).
				FillWidth().
				Height(28).
				Background(colorWhite).
				TextColor(colorBlack)
		})

		ui.Row(c).Gap(6).Children(func() {
			if customIconButton(c, 36, 28, false, drawMinimize).Clicked() {
				s.MinimizeApp()
			}
			if customIconButton(c, 36, 28, false, drawClose).Clicked() {
				s.CloseApp()
			}
		})
	})
}

func renderSidebar(c *ui.Context, s *AppState) {
	ui.Column(c).Width(220).FillHeight().MinHeight(0).Border(1, colorBlack).Radius(0).Padding(10).Gap(8).Children(func() {
		ui.Text(c, "NAVIGATION").FontSize(11).Bold().TextColor(colorBlack)

		if customButton(c, "ALL TRACKS", s.ViewMode == "tracks" && s.SelectedArtist == "" && s.SelectedAlbum == "").FillWidth().Height(30).Clicked() {
			s.ClearFilter()
		}
		if customButton(c, "BY ARTISTS", s.ViewMode == "artists").FillWidth().Height(30).Clicked() {
			s.ViewMode = "artists"
		}
		if customButton(c, "BY ALBUMS", s.ViewMode == "albums").FillWidth().Height(30).Clicked() {
			s.ViewMode = "albums"
		}

		ui.Row(c).FillWidth().Height(1).Border(1, colorBlack).Radius(0)

		ui.Text(c, "ACTIVE FILTER").FontSize(11).Bold().TextColor(colorBlack)
		filterText := "None (All Songs)"
		if s.SelectedArtist != "" {
			filterText = "Artist: " + s.SelectedArtist
		} else if s.SelectedAlbum != "" {
			filterText = "Album: " + s.SelectedAlbum
		}
		ui.Text(c, filterText).FontSize(10).TextColor(colorBlack)

		if s.SelectedArtist != "" || s.SelectedAlbum != "" || s.SearchQuery != "" {
			if customButton(c, "[RESET FILTER]", false).FillWidth().Height(26).Clicked() {
				s.ClearFilter()
			}
		}

		ui.Column(c).Fill().MinHeight(0)
	})
}

func renderMainContent(c *ui.Context, s *AppState) {
	ui.Column(c).Fill().FillHeight().MinHeight(0).Border(1, colorBlack).Radius(0).Padding(10).Gap(8).Children(func() {
		switch s.ViewMode {
		case "artists":
			renderArtistsList(c, s)
		case "albums":
			renderAlbumsList(c, s)
		default:
			renderTracksList(c, s)
		}
	})
}

// -------------------------------------------------------------
// By Artists: Grid 4 Card ke Samping
// -------------------------------------------------------------
func renderArtistsList(c *ui.Context, s *AppState) {
	cards := s.GetArtistCards()

	ui.Row(c).FillWidth().Height(24).Shrink(0).Children(func() {
		ui.Text(c, fmt.Sprintf("ARTISTS DIRECTORY (%d)", len(cards))).FontSize(12).Bold().TextColor(colorBlack)
	})

	ui.Scroll(c).Fill().MinHeight(0).Children(func() {
		ui.Column(c).FillWidth().Gap(10).Children(func() {
			// Susun 4 card per baris
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

						cardBg := colorWhite
						cardFg := colorBlack
						if isSelected {
							cardBg = colorBlack
							cardFg = colorWhite
						}

						displayName := card.Name
						if len([]rune(displayName)) > 22 {
							displayName = string([]rune(displayName)[:20]) + "..."
						}

						// Card Container
						cardBtn := ui.Button(c, "").
							Grow(1).
							Basis(0).
							MinWidth(0).
							Height(185).
							Border(1, colorBlack).
							Radius(8).
							Background(cardBg)

						cardBtn.Children(func() {
							ui.Column(c).Fill().Padding(8).Gap(6).Center().Children(func() {
								renderCoverBox(c, card.CoverData, 100, 6)
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

					// Padding pengisi agar baris terakhir tetap rata kiri jika kurang dari 4
					for k := 0; k < 4-len(rowCards); k++ {
						ui.Row(c).FillWidth()
					}
				})
			}
		})
	})
}

// -------------------------------------------------------------
// By Albums: Grid 4 Card ke Samping
// -------------------------------------------------------------
func renderAlbumsList(c *ui.Context, s *AppState) {
	cards := s.GetAlbumCards()

	ui.Row(c).FillWidth().Height(24).Shrink(0).Children(func() {
		ui.Text(c, fmt.Sprintf("ALBUMS DIRECTORY (%d)", len(cards))).FontSize(12).Bold().TextColor(colorBlack)
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

						cardBg := colorWhite
						cardFg := colorBlack
						if isSelected {
							cardBg = colorBlack
							cardFg = colorWhite
						}

						displayName := card.Name
						if len([]rune(displayName)) > 22 {
							displayName = string([]rune(displayName)[:20]) + "..."
						}

						// Card Container
						cardBtn := ui.Button(c, "").
							Grow(1).
							Basis(0).
							MinWidth(0).
							Height(185).
							Border(1, colorBlack).
							Radius(8).
							Background(cardBg)

						cardBtn.Children(func() {
							ui.Column(c).Fill().Padding(8).Gap(6).Center().Children(func() {
								renderCoverBox(c, card.CoverData, 100, 6)
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
						ui.Row(c).FillWidth()
					}
				})
			}
		})
	})
}

// -------------------------------------------------------------
// All Tracks: Daftar Lagu
// -------------------------------------------------------------
func renderTracksList(c *ui.Context, s *AppState) {
	tracks := s.FilteredTracks()

	ui.Row(c).FillWidth().Height(24).Shrink(0).Center().Children(func() {
		ui.Text(c, fmt.Sprintf("TRACKS (%d items)", len(tracks))).FontSize(12).Bold().TextColor(colorBlack)
	})

	ui.Scroll(c).Fill().MinHeight(0).Children(func() {
		ui.Column(c).FillWidth().Gap(4).Children(func() {
			if len(tracks) == 0 {
				ui.Row(c).Fill().Center().Children(func() {
					ui.Text(c, "NO TRACKS FOUND").FontSize(12).Bold().TextColor(colorBlack)
				})
				return
			}

			for idx, t := range tracks {
				track := t
				isSelected := s.CurrentTrack() != nil && s.CurrentTrack().Path == track.Path
				rowNum := idx + 1

				rowBg := colorWhite
				textColor := colorBlack
				btnBorderColor := colorBlack
				if isSelected {
					rowBg = colorBlack
					textColor = colorWhite
					btnBorderColor = colorWhite
				}

				ui.Row(c).FillWidth().Height(36).Padding(4, 10).Border(1, colorBlack).Radius(0).Background(rowBg).Center().Children(func() {
					// 1. Nomor Indeks [01], [*] jika aktif
					prefix := fmt.Sprintf("[%02d]", rowNum)
					if isSelected {
						prefix = "[*]"
					}
					ui.Text(c, prefix).Width(38).FontSize(11).Bold().TextColor(textColor)

					// 2. Judul Lagu & Artis
					ui.Column(c).FillWidth().Children(func() {
						ui.Text(c, track.Title).FontSize(12).Bold().TextColor(textColor)
						ui.Text(c, track.Artist).FontSize(10).TextColor(textColor)
					})

					// 3. Nama Album
					ui.Text(c, track.Album).Width(180).FontSize(10).TextColor(textColor)

					// 4. Tombol Aksi Ikon Play/Pause
					if isSelected {
						playIconDraw := drawPause
						if !s.Engine.IsPlaying() {
							playIconDraw = drawPlay
						}

						btn := ui.Button(c, "").Width(36).Height(26).Border(1, btnBorderColor).Radius(0).Background(colorWhite)
						btn.Draw(func(p *ui.Painter, r ui.Rect) {
							playIconDraw(p, r, colorBlack)
						})
						if btn.Clicked() {
							s.TogglePlay()
						}
					} else {
						btn := customIconButton(c, 36, 26, false, drawPlay)
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
// Playback Bar di Bawah (Cover Art Terhubung)
// -------------------------------------------------------------
func renderBottomPlaybackBar(c *ui.Context, s *AppState) {
	if s.Engine.IsPlaying() {
		s.Progress = s.Engine.Progress()
	}

	current := s.CurrentTrack()

	ui.Row(c).FillWidth().Height(72).Shrink(0).Padding(8, 14).Border(1, colorBlack).Radius(0).Background(colorWhite).Center().Children(func() {
		// --- KIRI: Cover Art + Info Lagu ---
		ui.Row(c).Width(250).Center().Gap(10).Children(func() {
			var cover []byte
			if current != nil {
				cover = current.CoverData
			}
			renderCoverBox(c, cover, 42, 6)

			ui.Column(c).FillWidth().Children(func() {
				title := "IDLE"
				artist := "READY"
				if current != nil {
					title = current.Title
					artist = current.Artist
				}
				ui.Text(c, title).FontSize(12).Bold().TextColor(colorBlack)
				ui.Text(c, artist).FontSize(10).TextColor(colorBlack)
			})
		})

		// --- TENGAH: Kontrol Navigasi & Timeline Seekbar ---
		ui.Column(c).FillWidth().Center().Gap(4).Children(func() {
			ui.Row(c).Center().Gap(6).Children(func() {
				if customIconButton(c, 36, 26, s.Shuffle, drawShuffle).Clicked() {
					s.Shuffle = !s.Shuffle
				}
				if customIconButton(c, 36, 26, false, drawPrev).Clicked() {
					s.PrevTrack()
				}

				playDraw := drawPlay
				if s.Engine.IsPlaying() {
					playDraw = drawPause
				}
				if customIconButton(c, 44, 26, s.Engine.IsPlaying(), playDraw).Clicked() {
					s.TogglePlay()
				}

				if customIconButton(c, 36, 26, false, drawNext).Clicked() {
					s.NextTrack()
				}
				if customIconButton(c, 36, 26, s.Repeat, drawRepeat).Clicked() {
					s.Repeat = !s.Repeat
				}
			})

			ui.Row(c).FillWidth().Center().Padding(0, 16).Gap(8).Children(func() {
				ui.Text(c, s.CurrentTimeString()).FontSize(10).Bold().TextColor(colorBlack)

				progSlider := customSlider(c, &s.Progress, 0.0, 1.0).FillWidth()
				if progSlider.Changed() {
					s.Seek(s.Progress)
				}

				ui.Text(c, s.TotalTimeString()).FontSize(10).Bold().TextColor(colorBlack)
			})
		})

		// --- KANAN: Volume Slider ---
		ui.Row(c).Width(160).Right(0).Center().Gap(6).Children(func() {
			ui.Text(c, "VOL:").FontSize(10).Bold().TextColor(colorBlack)

			volSlider := customSlider(c, &s.Volume, 0.0, 1.0).Width(70)
			if volSlider.Changed() {
				s.Engine.SetVolume(s.Volume)
			}

			ui.Text(c, fmt.Sprintf("%d%%", int(s.Volume*100))).FontSize(10).TextColor(colorBlack)
		})
	})
}
