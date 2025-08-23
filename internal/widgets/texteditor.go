package widgets

import (
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"simple-editor/internal/interaction"
)

// TextEditor widget that handles text editing and keyboard shortcuts
type TextEditor struct {
	theme    *material.Theme
	editor   widget.Editor
	fileOps  *interaction.FileOperations
}

// NewTextEditor creates a new TextEditor widget
func NewTextEditor(theme *material.Theme) *TextEditor {
	return &TextEditor{
		theme:   theme,
		fileOps: interaction.NewFileOperations(),
	}
}

// HandleEvents processes keyboard shortcuts for the text editor
func (te *TextEditor) HandleEvents(gtx layout.Context) {
	// Handle Ctrl+S (Save)
	for {
		evt, ok := gtx.Event(key.Filter{Focus: &te.editor, Name: "S", Required: key.ModCtrl})
		if !ok {
			break
		}
		if e, ok := evt.(key.Event); ok && e.State == key.Press {
			te.fileOps.SaveFile(&te.editor)
		}
	}

	// Handle Ctrl+O (Open)
	for {
		evt, ok := gtx.Event(key.Filter{Focus: &te.editor, Name: "O", Required: key.ModCtrl})
		if !ok {
			break
		}
		if e, ok := evt.(key.Event); ok && e.State == key.Press {
			te.fileOps.OpenFile(&te.editor)
		}
	}

	// Handle Ctrl+N (New)
	for {
		evt, ok := gtx.Event(key.Filter{Focus: &te.editor, Name: "N", Required: key.ModCtrl})
		if !ok {
			break
		}
		if e, ok := evt.(key.Event); ok && e.State == key.Press {
			te.fileOps.NewFile(&te.editor)
		}
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