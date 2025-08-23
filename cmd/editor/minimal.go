package main

import (
	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"simple-editor/internal/interaction"
)

func main() {
	go func() {
		w := &app.Window{}
		w.Option(app.Title("Minimal Editor"))
		w.Option(app.Size(unit.Dp(800), unit.Dp(600)))

		var editor widget.Editor
		theme := material.NewTheme()
		fileOps := interaction.NewFileOperations()

		var ops op.Ops

		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				return
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)

				// Handle keyboard shortcuts
				for {
					evt, ok := gtx.Event(key.Filter{Focus: &editor, Name: "S", Required: key.ModCtrl})
					if !ok {
						break
					}
					if e, ok := evt.(key.Event); ok && e.State == key.Press {
						fileOps.SaveFile(&editor)
					}
				}

				for {
					evt, ok := gtx.Event(key.Filter{Focus: &editor, Name: "O", Required: key.ModCtrl})
					if !ok {
						break
					}
					if e, ok := evt.(key.Event); ok && e.State == key.Press {
						fileOps.OpenFile(&editor)
					}
				}

				for {
					evt, ok := gtx.Event(key.Filter{Focus: &editor, Name: "N", Required: key.ModCtrl})
					if !ok {
						break
					}
					if e, ok := evt.(key.Event); ok && e.State == key.Press {
						fileOps.NewFile(&editor)
					}
				}

				layout.Background{}.Layout(gtx,
					func(gtx layout.Context) layout.Dimensions {
						return layout.Dimensions{Size: gtx.Constraints.Max}
					},
					func(gtx layout.Context) layout.Dimensions {
						inset := layout.Inset{
							Top:    unit.Dp(20),
							Bottom: unit.Dp(20),
							Left:   unit.Dp(20),
							Right:  unit.Dp(20),
						}
						return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return material.Editor(theme, &editor, "Type here...").Layout(gtx)
						})
					},
				)

				e.Frame(gtx.Ops)
			}
		}
	}()

	app.Main()
}