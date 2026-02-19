package ui

import (
	"image/color"

	"gioui.org/font/gofont"
	"gioui.org/text"
	"gioui.org/widget/material"
)

var (
	ColorBg        = color.NRGBA{R: 18, G: 18, B: 24, A: 255}
	ColorSidebar   = color.NRGBA{R: 24, G: 24, B: 32, A: 255}
	ColorCard      = color.NRGBA{R: 30, G: 30, B: 42, A: 255}
	ColorAccent    = color.NRGBA{R: 0, G: 200, B: 140, A: 255}
	ColorAccent2   = color.NRGBA{R: 100, G: 120, B: 255, A: 255}
	ColorText      = color.NRGBA{R: 230, G: 230, B: 240, A: 255}
	ColorSubtext   = color.NRGBA{R: 140, G: 140, B: 160, A: 255}
	ColorBorder    = color.NRGBA{R: 50, G: 50, B: 70, A: 255}
	ColorActive    = color.NRGBA{R: 0, G: 200, B: 140, A: 40}
	ColorGold      = color.NRGBA{R: 255, G: 200, B: 50, A: 255}
	ColorRed       = color.NRGBA{R: 255, G: 80, B: 80, A: 255}
	ColorChartLine = color.NRGBA{R: 0, G: 200, B: 140, A: 255}
	ColorChartVol  = color.NRGBA{R: 100, G: 120, B: 255, A: 255}
)

func NewTheme() *material.Theme {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	th.Palette.Bg = ColorBg
	th.Palette.Fg = ColorText
	th.Palette.ContrastBg = ColorAccent
	th.Palette.ContrastFg = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	return th
}
