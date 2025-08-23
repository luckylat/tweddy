package widgets

import (
	"image"
	"log"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"simple-editor/internal/style"
	"simple-editor/pkg/dialog"
)

type TextEditor struct {
	theme      *material.Theme
	tabManager *TabManager
	fileSystem FileSystem
	appTheme   style.Theme
}

func NewTextEditor(theme *material.Theme, appTheme style.Theme) *TextEditor {
	return &TextEditor{
		theme:      theme,
		tabManager: NewTabManager(),
		fileSystem: &RealFileSystem{},
		appTheme:   appTheme,
	}
}

func NewTextEditorWithFS(theme *material.Theme, fs FileSystem, appTheme style.Theme) *TextEditor {
	return &TextEditor{
		theme:      theme,
		tabManager: NewTabManager(),
		fileSystem: fs,
		appTheme:   appTheme,
	}
}

// drawBorder draws a border around the text editor
func (te *TextEditor) drawBorder(gtx layout.Context) {
	borderWidth := unit.Dp(1)
	
	// Top border
	rect := clip.Rect{
		Max: image.Point{
			X: gtx.Constraints.Max.X,
			Y: gtx.Dp(borderWidth),
		},
	}
	defer rect.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: te.appTheme.Application.BorderColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	
	// Left border
	defer op.Offset(image.Point{X: 0, Y: 0}).Push(gtx.Ops).Pop()
	rect = clip.Rect{
		Max: image.Point{
			X: gtx.Dp(borderWidth),
			Y: gtx.Constraints.Max.Y,
		},
	}
	defer rect.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: te.appTheme.Application.BorderColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	
	// Right border
	defer op.Offset(image.Point{X: gtx.Constraints.Max.X - gtx.Dp(borderWidth), Y: 0}).Push(gtx.Ops).Pop()
	rect = clip.Rect{
		Max: image.Point{
			X: gtx.Dp(borderWidth),
			Y: gtx.Constraints.Max.Y,
		},
	}
	defer rect.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: te.appTheme.Application.BorderColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	
	// Bottom border
	defer op.Offset(image.Point{X: 0, Y: gtx.Constraints.Max.Y - gtx.Dp(borderWidth)}).Push(gtx.Ops).Pop()
	rect = clip.Rect{
		Max: image.Point{
			X: gtx.Constraints.Max.X,
			Y: gtx.Dp(borderWidth),
		},
	}
	defer rect.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: te.appTheme.Application.BorderColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

func (te *TextEditor) Layout(gtx layout.Context) layout.Dimensions {
	activeTab := te.tabManager.GetActiveTab()
	if activeTab == nil {
		return layout.Dimensions{}
	}

	// Update tab content from editor
	currentText := activeTab.Editor.Text()
	if currentText != activeTab.Content {
		activeTab.Content = currentText
		activeTab.Modified = true
	}

	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Background{}.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					paint.ColorOp{Color: te.appTheme.Application.EditorBackground}.Add(gtx.Ops)
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
						return material.Editor(te.theme, &activeTab.Editor, "Enter your text here...").Layout(gtx)
					})
				},
			)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			te.drawBorder(gtx)
			return layout.Dimensions{}
		}),
	)
}

func (te *TextEditor) GetTabManager() *TabManager {
	return te.tabManager
}

func (te *TextEditor) HandleKeyboardShortcuts(gtx layout.Context) (newFile, openFile, saveFile bool) {
	for {
		event, ok := gtx.Event(key.Filter{Name: "N"}, key.Filter{Name: "O"}, key.Filter{Name: "S"})
		if !ok {
			break
		}
		if keyEvent, ok := event.(key.Event); ok && keyEvent.State == key.Press {
			switch {
			case keyEvent.Name == "N" && keyEvent.Modifiers.Contain(key.ModCtrl):
				newFile = true
			case keyEvent.Name == "O" && keyEvent.Modifiers.Contain(key.ModCtrl):
				openFile = true
			case keyEvent.Name == "S" && keyEvent.Modifiers.Contain(key.ModCtrl):
				saveFile = true
			}
		}
	}
	return
}

func (te *TextEditor) NewFile() {
	newTab := te.tabManager.NewTab("")
	newTab.Editor.SetText("")
}

func (te *TextEditor) OpenFile() {
	filename, err := dialog.OpenFile()
	if err != nil {
		if err.Error() != "Cancelled" {
			log.Printf("Failed to open file dialog: %v", err)
		}
		return
	}

	content, err := te.fileSystem.ReadFile(filename)
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		return
	}

	// Check if file is already open in a tab
	for _, tab := range te.tabManager.GetTabs() {
		if tab.FilePath == filename {
			te.tabManager.SetActiveTabByID(tab.ID)
			return
		}
	}

	// Create new tab for the file
	newTab := te.tabManager.NewTab(filename)
	newTab.Content = string(content)
	newTab.Editor.SetText(string(content))
	newTab.MarkSaved()
}

func (te *TextEditor) SaveFile() {
	activeTab := te.tabManager.GetActiveTab()
	if activeTab == nil {
		return
	}

	fileName := activeTab.FilePath
	if fileName == "" {
		var err error
		fileName, err = dialog.SaveFile()
		if err != nil {
			if err.Error() != "Cancelled" {
				log.Printf("Failed to open save dialog: %v", err)
			}
			return
		}
		activeTab.FilePath = fileName
		activeTab.FileName = fileName
	}

	content := activeTab.Editor.Text()
	err := te.fileSystem.WriteFile(fileName, []byte(content), 0644)
	if err != nil {
		log.Printf("Failed to save file: %v", err)
		return
	}

	activeTab.Content = content
	activeTab.MarkSaved()
	log.Printf("File saved as %s", fileName)
}

func (te *TextEditor) SaveAsFile() {
	activeTab := te.tabManager.GetActiveTab()
	if activeTab == nil {
		return
	}

	fileName, err := dialog.SaveFile()
	if err != nil {
		if err.Error() != "Cancelled" {
			log.Printf("Failed to open save dialog: %v", err)
		}
		return
	}

	content := activeTab.Editor.Text()
	err = te.fileSystem.WriteFile(fileName, []byte(content), 0644)
	if err != nil {
		log.Printf("Failed to save file: %v", err)
		return
	}

	activeTab.FilePath = fileName
	activeTab.FileName = fileName
	activeTab.Content = content
	activeTab.MarkSaved()
	log.Printf("File saved as %s", fileName)
}
