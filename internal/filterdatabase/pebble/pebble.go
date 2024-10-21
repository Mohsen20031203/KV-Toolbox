package Filterpebbledb

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"testgui/internal/filterdatabase"
	sharedfunc "testgui/internal/filterdatabase/SharedFunc"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type NameDatabasePebble struct{}

func NewFileterLeveldb() filterdatabase.FilterData {
	return &NameDatabasePebble{}
}

func (l *NameDatabasePebble) FilterFile(path string) bool {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal("Error reading folder:", err)
		return false
	}
	var count uint8
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "MANIFEST-") || filepath.Ext(file.Name()) == ".log" {
			count++
		}

		if count == 2 {
			return true
		}
	}
	return false
}

func (l *NameDatabasePebble) FilterFormat(folderDialog *dialog.FileDialog) {
	folderDialog.SetFilter(storage.NewExtensionFileFilter([]string{".log"}))
}

func (l *NameDatabasePebble) FormCreate(a fyne.App, title string, lastColumnContent *fyne.Container, rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button) {
	sharedfunc.FormPasteDatabase(a, title, lastColumnContent, rightColumnContentORG, nameButtonProject, buttonAdd)

}
