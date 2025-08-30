package widgets

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// TextEditor widget that handles text editing
type TextEditor struct {
	theme  *material.Theme
	editor widget.Editor
}

// NewTextEditor creates a new TextEditor widget
func NewTextEditor(theme *material.Theme) *TextEditor {
	return &TextEditor{
		theme: theme,
	}
}


// Layout renders the text editor with proper padding
func (te *TextEditor) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Background{}.Layout(gtx,
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
				return material.Editor(te.theme, &te.editor, "Type here...").Layout(gtx)
			})
		},
	)
}