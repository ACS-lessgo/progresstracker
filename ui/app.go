package ui

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"time"

	"progresstracker/data"
	"progresstracker/logic"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type NavTab int

const (
	TabLog NavTab = iota
	TabHistory
	TabAnalytics
)

type App struct {
	th      *material.Theme
	tracker *logic.Tracker
	anal    *logic.Analytics
	repo    *data.Repository

	activeTab NavTab

	dayBtns   [5]widget.Clickable
	activeDay int
	exBtns    []widget.Clickable
	activeEx  int
	exScroll  widget.List

	weightEdit widget.Editor
	repsEdit   widget.Editor
	setsEdit   widget.Editor
	notesEdit  widget.Editor
	dateEdit   widget.Editor
	saveBtn    widget.Clickable
	statusMsg  string
	statusOK   bool

	navBtns [3]widget.Clickable

	histList    widget.List
	histEntries []*data.Entry
	histPB      *data.PersonalBest

	chartWeight []logic.ChartPoint
	chartVolume []logic.ChartPoint
	chartScroll widget.List
	logScroll   widget.List

	currentWeek   int
	weekToggleBtn widget.Clickable
}

func NewApp(repo *data.Repository, tracker *logic.Tracker, anal *logic.Analytics) *App {
	a := &App{
		th:        NewTheme(),
		tracker:   tracker,
		anal:      anal,
		repo:      repo,
		activeDay: 0,
		activeEx:  0,
	}
	a.exScroll.Axis = layout.Vertical
	a.histList.Axis = layout.Vertical
	a.chartScroll.Axis = layout.Vertical
	a.currentWeek = repo.GetCurrentWeek()
	a.logScroll.Axis = layout.Vertical
	a.dateEdit.SetText(time.Now().Format("2006-01-02"))
	a.dateEdit.SingleLine = true
	a.weightEdit.SingleLine = true
	a.repsEdit.SingleLine = true
	a.setsEdit.SingleLine = true
	a.notesEdit.SingleLine = true
	a.rebuildExBtns()
	return a
}

func (a *App) rebuildExBtns() {
	exs := data.WorkoutDays(data.DayOrder[a.activeDay], a.currentWeek)
	a.exBtns = make([]widget.Clickable, len(exs))
	if a.activeEx >= len(exs) {
		a.activeEx = 0
	}
}

func (a *App) currentExercise() string {
	exs := data.WorkoutDays(data.DayOrder[a.activeDay], a.currentWeek)
	if len(exs) == 0 {
		return ""
	}
	if a.activeEx >= len(exs) {
		a.activeEx = 0
	}
	return exs[a.activeEx]
}

func (a *App) loadHistory() {
	ex := a.currentExercise()
	if ex == "" {
		return
	}
	entries, err := a.tracker.GetHistory(ex)
	if err == nil {
		a.histEntries = entries
	}
	pb, err := a.tracker.GetPersonalBest(ex)
	if err == nil {
		a.histPB = pb
	}
}

func (a *App) loadCharts() {
	ex := a.currentExercise()
	if ex == "" {
		return
	}
	pts, _ := a.anal.WeightOverTime(ex)
	a.chartWeight = pts
	pts2, _ := a.anal.VolumeOverTime(ex)
	a.chartVolume = pts2
}

func (a *App) Run(w *app.Window) error {
	var ops op.Ops
	for {
		e := w.Event()
		switch ev := e.(type) {
		case app.DestroyEvent:
			return ev.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, ev)
			a.update(gtx)
			a.layout(gtx)
			ev.Frame(&ops)
		}
	}
}

func (a *App) update(gtx layout.Context) {
	if a.weekToggleBtn.Clicked(gtx) {
		if a.currentWeek == 1 {
			a.currentWeek = 2
		} else {
			a.currentWeek = 1
		}
		a.repo.SetCurrentWeek(a.currentWeek)
		a.activeEx = 0
		a.rebuildExBtns()
		a.statusMsg = ""
	}

	for i := range a.navBtns {
		if a.navBtns[i].Clicked(gtx) {
			a.activeTab = NavTab(i)
			if a.activeTab == TabHistory {
				a.loadHistory()
			} else if a.activeTab == TabAnalytics {
				a.loadCharts()
			}
		}
	}

	for i := range a.dayBtns {
		if a.dayBtns[i].Clicked(gtx) {
			if a.activeDay != i {
				a.activeDay = i
				a.activeEx = 0
				a.rebuildExBtns()
				a.statusMsg = ""
			}
		}
	}

	day := data.DayOrder[a.activeDay]
	exs := data.WorkoutDays(day, a.currentWeek)
	for i := range a.exBtns {
		if i < len(exs) && a.exBtns[i].Clicked(gtx) {
			if a.activeEx != i {
				a.activeEx = i
				a.statusMsg = ""
				if a.activeTab == TabHistory {
					a.loadHistory()
				} else if a.activeTab == TabAnalytics {
					a.loadCharts()
				}
			}
		}
	}

	if a.saveBtn.Clicked(gtx) {
		ex := a.currentExercise()
		entry, err := a.tracker.AddEntry(
			ex,
			a.weightEdit.Text(),
			a.repsEdit.Text(),
			a.setsEdit.Text(),
			a.notesEdit.Text(),
			a.dateEdit.Text(),
		)
		if err != nil {
			a.statusMsg = "Error: " + err.Error()
			a.statusOK = false
		} else {
			a.statusMsg = fmt.Sprintf("Saved! Volume: %.0f kg", entry.Volume)
			a.statusOK = true
			a.weightEdit.SetText("")
			a.repsEdit.SetText("")
			a.setsEdit.SetText("")
			a.notesEdit.SetText("")
			a.dateEdit.SetText(time.Now().Format("2006-01-02"))
		}
	}
}

func (a *App) layout(gtx layout.Context) layout.Dimensions {
	fillRect(gtx, ColorBg, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(a.layoutSidebar),
		layout.Flexed(1, a.layoutMain),
	)
}

func (a *App) layoutSidebar(gtx layout.Context) layout.Dimensions {
	w := gtx.Dp(unit.Dp(210))
	gtx.Constraints.Min.X = w
	gtx.Constraints.Max.X = w
	fillRect(gtx, ColorSidebar, w, gtx.Constraints.Max.Y)

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// App title
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				lbl := material.H6(a.th, "ProgressTracker")
				lbl.Color = ColorAccent
				return lbl.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			label := "Switch to Week 2"
			if a.currentWeek == 2 {
				label = "Switch to Week 1"
			}
			weekLabel := fmt.Sprintf("WEEK %d  ·  %s", a.currentWeek, label)
			return layout.Inset{
				Left: unit.Dp(16), Right: unit.Dp(16),
				Top: unit.Dp(10), Bottom: unit.Dp(10),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						t := material.Caption(a.th, fmt.Sprintf("CURRENT: WEEK %d", a.currentWeek))
						t.Color = ColorAccent
						return t.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						_ = weekLabel
						btn := material.Button(a.th, &a.weekToggleBtn, label)
						btn.Background = color.NRGBA{R: 0, G: 80, B: 60, A: 255}
						btn.Color = ColorAccent
						btn.TextSize = unit.Sp(11)
						return btn.Layout(gtx)
					}),
				)
			})
		}),
		layout.Rigid(a.sidebarDivider),

		// Navigation section
		layout.Rigid(a.sectionLabel("NAVIGATION")),
		layout.Rigid(a.navBtn(0, "Log Workout")),
		layout.Rigid(a.navBtn(1, "History")),
		layout.Rigid(a.navBtn(2, "Analytics")),
		layout.Rigid(a.sidebarDivider),

		// Workout Day section
		layout.Rigid(a.sectionLabel("WORKOUT DAY")),
		layout.Rigid(a.dayBtn(0, "Mon · Chest")),
		layout.Rigid(a.dayBtn(1, "Tue · Back")),
		layout.Rigid(a.dayBtn(2, "Wed · Shoulders")),
		layout.Rigid(a.dayBtn(3, "Thu · Arms")),
		layout.Rigid(a.dayBtn(4, "Fri · Legs")),
		layout.Rigid(a.sidebarDivider),

		// Exercise section
		layout.Rigid(a.sectionLabel("EXERCISES")),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			day := data.DayOrder[a.activeDay]
			exs := data.WorkoutDays(day, a.currentWeek)
			return a.exScroll.Layout(gtx, len(exs), func(gtx layout.Context, idx int) layout.Dimensions {
				if idx >= len(a.exBtns) {
					return layout.Dimensions{}
				}
				return a.exBtn(idx, exs[idx])(gtx)
			})
		}),
	)
}

func (a *App) sidebarDivider(gtx layout.Context) layout.Dimensions {
	r := clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, 1)}.Push(gtx.Ops)
	paint.ColorOp{Color: ColorBorder}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	r.Pop()
	return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, 1)}
}

func (a *App) sectionLabel(text string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Left: unit.Dp(16), Top: unit.Dp(10), Bottom: unit.Dp(4), Right: unit.Dp(8)}.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				lbl := material.Caption(a.th, text)
				lbl.Color = ColorSubtext
				return lbl.Layout(gtx)
			})
	}
}

func (a *App) navBtn(idx int, label string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		active := a.activeTab == NavTab(idx)
		return a.sidebarClickable(gtx, &a.navBtns[idx], label, active)
	}
}

func (a *App) dayBtn(idx int, label string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		active := a.activeDay == idx
		return a.sidebarClickable(gtx, &a.dayBtns[idx], label, active)
	}
}

func (a *App) exBtn(idx int, label string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		active := a.activeEx == idx
		return a.sidebarClickable(gtx, &a.exBtns[idx], label, active)
	}
}

func (a *App) sidebarClickable(gtx layout.Context, btn *widget.Clickable, label string, active bool) layout.Dimensions {
	return btn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		bg := ColorSidebar
		if active {
			bg = color.NRGBA{R: 0, G: 60, B: 45, A: 255}
		}
		// fill background
		r := clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)}.Push(gtx.Ops)
		paint.ColorOp{Color: bg}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		r.Pop()

		return layout.Inset{
			Top: unit.Dp(10), Bottom: unit.Dp(10),
			Left: unit.Dp(16), Right: unit.Dp(8),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			// accent bar on left if active
			if active {
				bar := clip.Rect{
					Min: image.Pt(-16, -10),
					Max: image.Pt(-13, gtx.Dp(unit.Dp(40))),
				}.Push(gtx.Ops)
				paint.ColorOp{Color: ColorAccent}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				bar.Pop()
			}
			lbl := material.Body2(a.th, label)
			if active {
				lbl.Color = ColorAccent
			} else {
				lbl.Color = ColorText
			}
			return lbl.Layout(gtx)
		})
	})
}

func (a *App) layoutMain(gtx layout.Context) layout.Dimensions {
	switch a.activeTab {
	case TabLog:
		return a.layoutLog(gtx)
	case TabHistory:
		return a.layoutHistory(gtx)
	case TabAnalytics:
		return a.layoutAnalytics(gtx)
	}
	return layout.Dimensions{}
}

func (a *App) layoutLog(gtx layout.Context) layout.Dimensions {
	return layout.UniformInset(unit.Dp(28)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return a.logScroll.Layout(gtx, 3, func(gtx layout.Context, idx int) layout.Dimensions {
			switch idx {
			case 0:
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						t := material.H5(a.th, a.currentExercise())
						t.Color = ColorText
						return t.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						t := material.Body2(a.th, data.DayOrder[a.activeDay]+" Workout")
						t.Color = ColorSubtext
						return t.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
					layout.Rigid(a.layoutLastCard),
					layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
				)
			case 1:
				return a.layoutFormCard(gtx)
			case 2:
				return layout.Spacer{Height: unit.Dp(40)}.Layout(gtx)
			}
			return layout.Dimensions{}
		})
	})
}

func (a *App) layoutLastCard(gtx layout.Context) layout.Dimensions {
	last, _ := a.tracker.GetLastEntry(a.currentExercise())
	pb, _ := a.tracker.GetPersonalBest(a.currentExercise())

	// draw card background based on content height
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// card bg
			return drawCard(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						t := material.Body1(a.th, "Last Session & Personal Bests")
						t.Color = ColorAccent
						return t.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						s := "No previous entries"
						if last != nil {
							s = fmt.Sprintf("Last: %.1f kg × %d reps × %d sets = %.0f vol  (%s)", last.Weight, last.Reps, last.Sets, last.Volume, last.Date)
						}
						t := material.Body2(a.th, s)
						t.Color = ColorText
						return t.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						s := "No personal bests yet"
						if pb != nil && pb.MaxWeight > 0 {
							s = fmt.Sprintf("🏆  Best Weight: %.1f kg     Best Volume: %.0f", pb.MaxWeight, pb.MaxVolume)
						}
						t := material.Body2(a.th, s)
						t.Color = ColorGold
						return t.Layout(gtx)
					}),
				)
			})
		}),
	)
}

func (a *App) layoutFormCard(gtx layout.Context) layout.Dimensions {
	return drawCard(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				t := material.H6(a.th, "Log Entry")
				t.Color = ColorAccent
				return t.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				t := material.Caption(a.th, "WEIGHT (kg)")
				t.Color = ColorSubtext
				return t.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return plainEditor(gtx, a.th, &a.weightEdit, "e.g. 60")
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(14)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				t := material.Caption(a.th, "REPS")
				t.Color = ColorSubtext
				return t.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return plainEditor(gtx, a.th, &a.repsEdit, "e.g. 10")
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(14)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				t := material.Caption(a.th, "SETS")
				t.Color = ColorSubtext
				return t.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return plainEditor(gtx, a.th, &a.setsEdit, "e.g. 3")
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(14)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				t := material.Caption(a.th, "DATE (YYYY-MM-DD)")
				t.Color = ColorSubtext
				return t.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return plainEditor(gtx, a.th, &a.dateEdit, "2025-01-01")
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(14)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				t := material.Caption(a.th, "NOTES (optional)")
				t.Color = ColorSubtext
				return t.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return plainEditor(gtx, a.th, &a.notesEdit, "e.g. felt strong today")
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(22)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(a.th, &a.saveBtn, "SAVE ENTRY")
				btn.Background = ColorAccent
				btn.Color = color.NRGBA{A: 255}
				return btn.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if a.statusMsg == "" {
					return layout.Dimensions{}
				}
				t := material.Body2(a.th, a.statusMsg)
				if a.statusOK {
					t.Color = ColorAccent
				} else {
					t.Color = ColorRed
				}
				return t.Layout(gtx)
			}),
		)
	})
}

func (a *App) fieldOf(label string, ed *widget.Editor) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				t := material.Caption(a.th, label)
				t.Color = ColorSubtext
				return t.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return styledEditor(gtx, a.th, ed)
			}),
		)
	}
}

func (a *App) layoutHistory(gtx layout.Context) layout.Dimensions {
	a.loadHistory()
	return layout.UniformInset(unit.Dp(28)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				t := material.H5(a.th, "History: "+a.currentExercise())
				t.Color = ColorText
				return t.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if a.histPB == nil || a.histPB.MaxWeight == 0 {
					return layout.Dimensions{}
				}
				return cardLayout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							t := material.Body1(a.th, fmt.Sprintf("🏆  Best Weight: %.1f kg", a.histPB.MaxWeight))
							t.Color = ColorGold
							return t.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(40)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							t := material.Body1(a.th, fmt.Sprintf("🔥  Best Volume: %.0f", a.histPB.MaxVolume))
							t.Color = ColorAccent2
							return t.Layout(gtx)
						}),
					)
				})
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				if len(a.histEntries) == 0 {
					t := material.Body1(a.th, "No entries yet. Log your first workout!")
					t.Color = ColorSubtext
					return t.Layout(gtx)
				}
				return a.histList.Layout(gtx, len(a.histEntries), func(gtx layout.Context, idx int) layout.Dimensions {
					e := a.histEntries[idx]
					isPB := a.histPB != nil && e.Weight == a.histPB.MaxWeight
					return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return cardLayout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									t := material.Body2(a.th, e.Date)
									t.Color = ColorSubtext
									return layout.Inset{Right: unit.Dp(16)}.Layout(gtx, t.Layout)
								}),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									s := fmt.Sprintf("%.1f kg × %d reps × %d sets", e.Weight, e.Reps, e.Sets)
									t := material.Body1(a.th, s)
									if isPB {
										t.Color = ColorGold
									} else {
										t.Color = ColorText
									}
									return t.Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									t := material.Body2(a.th, fmt.Sprintf("Vol: %.0f", e.Volume))
									t.Color = ColorAccent2
									return t.Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if isPB {
										t := material.Body2(a.th, "  🏆")
										return t.Layout(gtx)
									}
									return layout.Dimensions{}
								}),
							)
						})
					})
				})
			}),
		)
	})
}

func (a *App) layoutAnalytics(gtx layout.Context) layout.Dimensions {
	a.loadCharts()
	return layout.UniformInset(unit.Dp(28)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return a.chartScroll.Layout(gtx, 5, func(gtx layout.Context, idx int) layout.Dimensions {
			switch idx {
			case 0:
				t := material.H5(a.th, "Analytics: "+a.currentExercise())
				t.Color = ColorText
				return t.Layout(gtx)
			case 1:
				return layout.Spacer{Height: unit.Dp(20)}.Layout(gtx)
			case 2:
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						t := material.Body1(a.th, "Weight Over Time")
						t.Color = ColorAccent
						return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, t.Layout)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if len(a.chartWeight) == 0 {
							t := material.Body2(a.th, "No data yet — log some entries first.")
							t.Color = ColorSubtext
							return t.Layout(gtx)
						}
						return drawLineChart(gtx, a.chartWeight, ColorChartLine, "Weight")
					}),
				)
			case 3:
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(layout.Spacer{Height: unit.Dp(24)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						t := material.Body1(a.th, "Volume Over Time")
						t.Color = ColorAccent2
						return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, t.Layout)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if len(a.chartVolume) == 0 {
							t := material.Body2(a.th, "No data yet — log some entries first.")
							t.Color = ColorSubtext
							return t.Layout(gtx)
						}
						return drawLineChart(gtx, a.chartVolume, ColorChartVol, "Volume")
					}),
				)
			case 4:
				return a.layoutStatsTable(gtx)
			}
			return layout.Dimensions{}
		})
	})
}

func (a *App) layoutStatsTable(gtx layout.Context) layout.Dimensions {
	return layout.Inset{Top: unit.Dp(24)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return cardLayout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					t := material.Body1(a.th, "Summary Statistics")
					t.Color = ColorAccent
					return layout.Inset{Bottom: unit.Dp(12)}.Layout(gtx, t.Layout)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if len(a.chartWeight) == 0 {
						t := material.Body2(a.th, "No data available")
						t.Color = ColorSubtext
						return t.Layout(gtx)
					}
					first := a.chartWeight[0]
					last := a.chartWeight[len(a.chartWeight)-1]
					diff := last.Value - first.Value
					diffStr := fmt.Sprintf("+%.1f", diff)
					if diff < 0 {
						diffStr = fmt.Sprintf("%.1f", diff)
					}
					rows := []struct{ label, value string }{
						{"Total Sessions", fmt.Sprintf("%d", len(a.chartWeight))},
						{"Starting Weight", fmt.Sprintf("%.1f kg", first.Value)},
						{"Current Weight", fmt.Sprintf("%.1f kg", last.Value)},
						{"Weight Change", diffStr + " kg"},
					}
					var children []layout.FlexChild
					for _, row := range rows {
						r := row
						children = append(children,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
									layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
										t := material.Body2(a.th, r.label)
										t.Color = ColorSubtext
										return t.Layout(gtx)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										t := material.Body2(a.th, r.value)
										t.Color = ColorText
										return t.Layout(gtx)
									}),
								)
							}),
							layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
						)
					}
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
				}),
			)
		})
	})
}

func Run(repo *data.Repository, tracker *logic.Tracker, anal *logic.Analytics) {
	a := NewApp(repo, tracker, anal)
	w := new(app.Window)
	w.Option(
		app.Title("ProgressTracker"),
		app.Size(unit.Dp(1100), unit.Dp(750)),
	)
	if err := a.Run(w); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
