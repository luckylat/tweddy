package widgets

import (
	"image"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"simple-editor/internal/style"
)

type Header struct {
	theme        *material.Theme
	newButton    widget.Clickable
	openButton   widget.Clickable
	saveButton   widget.Clickable
	saveAsButton widget.Clickable
	appTheme     style.Theme
}

func NewHeader(theme *material.Theme, appTheme style.Theme) *Header {
	return &Header{
		theme:    theme,
		appTheme: appTheme,
	}
}

// drawBottomBorder draws a bottom border line
func (h *Header) drawBottomBorder(gtx layout.Context) {
	borderHeight := unit.Dp(1)
	
	defer op.Offset(image.Point{X: 0, Y: gtx.Dp(borderHeight)}).Push(gtx.Ops).Pop()
	
	rect := clip.Rect{
		Max: image.Point{
			X: gtx.Constraints.Max.X,
			Y: gtx.Dp(borderHeight),
		},
	}
	
	defer rect.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: h.appTheme.Application.BorderColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

// customButton creates a button with custom colors from the theme
func (h *Header) customButton(gtx layout.Context, button *widget.Clickable, text string) layout.Dimensions {
	return material.Clickable(gtx, button, func(gtx layout.Context) layout.Dimensions {
		return layout.Background{}.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				// Button background color with hover effect
				bgColor := h.appTheme.Application.ButtonBackground
				if button.Hovered() {
					bgColor = h.appTheme.Application.ButtonHover
				}
				paint.ColorOp{Color: bgColor}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				return layout.Dimensions{Size: gtx.Constraints.Min}
			},
			func(gtx layout.Context) layout.Dimensions {
				inset := layout.Inset{
					Top:    unit.Dp(8),
					Bottom: unit.Dp(8),
					Left:   unit.Dp(12),
					Right:  unit.Dp(12),
				}
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					label := material.Body1(h.theme, text)
					label.Color = h.appTheme.Application.ButtonText
					return label.Layout(gtx)
				})
			},
		)
	})
}

func (h *Header) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Background{}.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					paint.ColorOp{Color: h.appTheme.Application.HeaderBackground}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)
					return layout.Dimensions{Size: gtx.Constraints.Max}
				},
				func(gtx layout.Context) layout.Dimensions {
					inset := layout.Inset{
						Top:    unit.Dp(8),
						Bottom: unit.Dp(8),
						Left:   unit.Dp(8),
						Right:  unit.Dp(8),
					}
					return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return h.customButton(gtx, &h.newButton, "New")
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Width: unit.Dp(10)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return h.customButton(gtx, &h.openButton, "Open")
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Width: unit.Dp(10)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return h.customButton(gtx, &h.saveButton, "Save")
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Width: unit.Dp(10)}.Layout(gtx)
		}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return h.customButton(gtx, &h.saveAsButton, "Save As")
							}),
						)
					})
				},
			)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			h.drawBottomBorder(gtx)
			return layout.Dimensions{}
		}),
	)
}

func (h *Header) NewClicked(gtx layout.Context) bool {
	return h.newButton.Clicked(gtx)
}

func (h *Header) OpenClicked(gtx layout.Context) bool {
	return h.openButton.Clicked(gtx)
}

func (h *Header) SaveClicked(gtx layout.Context) bool {
	return h.saveButton.Clicked(gtx)
}

func (h *Header) SaveAsClicked(gtx layout.Context) bool {
	return h.saveAsButton.Clicked(gtx)
}
