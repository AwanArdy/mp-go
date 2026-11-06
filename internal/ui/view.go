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
		// 1. Header Atas (Tinggi tetap, tidak boleh menyusut)
		renderTopBar(c, s)

		// 2. Konten Tengah (MinHeight 0 agar scroll aktif dan tidak mendorong footer)
		ui.Row(c).Fill().MinHeight(0).Padding(10).Gap(10).Children(func() {
			renderSidebar(c, s)
			renderMainContent(c, s)
		})

		// 3. Playback Bar Bawah (Tinggi tetap, terkunci di bawah layar)
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

		// Spacer pengisi sisa ruang ke bawah
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
// Daftar Artis (Scrollable)
// -------------------------------------------------------------
func renderArtistsList(c *ui.Context, s *AppState) {
	artists := s.GetArtists()

	ui.Row(c).FillWidth().Height(24).Children(func() {
		ui.Text(c, fmt.Sprintf("ARTISTS DIRECTORY (%d)", len(artists))).FontSize(12).Bold().TextColor(colorBlack)
	})

	// ui.Scroll sebagai pembungkus scroll vertikal
	ui.Scroll(c).Fill().Children(func() {
		ui.Column(c).FillWidth().Gap(4).Children(func() {
			for _, art := range artists {
				name := art
				isSelected := s.SelectedArtist == name
				if customButton(c, name, isSelected).FillWidth().Height(32).Clicked() {
					s.SelectedArtist = name
					s.SelectedAlbum = ""
					s.ViewMode = "tracks"
				}
			}
		})
	})
}

// -------------------------------------------------------------
// Daftar Album (Scrollable)
// -------------------------------------------------------------
func renderAlbumsList(c *ui.Context, s *AppState) {
	albums := s.GetAlbums()

	ui.Row(c).FillWidth().Height(24).Children(func() {
		ui.Text(c, fmt.Sprintf("ALBUMS DIRECTORY (%d)", len(albums))).FontSize(12).Bold().TextColor(colorBlack)
	})

	// ui.Scroll sebagai pembungkus scroll vertikal
	ui.Scroll(c).Fill().Children(func() {
		ui.Column(c).FillWidth().Gap(4).Children(func() {
			for _, alb := range albums {
				name := alb
				isSelected := s.SelectedAlbum == name
				if customButton(c, name, isSelected).FillWidth().Height(32).Clicked() {
					s.SelectedAlbum = name
					s.SelectedArtist = ""
					s.ViewMode = "tracks"
				}
			}
		})
	})
}

// -------------------------------------------------------------
// Daftar Seluruh Lagu (Scrollable)
// -------------------------------------------------------------
func renderTracksList(c *ui.Context, s *AppState) {
	tracks := s.FilteredTracks()

	ui.Row(c).FillWidth().Height(24).Center().Children(func() {
		ui.Text(c, fmt.Sprintf("TRACKS (%d items)", len(tracks))).FontSize(12).Bold().TextColor(colorBlack)
	})

	// ui.Scroll sebagai pembungkus scroll vertikal daftar lagu
	ui.Scroll(c).Fill().Children(func() {
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

				ui.Row(c).FillWidth().Height(36).Padding(4, 8).Border(1, colorBlack).Radius(0).Background(rowBg).Center().Children(func() {
					prefix := fmt.Sprintf("[%02d]", rowNum)
					if isSelected {
						prefix = "[*]"
					}
					ui.Text(c, prefix).Width(38).FontSize(11).Bold().TextColor(textColor)

					ui.Column(c).FillWidth().Children(func() {
						ui.Text(c, track.Title).FontSize(12).Bold().TextColor(textColor)
						ui.Text(c, track.Artist).FontSize(10).TextColor(textColor)
					})

					ui.Text(c, track.Album).Width(160).FontSize(10).TextColor(textColor)

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

func renderBottomPlaybackBar(c *ui.Context, s *AppState) {
	if s.Engine.IsPlaying() {
		s.Progress = s.Engine.Progress()
	}

	current := s.CurrentTrack()

	ui.Row(c).FillWidth().Height(72).Shrink(0).Padding(8, 14).Border(1, colorBlack).Radius(0).Background(colorWhite).Center().Children(func() {
		// --- SISI KIRI: Info Lagu (Kotak Cover + Judul + Artis Tanpa Label & Tanpa Tombol Like) ---
		ui.Row(c).Width(240).Center().Gap(10).Children(func() {
			// Kotak Ikon Not Musik
			artBox := ui.Row(c).Width(40).Height(40).Border(1, colorBlack).Radius(0).Center()
			artBox.Draw(func(p *ui.Painter, r ui.Rect) {
				drawNoteIcon(p, r, colorBlack)
			})

			// Detail Judul Lagu & Artis
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

		// --- SISI TENGAH: Kontrol Navigasi & Timeline Seekbar ---
		ui.Column(c).FillWidth().Center().Gap(4).Children(func() {
			// Tombol-tombol Playback
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

			// Seekbar Timeline
			ui.Row(c).FillWidth().Center().Padding(0, 16).Gap(8).Children(func() {
				ui.Text(c, s.CurrentTimeString()).FontSize(10).Bold().TextColor(colorBlack)

				progSlider := customSlider(c, &s.Progress, 0.0, 1.0).FillWidth()
				if progSlider.Changed() {
					s.Seek(s.Progress)
				}

				ui.Text(c, s.TotalTimeString()).FontSize(10).Bold().TextColor(colorBlack)
			})
		})

		// --- SISI KANAN: Kontrol Volume Slider ---
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
