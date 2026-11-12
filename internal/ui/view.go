package ui

import (
	"fmt"
	"math"
	"music-player/internal/model"
	"time"

	"github.com/egoist/mygo/ui"
)

const appFont = "Plus Jakarta Sans"

// -------------------------------------------------------------
// Token Warna M3 (Material 3 Expressive Tokens)
// -------------------------------------------------------------
type M3Theme struct {
	Background         ui.Color
	Surface            ui.Color
	SurfaceContainer   ui.Color
	CardBg             ui.Color
	CardHover          ui.Color
	Outline            ui.Color
	OutlineVariant     ui.Color
	Primary            ui.Color
	PrimaryContainer   ui.Color
	OnPrimary          ui.Color
	OnPrimaryContainer ui.Color
	TextPrimary        ui.Color
	TextSecondary      ui.Color
	SliderTrack        ui.Color
	Danger             ui.Color
}

func getM3Theme(isDark bool) M3Theme {
	if isDark {
		return M3Theme{
			Background:         ui.RGB(15, 17, 23),
			Surface:            ui.RGB(23, 27, 36),
			SurfaceContainer:   ui.RGB(32, 38, 51),
			CardBg:             ui.RGB(28, 34, 46),
			CardHover:          ui.RGB(40, 48, 64),
			Outline:            ui.RGBA(148, 163, 184, 0.22),
			OutlineVariant:     ui.RGBA(148, 163, 184, 0.12),
			Primary:            ui.RGB(168, 199, 250),
			PrimaryContainer:   ui.RGB(37, 56, 88),
			OnPrimary:          ui.RGB(10, 25, 50),
			OnPrimaryContainer: ui.RGB(218, 234, 254),
			TextPrimary:        ui.RGB(241, 245, 249),
			TextSecondary:      ui.RGB(148, 163, 184),
			SliderTrack:        ui.RGB(47, 57, 74),
			Danger:             ui.RGB(248, 113, 113),
		}
	}
	return M3Theme{
		Background:         ui.RGB(246, 248, 252),
		Surface:            ui.RGB(255, 255, 255),
		SurfaceContainer:   ui.RGB(241, 245, 249),
		CardBg:             ui.RGB(255, 255, 255),
		CardHover:          ui.RGB(239, 246, 255),
		Outline:            ui.RGBA(30, 41, 59, 0.14),
		OutlineVariant:     ui.RGBA(30, 41, 59, 0.08),
		Primary:            ui.RGB(37, 99, 235),
		PrimaryContainer:   ui.RGB(219, 234, 254),
		OnPrimary:          ui.RGB(255, 255, 255),
		OnPrimaryContainer: ui.RGB(30, 58, 138),
		TextPrimary:        ui.RGB(15, 23, 42),
		TextSecondary:      ui.RGB(100, 116, 139),
		SliderTrack:        ui.RGB(226, 232, 240),
		Danger:             ui.RGB(220, 38, 38),
	}
}

func m3Text(c *ui.Context, text string) ui.Element {
	return ui.Text(c, text).Font(appFont)
}

func m3NavPill(c *ui.Context, label string, active bool, theme M3Theme) ui.Element {
	btn := ui.Button(c, label).Font(appFont).Radius(20)
	bg := theme.CardBg
	fg := theme.TextSecondary

	if active {
		bg = theme.PrimaryContainer
		fg = theme.OnPrimaryContainer
		btn.Border(1, theme.Primary)
	} else {
		btn.Border(1, theme.OutlineVariant)
		if btn.Hovered() {
			bg = theme.CardHover
			fg = theme.TextPrimary
		}
	}

	btn.Background(bg).TextColor(fg)
	return btn
}

func m3IconButton(c *ui.Context, size float32, active bool, theme M3Theme, drawIcon func(p *ui.Painter, r ui.Rect, fg ui.Color)) ui.Element {
	btn := ui.Button(c, "").Width(size).Height(size).Radius(size / 2)
	bg := theme.CardBg
	fg := theme.TextPrimary

	if active {
		bg = theme.Primary
		fg = theme.OnPrimary
		btn.Border(1, theme.Primary)
	} else {
		btn.Border(1, theme.OutlineVariant)
		if btn.Hovered() {
			bg = theme.CardHover
		}
	}

	btn.Background(bg)
	btn.Draw(func(p *ui.Painter, r ui.Rect) {
		if drawIcon != nil {
			drawIcon(p, r, fg)
		}
	})
	return btn
}

func m3Slider(c *ui.Context, val *float64, min, max float64, theme M3Theme) ui.Element {
	s := ui.SliderBase(c, val, min, max).Height(20).Radius(10)
	ratio := float32((*val - min) / (max - min))
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	s.Draw(func(p *ui.Painter, r ui.Rect) {
		trackH := float32(6)
		trackY := r.Y + (r.H-trackH)/2

		p.Fill(ui.Rect{X: r.X, Y: trackY, W: r.W, H: trackH}, theme.SliderTrack, trackH/2)
		p.Fill(ui.Rect{X: r.X, Y: trackY, W: r.W * ratio, H: trackH}, theme.Primary, trackH/2)

		thumbSize := float32(14)
		thumbX := r.X + (r.W-thumbSize)*ratio
		thumbY := r.Y + (r.H-thumbSize)/2
		p.Fill(ui.Rect{X: thumbX, Y: thumbY, W: thumbSize, H: thumbSize}, theme.Primary, thumbSize/2)
	})

	return s
}

func renderCoverBox(c *ui.Context, coverData []byte, size float32, radius float32, theme M3Theme) {
	box := ui.Row(c).Width(size).Height(size).Shrink(0).Background(theme.SurfaceContainer).Border(1, theme.OutlineVariant).Radius(radius).Clip().Center()

	if img, ok := GetOrDecodeImage(coverData); ok {
		box.Children(func() {
			ui.Image(c, img).Width(size).Height(size)
		})
		return
	}

	box.Draw(func(p *ui.Painter, r ui.Rect) {
		drawNoteIcon(p, r, theme.TextSecondary)
	})
}

// -------------------------------------------------------------
// Vektor Ikon
// -------------------------------------------------------------
func drawTrash(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2

	// Tutup tempat sampah
	p.Fill(ui.Rect{X: cx - 5, Y: cy - 6, W: 10, H: 1.5}, fg, 0.5)
	p.Fill(ui.Rect{X: cx - 2, Y: cy - 7.5, W: 4, H: 1.5}, fg, 0.5)

	// Bodi tempat sampah
	p.Fill(ui.Rect{X: cx - 4, Y: cy - 4, W: 8, H: 10}, fg, 1)

	// Garis vertikal di dalam bodi
	p.Fill(ui.Rect{X: cx - 2, Y: cy - 2, W: 1, H: 6}, ui.RGB(255, 255, 255), 0)
	p.Fill(ui.Rect{X: cx + 1, Y: cy - 2, W: 1, H: 6}, ui.RGB(255, 255, 255), 0)
}

func drawCheck(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	p.Fill(ui.Rect{X: cx - 5, Y: cy, W: 2, H: 4}, fg, 0.5)
	p.Fill(ui.Rect{X: cx - 3, Y: cy + 3, W: 3, H: 2}, fg, 0.5)
	p.Fill(ui.Rect{X: cx, Y: cy + 1, W: 2, H: 4}, fg, 0.5)
	p.Fill(ui.Rect{X: cx + 2, Y: cy - 3, W: 2, H: 5}, fg, 0.5)
	p.Fill(ui.Rect{X: cx + 4, Y: cy - 5, W: 2, H: 3}, fg, 0.5)
}

func drawCross(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	size := float32(5)
	for i := float32(0); i <= size; i++ {
		p.Fill(ui.Rect{X: cx - size/2 + i, Y: cy - size/2 + i, W: 1.5, H: 1.5}, fg, 0)
		p.Fill(ui.Rect{X: cx + size/2 - i, Y: cy - size/2 + i, W: 1.5, H: 1.5}, fg, 0)
	}
}

func drawPlus(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	p.Fill(ui.Rect{X: cx - 5, Y: cy - 1, W: 10, H: 2}, fg, 1)
	p.Fill(ui.Rect{X: cx - 1, Y: cy - 5, W: 2, H: 10}, fg, 1)
}

func drawPlaylistIcon(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	p.Fill(ui.Rect{X: cx - 6, Y: cy - 5, W: 12, H: 1.5}, fg, 0.75)
	p.Fill(ui.Rect{X: cx - 6, Y: cy - 1, W: 12, H: 1.5}, fg, 0.75)
	p.Fill(ui.Rect{X: cx - 6, Y: cy + 3, W: 7, H: 1.5}, fg, 0.75)
}

func drawMoon(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	radius := float32(5)
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			d1 := x*x + y*y
			dx2 := x - 2.5
			dy2 := y + 1.0
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
	p.Fill(ui.Rect{X: cx - 2.5, Y: cy - 2.5, W: 5, H: 5}, fg, 2.5)
	p.Fill(ui.Rect{X: cx - 1, Y: cy - 6, W: 2, H: 2.5}, fg, 1)
	p.Fill(ui.Rect{X: cx - 1, Y: cy + 3.5, W: 2, H: 2.5}, fg, 1)
	p.Fill(ui.Rect{X: cx - 6, Y: cy - 1, W: 2.5, H: 2}, fg, 1)
	p.Fill(ui.Rect{X: cx + 3.5, Y: cy - 1, W: 2.5, H: 2}, fg, 1)
}

func drawSearchIcon(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2 - 1
	cy := r.Y + r.H/2 - 1
	p.Fill(ui.Rect{X: cx - 4, Y: cy - 4, W: 8, H: 1.5}, fg, 0.75)
	p.Fill(ui.Rect{X: cx - 4, Y: cy + 3, W: 8, H: 1.5}, fg, 0.75)
	p.Fill(ui.Rect{X: cx - 4, Y: cy - 3, W: 1.5, H: 6}, fg, 0.75)
	p.Fill(ui.Rect{X: cx + 3, Y: cy - 3, W: 1.5, H: 6}, fg, 0.75)
	for i := float32(0); i < 4; i++ {
		p.Fill(ui.Rect{X: cx + 3 + i, Y: cy + 3 + i, W: 1.5, H: 1.5}, fg, 0.75)
	}
}

func drawMinimize(p *ui.Painter, r ui.Rect, fg ui.Color) {
	w := float32(10)
	h := float32(2)
	p.Fill(ui.Rect{X: r.X + (r.W-w)/2, Y: r.Y + (r.H-h)/2 + 2, W: w, H: h}, fg, 1)
}

func drawClose(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	size := float32(7)
	half := size / 2
	for i := float32(0); i <= size; i++ {
		p.Fill(ui.Rect{X: cx - half + i, Y: cy - half + i, W: 1.5, H: 1.5}, fg, 0.5)
		p.Fill(ui.Rect{X: cx + half - i, Y: cy - half + i, W: 1.5, H: 1.5}, fg, 0.5)
	}
}

func drawPrev(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	p.Fill(ui.Rect{X: cx - 5, Y: cy - 5, W: 2, H: 10}, fg, 1)
	for col := float32(0); col < 6; col++ {
		h := 10 - col*1.5
		p.Fill(ui.Rect{X: cx + 3 - col, Y: cy - h/2, W: 1.5, H: h}, fg, 0)
	}
}

func drawNext(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	for col := float32(0); col < 6; col++ {
		h := 10 - col*1.5
		p.Fill(ui.Rect{X: cx - 3 + col, Y: cy - h/2, W: 1.5, H: h}, fg, 0)
	}
	p.Fill(ui.Rect{X: cx + 4, Y: cy - 5, W: 2, H: 10}, fg, 1)
}

func drawPlay(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2 - 1
	cy := r.Y + r.H/2
	for col := float32(0); col < 9; col++ {
		h := 12 - col*1.3
		p.Fill(ui.Rect{X: cx - 4 + col, Y: cy - h/2, W: 1.5, H: h}, fg, 0)
	}
}

func drawPause(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	p.Fill(ui.Rect{X: cx - 4, Y: cy - 6, W: 2.5, H: 12}, fg, 1)
	p.Fill(ui.Rect{X: cx + 2, Y: cy - 6, W: 2.5, H: 12}, fg, 1)
}

func drawShuffle(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	p.Fill(ui.Rect{X: cx - 6, Y: cy - 4, W: 3, H: 1.5}, fg, 0)
	p.Fill(ui.Rect{X: cx - 6, Y: cy + 3, W: 3, H: 1.5}, fg, 0)
	for i := float32(0); i < 6; i++ {
		p.Fill(ui.Rect{X: cx - 3 + i, Y: cy - 3 + i, W: 1.2, H: 1.2}, fg, 0)
		p.Fill(ui.Rect{X: cx - 3 + i, Y: cy + 3 - i, W: 1.2, H: 1.2}, fg, 0)
	}
	p.Fill(ui.Rect{X: cx + 3, Y: cy + 3, W: 3, H: 1.5}, fg, 0)
	p.Fill(ui.Rect{X: cx + 3, Y: cy - 4, W: 3, H: 1.5}, fg, 0)
}

func drawRepeat(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	p.Fill(ui.Rect{X: cx - 5, Y: cy - 4, W: 10, H: 1.5}, fg, 0)
	p.Fill(ui.Rect{X: cx - 5, Y: cy + 3, W: 10, H: 1.5}, fg, 0)
	p.Fill(ui.Rect{X: cx - 5, Y: cy - 4, W: 1.5, H: 8}, fg, 0)
	p.Fill(ui.Rect{X: cx + 4, Y: cy - 4, W: 1.5, H: 8}, fg, 0)
}

func drawMic(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2 - 2
	p.Fill(ui.Rect{X: cx - 2.5, Y: cy - 6, W: 5, H: 8}, fg, 2.5)
	p.Fill(ui.Rect{X: cx - 4.5, Y: cy - 2, W: 1.5, H: 5}, fg, 0)
	p.Fill(ui.Rect{X: cx + 3, Y: cy - 2, W: 1.5, H: 5}, fg, 0)
	p.Fill(ui.Rect{X: cx - 4.5, Y: cy + 3, W: 9, H: 1.5}, fg, 0)
	p.Fill(ui.Rect{X: cx - 0.75, Y: cy + 4.5, W: 1.5, H: 3.5}, fg, 0)
	p.Fill(ui.Rect{X: cx - 3.5, Y: cy + 8, W: 7, H: 1.5}, fg, 0.75)
}

func drawNoteIcon(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	p.Fill(ui.Rect{X: cx - 3, Y: cy - 6, W: 2, H: 9}, fg, 0)
	p.Fill(ui.Rect{X: cx - 6, Y: cy + 1, W: 5, H: 4}, fg, 2)
	p.Fill(ui.Rect{X: cx + 3, Y: cy - 4, W: 2, H: 7}, fg, 0)
	p.Fill(ui.Rect{X: cx, Y: cy + 1, W: 5, H: 4}, fg, 2)
	p.Fill(ui.Rect{X: cx - 3, Y: cy - 6, W: 8, H: 2.5}, fg, 0)
}

// -------------------------------------------------------------
// Root View
// -------------------------------------------------------------
func (s *AppState) View(c *ui.Context) {
	theme := getM3Theme(s.DarkMode)

	if s.Engine.IsPlaying() {
		c.After(250 * time.Millisecond)
		s.Progress = s.Engine.Progress()
	}

	ui.Column(c).Fill().Background(theme.Background).Padding(14).Gap(12).Children(func() {
		renderTopBar(c, s, theme)

		ui.Row(c).Fill().MinHeight(0).Gap(14).Children(func() {
			renderSidebar(c, s, theme)
			renderMainContent(c, s, theme)
		})

		renderBottomPlaybackBar(c, s, theme)
	})
}

// -------------------------------------------------------------
// 1. Top Bar
// -------------------------------------------------------------
func renderTopBar(c *ui.Context, s *AppState, theme M3Theme) {
	ui.Row(c).FillWidth().Height(48).Shrink(0).Center().Gap(14).Children(func() {
		ui.Row(c).Center().Gap(10).Children(func() {
			m3Text(c, "Music").FontSize(18).Bold().TextColor(theme.TextPrimary)

			themeIcon := drawMoon
			if s.DarkMode {
				themeIcon = drawSun
			}
			if m3IconButton(c, 34, false, theme, themeIcon).Clicked() {
				s.ToggleDarkMode()
			}
		})

		ui.Row(c).FillWidth().Height(40).Radius(20).Background(theme.SurfaceContainer).Border(1, theme.OutlineVariant).Padding(0, 14).Center().Gap(10).Children(func() {
			iconBox := ui.Row(c).Width(20).Height(20).Center()
			iconBox.Draw(func(p *ui.Painter, r ui.Rect) {
				drawSearchIcon(p, r, theme.TextSecondary)
			})

			ui.TextInput(c, &s.SearchQuery).
				Font(appFont).
				Border(0, theme.SurfaceContainer).
				Radius(0).
				FillWidth().
				Height(36).
				Placeholder("Search songs, artists, albums...").
				PlaceholderColor(theme.TextSecondary).
				Background(theme.SurfaceContainer).
				TextColor(theme.TextPrimary)
		})

		ui.Row(c).Gap(6).Children(func() {
			if m3IconButton(c, 34, false, theme, drawMinimize).Clicked() {
				s.MinimizeApp()
			}
			if m3IconButton(c, 34, false, theme, drawClose).Clicked() {
				s.CloseApp()
			}
		})
	})
}

// -------------------------------------------------------------
// 2. Sidebar: Playlist Inline (Input langsung di list + Ikon Trash)
// -------------------------------------------------------------
func renderSidebar(c *ui.Context, s *AppState, theme M3Theme) {
	ui.Column(c).Width(230).FillHeight().MinHeight(0).Radius(18).Background(theme.Surface).Border(1, theme.OutlineVariant).Padding(14).Gap(10).Children(func() {
		m3Text(c, "LIBRARY").FontSize(11).Bold().TextColor(theme.TextSecondary)

		if m3NavPill(c, "All Tracks", s.ViewMode == "tracks" && s.SelectedArtist == "" && s.SelectedAlbum == "" && s.SelectedPlaylistID == "", theme).FillWidth().Height(38).Clicked() {
			s.ClearFilter()
		}
		if m3NavPill(c, "By Artists", s.ViewMode == "artists", theme).FillWidth().Height(38).Clicked() {
			s.ViewMode = "artists"
		}
		if m3NavPill(c, "By Albums", s.ViewMode == "albums", theme).FillWidth().Height(38).Clicked() {
			s.ViewMode = "albums"
		}

		ui.Row(c).FillWidth().Height(1).Background(theme.OutlineVariant)

		ui.Row(c).FillWidth().Height(22).Center().Children(func() {
			m3Text(c, "PLAYLISTS").FontSize(11).Bold().TextColor(theme.TextSecondary)
			ui.Row(c).FillWidth().Right(0).Children(func() {
				m3Text(c, fmt.Sprintf("(%d)", len(s.Playlists))).FontSize(11).TextColor(theme.TextSecondary)
			})
		})

		// Area Scroll Playlist
		ui.Scroll(c).Fill().MinHeight(0).Children(func() {
			ui.Column(c).FillWidth().Gap(4).Children(func() {
				if len(s.Playlists) == 0 && !s.IsAddingPlaylist {
					ui.Row(c).FillWidth().Padding(6, 4).Children(func() {
						m3Text(c, "No playlists yet").FontSize(12).TextColor(theme.TextSecondary)
					})
				}

				for _, pl := range s.Playlists {
					playlist := pl
					isPlActive := s.ViewMode == "playlist" && s.SelectedPlaylistID == playlist.ID
					isEditing := s.EditingPlaylistID == playlist.ID

					// 1. TAMPILAN JIKA SEDANG DI-RENAME SECARA INLINE
					if isEditing {
						ui.Row(c).FillWidth().Height(36).Radius(10).Background(theme.CardBg).Border(1, theme.Primary).Padding(2, 6).Center().Gap(4).Children(func() {
							ui.TextInput(c, &s.InlinePlaylistInput).
								Font(appFont).
								Border(0, theme.CardBg).
								Radius(6).
								FillWidth().
								Height(28).
								Background(theme.CardBg).
								TextColor(theme.TextPrimary)

							// Tombol Simpan Centang
							saveBtn := ui.Button(c, "").Width(22).Height(22).Radius(11).Background(theme.Primary).Border(0, theme.Primary)
							saveBtn.Draw(func(p *ui.Painter, r ui.Rect) {
								drawCheck(p, r, theme.OnPrimary)
							})
							if saveBtn.Clicked() {
								s.CommitRenamePlaylist()
							}

							// Tombol Batal Silang
							cancelBtn := ui.Button(c, "").Width(22).Height(22).Radius(11).Background(theme.SurfaceContainer).Border(0, theme.SurfaceContainer)
							cancelBtn.Draw(func(p *ui.Painter, r ui.Rect) {
								drawCross(p, r, theme.TextSecondary)
							})
							if cancelBtn.Clicked() {
								s.CancelInlinePlaylist()
							}
						})
						continue
					}

					// 2. TAMPILAN NORMAL ITEM PLAYLIST DENGAN IKON TRASH DI KANAN
					plRow := ui.Row(c).FillWidth().Height(36).Radius(18).Center().Padding(0, 6)
					rowBg := theme.Surface
					txtColor := theme.TextSecondary

					if isPlActive {
						rowBg = theme.PrimaryContainer
						txtColor = theme.OnPrimaryContainer
						plRow.Border(1, theme.Primary)
					} else {
						plRow.Border(1, theme.OutlineVariant)
						if plRow.Hovered() {
							rowBg = theme.CardHover
							txtColor = theme.TextPrimary
						}
					}
					plRow.Background(rowBg)

					// Klik Kanan pada item playlist langsung memicu rename inline
					if plRow.RightClicked() {
						s.SelectedPlaylistID = playlist.ID
						s.SelectedArtist = ""
						s.SelectedAlbum = ""
						s.ViewMode = "playlist"
						s.StartRenamePlaylist(playlist.ID, playlist.Name)
					}

					plRow.Children(func() {
						// Area klik nama playlist
						itemBtn := ui.Button(c, "").Grow(1).Basis(0).Height(34).Radius(16).Border(0, rowBg).Background(rowBg)
						itemBtn.Children(func() {
							ui.Row(c).Fill().Padding(0, 4).Center().Children(func() {
								iconRow := ui.Row(c).Width(16).Height(16).Center()
								iconRow.Draw(func(p *ui.Painter, r ui.Rect) {
									drawPlaylistIcon(p, r, txtColor)
								})

								ui.Row(c).Width(8)

								displayName := playlist.Name
								if len([]rune(displayName)) > 11 {
									displayName = string([]rune(displayName)[:9]) + "..."
								}
								m3Text(c, displayName).FontSize(12).TextColor(txtColor)
							})
						})

						if itemBtn.Clicked() {
							s.SelectedPlaylistID = playlist.ID
							s.SelectedArtist = ""
							s.SelectedAlbum = ""
							s.ViewMode = "playlist"
						}

						if itemBtn.RightClicked() {
							s.SelectedPlaylistID = playlist.ID
							s.SelectedArtist = ""
							s.SelectedAlbum = ""
							s.ViewMode = "playlist"
							s.StartRenamePlaylist(playlist.ID, playlist.Name)
						}

						// IKON TRASH DI SISI KANAN TOMBOL PLAYLIST
						trashBtn := ui.Button(c, "").Width(26).Height(26).Radius(13).Border(0, rowBg).Background(rowBg)
						trashFg := theme.TextSecondary
						if trashBtn.Hovered() {
							trashBtn.Background(ui.RGBA(239, 68, 68, 0.15))
							trashFg = theme.Danger
						}
						trashBtn.Draw(func(p *ui.Painter, r ui.Rect) {
							drawTrash(p, r, trashFg)
						})

						// Sekali klik langsung menghapus tanpa konfirmasi modal
						if trashBtn.Clicked() {
							s.DeletePlaylist(playlist.ID)
						}
					})
				}

				// 3. BARIS INPUT INLINE UNTUK MEMBUAT PLAYLIST BARU
				if s.IsAddingPlaylist {
					ui.Row(c).FillWidth().Height(36).Radius(10).Background(theme.CardBg).Border(1, theme.Primary).Padding(2, 6).Center().Gap(4).Children(func() {
						ui.TextInput(c, &s.InlinePlaylistInput).
							Font(appFont).
							Border(0, theme.CardBg).
							Radius(6).
							FillWidth().
							Height(28).
							Background(theme.CardBg).
							TextColor(theme.TextPrimary)

						saveBtn := ui.Button(c, "").Width(22).Height(22).Radius(11).Background(theme.Primary).Border(0, theme.Primary)
						saveBtn.Draw(func(p *ui.Painter, r ui.Rect) {
							drawCheck(p, r, theme.OnPrimary)
						})
						if saveBtn.Clicked() {
							s.CommitAddPlaylist()
						}

						cancelBtn := ui.Button(c, "").Width(22).Height(22).Radius(11).Background(theme.SurfaceContainer).Border(0, theme.SurfaceContainer)
						cancelBtn.Draw(func(p *ui.Painter, r ui.Rect) {
							drawCross(p, r, theme.TextSecondary)
						})
						if cancelBtn.Clicked() {
							s.CancelInlinePlaylist()
						}
					})
				}
			})
		})

		if s.SelectedArtist != "" || s.SelectedAlbum != "" || s.SearchQuery != "" {
			ui.Row(c).FillWidth().Height(1).Background(theme.OutlineVariant)
			m3Text(c, "ACTIVE FILTER").FontSize(10).Bold().TextColor(theme.TextSecondary)
			if m3NavPill(c, "Clear Filter", false, theme).FillWidth().Height(28).Clicked() {
				s.ClearFilter()
			}
		}

		ui.Row(c).FillWidth().Height(1).Background(theme.OutlineVariant)

		// Tombol "+ Add Playlist" di Bagian Bawah Sidebar
		addBtn := ui.Button(c, "").FillWidth().Height(40).Radius(20).Background(theme.Primary).Border(0, theme.Primary)
		if addBtn.Hovered() {
			addBtn.Background(theme.PrimaryContainer)
		}
		addBtn.Children(func() {
			ui.Row(c).Fill().Center().Gap(8).Children(func() {
				iconBox := ui.Row(c).Width(16).Height(16).Center()
				iconBox.Draw(func(p *ui.Painter, r ui.Rect) {
					drawPlus(p, r, theme.OnPrimary)
				})
				m3Text(c, "Add Playlist").FontSize(12).Bold().TextColor(theme.OnPrimary)
			})
		})

		if addBtn.Clicked() {
			s.StartAddPlaylist()
		}
	})
}

// -------------------------------------------------------------
// 3. Main Content
// -------------------------------------------------------------
func renderMainContent(c *ui.Context, s *AppState, theme M3Theme) {
	ui.Column(c).Fill().FillHeight().MinHeight(0).Radius(18).Background(theme.Surface).Border(1, theme.OutlineVariant).Padding(16).Gap(12).Children(func() {
		switch s.ViewMode {
		case "artists":
			renderArtistsList(c, s, theme)
		case "albums":
			renderAlbumsList(c, s, theme)
		case "lyrics":
			renderLyricsView(c, s, theme)
		case "playlist":
			renderPlaylistView(c, s, theme)
		default:
			renderTracksList(c, s, theme)
		}
	})
}

// -------------------------------------------------------------
// Detail Playlist View
// -------------------------------------------------------------
func renderPlaylistView(c *ui.Context, s *AppState, theme M3Theme) {
	pl := s.GetSelectedPlaylist()
	if pl == nil {
		ui.Column(c).Fill().Center().Children(func() {
			m3Text(c, "Playlist not found or deleted").FontSize(14).TextColor(theme.TextSecondary)
		})
		return
	}

	tracks := s.PlaylistTracks(pl)

	// Playlist Header Card
	ui.Row(c).FillWidth().Height(64).Shrink(0).Radius(14).Background(theme.SurfaceContainer).Border(1, theme.OutlineVariant).Padding(10, 16).Center().Children(func() {
		iconBox := ui.Row(c).Width(40).Height(40).Radius(10).Background(theme.PrimaryContainer).Center()
		iconBox.Draw(func(p *ui.Painter, r ui.Rect) {
			drawPlaylistIcon(p, r, theme.OnPrimaryContainer)
		})

		ui.Row(c).Width(12)

		ui.Column(c).FillWidth().Children(func() {
			m3Text(c, pl.Name).FontSize(15).Bold().TextColor(theme.TextPrimary)
			m3Text(c, fmt.Sprintf("%d songs", len(tracks))).FontSize(11).TextColor(theme.TextSecondary)
		})

		trashHeaderBtn := ui.Button(c, "").Width(34).Height(34).Radius(17).Background(theme.CardBg).Border(1, theme.OutlineVariant)
		if trashHeaderBtn.Hovered() {
			trashHeaderBtn.Background(ui.RGBA(239, 68, 68, 0.15))
		}
		trashHeaderBtn.Draw(func(p *ui.Painter, r ui.Rect) {
			drawTrash(p, r, theme.Danger)
		})
		if trashHeaderBtn.Clicked() {
			s.DeletePlaylist(pl.ID)
		}
	})

	renderTrackTable(c, s, tracks, theme)
}

// -------------------------------------------------------------
// By Artists: M3 Outlined Cards
// -------------------------------------------------------------
func renderArtistsList(c *ui.Context, s *AppState, theme M3Theme) {
	cards := s.GetArtistCards()

	ui.Row(c).FillWidth().Height(26).Shrink(0).Children(func() {
		m3Text(c, fmt.Sprintf("Artists (%d)", len(cards))).FontSize(14).Bold().TextColor(theme.TextPrimary)
	})

	ui.Scroll(c).Fill().MinHeight(0).Children(func() {
		ui.Column(c).FillWidth().Gap(14).Children(func() {
			total := len(cards)
			for i := 0; i < total; i += 4 {
				end := i + 4
				if end > total {
					end = total
				}
				rowCards := cards[i:end]

				ui.Row(c).FillWidth().Gap(14).Children(func() {
					for _, item := range rowCards {
						card := item
						isSelected := s.SelectedArtist == card.Name

						cardBtn := ui.Button(c, "").
							Grow(1).
							Basis(0).
							MinWidth(0).
							Height(185).
							Radius(16)

						bg := theme.CardBg
						if isSelected {
							bg = theme.PrimaryContainer
							cardBtn.Border(1, theme.Primary)
						} else {
							cardBtn.Border(1, theme.OutlineVariant)
							if cardBtn.Hovered() {
								bg = theme.CardHover
							}
						}
						cardBtn.Background(bg)

						displayName := card.Name
						if len([]rune(displayName)) > 18 {
							displayName = string([]rune(displayName)[:16]) + "..."
						}

						titleColor := theme.TextPrimary
						subColor := theme.TextSecondary
						if isSelected {
							titleColor = theme.OnPrimaryContainer
							subColor = theme.OnPrimaryContainer
						}

						cardBtn.Children(func() {
							ui.Column(c).Fill().Padding(10).Gap(8).Center().Children(func() {
								renderCoverBox(c, card.CoverData, 96, 12, theme)
								m3Text(c, displayName).FontSize(12).Bold().TextColor(titleColor)
								m3Text(c, card.Subtext).FontSize(10).TextColor(subColor)
							})
						})

						if cardBtn.Clicked() {
							s.SelectedArtist = card.Name
							s.SelectedAlbum = ""
							s.SelectedPlaylistID = ""
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

// -------------------------------------------------------------
// By Albums: M3 Outlined Cards
// -------------------------------------------------------------
func renderAlbumsList(c *ui.Context, s *AppState, theme M3Theme) {
	cards := s.GetAlbumCards()

	ui.Row(c).FillWidth().Height(26).Shrink(0).Children(func() {
		m3Text(c, fmt.Sprintf("Albums (%d)", len(cards))).FontSize(14).Bold().TextColor(theme.TextPrimary)
	})

	ui.Scroll(c).Fill().MinHeight(0).Children(func() {
		ui.Column(c).FillWidth().Gap(14).Children(func() {
			total := len(cards)
			for i := 0; i < total; i += 4 {
				end := i + 4
				if end > total {
					end = total
				}
				rowCards := cards[i:end]

				ui.Row(c).FillWidth().Gap(14).Children(func() {
					for _, item := range rowCards {
						card := item
						isSelected := s.SelectedAlbum == card.Name

						cardBtn := ui.Button(c, "").
							Grow(1).
							Basis(0).
							MinWidth(0).
							Height(185).
							Radius(16)

						bg := theme.CardBg
						if isSelected {
							bg = theme.PrimaryContainer
							cardBtn.Border(1, theme.Primary)
						} else {
							cardBtn.Border(1, theme.OutlineVariant)
							if cardBtn.Hovered() {
								bg = theme.CardHover
							}
						}
						cardBtn.Background(bg)

						displayName := card.Name
						if len([]rune(displayName)) > 18 {
							displayName = string([]rune(displayName)[:16]) + "..."
						}

						titleColor := theme.TextPrimary
						subColor := theme.TextSecondary
						if isSelected {
							titleColor = theme.OnPrimaryContainer
							subColor = theme.OnPrimaryContainer
						}

						cardBtn.Children(func() {
							ui.Column(c).Fill().Padding(10).Gap(8).Center().Children(func() {
								renderCoverBox(c, card.CoverData, 96, 12, theme)
								m3Text(c, displayName).FontSize(12).Bold().TextColor(titleColor)
								m3Text(c, card.Subtext).FontSize(10).TextColor(subColor)
							})
						})

						if cardBtn.Clicked() {
							s.SelectedAlbum = card.Name
							s.SelectedArtist = ""
							s.SelectedPlaylistID = ""
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

// -------------------------------------------------------------
// All Tracks M3
// -------------------------------------------------------------
func formatTrackDuration(s *AppState, t model.Track) string {
	if s.CurrentTrack() != nil && s.CurrentTrack().Path == t.Path {
		d := s.Engine.Duration()
		if d > 0 {
			sec := int(d.Seconds())
			return fmt.Sprintf("%d:%02d", sec/60, sec%60)
		}
	}
	if t.Duration > 0 {
		sec := int(t.Duration.Seconds())
		return fmt.Sprintf("%d:%02d", sec/60, sec%60)
	}
	return "--:--"
}

func renderTrackTable(c *ui.Context, s *AppState, tracks []model.Track, theme M3Theme) {
	// 1. HEADER TABEL (Track, Title, Artist, Album, Time)
	ui.Row(c).FillWidth().Height(32).Radius(8).Background(theme.SurfaceContainer).Border(1, theme.OutlineVariant).Padding(0, 12).Center().Children(func() {
		m3Text(c, "#").Width(44).FontSize(11).Bold().TextColor(theme.TextSecondary)
		m3Text(c, "Title").Grow(1).Basis(0).MinWidth(0).FontSize(11).Bold().TextColor(theme.TextSecondary)
		m3Text(c, "Artist").Width(190).FontSize(11).Bold().TextColor(theme.TextSecondary)
		m3Text(c, "Album").Width(210).FontSize(11).Bold().TextColor(theme.TextSecondary)
		m3Text(c, "Time").Width(55).Right(0).FontSize(11).Bold().TextColor(theme.TextSecondary)
	})

	// 2. DAFTAR BARIS TABEL
	ui.Scroll(c).Fill().MinHeight(0).Children(func() {
		ui.Column(c).FillWidth().Gap(2).Children(func() {
			if len(tracks) == 0 {
				ui.Row(c).FillWidth().Height(80).Center().Children(func() {
					m3Text(c, "No tracks found").FontSize(13).TextColor(theme.TextSecondary)
				})
				return
			}

			for idx, t := range tracks {
				track := t
				rowNum := idx + 1
				isPlaying := s.CurrentTrack() != nil && s.CurrentTrack().Path == track.Path
				isSelected := s.SelectedTrackPath == track.Path || (s.SelectedTrackPath == "" && isPlaying)

				rowBg := theme.Surface
				textColor := theme.TextPrimary
				subTextColor := theme.TextSecondary

				rowBtn := ui.Button(c, "").FillWidth().Height(34).Radius(6).Border(0, rowBg)

				if isSelected {
					rowBg = theme.PrimaryContainer
					textColor = theme.OnPrimaryContainer
					subTextColor = theme.OnPrimaryContainer
					rowBtn.Border(1, theme.Primary)
				} else if rowBtn.Hovered() {
					rowBg = theme.CardHover
				}
				rowBtn.Background(rowBg)

				rowBtn.Children(func() {
					ui.Row(c).Fill().Padding(0, 12).Center().Children(func() {
						// Kolom 1: No / Track
						prefix := fmt.Sprintf("%02d", rowNum)
						numColor := subTextColor
						if isPlaying {
							prefix = "▶ " + prefix
							numColor = theme.Primary
							if isSelected {
								numColor = theme.OnPrimaryContainer
							}
						}
						m3Text(c, prefix).Width(44).FontSize(11).Bold().TextColor(numColor)

						// Kolom 2: Title
						titleTxt := m3Text(c, track.Title).Grow(1).Basis(0).MinWidth(0).FontSize(12).TextColor(textColor)
						if isPlaying {
							titleTxt.Bold()
						}

						// Kolom 3: Artist
						m3Text(c, track.Artist).Width(190).FontSize(11).TextColor(subTextColor)

						// Kolom 4: Album
						m3Text(c, track.Album).Width(210).FontSize(11).TextColor(subTextColor)

						// Kolom 5: Time
						m3Text(c, formatTrackDuration(s, track)).Width(55).Right(0).FontSize(11).TextColor(subTextColor)
					})
				})

				// Double-Click: jeda klik < 350ms langsung memutar lagu
				if rowBtn.Clicked() {
					now := time.Now()
					if s.lastClickPath == track.Path && now.Sub(s.lastClickTime) < 350*time.Millisecond {
						s.PlayTrack(track)
						s.lastClickTime = time.Time{}
						s.lastClickPath = ""
					} else {
						s.SelectedTrackPath = track.Path
						s.lastClickTime = now
						s.lastClickPath = track.Path
					}
				}
			}
		})
	})
}

func renderTracksList(c *ui.Context, s *AppState, theme M3Theme) {
	tracks := s.FilteredTracks()

	ui.Row(c).FillWidth().Height(24).Shrink(0).Center().Children(func() {
		m3Text(c, fmt.Sprintf("Tracks (%d items)", len(tracks))).FontSize(13).Bold().TextColor(theme.TextPrimary)
	})

	renderTrackTable(c, s, tracks, theme)
}

// -------------------------------------------------------------
// Tampilan Lirik M3
// -------------------------------------------------------------
func renderLyricsView(c *ui.Context, s *AppState, theme M3Theme) {
	current := s.CurrentTrack()

	ui.Row(c).FillWidth().Height(26).Shrink(0).Center().Children(func() {
		title := "No track playing"
		if current != nil {
			title = fmt.Sprintf("%s — %s", current.Title, current.Artist)
		}
		m3Text(c, title).FontSize(12).Bold().TextColor(theme.TextPrimary)

		ui.Row(c).FillWidth().Right(0).Children(func() {
			if m3NavPill(c, "Back to tracks", false, theme).Height(26).Clicked() {
				s.ViewMode = "tracks"
			}
		})
	})

	if current == nil {
		ui.Column(c).Fill().Center().Children(func() {
			m3Text(c, "Putar lagu untuk melihat lirik").FontSize(13).TextColor(theme.TextSecondary)
		})
		return
	}

	if s.LyricsStatus == "loading" {
		ui.Column(c).Fill().Center().Gap(6).Children(func() {
			m3Text(c, "Mencari lirik tersinkronisasi...").FontSize(13).Bold().TextColor(theme.TextPrimary)
			m3Text(c, "YTMusic API > LRCLIB").FontSize(11).TextColor(theme.TextSecondary)
		})
		return
	}

	if s.LyricsStatus == "not_found" || len(s.LyricsLines) == 0 {
		ui.Column(c).Fill().Center().Gap(6).Children(func() {
			m3Text(c, "Lirik tidak ditemukan").FontSize(14).Bold().TextColor(theme.TextPrimary)
			m3Text(c, "Hanya lirik tersinkronisasi yang didukung").FontSize(11).TextColor(theme.TextSecondary)
		})
		return
	}

	currentPos := s.Engine.Position()
	activeIdx := 0
	for i, line := range s.LyricsLines {
		if line.Time <= currentPos {
			activeIdx = i
		} else {
			break
		}
	}

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
			if s.LyricVisualPos+step >= targetPos {
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

	const lineH float32 = 46
	const viewportH float32 = 7 * lineH

	baseIdx := int(math.Floor(s.LyricVisualPos))
	frac := float32(s.LyricVisualPos - float64(baseIdx))

	ui.Column(c).Fill().Center().Children(func() {
		ui.Column(c).FillWidth().Height(viewportH).Clip().Children(func() {
			for k := -3; k <= 4; k++ {
				lineIdx := baseIdx + k
				rowH := lineH
				if k == -3 {
					rowH = (1.0 - frac) * lineH
				}

				if lineIdx < 0 || lineIdx >= len(s.LyricsLines) {
					ui.Row(c).FillWidth().Height(rowH)
					continue
				}

				line := s.LyricsLines[lineIdx]
				displayText := line.Text
				if displayText == "" {
					displayText = "• • •"
				}

				dist := float64(lineIdx) - s.LyricVisualPos
				absDist := math.Abs(dist)

				fontSize := float32(14)
				textColor := theme.TextSecondary
				isBold := false

				if absDist < 0.6 {
					fontSize = 22
					textColor = theme.TextPrimary
					isBold = true
				} else if absDist < 1.6 {
					fontSize = 17
					textColor = theme.TextSecondary
				}

				row := ui.Row(c).FillWidth().Height(rowH).Clip().Center()
				if rowH > 14 {
					row.Children(func() {
						txt := m3Text(c, displayText).FontSize(fontSize).TextColor(textColor)
						if isBold {
							txt.Bold()
						}
					})
				}

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

// -------------------------------------------------------------
// 4. Playback Bar M3
// -------------------------------------------------------------
func renderBottomPlaybackBar(c *ui.Context, s *AppState, theme M3Theme) {
	current := s.CurrentTrack()

	ui.Row(c).FillWidth().Height(84).Shrink(0).Radius(24).Background(theme.Surface).Border(1, theme.OutlineVariant).Padding(10, 20).Center().Children(func() {
		ui.Row(c).Width(320).Center().Gap(14).Children(func() {
			var cover []byte
			if current != nil {
				cover = current.CoverData
			}
			renderCoverBox(c, cover, 54, 12, theme)

			ui.Column(c).FillWidth().Children(func() {
				title := "Ready"
				artist := "Select a song"
				if current != nil {
					title = current.Title
					artist = current.Artist
				}
				m3Text(c, title).FontSize(13).Bold().TextColor(theme.TextPrimary)
				m3Text(c, artist).FontSize(11).TextColor(theme.TextSecondary)
			})

			isLyricsMode := s.ViewMode == "lyrics"
			if m3IconButton(c, 36, isLyricsMode, theme, drawMic).Clicked() {
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

		ui.Column(c).FillWidth().Center().Gap(6).Children(func() {
			ui.Row(c).Center().Gap(12).Children(func() {
				if m3IconButton(c, 34, s.Shuffle, theme, drawShuffle).Clicked() {
					s.Shuffle = !s.Shuffle
				}
				if m3IconButton(c, 34, false, theme, drawPrev).Clicked() {
					s.PrevTrack()
				}

				playDraw := drawPlay
				if s.Engine.IsPlaying() {
					playDraw = drawPause
				}
				if m3IconButton(c, 46, true, theme, playDraw).Clicked() {
					s.TogglePlay()
				}

				if m3IconButton(c, 34, false, theme, drawNext).Clicked() {
					s.NextTrack()
				}
				if m3IconButton(c, 34, s.Repeat, theme, drawRepeat).Clicked() {
					s.Repeat = !s.Repeat
				}
			})

			ui.Row(c).FillWidth().Center().Padding(0, 16).Gap(10).Children(func() {
				m3Text(c, s.CurrentTimeString()).FontSize(11).Bold().TextColor(theme.TextSecondary)

				progSlider := m3Slider(c, &s.Progress, 0.0, 1.0, theme).FillWidth()
				if progSlider.Changed() {
					s.Seek(s.Progress)
				}

				m3Text(c, s.TotalTimeString()).FontSize(11).Bold().TextColor(theme.TextSecondary)
			})
		})

		ui.Row(c).Width(160).Right(0).Center().Gap(8).Children(func() {
			m3Text(c, "Vol").FontSize(11).Bold().TextColor(theme.TextSecondary)

			volSlider := m3Slider(c, &s.Volume, 0.0, 1.0, theme).Width(70)
			if volSlider.Changed() {
				s.Engine.SetVolume(s.Volume)
			}

			m3Text(c, fmt.Sprintf("%d%%", int(s.Volume*100))).FontSize(11).TextColor(theme.TextSecondary)
		})
	})
}
