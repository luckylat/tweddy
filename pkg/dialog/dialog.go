package dialog

import "github.com/sqweek/dialog"

func OpenFile() (string, error) {
	return dialog.File().Load()
}

func SaveFile() (string, error) {
	return dialog.File().Title("Save File").Save()
}
