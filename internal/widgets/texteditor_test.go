package widgets

import (
	"image"
	"testing"
	"time"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"simple-editor/internal/style"
)

func TestNewTextEditor(t *testing.T) {
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	editor := NewTextEditor(theme, style.DefaultTheme())

	if editor == nil {
		t.Fatal("NewTextEditor() returned nil")
	}

	if editor.theme != theme {
		t.Error("TextEditor theme not set correctly")
	}
}

func TestTextEditorLayout(t *testing.T) {
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	editor := NewTextEditor(theme, style.DefaultTheme())

	var ops op.Ops
	gtx := layout.Context{
		Ops:         &ops,
		Now:         time.Now(),
		Metric:      unit.Metric{},
		Constraints: layout.Exact(image.Pt(800, 600)),
	}

	// Test that Layout doesn't panic and returns valid dimensions
	dims := editor.Layout(gtx)

	if dims.Size.X <= 0 || dims.Size.Y <= 0 {
		t.Errorf("TextEditor layout returned invalid dimensions: %v", dims.Size)
	}
}

func TestTextEditorNewFile(t *testing.T) {
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	editor := NewTextEditor(theme, style.DefaultTheme())

	// Set some text in active tab
	activeTab := editor.GetTabManager().GetActiveTab()
	activeTab.Editor.SetText("test content")
	activeTab.FilePath = "test.txt"

	// Call NewFile
	editor.NewFile()

	// Verify new tab was created and is active
	newActiveTab := editor.GetTabManager().GetActiveTab()
	if newActiveTab.Editor.Text() != "" {
		t.Error("NewFile() should create a new tab with empty content")
	}

	if newActiveTab.FilePath != "" {
		t.Error("NewFile() should create a new tab with empty filepath")
	}
}

func TestTextEditorHandleKeyboardShortcuts(t *testing.T) {
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	editor := NewTextEditor(theme, style.DefaultTheme())

	var ops op.Ops
	gtx := layout.Context{
		Ops:         &ops,
		Now:         time.Now(),
		Metric:      unit.Metric{},
		Constraints: layout.Exact(image.Pt(800, 600)),
	}

	// Test that HandleKeyboardShortcuts doesn't panic
	newFile, openFile, saveFile := editor.HandleKeyboardShortcuts(gtx)

	// Without actual key events, all should be false
	if newFile || openFile || saveFile {
		t.Error("HandleKeyboardShortcuts should return false for all actions when no key events")
	}
}

func TestTextEditorSaveFileWithVirtualFS(t *testing.T) {
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	// Create text editor with virtual filesystem
	fs := NewTestFileSystem()
	editor := NewTextEditorWithFS(theme, fs, style.DefaultTheme())

	// Set content and filename in active tab
	testContent := "Hello, World!"
	testFile := "test.txt"
	activeTab := editor.GetTabManager().GetActiveTab()
	activeTab.Editor.SetText(testContent)
	activeTab.FilePath = testFile

	// Save the file
	editor.SaveFile()

	// Verify the file was created and has correct content
	content, err := fs.ReadFile(testFile)
	if err != nil {
		t.Errorf("Failed to read saved file: %v", err)
		return
	}

	if string(content) != testContent {
		t.Errorf("Saved file content mismatch. Expected: %s, Got: %s", testContent, string(content))
	}
}

func TestTextEditorOpenFileWithVirtualFS(t *testing.T) {
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	// Create text editor with virtual filesystem
	fs := NewTestFileSystem()
	editor := NewTextEditorWithFS(theme, fs, style.DefaultTheme())

	// Create a test file in virtual filesystem
	testFile := "test.txt"
	testContent := "Virtual file content"
	err := fs.CreateFile(testFile, testContent)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Manually test file reading without dialog (since OpenFile uses dialog)
	content, err := fs.ReadFile(testFile)
	if err != nil {
		t.Errorf("Failed to read file: %v", err)
		return
	}

	// Create a new tab and set content
	newTab := editor.GetTabManager().NewTab(testFile)
	newTab.Content = string(content)
	newTab.Editor.SetText(string(content))

	// Verify content was loaded
	if newTab.Editor.Text() != testContent {
		t.Errorf("File content mismatch. Expected: %s, Got: %s", testContent, newTab.Editor.Text())
	}
}

func TestTextEditorSaveAsFileWithoutDialog(t *testing.T) {
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	editor := NewTextEditor(theme, style.DefaultTheme())

	// Set content in active tab
	testContent := "Save As Test Content"
	activeTab := editor.GetTabManager().GetActiveTab()
	activeTab.Editor.SetText(testContent)

	// Note: SaveAsFile() will try to open a dialog, which will fail in test environment
	// This test mainly verifies the method doesn't panic
	editor.SaveAsFile()

	// The dialog will fail, so filename should remain empty
	if activeTab.FilePath != "" {
		t.Error("SaveAsFile() should not set filename when dialog fails")
	}
}
