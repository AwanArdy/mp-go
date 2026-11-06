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

// Helper tombol tajam monokrom (0px radius)
func customButton(c *ui.Context, label string, active bool) ui.Element {
	btn := ui.Button(c, label).Border(1, colorBlack).Radius(0)
	if active {
		btn.Background(colorBlack).TextColor(colorWhite)
	} else {
		btn.Background(colorWhite).TextColor(colorBlack)
	}
	return btn
}

// Custom Slider Monokrom Murni (Tanpa Aksen Oranye)
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

		// 1. Background Track (Abu-abu Terang)
		p.Fill(ui.Rect{X: r.X, Y: trackY, W: r.W, H: trackH}, colorLightGray, 0)

		// 2. Filled Track (Hitam Pekat)
		p.Fill(ui.Rect{X: r.X, Y: trackY, W: r.W * ratio, H: trackH}, colorBlack, 0)

		// 3. Thumb/Knob Slider (Balok Hitam Kotak 0px Radius)
		thumbW := float32(8)
		thumbH := float32(14)
		thumbX := r.X + (r.W-thumbW)*ratio
		thumbY := r.Y + (r.H-thumbH)/2
		p.Fill(ui.Rect{X: thumbX, Y: thumbY, W: thumbW, H: thumbH}, colorBlack, 0)
	})

	return s
}

func (s *AppState) View(c *ui.Context) {
	// Container Induk Utama (Latar Putih, Garis Outline 2px Hitam)
	ui.Column(c).Fill().Background(colorWhite).Border(2, colorBlack).Radius(0).Children(func() {
		// 1. TOP BAR TERPADU (Satu Garis Lurus)
		renderTopBar(c, s)

		// 2. KONTEN TENGAH (Sidebar Kiri + Konten Library Kanan)
		ui.Row(c).Fill().Padding(10).Gap(10).Children(func() {
			renderSidebar(c, s)
			renderMainContent(c, s)
		})

		// 3. PLAYBACK BAR (Fokus Murni Kontrol & Timeline)
		renderBottomPlaybackBar(c, s)
	})
}

// -------------------------------------------------------------
// 1. Header Bar Terpadu
// -------------------------------------------------------------
func renderTopBar(c *ui.Context, s *AppState) {
	ui.Row(c).FillWidth().Height(46).Padding(6, 12).Border(1, colorBlack).Radius(0).Background(colorWhite).Center().Children(func() {
		// Kiri: Title
		ui.Text(c, "MUSIC PLAYER").FontSize(13).Bold().TextColor(colorBlack).Width(150)

		// Tengah: Search Bar Fleksibel Melebar
		ui.Row(c).FillWidth().Height(30).Border(1, colorBlack).Radius(0).Padding(2, 8).Background(colorWhite).Center().Children(func() {
			ui.Text(c, "FIND: ").FontSize(11).Bold().TextColor(colorBlack)
			ui.TextInput(c, &s.SearchQuery).
				Placeholder("Filter tracks, artists, albums...").
				FillWidth().
				Background(colorWhite).
				TextColor(colorBlack)
		})

		// Kanan: Tombol Window Controls Presisi [ — ] dan [ ✕ ]
		ui.Row(c).Gap(6).Children(func() {
			if customButton(c, " — ", false).Width(40).Height(28).Clicked() {
				s.MinimizeApp()
			}
			if customButton(c, " ✕ ", false).Width(40).Height(28).Clicked() {
				s.CloseApp()
			}
		})
	})
}

// -------------------------------------------------------------
// 2. Sidebar Kiri (Navigasi + Active Filter + Now Playing Card)
// -------------------------------------------------------------
func renderSidebar(c *ui.Context, s *AppState) {
	ui.Column(c).Width(220).FillHeight().Border(1, colorBlack).Radius(0).Padding(10).Gap(8).Children(func() {
		// Menu Navigasi Direktori
		ui.Text(c, "NAVIGATION").FontSize(11).Bold().TextColor(colorBlack)

		if customButton(c, "ALL TRACKS", s.ViewMode == "tracks" && s.SelectedArtist == "" && s.SelectedAlbum == "").FillWidth().Height(30).Clicked() {
			s.ClearFilter()
		}
		if customButton(c, "BY ARTISTS", s.ViewMode == "artists").FillWidth().Height(30).Clicked() {
			s.ViewMode = "artists"
			s.CurrentPage = 0
		}
		if customButton(c, "BY ALBUMS", s.ViewMode == "albums").FillWidth().Height(30).Clicked() {
			s.ViewMode = "albums"
			s.CurrentPage = 0
		}

		// Garis Pemisah
		ui.Row(c).FillWidth().Height(1).Border(1, colorBlack).Radius(0)

		// Blok Status Filter Aktif
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

		// Spacer pembatas ke bagian bawah sidebar
		ui.Column(c).Fill()

		// --- NOW PLAYING CARD DI BAGIAN BAWAH SIDEBAR ---
		renderSidebarNowPlaying(c, s)
	})
}

func renderSidebarNowPlaying(c *ui.Context, s *AppState) {
	current := s.CurrentTrack()

	ui.Column(c).FillWidth().Border(1, colorBlack).Radius(0).Padding(8).Gap(6).Background(colorWhite).Children(func() {
		ui.Row(c).FillWidth().Center().Children(func() {
			ui.Text(c, "NOW PLAYING").FontSize(10).Bold().TextColor(colorBlack)
			ui.Row(c).FillWidth().Right(0).Children(func() {
				likeLabel := "[LIKE]"
				if s.IsLiked {
					likeLabel = "[*LIKED*]"
				}
				if customButton(c, likeLabel, s.IsLiked).Height(20).Clicked() {
					s.IsLiked = !s.IsLiked
				}
			})
		})

		ui.Row(c).FillWidth().Gap(8).Center().Children(func() {
			// Artwork Placeholder Kotak Tajam
			ui.Row(c).Width(38).Height(38).Border(1, colorBlack).Radius(0).Center().Children(func() {
				ui.Text(c, "♫").FontSize(16).TextColor(colorBlack)
			})

			// Detail Judul & Artis
			ui.Column(c).FillWidth().Children(func() {
				title := "IDLE"
				artist := "READY"
				if current != nil {
					title = current.Title
					artist = current.Artist
				}
				ui.Text(c, title).FontSize(11).Bold().TextColor(colorBlack)
				ui.Text(c, artist).FontSize(9).TextColor(colorBlack)
			})
		})
	})
}

// -------------------------------------------------------------
// 3. Konten Library Kanan (Daftar Lagu / Artis / Album)
// -------------------------------------------------------------
func renderMainContent(c *ui.Context, s *AppState) {
	ui.Column(c).Fill().FillHeight().Border(1, colorBlack).Radius(0).Padding(10).Gap(8).Children(func() {
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

func renderArtistsList(c *ui.Context, s *AppState) {
	artists := s.GetArtists()

	ui.Row(c).FillWidth().Height(24).Children(func() {
		ui.Text(c, fmt.Sprintf("ARTISTS DIRECTORY (%d)", len(artists))).FontSize(12).Bold().TextColor(colorBlack)
	})

	ui.Column(c).Fill().Gap(4).Children(func() {
		for _, art := range artists {
			name := art
			isSelected := s.SelectedArtist == name
			if customButton(c, name, isSelected).FillWidth().Height(32).Clicked() {
				s.SelectedArtist = name
				s.SelectedAlbum = ""
				s.ViewMode = "tracks"
				s.CurrentPage = 0
			}
		}
	})
}

func renderAlbumsList(c *ui.Context, s *AppState) {
	albums := s.GetAlbums()

	ui.Row(c).FillWidth().Height(24).Children(func() {
		ui.Text(c, fmt.Sprintf("ALBUMS DIRECTORY (%d)", len(albums))).FontSize(12).Bold().TextColor(colorBlack)
	})

	ui.Column(c).Fill().Gap(4).Children(func() {
		for _, alb := range albums {
			name := alb
			isSelected := s.SelectedAlbum == name
			if customButton(c, name, isSelected).FillWidth().Height(32).Clicked() {
				s.SelectedAlbum = name
				s.SelectedArtist = ""
				s.ViewMode = "tracks"
				s.CurrentPage = 0
			}
		}
	})
}

func renderTracksList(c *ui.Context, s *AppState) {
	tracks, currentPg, totalPg := s.PaginatedTracks()
	totalTracks := len(s.FilteredTracks())

	// Subheader & Paginasi
	ui.Row(c).FillWidth().Height(28).Center().Children(func() {
		ui.Text(c, fmt.Sprintf("TRACKS (%d items)", totalTracks)).FontSize(12).Bold().TextColor(colorBlack)

		ui.Row(c).FillWidth().Right(0).Gap(6).Children(func() {
			if customButton(c, " < PREV ", false).Height(24).Clicked() {
				s.PrevPage()
			}
			pageIndicator := fmt.Sprintf("PAGE %d / %d", currentPg, totalPg)
			if totalPg == 0 {
				pageIndicator = "PAGE 0 / 0"
			}
			ui.Text(c, pageIndicator).FontSize(11).Bold().TextColor(colorBlack)
			if customButton(c, " NEXT > ", false).Height(24).Clicked() {
				s.NextPage()
			}
		})
	})

	// Daftar Lagu (Active track mendapat background full hitam dan teks putih)
	ui.Column(c).Fill().Gap(4).Children(func() {
		if len(tracks) == 0 {
			ui.Row(c).Fill().Center().Children(func() {
				ui.Text(c, "NO TRACKS FOUND").FontSize(12).Bold().TextColor(colorBlack)
			})
			return
		}

		for idx, t := range tracks {
			track := t
			isSelected := s.CurrentTrack() != nil && s.CurrentTrack().Path == track.Path
			rowNum := s.CurrentPage*s.PageSize + idx + 1

			// Penentuan Warna Kontras Baris
			rowBg := colorWhite
			textColor := colorBlack
			if isSelected {
				rowBg = colorBlack
				textColor = colorWhite
			}

			ui.Row(c).FillWidth().Height(36).Padding(4, 8).Border(1, colorBlack).Radius(0).Background(rowBg).Center().Children(func() {
				// 1. Index Kolom
				prefix := fmt.Sprintf("[%02d]", rowNum)
				if isSelected {
					prefix = "[▶]"
				}
				ui.Text(c, prefix).Width(38).FontSize(11).Bold().TextColor(textColor)

				// 2. Judul & Artis (Flexible Grow)
				ui.Column(c).FillWidth().Children(func() {
					ui.Text(c, track.Title).FontSize(12).Bold().TextColor(textColor)
					ui.Text(c, track.Artist).FontSize(10).TextColor(textColor)
				})

				// 3. Kolom Nama Album (Lebar Pasti dengan Jarak Aman)
				ui.Text(c, track.Album).Width(160).FontSize(10).TextColor(textColor)

				// 4. Tombol Aksi Play / Pause
				btnText := "PLAY"
				if isSelected {
					if s.Engine.IsPlaying() {
						btnText = "❚❚ PAUSE"
					} else {
						btnText = "▶ RESUME"
					}
					// Tombol putih kontras di atas baris hitam
					btn := ui.Button(c, btnText).Width(76).Height(26).Border(1, colorWhite).Radius(0).Background(colorWhite).TextColor(colorBlack)
					if btn.Clicked() {
						s.TogglePlay()
					}
				} else {
					if customButton(c, btnText, false).Width(76).Height(26).Clicked() {
						s.PlayTrack(track)
					}
				}
			})
		}
	})
}

// -------------------------------------------------------------
// 4. Playback Bar di Bawah (Fokus Murni Kontrol & Timeline)
// -------------------------------------------------------------
func renderBottomPlaybackBar(c *ui.Context, s *AppState) {
	// Sinkronisasi posisi progress secara berkala
	if s.Engine.IsPlaying() {
		s.Progress = s.Engine.Progress()
	}

	ui.Row(c).FillWidth().Height(72).Padding(8, 14).Border(1, colorBlack).Radius(0).Background(colorWhite).Center().Children(func() {
		// --- 1. SISI KIRI: Tombol Kontrol Navigasi Utama ---
		ui.Row(c).Center().Gap(6).Children(func() {
			if customButton(c, "[SHUF]", s.Shuffle).Width(52).Height(28).Clicked() {
				s.Shuffle = !s.Shuffle
			}
			if customButton(c, "|◀◀", false).Width(44).Height(28).Clicked() {
				s.PrevTrack()
			}

			playLabel := "▶ PLAY"
			if s.Engine.IsPlaying() {
				playLabel = "❚❚ PAUS"
			}
			if customButton(c, playLabel, s.Engine.IsPlaying()).Width(74).Height(28).Clicked() {
				s.TogglePlay()
			}

			if customButton(c, "▶▶|", false).Width(44).Height(28).Clicked() {
				s.NextTrack()
			}
			if customButton(c, "[LOOP]", s.Repeat).Width(52).Height(28).Clicked() {
				s.Repeat = !s.Repeat
			}
		})

		// --- 2. BAGIAN TENGAH: Timeline Slider dengan Custom Monochrome Look ---
		ui.Row(c).FillWidth().Center().Padding(0, 16).Gap(8).Children(func() {
			ui.Text(c, s.CurrentTimeString()).FontSize(11).Bold().TextColor(colorBlack)

			progSlider := customSlider(c, &s.Progress, 0.0, 1.0).FillWidth()
			if progSlider.Changed() {
				s.Seek(s.Progress)
			}

			ui.Text(c, s.TotalTimeString()).FontSize(11).Bold().TextColor(colorBlack)
		})

		// --- 3. SISI KANAN: Volume Slider Monokrom ---
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
