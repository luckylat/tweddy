package main

import (
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"simple-editor/internal/style"
	"simple-editor/internal/widgets"
)

type Application struct {
	theme      *material.Theme
	header     *widgets.Header
	textEditor *widgets.TextEditor
	tabBar     *widgets.TabBar
	appTheme   style.Theme
}

func NewApplication() *Application {
	theme := material.NewTheme()
	appTheme := style.DefaultTheme()
	textEditor := widgets.NewTextEditor(theme, appTheme)
	tabBar := widgets.NewTabBar(theme, textEditor.GetTabManager())

	return &Application{
		theme:      theme,
		header:     widgets.NewHeader(theme, appTheme),
		textEditor: textEditor,
		tabBar:     tabBar,
		appTheme:   appTheme,
	}
}

func (a *Application) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Background{}.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			paint.ColorOp{Color: a.appTheme.Application.MainBackground}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: gtx.Constraints.Max}
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.header.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(10)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.tabBar.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(10)}.Layout(gtx)
		}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return a.textEditor.Layout(gtx)
			}),
			)
		},
	)
}

func (a *Application) HandleEvents(gtx layout.Context) {
	newFile, openFile, saveFile := a.textEditor.HandleKeyboardShortcuts(gtx)

	// Handle tab bar events (keyboard shortcuts for tab switching)
	a.tabBar.HandleEvents(gtx)

	if newFile || a.header.NewClicked(gtx) {
		a.textEditor.NewFile()
	}

	if openFile || a.header.OpenClicked(gtx) {
		a.textEditor.OpenFile()
	}

	if saveFile || a.header.SaveClicked(gtx) {
		a.textEditor.SaveFile()
	}

	if a.header.SaveAsClicked(gtx) {
		a.textEditor.SaveAsFile()
	}
}

func main() {
	go func() {
		w := &app.Window{}
		w.Option(app.Title("Simple Editor"))
		w.Option(app.Size(unit.Dp(800), unit.Dp(600)))

		application := NewApplication()

		var ops op.Ops

		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				return
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)

				application.HandleEvents(gtx)

				application.Layout(gtx)

				e.Frame(gtx.Ops)
			}
		}
	}()

	app.Main()
}
