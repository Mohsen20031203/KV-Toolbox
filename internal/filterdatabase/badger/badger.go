package Filterbadger

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"testgui/internal/filterdatabase"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

type NameDatabaseBadger struct {
	name string
}

func NewFileterBadger() filterdatabase.FilterData {
	return &NameDatabaseBadger{}
}

func (l *NameDatabaseBadger) FilterFile(path string) bool {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal("Error reading folder:", err)
		return false
	}
	var count uint8
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sst" || filepath.Ext(file.Name()) == ".vlog" {
			count++
		}

		if count == 2 {
			return true
		}
	}
	return false
}

func (l *NameDatabaseBadger) FilterFormat(folderDialog *dialog.FileDialog) {
	folderDialog.SetFilter(storage.NewExtensionFileFilter([]string{".sst", ".vlog"}))
}
