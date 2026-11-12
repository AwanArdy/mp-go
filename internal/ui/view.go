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
// Material Design 3 (Material You) Dark Scheme Tokens
// -------------------------------------------------------------
type M3Theme struct {
	Background              ui.Color
	Surface                 ui.Color
	SurfaceDim              ui.Color
	SurfaceBright           ui.Color
	SurfaceContainerLowest  ui.Color
	SurfaceContainerLow     ui.Color
	SurfaceContainer        ui.Color
	SurfaceContainerHigh    ui.Color
	SurfaceContainerHighest ui.Color
	CardBg                  ui.Color
	CardHover               ui.Color
	Outline                 ui.Color
	OutlineVariant          ui.Color
	Primary                 ui.Color
	PrimaryContainer        ui.Color
	OnPrimary               ui.Color
	OnPrimaryContainer      ui.Color
	Secondary               ui.Color
	SecondaryContainer      ui.Color
	OnSecondaryContainer    ui.Color
	Tertiary                ui.Color
	TertiaryContainer       ui.Color
	OnTertiaryContainer     ui.Color
	TextPrimary             ui.Color
	TextSecondary           ui.Color
	SliderTrack             ui.Color
	Danger                  ui.Color
	DangerContainer         ui.Color
	OnDanger                ui.Color
}

func getM3Theme(isDark bool) M3Theme {
	// Full Material 3 Dark Scheme (Deep Midnight Slate + Luminous Accent)
	return M3Theme{
		Background:              ui.RGB(15, 17, 22),       // #0F1116
		Surface:                 ui.RGB(18, 20, 26),       // #12141A
		SurfaceDim:              ui.RGB(15, 17, 22),       // #0F1116
		SurfaceBright:           ui.RGB(56, 59, 68),       // #383B44
		SurfaceContainerLowest:  ui.RGB(11, 13, 17),       // #0B0D11
		SurfaceContainerLow:     ui.RGB(21, 24, 31),       // #15181F
		SurfaceContainer:        ui.RGB(27, 30, 39),       // #1B1E27
		SurfaceContainerHigh:    ui.RGB(37, 41, 52),       // #252934
		SurfaceContainerHighest: ui.RGB(48, 53, 67),       // #303543
		CardBg:                  ui.RGB(27, 30, 39),
		CardHover:               ui.RGB(41, 46, 59),
		Outline:                 ui.RGBA(142, 145, 155, 0.35),
		OutlineVariant:          ui.RGBA(142, 145, 155, 0.14),
		Primary:                 ui.RGB(168, 199, 250),    // #A8C7FA M3 Cyan-Blue
		PrimaryContainer:        ui.RGB(18, 62, 102),      // #123E66 Deep Tonal Blue
		OnPrimary:               ui.RGB(6, 46, 111),       // #062E6F Contrast deep blue
		OnPrimaryContainer:      ui.RGB(210, 228, 255),    // #D2E4FF Bright tint
		Secondary:               ui.RGB(188, 199, 219),
		SecondaryContainer:      ui.RGB(60, 72, 88),
		OnSecondaryContainer:    ui.RGB(218, 226, 249),
		Tertiary:                ui.RGB(223, 189, 226),
		TertiaryContainer:       ui.RGB(85, 60, 90),
		OnTertiaryContainer:     ui.RGB(250, 216, 253),
		TextPrimary:             ui.RGB(230, 233, 239),    // #E6E9EF OnSurface
		TextSecondary:           ui.RGB(156, 163, 175),    // #9CA3AF OnSurfaceVariant
		SliderTrack:             ui.RGB(48, 53, 67),
		Danger:                  ui.RGB(255, 180, 171),    // #FFB4AB Error
		DangerContainer:         ui.RGB(147, 0, 10),       // #93000A ErrorContainer
		OnDanger:                ui.RGB(105, 0, 5),
	}
}

func m3Text(c *ui.Context, text string) ui.Element {
	return ui.Text(c, text).Font(appFont)
}

func m3NavPill(c *ui.Context, label string, active bool, theme M3Theme) ui.Element {
	btn := ui.Button(c, label).Font(appFont).Radius(18)
	bg := theme.SurfaceContainerHigh
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

// m3NavDrawerItem: Standard Material 3 Navigation Drawer Item (Stadium Pill)
func m3NavDrawerItem(c *ui.Context, label string, active bool, theme M3Theme, drawIcon func(p *ui.Painter, r ui.Rect, fg ui.Color)) ui.Element {
	btn := ui.Button(c, "").FillWidth().Height(42).Radius(21)
	bg := ui.RGBA(0, 0, 0, 0)
	fg := theme.TextSecondary

	if active {
		bg = theme.PrimaryContainer
		fg = theme.OnPrimaryContainer
	} else if btn.Hovered() {
		bg = theme.CardHover
		fg = theme.TextPrimary
	}

	btn.Background(bg)
	btn.Children(func() {
		ui.Row(c).Fill().Padding(0, 14).Center().Gap(12).Children(func() {
			if drawIcon != nil {
				iconBox := ui.Row(c).Width(20).Height(20).Center()
				iconBox.Draw(func(p *ui.Painter, r ui.Rect) {
					drawIcon(p, r, fg)
				})
			}
			txt := m3Text(c, label).FontSize(13).TextColor(fg)
			if active {
				txt.Bold()
			}
		})
	})
	return btn
}

func m3IconButton(c *ui.Context, size float32, active bool, theme M3Theme, drawIcon func(p *ui.Painter, r ui.Rect, fg ui.Color)) ui.Element {
	btn := ui.Button(c, "").Width(size).Height(size).Radius(size / 2)
	bg := theme.SurfaceContainerHigh
	fg := theme.TextSecondary

	if active {
		bg = theme.Primary
		fg = theme.OnPrimary
	} else if btn.Hovered() {
		bg = theme.CardHover
		fg = theme.TextPrimary
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
	s := ui.SliderBase(c, val, min, max).Height(22).Radius(11)
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

		// Inactive Track (M3 SurfaceContainerHighest)
		p.Fill(ui.Rect{X: r.X, Y: trackY, W: r.W, H: trackH}, theme.SliderTrack, trackH/2)
		// Active Track (M3 Primary)
		if r.W*ratio > 0 {
			p.Fill(ui.Rect{X: r.X, Y: trackY, W: r.W * ratio, H: trackH}, theme.Primary, trackH/2)
		}

		// Thumb: Material 3 handle
		thumbSize := float32(14)
		if s.Hovered() {
			thumbSize = 16
		}
		thumbX := r.X + (r.W-thumbSize)*ratio
		thumbY := r.Y + (r.H-thumbSize)/2
		p.Fill(ui.Rect{X: thumbX, Y: thumbY, W: thumbSize, H: thumbSize}, theme.Primary, thumbSize/2)
	})

	return s
}

func renderCoverBox(c *ui.Context, coverData []byte, size float32, radius float32, theme M3Theme) {
	renderCoverBoxWithFallback(c, coverData, size, radius, theme, drawNoteIcon)
}

func renderCoverBoxWithFallback(c *ui.Context, coverData []byte, size float32, radius float32, theme M3Theme, fallback func(p *ui.Painter, r ui.Rect, fg ui.Color)) {
	box := ui.Row(c).Width(size).Height(size).Shrink(0).Background(theme.SurfaceContainerHighest).Border(1, theme.OutlineVariant).Radius(radius).Clip().Center()

	if img, ok := GetOrDecodeImage(coverData); ok {
		box.Children(func() {
			ui.Image(c, img).Width(size).Height(size)
		})
		return
	}

	box.Draw(func(p *ui.Painter, r ui.Rect) {
		if fallback != nil {
			fallback(p, r, theme.TextSecondary)
		} else {
			drawNoteIcon(p, r, theme.TextSecondary)
		}
	})
}

// -------------------------------------------------------------
// Vector Icons (Material 3 Symbols)
// -------------------------------------------------------------
func drawSearchIcon(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	var circle ui.Path
	circle.Circle(cx-1.5, cy-1.5, 4.5)
	p.StrokePath(&circle, 1.8, fg)

	var handle ui.Path
	handle.Cap(ui.RoundCap)
	handle.MoveTo(cx+1.8, cy+1.8).LineTo(cx+6, cy+6)
	p.StrokePath(&handle, 1.8, fg)
}

func drawMinimize(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2 + 2
	var path ui.Path
	path.Cap(ui.RoundCap)
	path.MoveTo(cx-5, cy).LineTo(cx+5, cy)
	p.StrokePath(&path, 1.8, fg)
}

func drawClose(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	var path ui.Path
	path.Cap(ui.RoundCap)
	path.MoveTo(cx-4.5, cy-4.5).LineTo(cx+4.5, cy+4.5)
	path.MoveTo(cx+4.5, cy-4.5).LineTo(cx-4.5, cy+4.5)
	p.StrokePath(&path, 1.8, fg)
}

func drawPlay(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2 - 1
	cy := r.Y + r.H/2
	var path ui.Path
	path.Cap(ui.RoundCap).Join(ui.RoundJoin)
	path.MoveTo(cx-4, cy-6).LineTo(cx+6, cy).LineTo(cx-4, cy+6).Close()
	p.FillPath(&path, fg)
}

func drawPause(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	p.Fill(ui.Rect{X: cx - 4.5, Y: cy - 6, W: 3, H: 12}, fg, 1.5)
	p.Fill(ui.Rect{X: cx + 1.5, Y: cy - 6, W: 3, H: 12}, fg, 1.5)
}

func drawPrev(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	p.Fill(ui.Rect{X: cx - 6, Y: cy - 5.5, W: 2, H: 11}, fg, 1)
	var path ui.Path
	path.Cap(ui.RoundCap).Join(ui.RoundJoin)
	path.MoveTo(cx+5, cy-5.5).LineTo(cx-3.5, cy).LineTo(cx+5, cy+5.5).Close()
	p.FillPath(&path, fg)
}

func drawNext(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	var path ui.Path
	path.Cap(ui.RoundCap).Join(ui.RoundJoin)
	path.MoveTo(cx-5, cy-5.5).LineTo(cx+3.5, cy).LineTo(cx-5, cy+5.5).Close()
	p.FillPath(&path, fg)
	p.Fill(ui.Rect{X: cx + 4, Y: cy - 5.5, W: 2, H: 11}, fg, 1)
}

func drawShuffle(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	var path1, path2 ui.Path
	path1.Cap(ui.RoundCap).Join(ui.RoundJoin)
	path1.MoveTo(cx-6, cy-4).LineTo(cx-2, cy-4).LineTo(cx+2, cy+4).LineTo(cx+6, cy+4)
	p.StrokePath(&path1, 1.6, fg)

	path2.Cap(ui.RoundCap).Join(ui.RoundJoin)
	path2.MoveTo(cx-6, cy+4).LineTo(cx-2, cy+4).LineTo(cx+2, cy-4).LineTo(cx+6, cy-4)
	p.StrokePath(&path2, 1.6, fg)

	var ah1, ah2 ui.Path
	ah1.Cap(ui.RoundCap).Join(ui.RoundJoin)
	ah1.MoveTo(cx+4, cy-6).LineTo(cx+6.5, cy-4).LineTo(cx+4, cy-2)
	p.StrokePath(&ah1, 1.6, fg)

	ah2.Cap(ui.RoundCap).Join(ui.RoundJoin)
	ah2.MoveTo(cx+4, cy+2).LineTo(cx+6.5, cy+4).LineTo(cx+4, cy+6)
	p.StrokePath(&ah2, 1.6, fg)
}

func drawRepeat(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	var path ui.Path
	path.Cap(ui.RoundCap).Join(ui.RoundJoin)
	path.MoveTo(cx+2, cy-5).LineTo(cx-4, cy-5).QuadTo(cx-6, cy-5, cx-6, cy-3).LineTo(cx-6, cy+1)
	path.MoveTo(cx-2, cy+5).LineTo(cx+4, cy+5).QuadTo(cx+6, cy+5, cx+6, cy+3).LineTo(cx+6, cy-1)
	p.StrokePath(&path, 1.6, fg)

	var ah1, ah2 ui.Path
	ah1.Cap(ui.RoundCap).Join(ui.RoundJoin)
	ah1.MoveTo(cx, cy-7).LineTo(cx+3, cy-5).LineTo(cx, cy-3)
	p.StrokePath(&ah1, 1.6, fg)

	ah2.Cap(ui.RoundCap).Join(ui.RoundJoin)
	ah2.MoveTo(cx, cy+3).LineTo(cx-3, cy+5).LineTo(cx, cy+7)
	p.StrokePath(&ah2, 1.6, fg)
}

func drawMic(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2 - 1
	p.Fill(ui.Rect{X: cx - 2.5, Y: cy - 6, W: 5, H: 8}, fg, 2.5)

	var arc ui.Path
	arc.Cap(ui.RoundCap)
	arc.MoveTo(cx-4.5, cy-2).QuadTo(cx-4.5, cy+3, cx, cy+3).QuadTo(cx+4.5, cy+3, cx+4.5, cy-2)
	p.StrokePath(&arc, 1.6, fg)

	var stem ui.Path
	stem.Cap(ui.RoundCap)
	stem.MoveTo(cx, cy+3).LineTo(cx, cy+6.5)
	stem.MoveTo(cx-3.5, cy+6.5).LineTo(cx+3.5, cy+6.5)
	p.StrokePath(&stem, 1.6, fg)
}

func drawNoteIcon(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	var note1 ui.Path
	note1.Circle(cx-3, cy+3, 2.5)
	p.FillPath(&note1, fg)

	var note2 ui.Path
	note2.Circle(cx+3.5, cy+1, 2.5)
	p.FillPath(&note2, fg)

	var stems ui.Path
	stems.Cap(ui.RoundCap).Join(ui.RoundJoin)
	stems.MoveTo(cx-1, cy+3).LineTo(cx-1, cy-5).LineTo(cx+5.5, cy-7).LineTo(cx+5.5, cy+1)
	p.StrokePath(&stems, 1.8, fg)
}

func drawVolume(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	var speaker ui.Path
	speaker.Cap(ui.RoundCap).Join(ui.RoundJoin)
	speaker.MoveTo(cx-5, cy-2.5).LineTo(cx-3, cy-2.5).LineTo(cx, cy-5.5).LineTo(cx, cy+5.5).LineTo(cx-3, cy+2.5).LineTo(cx-5, cy+2.5).Close()
	p.FillPath(&speaker, fg)

	var wave ui.Path
	wave.Cap(ui.RoundCap)
	wave.MoveTo(cx+3, cy-3).QuadTo(cx+5.5, cy, cx+3, cy+3)
	p.StrokePath(&wave, 1.5, fg)
}

func drawLibraryIcon(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	var lines ui.Path
	lines.Cap(ui.RoundCap)
	lines.MoveTo(cx-6, cy-4).LineTo(cx+6, cy-4)
	lines.MoveTo(cx-6, cy).LineTo(cx+2, cy)
	lines.MoveTo(cx-6, cy+4).LineTo(cx, cy+4)
	p.StrokePath(&lines, 1.6, fg)

	var note ui.Path
	note.Circle(cx+4, cy+3.5, 1.8)
	p.FillPath(&note, fg)

	var stem ui.Path
	stem.Cap(ui.RoundCap)
	stem.MoveTo(cx+5.5, cy+3.5).LineTo(cx+5.5, cy-1)
	p.StrokePath(&stem, 1.4, fg)
}

func drawArtistIcon(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	var head ui.Path
	head.Circle(cx, cy-3.5, 3)
	p.FillPath(&head, fg)

	var body ui.Path
	body.Cap(ui.RoundCap)
	body.MoveTo(cx-5.5, cy+5.5).QuadTo(cx-5.5, cy+1.5, cx, cy+1.5).QuadTo(cx+5.5, cy+1.5, cx+5.5, cy+5.5)
	p.StrokePath(&body, 1.6, fg)
}

func drawAlbumIcon(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	var outer ui.Path
	outer.Circle(cx, cy, 6)
	p.StrokePath(&outer, 1.5, fg)

	var inner ui.Path
	inner.Circle(cx, cy, 2)
	p.StrokePath(&inner, 1.2, fg)
}

func drawPlaylistIcon(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	var lines ui.Path
	lines.Cap(ui.RoundCap)
	lines.MoveTo(cx-6, cy-4).LineTo(cx+6, cy-4)
	lines.MoveTo(cx-6, cy).LineTo(cx+6, cy)
	lines.MoveTo(cx-6, cy+4).LineTo(cx+1, cy+4)
	p.StrokePath(&lines, 1.6, fg)
}

func drawTrash(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	var lid ui.Path
	lid.Cap(ui.RoundCap)
	lid.MoveTo(cx-5.5, cy-4).LineTo(cx+5.5, cy-4)
	lid.MoveTo(cx-2, cy-4).LineTo(cx-2, cy-6).LineTo(cx+2, cy-6).LineTo(cx+2, cy-4)
	p.StrokePath(&lid, 1.5, fg)

	var bin ui.Path
	bin.Cap(ui.RoundCap).Join(ui.RoundJoin)
	bin.MoveTo(cx-4, cy-2.5).LineTo(cx-3.5, cy+5).QuadTo(cx-3.5, cy+6, cx-2, cy+6).LineTo(cx+2, cy+6).QuadTo(cx+3.5, cy+6, cx+3.5, cy+5).LineTo(cx+4, cy-2.5)
	p.StrokePath(&bin, 1.5, fg)

	var ribs ui.Path
	ribs.Cap(ui.RoundCap)
	ribs.MoveTo(cx-1.5, cy-0.5).LineTo(cx-1.5, cy+4)
	ribs.MoveTo(cx+1.5, cy-0.5).LineTo(cx+1.5, cy+4)
	p.StrokePath(&ribs, 1.2, fg)
}

func drawPlus(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	var path ui.Path
	path.Cap(ui.RoundCap)
	path.MoveTo(cx-5, cy).LineTo(cx+5, cy)
	path.MoveTo(cx, cy-5).LineTo(cx, cy+5)
	p.StrokePath(&path, 2.0, fg)
}

func drawCheck(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	var path ui.Path
	path.Cap(ui.RoundCap).Join(ui.RoundJoin)
	path.MoveTo(cx-4.5, cy).LineTo(cx-1.5, cy+3.5).LineTo(cx+5, cy-3.5)
	p.StrokePath(&path, 2.0, fg)
}

func drawCross(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	var path ui.Path
	path.Cap(ui.RoundCap)
	path.MoveTo(cx-3.5, cy-3.5).LineTo(cx+3.5, cy+3.5)
	path.MoveTo(cx+3.5, cy-3.5).LineTo(cx-3.5, cy+3.5)
	p.StrokePath(&path, 1.8, fg)
}

func drawEqualizer(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	t := float64(time.Now().UnixMilli()) / 180.0
	h1 := float32(4 + 5*math.Abs(math.Sin(t)))
	h2 := float32(4 + 7*math.Abs(math.Sin(t+1.1)))
	h3 := float32(4 + 6*math.Abs(math.Sin(t+2.2)))

	p.Fill(ui.Rect{X: cx - 4.5, Y: cy - h1/2, W: 2, H: h1}, fg, 1)
	p.Fill(ui.Rect{X: cx - 1, Y: cy - h2/2, W: 2, H: h2}, fg, 1)
	p.Fill(ui.Rect{X: cx + 2.5, Y: cy - h3/2, W: 2, H: h3}, fg, 1)
}

func drawArtistAvatarFallback(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	var head ui.Path
	head.Circle(cx, cy-12, 14)
	p.FillPath(&head, fg)

	var body ui.Path
	body.Cap(ui.RoundCap)
	body.MoveTo(cx-26, cy+26).QuadTo(cx-26, cy+6, cx, cy+6).QuadTo(cx+26, cy+6, cx+26, cy+26)
	p.StrokePath(&body, 6, fg)
}

func drawAlbumArtFallback(p *ui.Painter, r ui.Rect, fg ui.Color) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2
	var outer ui.Path
	outer.Circle(cx, cy, 32)
	p.StrokePath(&outer, 4, fg)

	var ring ui.Path
	ring.Circle(cx, cy, 22)
	p.StrokePath(&ring, 2, fg)

	var center ui.Path
	center.Circle(cx, cy, 10)
	p.FillPath(&center, fg)
}

// -------------------------------------------------------------
// Root View
// -------------------------------------------------------------
func (s *AppState) View(c *ui.Context) {
	// Full Material 3 Dark Scheme
	theme := getM3Theme(true)

	if s.Engine.IsPlaying() {
		c.After(120 * time.Millisecond)
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
// 1. Top Bar: Full M3 Dark with Search & Window Controls
// (Dark mode button removed as requested)
// -------------------------------------------------------------
func renderTopBar(c *ui.Context, s *AppState, theme M3Theme) {
	ui.Row(c).FillWidth().Height(48).Shrink(0).Center().Gap(14).Children(func() {
		// App Brand Badge
		ui.Row(c).Center().Gap(10).Children(func() {
			logoBadge := ui.Row(c).Width(36).Height(36).Radius(18).Background(theme.SurfaceContainerHigh).Border(1, theme.OutlineVariant).Center()
			logoBadge.Draw(func(p *ui.Painter, r ui.Rect) {
				drawNoteIcon(p, r, theme.Primary)
			})

			m3Text(c, "Music").FontSize(20).Bold().TextColor(theme.TextPrimary)

			// M3 Dark Pill Indicator
			ui.Row(c).Height(20).Radius(10).Background(theme.SurfaceContainerHighest).Padding(0, 8).Center().Children(func() {
				m3Text(c, "M3").FontSize(10).Bold().TextColor(theme.Primary)
			})
		})

		// M3 Floating Search Bar (Stadium shape)
		ui.Row(c).FillWidth().Height(44).Radius(22).Background(theme.SurfaceContainerHigh).Border(1, theme.OutlineVariant).Padding(0, 16).Center().Gap(10).Children(func() {
			iconBox := ui.Row(c).Width(20).Height(20).Center()
			iconBox.Draw(func(p *ui.Painter, r ui.Rect) {
				drawSearchIcon(p, r, theme.TextSecondary)
			})

			ui.TextInput(c, &s.SearchQuery).
				Font(appFont).
				Border(0, theme.SurfaceContainerHigh).
				Radius(0).
				FillWidth().
				Height(36).
				Placeholder("Search songs, artists, albums...").
				PlaceholderColor(theme.TextSecondary).
				Background(theme.SurfaceContainerHigh).
				TextColor(theme.TextPrimary)

			if s.SearchQuery != "" {
				clearBtn := ui.Button(c, "").Width(22).Height(22).Radius(11).Background(theme.SurfaceContainerHighest).Border(0, theme.SurfaceContainerHighest)
				clearBtn.Draw(func(p *ui.Painter, r ui.Rect) {
					drawCross(p, r, theme.TextSecondary)
				})
				if clearBtn.Clicked() {
					s.SearchQuery = ""
				}
			}
		})

		// Window Controls (M3 Circular Icon Buttons)
		ui.Row(c).Gap(8).Center().Children(func() {
			minBtn := ui.Button(c, "").Width(34).Height(34).Radius(17).Background(theme.SurfaceContainerHigh).Border(1, theme.OutlineVariant)
			if minBtn.Hovered() {
				minBtn.Background(theme.CardHover)
			}
			minBtn.Draw(func(p *ui.Painter, r ui.Rect) {
				drawMinimize(p, r, theme.TextPrimary)
			})
			if minBtn.Clicked() {
				s.MinimizeApp()
			}

			closeBtn := ui.Button(c, "").Width(34).Height(34).Radius(17).Background(theme.SurfaceContainerHigh).Border(1, theme.OutlineVariant)
			closeFg := theme.TextPrimary
			if closeBtn.Hovered() {
				closeBtn.Background(theme.DangerContainer)
				closeFg = theme.Danger
			}
			closeBtn.Draw(func(p *ui.Painter, r ui.Rect) {
				drawClose(p, r, closeFg)
			})
			if closeBtn.Clicked() {
				s.CloseApp()
			}
		})
	})
}

// -------------------------------------------------------------
// 2. Sidebar: Material 3 Navigation Drawer
// -------------------------------------------------------------
func renderSidebar(c *ui.Context, s *AppState, theme M3Theme) {
	ui.Column(c).Width(236).FillHeight().MinHeight(0).Radius(24).Background(theme.SurfaceContainerLow).Border(1, theme.OutlineVariant).Padding(16, 12).Gap(10).Children(func() {
		// Section: Library
		ui.Row(c).FillWidth().Padding(0, 10).Children(func() {
			m3Text(c, "LIBRARY").FontSize(11).Bold().TextColor(theme.TextSecondary)
		})

		isAllActive := s.ViewMode == "tracks" && s.SelectedArtist == "" && s.SelectedAlbum == "" && s.SelectedPlaylistID == ""
		if m3NavDrawerItem(c, "All Tracks", isAllActive, theme, drawLibraryIcon).Clicked() {
			s.ClearFilter()
		}

		if m3NavDrawerItem(c, "Artists", s.ViewMode == "artists", theme, drawArtistIcon).Clicked() {
			s.ViewMode = "artists"
		}

		if m3NavDrawerItem(c, "Albums", s.ViewMode == "albums", theme, drawAlbumIcon).Clicked() {
			s.ViewMode = "albums"
		}

		// Divider
		ui.Row(c).FillWidth().Height(1).Background(theme.OutlineVariant)

		// Section: Playlists
		ui.Row(c).FillWidth().Height(26).Padding(0, 10).Center().Children(func() {
			m3Text(c, "PLAYLISTS").FontSize(11).Bold().TextColor(theme.TextSecondary)
			ui.Row(c).FillWidth().Right(0).Children(func() {
				ui.Row(c).Height(18).Radius(9).Background(theme.SurfaceContainerHigh).Padding(0, 7).Center().Children(func() {
					m3Text(c, fmt.Sprintf("%d", len(s.Playlists))).FontSize(10).Bold().TextColor(theme.TextSecondary)
				})
			})
		})

		// Playlists List Scroll Area
		ui.Scroll(c).Fill().MinHeight(0).Children(func() {
			ui.Column(c).FillWidth().Gap(4).Children(func() {
				if len(s.Playlists) == 0 && !s.IsAddingPlaylist {
					ui.Row(c).FillWidth().Padding(8, 10).Children(func() {
						m3Text(c, "No playlists yet").FontSize(12).TextColor(theme.TextSecondary)
					})
				}

				for _, pl := range s.Playlists {
					playlist := pl
					isPlActive := s.ViewMode == "playlist" && s.SelectedPlaylistID == playlist.ID
					isEditing := s.EditingPlaylistID == playlist.ID

					// 1. TAMPILAN JIKA SEDANG DI-RENAME SECARA INLINE
					if isEditing {
						ui.Row(c).FillWidth().Height(38).Radius(12).Background(theme.SurfaceContainerHigh).Border(1.5, theme.Primary).Padding(2, 6).Center().Gap(4).Children(func() {
							ui.TextInput(c, &s.InlinePlaylistInput).
								Font(appFont).
								Border(0, theme.SurfaceContainerHigh).
								Radius(6).
								FillWidth().
								Height(30).
								Background(theme.SurfaceContainerHigh).
								TextColor(theme.TextPrimary)

							saveBtn := ui.Button(c, "").Width(24).Height(24).Radius(12).Background(theme.Primary).Border(0, theme.Primary)
							saveBtn.Draw(func(p *ui.Painter, r ui.Rect) {
								drawCheck(p, r, theme.OnPrimary)
							})
							if saveBtn.Clicked() {
								s.CommitRenamePlaylist()
							}

							cancelBtn := ui.Button(c, "").Width(24).Height(24).Radius(12).Background(theme.SurfaceContainerHighest).Border(0, theme.SurfaceContainerHighest)
							cancelBtn.Draw(func(p *ui.Painter, r ui.Rect) {
								drawCross(p, r, theme.TextSecondary)
							})
							if cancelBtn.Clicked() {
								s.CancelInlinePlaylist()
							}
						})
						continue
					}

					// 2. TAMPILAN NORMAL ITEM PLAYLIST
					plRow := ui.Row(c).FillWidth().Height(38).Radius(19).Center().Padding(0, 8)
					rowBg := ui.RGBA(0, 0, 0, 0)
					txtColor := theme.TextSecondary

					if isPlActive {
						rowBg = theme.PrimaryContainer
						txtColor = theme.OnPrimaryContainer
					} else if plRow.Hovered() {
						rowBg = theme.CardHover
						txtColor = theme.TextPrimary
					}
					plRow.Background(rowBg)

					if plRow.RightClicked() {
						s.SelectedPlaylistID = playlist.ID
						s.SelectedArtist = ""
						s.SelectedAlbum = ""
						s.ViewMode = "playlist"
						s.StartRenamePlaylist(playlist.ID, playlist.Name)
					}

					plRow.Children(func() {
						itemBtn := ui.Button(c, "").Grow(1).Basis(0).Height(34).Radius(17).Border(0, rowBg).Background(rowBg)
						itemBtn.Children(func() {
							ui.Row(c).Fill().Padding(0, 6).Center().Children(func() {
								iconRow := ui.Row(c).Width(16).Height(16).Center()
								iconRow.Draw(func(p *ui.Painter, r ui.Rect) {
									drawPlaylistIcon(p, r, txtColor)
								})

								ui.Row(c).Width(8)

								displayName := playlist.Name
								if len([]rune(displayName)) > 12 {
									displayName = string([]rune(displayName)[:10]) + "..."
								}
								txt := m3Text(c, displayName).FontSize(12).TextColor(txtColor)
								if isPlActive {
									txt.Bold()
								}
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

						trashBtn := ui.Button(c, "").Width(26).Height(26).Radius(13).Border(0, rowBg).Background(rowBg)
						trashFg := theme.TextSecondary
						if trashBtn.Hovered() {
							trashBtn.Background(theme.DangerContainer)
							trashFg = theme.Danger
						}
						trashBtn.Draw(func(p *ui.Painter, r ui.Rect) {
							drawTrash(p, r, trashFg)
						})

						if trashBtn.Clicked() {
							s.DeletePlaylist(playlist.ID)
						}
					})
				}

				// 3. BARIS INPUT INLINE UNTUK MEMBUAT PLAYLIST BARU
				if s.IsAddingPlaylist {
					ui.Row(c).FillWidth().Height(38).Radius(12).Background(theme.SurfaceContainerHigh).Border(1.5, theme.Primary).Padding(2, 6).Center().Gap(4).Children(func() {
						ui.TextInput(c, &s.InlinePlaylistInput).
							Font(appFont).
							Border(0, theme.SurfaceContainerHigh).
							Radius(6).
							FillWidth().
							Height(30).
							Background(theme.SurfaceContainerHigh).
							TextColor(theme.TextPrimary)

						saveBtn := ui.Button(c, "").Width(24).Height(24).Radius(12).Background(theme.Primary).Border(0, theme.Primary)
						saveBtn.Draw(func(p *ui.Painter, r ui.Rect) {
							drawCheck(p, r, theme.OnPrimary)
						})
						if saveBtn.Clicked() {
							s.CommitAddPlaylist()
						}

						cancelBtn := ui.Button(c, "").Width(24).Height(24).Radius(12).Background(theme.SurfaceContainerHighest).Border(0, theme.SurfaceContainerHighest)
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

		// Active Filter Pill in Sidebar if any
		if s.SelectedArtist != "" || s.SelectedAlbum != "" || s.SearchQuery != "" {
			ui.Row(c).FillWidth().Height(1).Background(theme.OutlineVariant)
			ui.Row(c).FillWidth().Padding(0, 6).Center().Children(func() {
				m3Text(c, "ACTIVE FILTER").FontSize(10).Bold().TextColor(theme.TextSecondary)
			})
			clearChip := ui.Button(c, "").FillWidth().Height(32).Radius(16).Background(theme.SurfaceContainerHigh).Border(1, theme.OutlineVariant)
			if clearChip.Hovered() {
				clearChip.Background(theme.CardHover)
			}
			clearChip.Children(func() {
				ui.Row(c).Fill().Padding(0, 10).Center().Gap(6).Children(func() {
					iconBox := ui.Row(c).Width(14).Height(14).Center()
					iconBox.Draw(func(p *ui.Painter, r ui.Rect) {
						drawCross(p, r, theme.Primary)
					})
					m3Text(c, "Clear Filter").FontSize(11).Bold().TextColor(theme.Primary)
				})
			})
			if clearChip.Clicked() {
				s.ClearFilter()
			}
		}

		ui.Row(c).FillWidth().Height(1).Background(theme.OutlineVariant)

		// M3 Extended FAB ("+ New Playlist")
		addBtn := ui.Button(c, "").FillWidth().Height(46).Radius(23).Background(theme.Primary).Border(0, theme.Primary)
		if addBtn.Hovered() {
			addBtn.Background(ui.RGB(190, 215, 255))
		}
		addBtn.Children(func() {
			ui.Row(c).Fill().Center().Gap(8).Children(func() {
				iconBox := ui.Row(c).Width(18).Height(18).Center()
				iconBox.Draw(func(p *ui.Painter, r ui.Rect) {
					drawPlus(p, r, theme.OnPrimary)
				})
				m3Text(c, "New Playlist").FontSize(13).Bold().TextColor(theme.OnPrimary)
			})
		})

		if addBtn.Clicked() {
			s.StartAddPlaylist()
		}
	})
}

// -------------------------------------------------------------
// 3. Main Content Panel
// -------------------------------------------------------------
func renderMainContent(c *ui.Context, s *AppState, theme M3Theme) {
	ui.Column(c).Fill().FillHeight().MinHeight(0).Radius(24).Background(theme.SurfaceContainerLow).Border(1, theme.OutlineVariant).Padding(18).Gap(14).Children(func() {
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

	// Material 3 Header Card
	ui.Row(c).FillWidth().Height(74).Shrink(0).Radius(18).Background(theme.SurfaceContainer).Border(1, theme.OutlineVariant).Padding(12, 18).Center().Children(func() {
		iconBox := ui.Row(c).Width(48).Height(48).Radius(12).Background(theme.PrimaryContainer).Center()
		iconBox.Draw(func(p *ui.Painter, r ui.Rect) {
			drawPlaylistIcon(p, r, theme.OnPrimaryContainer)
		})

		ui.Row(c).Width(14)

		ui.Column(c).FillWidth().Children(func() {
			m3Text(c, pl.Name).FontSize(17).Bold().TextColor(theme.TextPrimary)
			m3Text(c, fmt.Sprintf("%d tracks in playlist", len(tracks))).FontSize(12).TextColor(theme.TextSecondary)
		})

		trashHeaderBtn := ui.Button(c, "").Width(36).Height(36).Radius(18).Background(theme.SurfaceContainerHigh).Border(1, theme.OutlineVariant)
		trashFg := theme.Danger
		if trashHeaderBtn.Hovered() {
			trashHeaderBtn.Background(theme.DangerContainer)
		}
		trashHeaderBtn.Draw(func(p *ui.Painter, r ui.Rect) {
			drawTrash(p, r, trashFg)
		})
		if trashHeaderBtn.Clicked() {
			s.DeletePlaylist(pl.ID)
		}
	})

	renderTrackTable(c, s, tracks, theme)
}

// -------------------------------------------------------------
// By Artists: M3 Outlined Cards with Circular Avatars
// -------------------------------------------------------------
func renderArtistsList(c *ui.Context, s *AppState, theme M3Theme) {
	cards := s.GetArtistCards()

	ui.Row(c).FillWidth().Height(30).Shrink(0).Center().Children(func() {
		m3Text(c, "Artists").FontSize(18).Bold().TextColor(theme.TextPrimary)
		ui.Row(c).Width(10)
		ui.Row(c).Height(20).Radius(10).Background(theme.SurfaceContainerHigh).Padding(0, 8).Center().Children(func() {
			m3Text(c, fmt.Sprintf("%d", len(cards))).FontSize(11).Bold().TextColor(theme.TextSecondary)
		})
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
							Height(215).
							Radius(20)

						bg := theme.CardBg
						if isSelected {
							bg = theme.PrimaryContainer
							cardBtn.Border(1.5, theme.Primary)
						} else {
							cardBtn.Border(1, theme.OutlineVariant)
							if cardBtn.Hovered() {
								bg = theme.CardHover
							}
						}
						cardBtn.Background(bg)

						displayName := card.Name
						if len([]rune(displayName)) > 16 {
							displayName = string([]rune(displayName)[:14]) + "..."
						}

						titleColor := theme.TextPrimary
						subColor := theme.TextSecondary
						if isSelected {
							titleColor = theme.OnPrimaryContainer
							subColor = theme.OnPrimaryContainer
						}

						cardBtn.Children(func() {
							ui.Column(c).Fill().Padding(12).Gap(10).Center().Children(func() {
								// M3 Circular Artist Avatar
								renderCoverBoxWithFallback(c, card.CoverData, 110, 55, theme, drawArtistAvatarFallback)
								m3Text(c, displayName).FontSize(13).Bold().TextColor(titleColor)
								m3Text(c, card.Subtext).FontSize(11).TextColor(subColor)
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
// By Albums: M3 Outlined Cards with Rounded Squircles
// -------------------------------------------------------------
func renderAlbumsList(c *ui.Context, s *AppState, theme M3Theme) {
	cards := s.GetAlbumCards()

	ui.Row(c).FillWidth().Height(30).Shrink(0).Center().Children(func() {
		m3Text(c, "Albums").FontSize(18).Bold().TextColor(theme.TextPrimary)
		ui.Row(c).Width(10)
		ui.Row(c).Height(20).Radius(10).Background(theme.SurfaceContainerHigh).Padding(0, 8).Center().Children(func() {
			m3Text(c, fmt.Sprintf("%d", len(cards))).FontSize(11).Bold().TextColor(theme.TextSecondary)
		})
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
							Height(215).
							Radius(20)

						bg := theme.CardBg
						if isSelected {
							bg = theme.PrimaryContainer
							cardBtn.Border(1.5, theme.Primary)
						} else {
							cardBtn.Border(1, theme.OutlineVariant)
							if cardBtn.Hovered() {
								bg = theme.CardHover
							}
						}
						cardBtn.Background(bg)

						displayName := card.Name
						if len([]rune(displayName)) > 16 {
							displayName = string([]rune(displayName)[:14]) + "..."
						}

						titleColor := theme.TextPrimary
						subColor := theme.TextSecondary
						if isSelected {
							titleColor = theme.OnPrimaryContainer
							subColor = theme.OnPrimaryContainer
						}

						cardBtn.Children(func() {
							ui.Column(c).Fill().Padding(12).Gap(10).Center().Children(func() {
								// M3 Squircle Album Artwork
								renderCoverBoxWithFallback(c, card.CoverData, 110, 16, theme, drawAlbumArtFallback)
								m3Text(c, displayName).FontSize(13).Bold().TextColor(titleColor)
								m3Text(c, card.Subtext).FontSize(11).TextColor(subColor)
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
// All Tracks M3 Table
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
	// 1. M3 Table Header
	ui.Row(c).FillWidth().Height(34).Radius(10).Background(theme.SurfaceContainerHigh).Border(1, theme.OutlineVariant).Padding(0, 14).Center().Children(func() {
		m3Text(c, "#").Width(48).FontSize(11).Bold().TextColor(theme.TextSecondary)
		m3Text(c, "Title").Grow(1).Basis(0).MinWidth(0).FontSize(11).Bold().TextColor(theme.TextSecondary)
		m3Text(c, "Artist").Width(200).FontSize(11).Bold().TextColor(theme.TextSecondary)
		m3Text(c, "Album").Width(220).FontSize(11).Bold().TextColor(theme.TextSecondary)
		m3Text(c, "Duration").Width(60).Right(0).FontSize(11).Bold().TextColor(theme.TextSecondary)
	})

	// 2. M3 Track Rows List
	ui.Scroll(c).Fill().MinHeight(0).Children(func() {
		ui.Column(c).FillWidth().Gap(3).Children(func() {
			if len(tracks) == 0 {
				ui.Row(c).FillWidth().Height(100).Center().Children(func() {
					m3Text(c, "No tracks found").FontSize(14).TextColor(theme.TextSecondary)
				})
				return
			}

			for idx, t := range tracks {
				track := t
				rowNum := idx + 1
				isPlaying := s.CurrentTrack() != nil && s.CurrentTrack().Path == track.Path
				isSelected := s.SelectedTrackPath == track.Path || (s.SelectedTrackPath == "" && isPlaying)

				rowBg := ui.RGBA(0, 0, 0, 0)
				textColor := theme.TextPrimary
				subTextColor := theme.TextSecondary

				rowBtn := ui.Button(c, "").FillWidth().Height(40).Radius(10).Border(0, rowBg)

				if isSelected {
					rowBg = theme.PrimaryContainer
					textColor = theme.OnPrimaryContainer
					subTextColor = theme.OnPrimaryContainer
					rowBtn.Border(1, theme.Primary)
				} else if isPlaying {
					rowBg = theme.SurfaceContainerHighest
					textColor = theme.Primary
					subTextColor = theme.TextPrimary
				} else if rowBtn.Hovered() {
					rowBg = theme.CardHover
				}
				rowBtn.Background(rowBg)

				rowBtn.Children(func() {
					ui.Row(c).Fill().Padding(0, 14).Center().Children(func() {
						// Column 1: Equalizer or Number
						if isPlaying {
							eqBox := ui.Row(c).Width(48).Height(20).Center()
							eqColor := theme.Primary
							if isSelected {
								eqColor = theme.OnPrimaryContainer
							}
							eqBox.Draw(func(p *ui.Painter, r ui.Rect) {
								drawEqualizer(p, r, eqColor)
							})
						} else {
							prefix := fmt.Sprintf("%02d", rowNum)
							numColor := subTextColor
							m3Text(c, prefix).Width(48).FontSize(11).Bold().TextColor(numColor)
						}

						// Column 2: Title
						titleTxt := m3Text(c, track.Title).Grow(1).Basis(0).MinWidth(0).FontSize(13).TextColor(textColor)
						if isPlaying || isSelected {
							titleTxt.Bold()
						}

						// Column 3: Artist
						m3Text(c, track.Artist).Width(200).FontSize(12).TextColor(subTextColor)

						// Column 4: Album
						m3Text(c, track.Album).Width(220).FontSize(12).TextColor(subTextColor)

						// Column 5: Duration
						m3Text(c, formatTrackDuration(s, track)).Width(60).Right(0).FontSize(11).Bold().TextColor(subTextColor)
					})
				})

				// Double-Click plays song (< 350ms)
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

	ui.Row(c).FillWidth().Height(32).Shrink(0).Center().Children(func() {
		m3Text(c, "Tracks").FontSize(18).Bold().TextColor(theme.TextPrimary)
		ui.Row(c).Width(10)
		ui.Row(c).Height(20).Radius(10).Background(theme.SurfaceContainerHigh).Padding(0, 8).Center().Children(func() {
			m3Text(c, fmt.Sprintf("%d songs", len(tracks))).FontSize(11).Bold().TextColor(theme.TextSecondary)
		})

		// M3 Active Filter Chip if artist/album/search active
		if s.SelectedArtist != "" || s.SelectedAlbum != "" || s.SearchQuery != "" {
			ui.Row(c).Width(14)
			filterLabel := ""
			if s.SelectedArtist != "" {
				filterLabel = "Artist: " + s.SelectedArtist
			} else if s.SelectedAlbum != "" {
				filterLabel = "Album: " + s.SelectedAlbum
			} else if s.SearchQuery != "" {
				filterLabel = "Search: \"" + s.SearchQuery + "\""
			}

			filterChip := ui.Button(c, "").Height(26).Radius(13).Background(theme.SurfaceContainerHighest).Border(1, theme.Primary)
			filterChip.Children(func() {
				ui.Row(c).Fill().Padding(0, 10).Center().Gap(6).Children(func() {
					m3Text(c, filterLabel).FontSize(11).Bold().TextColor(theme.Primary)
					iconBox := ui.Row(c).Width(12).Height(12).Center()
					iconBox.Draw(func(p *ui.Painter, r ui.Rect) {
						drawCross(p, r, theme.Primary)
					})
				})
			})
			if filterChip.Clicked() {
				s.ClearFilter()
			}
		}
	})

	renderTrackTable(c, s, tracks, theme)
}

// -------------------------------------------------------------
// Lyrics View (Material 3 Synchronized Lyrics)
// -------------------------------------------------------------
func renderLyricsView(c *ui.Context, s *AppState, theme M3Theme) {
	current := s.CurrentTrack()

	ui.Row(c).FillWidth().Height(36).Shrink(0).Center().Children(func() {
		title := "No track playing"
		if current != nil {
			title = fmt.Sprintf("%s — %s", current.Title, current.Artist)
		}
		m3Text(c, title).FontSize(14).Bold().TextColor(theme.TextPrimary)

		ui.Row(c).FillWidth().Right(0).Children(func() {
			backBtn := ui.Button(c, "").Height(32).Radius(16).Background(theme.SurfaceContainerHigh).Border(1, theme.OutlineVariant)
			if backBtn.Hovered() {
				backBtn.Background(theme.CardHover)
			}
			backBtn.Children(func() {
				ui.Row(c).Fill().Padding(0, 12).Center().Gap(6).Children(func() {
					m3Text(c, "← Back to tracks").FontSize(11).Bold().TextColor(theme.Primary)
				})
			})
			if backBtn.Clicked() {
				s.ViewMode = "tracks"
			}
		})
	})

	if current == nil {
		ui.Column(c).Fill().Center().Children(func() {
			m3Text(c, "Play a track to view synchronized lyrics").FontSize(14).TextColor(theme.TextSecondary)
		})
		return
	}

	if s.LyricsStatus == "loading" {
		ui.Column(c).Fill().Center().Gap(8).Children(func() {
			m3Text(c, "Fetching synchronized lyrics...").FontSize(15).Bold().TextColor(theme.Primary)
			m3Text(c, "Connecting to lyrics provider").FontSize(12).TextColor(theme.TextSecondary)
		})
		return
	}

	if s.LyricsStatus == "not_found" || len(s.LyricsLines) == 0 {
		ui.Column(c).Fill().Center().Gap(8).Children(func() {
			m3Text(c, "No lyrics available for this song").FontSize(15).Bold().TextColor(theme.TextPrimary)
			m3Text(c, "Synchronized LRC lyrics not found").FontSize(12).TextColor(theme.TextSecondary)
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

	const lineH float32 = 48
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
					fontSize = 24
					textColor = theme.Primary
					isBold = true
				} else if absDist < 1.6 {
					fontSize = 17
					textColor = theme.TextPrimary
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
// 4. Playback Bar: Material 3 Persistent Bottom Dock
// -------------------------------------------------------------
func renderBottomPlaybackBar(c *ui.Context, s *AppState, theme M3Theme) {
	current := s.CurrentTrack()

	ui.Row(c).FillWidth().Height(88).Shrink(0).Radius(28).Background(theme.SurfaceContainer).Border(1, theme.OutlineVariant).Padding(10, 22).Center().Children(func() {
		// Left: Track Info & Artwork
		ui.Row(c).Width(320).Center().Gap(14).Children(func() {
			var cover []byte
			if current != nil {
				cover = current.CoverData
			}
			renderCoverBox(c, cover, 56, 14, theme)

			ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(2).Children(func() {
				title := "Ready"
				artist := "Select a track to play"
				if current != nil {
					title = current.Title
					artist = current.Artist
				}
				if len([]rune(title)) > 24 {
					title = string([]rune(title)[:22]) + "..."
				}
				if len([]rune(artist)) > 26 {
					artist = string([]rune(artist)[:24]) + "..."
				}
				m3Text(c, title).FontSize(13).Bold().TextColor(theme.TextPrimary)
				m3Text(c, artist).FontSize(11).TextColor(theme.TextSecondary)
			})

			// Lyrics Toggle Button
			isLyricsMode := s.ViewMode == "lyrics"
			lyricsBtn := ui.Button(c, "").Width(36).Height(36).Radius(18)
			lyBg := theme.SurfaceContainerHigh
			lyFg := theme.TextSecondary
			if isLyricsMode {
				lyBg = theme.Primary
				lyFg = theme.OnPrimary
			} else if lyricsBtn.Hovered() {
				lyBg = theme.CardHover
				lyFg = theme.TextPrimary
			}
			lyricsBtn.Background(lyBg)
			lyricsBtn.Draw(func(p *ui.Painter, r ui.Rect) {
				drawMic(p, r, lyFg)
			})
			if lyricsBtn.Clicked() {
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

		// Center: Playback Controls & Progress Bar
		ui.Column(c).FillWidth().Center().Gap(6).Children(func() {
			// Controls Row
			ui.Row(c).Center().Gap(12).Children(func() {
				// Shuffle Button
				shufBtn := ui.Button(c, "").Width(36).Height(36).Radius(18)
				shufBg := ui.RGBA(0, 0, 0, 0)
				shufFg := theme.TextSecondary
				if s.Shuffle {
					shufBg = theme.PrimaryContainer
					shufFg = theme.OnPrimaryContainer
				} else if shufBtn.Hovered() {
					shufBg = theme.SurfaceContainerHigh
					shufFg = theme.TextPrimary
				}
				shufBtn.Background(shufBg)
				shufBtn.Draw(func(p *ui.Painter, r ui.Rect) {
					drawShuffle(p, r, shufFg)
					if s.Shuffle {
						cx := r.X + r.W/2
						p.Fill(ui.Rect{X: cx - 1.5, Y: r.Y + r.H - 5, W: 3, H: 3}, shufFg, 1.5)
					}
				})
				if shufBtn.Clicked() {
					s.Shuffle = !s.Shuffle
				}

				// Prev Button
				prevBtn := ui.Button(c, "").Width(36).Height(36).Radius(18).Background(ui.RGBA(0, 0, 0, 0))
				prevFg := theme.TextPrimary
				if prevBtn.Hovered() {
					prevBtn.Background(theme.SurfaceContainerHigh)
				}
				prevBtn.Draw(func(p *ui.Painter, r ui.Rect) {
					drawPrev(p, r, prevFg)
				})
				if prevBtn.Clicked() {
					s.PrevTrack()
				}

				// Play/Pause FAB Button (50x50 Primary Circle)
				playBtn := ui.Button(c, "").Width(50).Height(50).Radius(25).Background(theme.Primary).Border(0, theme.Primary)
				if playBtn.Hovered() {
					playBtn.Background(ui.RGB(190, 215, 255))
				}
				playDraw := drawPlay
				if s.Engine.IsPlaying() {
					playDraw = drawPause
				}
				playBtn.Draw(func(p *ui.Painter, r ui.Rect) {
					playDraw(p, r, theme.OnPrimary)
				})
				if playBtn.Clicked() {
					s.TogglePlay()
				}

				// Next Button
				nextBtn := ui.Button(c, "").Width(36).Height(36).Radius(18).Background(ui.RGBA(0, 0, 0, 0))
				nextFg := theme.TextPrimary
				if nextBtn.Hovered() {
					nextBtn.Background(theme.SurfaceContainerHigh)
				}
				nextBtn.Draw(func(p *ui.Painter, r ui.Rect) {
					drawNext(p, r, nextFg)
				})
				if nextBtn.Clicked() {
					s.NextTrack()
				}

				// Repeat Button
				repBtn := ui.Button(c, "").Width(36).Height(36).Radius(18)
				repBg := ui.RGBA(0, 0, 0, 0)
				repFg := theme.TextSecondary
				if s.Repeat {
					repBg = theme.PrimaryContainer
					repFg = theme.OnPrimaryContainer
				} else if repBtn.Hovered() {
					repBg = theme.SurfaceContainerHigh
					repFg = theme.TextPrimary
				}
				repBtn.Background(repBg)
				repBtn.Draw(func(p *ui.Painter, r ui.Rect) {
					drawRepeat(p, r, repFg)
					if s.Repeat {
						cx := r.X + r.W/2
						p.Fill(ui.Rect{X: cx - 1.5, Y: r.Y + r.H - 5, W: 3, H: 3}, repFg, 1.5)
					}
				})
				if repBtn.Clicked() {
					s.Repeat = !s.Repeat
				}
			})

			// Timeline Slider Row
			ui.Row(c).FillWidth().Center().Padding(0, 16).Gap(10).Children(func() {
				m3Text(c, s.CurrentTimeString()).FontSize(11).Bold().TextColor(theme.TextSecondary)

				progSlider := m3Slider(c, &s.Progress, 0.0, 1.0, theme).FillWidth()
				if progSlider.Changed() {
					s.Seek(s.Progress)
				}

				m3Text(c, s.TotalTimeString()).FontSize(11).Bold().TextColor(theme.TextSecondary)
			})
		})

		// Right: Volume Controls
		ui.Row(c).Width(170).Right(0).Center().Gap(8).Children(func() {
			volIconBox := ui.Row(c).Width(20).Height(20).Center()
			volIconBox.Draw(func(p *ui.Painter, r ui.Rect) {
				drawVolume(p, r, theme.TextSecondary)
			})

			volSlider := m3Slider(c, &s.Volume, 0.0, 1.0, theme).Width(80)
			if volSlider.Changed() {
				s.Engine.SetVolume(s.Volume)
			}

			m3Text(c, fmt.Sprintf("%d%%", int(s.Volume*100))).Width(36).FontSize(11).Bold().TextColor(theme.TextSecondary)
		})
	})
}
