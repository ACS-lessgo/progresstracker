package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func fillRect(gtx layout.Context, col color.NRGBA, w, h int) {
	r := clip.Rect{Max: image.Pt(w, h)}.Push(gtx.Ops)
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	r.Pop()
}

func cardLayout(gtx layout.Context, content layout.Widget) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			r := clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Max}, 6)
			stack := r.Push(gtx.Ops)
			paint.ColorOp{Color: ColorCard}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			stack.Pop()
			return layout.Dimensions{Size: gtx.Constraints.Max}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, content)
		}),
	)
}

func styledEditor(gtx layout.Context, th *material.Theme, ed *widget.Editor) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			r := clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Max}, 4)
			stack := r.Push(gtx.Ops)
			paint.ColorOp{Color: color.NRGBA{R: 40, G: 40, B: 55, A: 255}}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			stack.Pop()
			borderStack := clip.Stroke{
				Path:  clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Max}, 4).Path(gtx.Ops),
				Width: 1,
			}.Op().Push(gtx.Ops)
			paint.ColorOp{Color: ColorBorder}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			borderStack.Pop()
			return layout.Dimensions{Size: gtx.Constraints.Max}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				e := material.Editor(th, ed, "")
				e.Color = ColorText
				e.HintColor = ColorSubtext
				return e.Layout(gtx)
			})
		}),
	)
}

// drawCard draws a rounded card whose height matches its content exactly.
// Unlike cardLayout (which uses Stack+Expanded), this works correctly inside
// scrollable lists because the height is driven by the Stacked child.
func drawCard(gtx layout.Context, content layout.Widget) layout.Dimensions {
	// First measure content
	m := op.Record(gtx.Ops)
	dims := layout.UniformInset(unit.Dp(16)).Layout(gtx, content)
	call := m.Stop()

	// Draw background at measured size
	r := clip.UniformRRect(image.Rectangle{Max: dims.Size}, 8)
	stack := r.Push(gtx.Ops)
	paint.ColorOp{Color: ColorCard}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	stack.Pop()

	// Replay content on top
	call.Add(gtx.Ops)
	return dims
}

// plainEditor renders a simple single-line text input with visible background.
func plainEditor(gtx layout.Context, th *material.Theme, ed *widget.Editor, hint string) layout.Dimensions {
	// fixed height input box
	gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(36))
	gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(36))

	r := clip.UniformRRect(image.Rectangle{Max: image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Min.Y)}, 4)
	stack := r.Push(gtx.Ops)
	paint.ColorOp{Color: color.NRGBA{R: 45, G: 45, B: 60, A: 255}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	stack.Pop()

	// border
	borderStack := clip.Stroke{
		Path:  clip.UniformRRect(image.Rectangle{Max: image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Min.Y)}, 4).Path(gtx.Ops),
		Width: 1,
	}.Op().Push(gtx.Ops)
	paint.ColorOp{Color: ColorBorder}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	borderStack.Pop()

	return layout.Inset{Left: unit.Dp(10), Right: unit.Dp(10), Top: unit.Dp(8), Bottom: unit.Dp(8)}.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			e := material.Editor(th, ed, hint)
			e.Color = ColorText
			e.HintColor = ColorSubtext
			return e.Layout(gtx)
		},
	)
}
