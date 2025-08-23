package interaction

import (
	"fmt"
	"os"
	"gioui.org/widget"
	"simple-editor/pkg/dialog"
)

// FileOperations handles file-related operations for the editor
type FileOperations struct {
	currentFile string
}

// NewFileOperations creates a new FileOperations instance
func NewFileOperations() *FileOperations {
	return &FileOperations{}
}

// NewFile clears the editor content and resets the current file
func (fo *FileOperations) NewFile(editor *widget.Editor) {
	editor.SetText("")
	fo.currentFile = ""
	fmt.Println("New file created")
}

// SaveFile saves the current editor content
func (fo *FileOperations) SaveFile(editor *widget.Editor) error {
	content := editor.Text()
	
	if fo.currentFile == "" {
		// If no current file, prompt for save location
		return fo.SaveAsFile(editor, "")
	}
	
	err := os.WriteFile(fo.currentFile, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error saving file: %v\n", err)
		return err
	}
	
	fmt.Printf("File saved: %s\n", fo.currentFile)
	return nil
}

// OpenFile opens a file using OS file dialog and loads its content into the editor
func (fo *FileOperations) OpenFile(editor *widget.Editor) error {
	filename, err := dialog.OpenFile()
	if err != nil {
		fmt.Printf("Error opening file dialog: %v\n", err)
		return err
	}
	
	if filename == "" {
		// User cancelled the dialog
		fmt.Println("Open file cancelled")
		return nil
	}
	
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return err
	}
	
	editor.SetText(string(content))
	fo.currentFile = filename
	fmt.Printf("File opened: %s\n", filename)
	return nil
}

// SaveAsFile saves the current editor content using OS file dialog
func (fo *FileOperations) SaveAsFile(editor *widget.Editor, suggestedFilename string) error {
	content := editor.Text()
	
	filename, err := dialog.SaveFile()
	if err != nil {
		fmt.Printf("Error opening save dialog: %v\n", err)
		return err
	}
	
	if filename == "" {
		// User cancelled the dialog
		fmt.Println("Save file cancelled")
		return nil
	}
	
	err = os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error saving file: %v\n", err)
		return err
	}
	
	fo.currentFile = filename
	fmt.Printf("File saved as: %s\n", filename)
	return nil
}

// GetCurrentFile returns the current file path
func (fo *FileOperations) GetCurrentFile() string {
	return fo.currentFile
}