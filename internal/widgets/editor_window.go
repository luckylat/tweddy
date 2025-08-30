package widgets

import (
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"simple-editor/internal/interaction"
)

// EditorWindow is the main window widget that contains the text editor and handles global shortcuts
type EditorWindow struct {
	theme      *material.Theme
	textEditor *TextEditor
	fileOps    *interaction.FileOperations
}

// NewEditorWindow creates a new EditorWindow widget
func NewEditorWindow(theme *material.Theme) *EditorWindow {
	return &EditorWindow{
		theme:      theme,
		textEditor: NewTextEditor(theme),
		fileOps:    interaction.NewFileOperations(),
	}
}

// HandleEvents processes window-level keyboard shortcuts
func (ew *EditorWindow) HandleEvents(gtx layout.Context) {
	// Handle Ctrl+S (Save)
	for {
		evt, ok := gtx.Event(key.Filter{Name: "S", Required: key.ModCtrl})
		if !ok {
			break
		}
		if e, ok := evt.(key.Event); ok && e.State == key.Press {
			ew.fileOps.SaveFile(&ew.textEditor.editor)
		}
	}

	// Handle Ctrl+O (Open)
	for {
		evt, ok := gtx.Event(key.Filter{Name: "O", Required: key.ModCtrl})
		if !ok {
			break
		}
		if e, ok := evt.(key.Event); ok && e.State == key.Press {
			ew.fileOps.OpenFile(&ew.textEditor.editor)
		}
	}

	// Handle Ctrl+N (New)
	for {
		evt, ok := gtx.Event(key.Filter{Name: "N", Required: key.ModCtrl})
		if !ok {
			break
		}
		if e, ok := evt.(key.Event); ok && e.State == key.Press {
			ew.fileOps.NewFile(&ew.textEditor.editor)
		}
	}
}

// Layout renders the editor window
func (ew *EditorWindow) Layout(gtx layout.Context) layout.Dimensions {
	return ew.textEditor.Layout(gtx)
}