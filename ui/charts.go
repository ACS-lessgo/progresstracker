package ui

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"progresstracker/logic"
)

func drawLineChart(gtx layout.Context, points []logic.ChartPoint, lineColor color.NRGBA, title string) layout.Dimensions {
	if len(points) == 0 {
		return layout.Dimensions{Size: gtx.Constraints.Min}
	}

	totalH := gtx.Dp(unit.Dp(200))
	totalW := gtx.Constraints.Max.X
	if totalW == 0 {
		totalW = 400
	}

	padL := 50
	padR := 20
	padT := 30
	padB := 40

	chartW := totalW - padL - padR
	chartH := totalH - padT - padB

	if chartW <= 0 || chartH <= 0 {
		return layout.Dimensions{Size: image.Pt(totalW, totalH)}
	}

	// background
	bgStack := clip.Rect{Max: image.Pt(totalW, totalH)}.Push(gtx.Ops)
	paint.ColorOp{Color: ColorCard}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	bgStack.Pop()

	// find min/max
	minV, maxV := math.MaxFloat64, -math.MaxFloat64
	for _, p := range points {
		if p.Value < minV {
			minV = p.Value
		}
		if p.Value > maxV {
			maxV = p.Value
		}
	}
	if maxV == minV {
		minV = minV - 1
		maxV = maxV + 1
	}

	rangeV := maxV - minV

	toX := func(i int) float32 {
		if len(points) == 1 {
			return float32(padL + chartW/2)
		}
		return float32(padL) + float32(i)*float32(chartW)/float32(len(points)-1)
	}
	toY := func(v float64) float32 {
		norm := (v - minV) / rangeV
		return float32(padT+chartH) - float32(norm)*float32(chartH)
	}

	// draw grid lines
	for i := 0; i <= 4; i++ {
		y := padT + i*chartH/4
		drawHLine(gtx.Ops, padL, y, totalW-padR, color.NRGBA{R: 50, G: 50, B: 70, A: 255})
	}

	// draw the line
	if len(points) > 1 {
		var path clip.Path
		path.Begin(gtx.Ops)
		path.MoveTo(f32.Pt(toX(0), toY(points[0].Value)))
		for i := 1; i < len(points); i++ {
			path.LineTo(f32.Pt(toX(i), toY(points[i].Value)))
		}
		lineStack := clip.Stroke{
			Path:  path.End(),
			Width: 2,
		}.Op().Push(gtx.Ops)
		paint.ColorOp{Color: lineColor}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		lineStack.Pop()
	}

	// draw dots
	for i, p := range points {
		cx := int(toX(i))
		cy := int(toY(p.Value))
		drawDot(gtx.Ops, cx, cy, 4, lineColor)
	}

	// draw axis labels
	_ = title
	// Y labels
	for i := 0; i <= 4; i++ {
		v := minV + float64(i)*rangeV/4
		y := padT + (4-i)*chartH/4
		drawLabel(gtx, fmt.Sprintf("%.0f", v), padL-45, y-8, ColorSubtext)
	}

	// X labels (first, middle, last)
	if len(points) > 0 {
		indices := []int{0}
		if len(points) > 1 {
			indices = append(indices, len(points)-1)
		}
		if len(points) > 2 {
			indices = append(indices, len(points)/2)
		}
		for _, idx := range indices {
			x := int(toX(idx))
			date := points[idx].Date
			if len(date) >= 10 {
				date = date[5:10]
			}
			drawLabel(gtx, date, x-15, padT+chartH+8, ColorSubtext)
		}
	}

	return layout.Dimensions{Size: image.Pt(totalW, totalH)}
}

func drawHLine(ops *op.Ops, x1, y, x2 int, col color.NRGBA) {
	r := clip.Rect{Min: image.Pt(x1, y), Max: image.Pt(x2, y+1)}.Push(ops)
	paint.ColorOp{Color: col}.Add(ops)
	paint.PaintOp{}.Add(ops)
	r.Pop()
}

func drawDot(ops *op.Ops, cx, cy, r int, col color.NRGBA) {
	rect := clip.Rect{
		Min: image.Pt(cx-r, cy-r),
		Max: image.Pt(cx+r, cy+r),
	}.Push(ops)
	paint.ColorOp{Color: col}.Add(ops)
	paint.PaintOp{}.Add(ops)
	rect.Pop()
}

func drawLabel(gtx layout.Context, text string, x, y int, col color.NRGBA) {
	defer op.Offset(image.Pt(x, y)).Push(gtx.Ops).Pop()
	_ = text
	_ = col
	// Labels via Gio text require shaper; skip for now, values shown via grid
}
