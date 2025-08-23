package widgets

import (
	"path/filepath"

	"gioui.org/widget"
)

// Tab represents a single file tab
type Tab struct {
	ID       int
	FileName string
	FilePath string
	Content  string
	Modified bool
	Editor   widget.Editor
}

// NewTab creates a new tab with given file information
func NewTab(id int, filePath string) *Tab {
	fileName := "Untitled"
	if filePath != "" {
		fileName = filepath.Base(filePath)
	}

	tab := &Tab{
		ID:       id,
		FileName: fileName,
		FilePath: filePath,
		Content:  "",
		Modified: false,
	}

	return tab
}

// GetDisplayName returns the tab display name with modification indicator
func (t *Tab) GetDisplayName() string {
	name := t.FileName
	if t.Modified {
		name += "*"
	}
	return name
}

// UpdateContent updates the tab content and marks as modified if changed
func (t *Tab) UpdateContent(content string) {
	if t.Content != content {
		t.Content = content
		t.Modified = true
		t.Editor.SetText(content)
	}
}

// MarkSaved marks the tab as saved (not modified)
func (t *Tab) MarkSaved() {
	t.Modified = false
}

// TabManager manages multiple tabs with browser-like behavior
type TabManager struct {
	tabs         []*Tab
	activeTab    int
	nextID       int
	accessHistory []int  // Stack of recently accessed tab IDs
	maxTabs      int     // Maximum number of tabs (0 = unlimited)
}

// NewTabManager creates a new tab manager with browser-like behavior
func NewTabManager() *TabManager {
	tm := &TabManager{
		tabs:         make([]*Tab, 0),
		activeTab:    -1,
		nextID:       1,
		accessHistory: make([]int, 0),
		maxTabs:      20, // Default limit of 20 tabs like most browsers
	}

	// Create initial untitled tab
	tm.NewTab("")

	return tm
}

// NewTab creates a new tab with browser-like positioning
func (tm *TabManager) NewTab(filePath string) *Tab {
	// Check tab limit
	if tm.maxTabs > 0 && len(tm.tabs) >= tm.maxTabs {
		// Close oldest tab if we're at the limit
		tm.closeOldestTab()
	}
	
	tab := NewTab(tm.nextID, filePath)
	tm.nextID++
	
	// Insert new tab next to current active tab (browser behavior)
	insertPos := tm.activeTab + 1
	if insertPos > len(tm.tabs) {
		insertPos = len(tm.tabs)
	}
	
	// Insert at position
	tm.tabs = append(tm.tabs[:insertPos], append([]*Tab{tab}, tm.tabs[insertPos:]...)...)
	tm.activeTab = insertPos
	tm.updateAccessHistory(tab.ID)
	
	return tab
}

// GetActiveTab returns the currently active tab
func (tm *TabManager) GetActiveTab() *Tab {
	if tm.activeTab >= 0 && tm.activeTab < len(tm.tabs) {
		return tm.tabs[tm.activeTab]
	}
	return nil
}

// GetTabs returns all tabs
func (tm *TabManager) GetTabs() []*Tab {
	return tm.tabs
}

// SetActiveTab sets the active tab by index and updates access history
func (tm *TabManager) SetActiveTab(index int) {
	if index >= 0 && index < len(tm.tabs) {
		tm.activeTab = index
		tm.updateAccessHistory(tm.tabs[index].ID)
	}
}

// SetActiveTabByID sets the active tab by ID and updates access history
func (tm *TabManager) SetActiveTabByID(id int) {
	for i, tab := range tm.tabs {
		if tab.ID == id {
			tm.activeTab = i
			tm.updateAccessHistory(id)
			break
		}
	}
}

// CloseTab closes a tab by index with browser-like behavior
func (tm *TabManager) CloseTab(index int) {
	if index >= 0 && index < len(tm.tabs) {
		closingActiveTab := (index == tm.activeTab)
		closedTabID := tm.tabs[index].ID
		
		// Remove tab from slice
		tm.tabs = append(tm.tabs[:index], tm.tabs[index+1:]...)
		
		// Remove from access history
		tm.removeFromAccessHistory(closedTabID)
		
		// If no tabs left, create a new untitled tab
		if len(tm.tabs) == 0 {
			tm.NewTab("")
			return
		}
		
		// Browser-like tab switching behavior
		if closingActiveTab {
			// Switch to most recently accessed tab
			tm.switchToMostRecentTab()
		} else if tm.activeTab > index {
			// Adjust active tab index if needed
			tm.activeTab--
		}
	}
}

// CloseTabByID closes a tab by ID
func (tm *TabManager) CloseTabByID(id int) {
	for i, tab := range tm.tabs {
		if tab.ID == id {
			tm.CloseTab(i)
			break
		}
	}
}

// GetActiveTabIndex returns the index of the active tab
func (tm *TabManager) GetActiveTabIndex() int {
	return tm.activeTab
}

// HasUnsavedChanges returns true if any tab has unsaved changes
func (tm *TabManager) HasUnsavedChanges() bool {
	for _, tab := range tm.tabs {
		if tab.Modified {
			return true
		}
	}
	return false
}

// updateAccessHistory updates the access history stack
func (tm *TabManager) updateAccessHistory(tabID int) {
	// Remove existing entry if present
	tm.removeFromAccessHistory(tabID)
	
	// Add to front of history
	tm.accessHistory = append([]int{tabID}, tm.accessHistory...)
	
	// Limit history size to prevent memory growth
	maxHistory := 10
	if len(tm.accessHistory) > maxHistory {
		tm.accessHistory = tm.accessHistory[:maxHistory]
	}
}

// removeFromAccessHistory removes a tab ID from access history
func (tm *TabManager) removeFromAccessHistory(tabID int) {
	for i, id := range tm.accessHistory {
		if id == tabID {
			tm.accessHistory = append(tm.accessHistory[:i], tm.accessHistory[i+1:]...)
			break
		}
	}
}

// switchToMostRecentTab switches to the most recently accessed tab
func (tm *TabManager) switchToMostRecentTab() {
	// Find the most recent tab that still exists
	for _, historyID := range tm.accessHistory {
		for i, tab := range tm.tabs {
			if tab.ID == historyID {
				tm.activeTab = i
				return
			}
		}
	}
	
	// Fallback: if no history or no valid tabs in history, use adjacent tab
	if tm.activeTab >= len(tm.tabs) {
		tm.activeTab = len(tm.tabs) - 1
	}
	if tm.activeTab < 0 && len(tm.tabs) > 0 {
		tm.activeTab = 0
	}
}

// closeOldestTab closes the least recently accessed tab
func (tm *TabManager) closeOldestTab() {
	if len(tm.tabs) == 0 {
		return
	}
	
	// Find the oldest tab (last in access history or first in tabs if no history)
	var oldestIndex int = 0
	
	if len(tm.accessHistory) > 0 {
		// Find tab that's been accessed longest ago
		for i := len(tm.accessHistory) - 1; i >= 0; i-- {
			historyID := tm.accessHistory[i]
			for j, tab := range tm.tabs {
				if tab.ID == historyID {
					oldestIndex = j
					goto foundOldest
				}
			}
		}
	}
	
	foundOldest:
	tm.CloseTab(oldestIndex)
}

// MoveTab moves a tab from one position to another (for drag and drop)
func (tm *TabManager) MoveTab(fromIndex, toIndex int) {
	if fromIndex < 0 || fromIndex >= len(tm.tabs) || toIndex < 0 || toIndex >= len(tm.tabs) {
		return
	}
	
	if fromIndex == toIndex {
		return
	}
	
	// Remove tab from current position
	tab := tm.tabs[fromIndex]
	tm.tabs = append(tm.tabs[:fromIndex], tm.tabs[fromIndex+1:]...)
	
	// Adjust toIndex if needed
	if fromIndex < toIndex {
		toIndex--
	}
	
	// Insert at new position
	tm.tabs = append(tm.tabs[:toIndex], append([]*Tab{tab}, tm.tabs[toIndex:]...)...)
	
	// Update active tab index
	if tm.activeTab == fromIndex {
		tm.activeTab = toIndex
	} else if tm.activeTab > fromIndex && tm.activeTab <= toIndex {
		tm.activeTab--
	} else if tm.activeTab < fromIndex && tm.activeTab >= toIndex {
		tm.activeTab++
	}
}

// GetTabCount returns the number of open tabs
func (tm *TabManager) GetTabCount() int {
	return len(tm.tabs)
}

// SetMaxTabs sets the maximum number of tabs (0 = unlimited)
func (tm *TabManager) SetMaxTabs(max int) {
	tm.maxTabs = max
}
