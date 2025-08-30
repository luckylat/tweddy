package widgets

import (
	"fmt"
	"image"
	"path/filepath"
	"strconv"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// Tab represents a single tab containing a text editor
type Tab struct {
	TextEditor *TextEditor
	Title      string
	FilePath   string
	Modified   bool
	CloseBtn   widget.Clickable
	TabBtn     widget.Clickable
}

// TabbedTextEditor manages multiple text editor tabs
type TabbedTextEditor struct {
	theme     *material.Theme
	tabs      []*Tab
	activeTab int
	newTabBtn widget.Clickable
}

// NewTabbedTextEditor creates a new tabbed text editor with one initial tab
func NewTabbedTextEditor(theme *material.Theme) *TabbedTextEditor {
	tte := &TabbedTextEditor{
		theme:     theme,
		tabs:      make([]*Tab, 0),
		activeTab: 0,
	}
	
	// Create initial tab
	tte.AddNewTab()
	
	return tte
}

// AddNewTab creates and adds a new tab
func (tte *TabbedTextEditor) AddNewTab() {
	tab := &Tab{
		TextEditor: NewTextEditor(tte.theme),
		Title:      fmt.Sprintf("Untitled-%d", len(tte.tabs)+1),
		FilePath:   "",
		Modified:   false,
	}
	
	tte.tabs = append(tte.tabs, tab)
	tte.activeTab = len(tte.tabs) - 1
}

// CloseTab removes a tab at the specified index
func (tte *TabbedTextEditor) CloseTab(index int) {
	if len(tte.tabs) <= 1 {
		// Don't close the last tab, just create a new one
		tte.tabs[0] = &Tab{
			TextEditor: NewTextEditor(tte.theme),
			Title:      "Untitled-1",
			FilePath:   "",
			Modified:   false,
		}
		return
	}
	
	// Remove tab
	tte.tabs = append(tte.tabs[:index], tte.tabs[index+1:]...)
	
	// Adjust active tab index
	if tte.activeTab >= len(tte.tabs) {
		tte.activeTab = len(tte.tabs) - 1
	} else if tte.activeTab > index {
		tte.activeTab--
	}
}

// SwitchToTab changes the active tab
func (tte *TabbedTextEditor) SwitchToTab(index int) {
	if index >= 0 && index < len(tte.tabs) {
		tte.activeTab = index
	}
}

// GetActiveTab returns the currently active tab
func (tte *TabbedTextEditor) GetActiveTab() *Tab {
	if len(tte.tabs) == 0 {
		return nil
	}
	return tte.tabs[tte.activeTab]
}

// HandleEvents processes keyboard shortcuts and tab interactions
func (tte *TabbedTextEditor) HandleEvents(gtx layout.Context) {
	// Handle tab switching with Ctrl+1-9
	for i := 0; i < 9 && i < len(tte.tabs); i++ {
		keyName := strconv.Itoa(i + 1)
		for {
			evt, ok := gtx.Event(key.Filter{Name: key.Name(keyName), Required: key.ModCtrl})
			if !ok {
				break
			}
			if e, ok := evt.(key.Event); ok && e.State == key.Press {
				tte.SwitchToTab(i)
			}
		}
	}
	
	// Handle Ctrl+T (New Tab)
	for {
		evt, ok := gtx.Event(key.Filter{Name: key.Name("T"), Required: key.ModCtrl})
		if !ok {
			break
		}
		if e, ok := evt.(key.Event); ok && e.State == key.Press {
			tte.AddNewTab()
		}
	}
	
	// Handle Ctrl+W (Close Tab)
	for {
		evt, ok := gtx.Event(key.Filter{Name: key.Name("W"), Required: key.ModCtrl})
		if !ok {
			break
		}
		if e, ok := evt.(key.Event); ok && e.State == key.Press {
			tte.CloseTab(tte.activeTab)
		}
	}
	
	// Handle new tab button click
	if tte.newTabBtn.Clicked(gtx) {
		tte.AddNewTab()
	}
	
	// Handle tab clicks for switching
	for i, tab := range tte.tabs {
		if tab.TabBtn.Clicked(gtx) {
			tte.SwitchToTab(i)
		}
	}
	
	// Handle tab close button clicks
	for i, tab := range tte.tabs {
		if tab.CloseBtn.Clicked(gtx) {
			tte.CloseTab(i)
			break
		}
	}
	
	// Handle events for the active tab
	if activeTab := tte.GetActiveTab(); activeTab != nil {
		activeTab.TextEditor.HandleEvents(gtx)
		
		// Update tab title based on file operations
		currentFile := activeTab.TextEditor.fileOps.GetCurrentFile()
		if currentFile != "" && currentFile != activeTab.FilePath {
			activeTab.FilePath = currentFile
			activeTab.Title = filepath.Base(currentFile)
		}
	}
}

// layoutTabBar renders the tab bar at the top
func (tte *TabbedTextEditor) layoutTabBar(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		// Tab buttons
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			if len(tte.tabs) == 0 {
				return layout.Dimensions{}
			}
			
			children := make([]layout.FlexChild, 0, len(tte.tabs))
			for i, tab := range tte.tabs {
				tabIndex := i // Capture for closure
				tabRef := tab  // Capture for closure
				children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return tte.layoutTab(gtx, tabRef, tabIndex)
				}))
			}
			
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, children...)
		}),
		// New tab button
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(tte.theme, &tte.newTabBtn, "+")
			btn.Background = tte.theme.Bg
			return btn.Layout(gtx)
		}),
	)
}

// layoutTab renders a single tab button
func (tte *TabbedTextEditor) layoutTab(gtx layout.Context, tab *Tab, index int) layout.Dimensions {
	isActive := index == tte.activeTab
	
	return material.ButtonLayout(tte.theme, &tab.TabBtn).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		// Set background color based on active state
		bgColor := tte.theme.Bg
		if isActive {
			bgColor = tte.theme.ContrastBg
		}
		
		// Draw background
		defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
		paint.Fill(gtx.Ops, bgColor)
		
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			// Tab title
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				title := tab.Title
				if tab.Modified {
					title += "*"
				}
				label := material.Label(tte.theme, unit.Sp(14), title)
				if isActive {
					label.Color = tte.theme.ContrastFg
				} else {
					label.Color = tte.theme.Fg
				}
				return layout.Inset{
					Top: unit.Dp(8), Bottom: unit.Dp(8),
					Left: unit.Dp(12), Right: unit.Dp(4),
				}.Layout(gtx, label.Layout)
			}),
			// Close button (only show if more than one tab)
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if len(tte.tabs) <= 1 {
					return layout.Dimensions{}
				}
				
				closeBtn := material.Button(tte.theme, &tab.CloseBtn, "×")
				closeBtn.Background = bgColor
				closeBtn.Color = tte.theme.Fg
				return layout.Inset{Right: unit.Dp(4)}.Layout(gtx, closeBtn.Layout)
			}),
		)
	})
}

// Layout renders the complete tabbed text editor
func (tte *TabbedTextEditor) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Tab bar
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return tte.layoutTabBar(gtx)
		}),
		// Separator line
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{
				Size: image.Pt(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(1))),
			}
		}),
		// Active tab content
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			if activeTab := tte.GetActiveTab(); activeTab != nil {
				return activeTab.TextEditor.Layout(gtx)
			}
			return layout.Dimensions{}
		}),
	)
}