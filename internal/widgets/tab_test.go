package widgets

import (
	"testing"
)

func TestNewTab(t *testing.T) {
	tab := NewTab(1, "/path/to/file.txt")

	if tab.ID != 1 {
		t.Errorf("Expected tab ID to be 1, got %d", tab.ID)
	}

	if tab.FileName != "file.txt" {
		t.Errorf("Expected filename to be 'file.txt', got '%s'", tab.FileName)
	}

	if tab.FilePath != "/path/to/file.txt" {
		t.Errorf("Expected filepath to be '/path/to/file.txt', got '%s'", tab.FilePath)
	}

	if tab.Modified {
		t.Error("New tab should not be marked as modified")
	}
}

func TestNewTabUntitled(t *testing.T) {
	tab := NewTab(1, "")

	if tab.FileName != "Untitled" {
		t.Errorf("Expected filename to be 'Untitled', got '%s'", tab.FileName)
	}
}

func TestTabGetDisplayName(t *testing.T) {
	tab := NewTab(1, "test.txt")

	// Test unmodified
	if tab.GetDisplayName() != "test.txt" {
		t.Errorf("Expected display name to be 'test.txt', got '%s'", tab.GetDisplayName())
	}

	// Test modified
	tab.Modified = true
	if tab.GetDisplayName() != "test.txt*" {
		t.Errorf("Expected display name to be 'test.txt*', got '%s'", tab.GetDisplayName())
	}
}

func TestTabUpdateContent(t *testing.T) {
	tab := NewTab(1, "test.txt")

	// Update content
	tab.UpdateContent("new content")

	if tab.Content != "new content" {
		t.Errorf("Expected content to be 'new content', got '%s'", tab.Content)
	}

	if !tab.Modified {
		t.Error("Tab should be marked as modified after content update")
	}
}

func TestTabMarkSaved(t *testing.T) {
	tab := NewTab(1, "test.txt")
	tab.Modified = true

	tab.MarkSaved()

	if tab.Modified {
		t.Error("Tab should not be marked as modified after MarkSaved()")
	}
}

func TestNewTabManager(t *testing.T) {
	tm := NewTabManager()

	if len(tm.GetTabs()) != 1 {
		t.Errorf("Expected 1 initial tab, got %d", len(tm.GetTabs()))
	}

	if tm.GetActiveTabIndex() != 0 {
		t.Errorf("Expected active tab index to be 0, got %d", tm.GetActiveTabIndex())
	}
}

func TestTabManagerNewTab(t *testing.T) {
	tm := NewTabManager()
	initialCount := len(tm.GetTabs())

	newTab := tm.NewTab("test.txt")

	if len(tm.GetTabs()) != initialCount+1 {
		t.Errorf("Expected %d tabs, got %d", initialCount+1, len(tm.GetTabs()))
	}

	if newTab.FileName != "test.txt" {
		t.Errorf("Expected filename to be 'test.txt', got '%s'", newTab.FileName)
	}

	if tm.GetActiveTabIndex() != len(tm.GetTabs())-1 {
		t.Error("New tab should be active")
	}
}

func TestTabManagerSetActiveTab(t *testing.T) {
	tm := NewTabManager()
	tm.NewTab("test1.txt")
	tm.NewTab("test2.txt")

	tm.SetActiveTab(0)

	if tm.GetActiveTabIndex() != 0 {
		t.Errorf("Expected active tab index to be 0, got %d", tm.GetActiveTabIndex())
	}
}

func TestTabManagerCloseTab(t *testing.T) {
	tm := NewTabManager()
	tm.NewTab("test1.txt")
	tm.NewTab("test2.txt")

	initialCount := len(tm.GetTabs())
	tm.CloseTab(1)

	if len(tm.GetTabs()) != initialCount-1 {
		t.Errorf("Expected %d tabs after closing, got %d", initialCount-1, len(tm.GetTabs()))
	}
}

func TestTabManagerCloseLastTab(t *testing.T) {
	tm := NewTabManager()

	// Close the only tab
	tm.CloseTab(0)

	// Should create a new untitled tab
	if len(tm.GetTabs()) != 1 {
		t.Errorf("Expected 1 tab after closing last tab, got %d", len(tm.GetTabs()))
	}
}

func TestTabManagerHasUnsavedChanges(t *testing.T) {
	tm := NewTabManager()

	if tm.HasUnsavedChanges() {
		t.Error("New tab manager should not have unsaved changes")
	}

	// Mark a tab as modified
	activeTab := tm.GetActiveTab()
	activeTab.Modified = true

	if !tm.HasUnsavedChanges() {
		t.Error("Tab manager should have unsaved changes after modifying a tab")
	}
}
