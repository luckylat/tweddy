package widgets

import (
	"fmt"
	"strconv"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// TabManager handles tab creation, deletion, and switching operations
type TabManager struct {
	theme     *material.Theme
	tabs      []*Tab
	activeTab int
	newTabBtn widget.Clickable
}

// NewTabManager creates a new tab manager
func NewTabManager(theme *material.Theme) *TabManager {
	tm := &TabManager{
		theme:     theme,
		tabs:      make([]*Tab, 0),
		activeTab: 0,
	}
	
	// Create initial tab
	tm.AddNewTab()
	
	return tm
}

// AddNewTab creates and adds a new tab
func (tm *TabManager) AddNewTab() {
	tab := &Tab{
		TextEditor: NewTextEditor(tm.theme),
		Title:      fmt.Sprintf("Untitled-%d", len(tm.tabs)+1),
		FilePath:   "",
		Modified:   false,
	}
	
	tm.tabs = append(tm.tabs, tab)
	tm.activeTab = len(tm.tabs) - 1
}

// CloseTab removes a tab at the specified index
func (tm *TabManager) CloseTab(index int) {
	if len(tm.tabs) <= 1 {
		// Don't close the last tab, just create a new one
		tm.tabs[0] = &Tab{
			TextEditor: NewTextEditor(tm.theme),
			Title:      "Untitled-1",
			FilePath:   "",
			Modified:   false,
		}
		return
	}
	
	// Remove tab
	tm.tabs = append(tm.tabs[:index], tm.tabs[index+1:]...)
	
	// Adjust active tab index
	if tm.activeTab >= len(tm.tabs) {
		tm.activeTab = len(tm.tabs) - 1
	} else if tm.activeTab > index {
		tm.activeTab--
	}
}

// SwitchToTab changes the active tab
func (tm *TabManager) SwitchToTab(index int) {
	if index >= 0 && index < len(tm.tabs) {
		tm.activeTab = index
	}
}

// GetActiveTab returns the currently active tab
func (tm *TabManager) GetActiveTab() *Tab {
	if len(tm.tabs) == 0 {
		return nil
	}
	return tm.tabs[tm.activeTab]
}

// GetTabs returns all tabs
func (tm *TabManager) GetTabs() []*Tab {
	return tm.tabs
}

// GetActiveTabIndex returns the index of the active tab
func (tm *TabManager) GetActiveTabIndex() int {
	return tm.activeTab
}

// HandleEvents processes keyboard shortcuts for tab management
func (tm *TabManager) HandleEvents(gtx layout.Context) {
	// Handle tab switching with Ctrl+1-9
	for i := 0; i < 9 && i < len(tm.tabs); i++ {
		keyName := strconv.Itoa(i + 1)
		for {
			evt, ok := gtx.Event(key.Filter{Name: key.Name(keyName), Required: key.ModCtrl})
			if !ok {
				break
			}
			if e, ok := evt.(key.Event); ok && e.State == key.Press {
				tm.SwitchToTab(i)
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
			tm.AddNewTab()
		}
	}
	
	// Handle Ctrl+W (Close Tab)
	for {
		evt, ok := gtx.Event(key.Filter{Name: key.Name("W"), Required: key.ModCtrl})
		if !ok {
			break
		}
		if e, ok := evt.(key.Event); ok && e.State == key.Press {
			tm.CloseTab(tm.activeTab)
		}
	}
	
	// Handle new tab button click
	if tm.newTabBtn.Clicked(gtx) {
		tm.AddNewTab()
	}
	
	// Handle tab clicks for switching
	for i, tab := range tm.tabs {
		if tab.TabBtn.Clicked(gtx) {
			tm.SwitchToTab(i)
		}
	}
	
	// Handle tab close button clicks
	for i, tab := range tm.tabs {
		if tab.CloseBtn.Clicked(gtx) {
			tm.CloseTab(i)
			break
		}
	}
}

// GetNewTabButton returns the new tab button for layout purposes
func (tm *TabManager) GetNewTabButton() *widget.Clickable {
	return &tm.newTabBtn
}