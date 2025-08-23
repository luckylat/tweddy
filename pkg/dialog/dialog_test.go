package dialog

import (
	"testing"
)

func TestOpenFile(t *testing.T) {
	// Note: OpenFile() will try to open a native dialog, which will fail in test environment
	// This test mainly verifies the function doesn't panic and returns appropriate error
	
	filename, err := OpenFile()
	
	// In test environment without GUI, this should return an error or cancelled
	if err == nil && filename == "" {
		t.Error("OpenFile() should return error or filename")
	}
	
	// If error is returned, it should be meaningful
	if err != nil && err.Error() == "" {
		t.Error("OpenFile() error should have a message")
	}
}

func TestSaveFile(t *testing.T) {
	// Note: SaveFile() will try to open a native dialog, which will fail in test environment
	// This test mainly verifies the function doesn't panic and returns appropriate error
	
	filename, err := SaveFile()
	
	// In test environment without GUI, this should return an error or cancelled
	if err == nil && filename == "" {
		t.Error("SaveFile() should return error or filename")
	}
	
	// If error is returned, it should be meaningful
	if err != nil && err.Error() == "" {
		t.Error("SaveFile() error should have a message")
	}
}

func TestDialogFunctions(t *testing.T) {
	// Test that both functions exist and can be called
	// This is mainly a compilation test
	
	_, err1 := OpenFile()
	_, err2 := SaveFile()
	
	// Both should return errors in test environment
	if err1 == nil {
		t.Log("OpenFile() succeeded unexpectedly (might be in GUI environment)")
	}
	
	if err2 == nil {
		t.Log("SaveFile() succeeded unexpectedly (might be in GUI environment)")
	}
}