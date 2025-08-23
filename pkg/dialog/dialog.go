package dialog

import "github.com/sqweek/dialog"

func OpenFile() (string, error) {
	return dialog.File().Filter("Text files", "txt").Filter("All files", "*").Load()
}

func SaveFile() (string, error) {
	return dialog.File().Filter("Text files", "txt").Filter("All files", "*").Title("Save File").Save()
}