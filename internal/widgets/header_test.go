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

func TestNewHeader(t *testing.T) {
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	header := NewHeader(theme, style.DefaultTheme())

	if header == nil {
		t.Fatal("NewHeader() returned nil")
	}

	if header.theme != theme {
		t.Error("Header theme not set correctly")
	}
}

func TestHeaderLayout(t *testing.T) {
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	header := NewHeader(theme, style.DefaultTheme())

	var ops op.Ops
	gtx := layout.Context{
		Ops:         &ops,
		Now:         time.Now(),
		Metric:      unit.Metric{},
		Constraints: layout.Exact(image.Pt(800, 600)),
	}

	// Test that Layout doesn't panic and returns valid dimensions
	dims := header.Layout(gtx)

	if dims.Size.X <= 0 || dims.Size.Y <= 0 {
		t.Errorf("Header layout returned invalid dimensions: %v", dims.Size)
	}
}

func TestHeaderButtonClicks(t *testing.T) {
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	header := NewHeader(theme, style.DefaultTheme())

	var ops op.Ops
	gtx := layout.Context{
		Ops:         &ops,
		Now:         time.Now(),
		Metric:      unit.Metric{},
		Constraints: layout.Exact(image.Pt(800, 600)),
	}

	// Test that button click methods don't panic and return false by default
	if header.NewClicked(gtx) {
		t.Error("NewClicked should return false when not clicked")
	}

	if header.OpenClicked(gtx) {
		t.Error("OpenClicked should return false when not clicked")
	}

	if header.SaveClicked(gtx) {
		t.Error("SaveClicked should return false when not clicked")
	}

	if header.SaveAsClicked(gtx) {
		t.Error("SaveAsClicked should return false when not clicked")
	}
}
