package main

import (
	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"simple-editor/internal/widgets"
)

func main() {
	go func() {
		w := &app.Window{}
		w.Option(app.Title("Minimal Editor"))
		w.Option(app.Size(unit.Dp(800), unit.Dp(600)))

		theme := material.NewTheme()
		textEditor := widgets.NewTextEditor(theme)

		var ops op.Ops

		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				return
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)

				// Handle events and render
				textEditor.HandleEvents(gtx)
				textEditor.Layout(gtx)

				e.Frame(gtx.Ops)
			}
		}
	}()

	app.Main()
}