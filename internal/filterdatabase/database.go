package filterdatabase

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type FilterData interface {
	FilterFile(path string) bool
	FilterFormat(folderDialog *dialog.FileDialog)
	FormCreate(a fyne.App, title string, lastColumnContent *fyne.Container, rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button, inputEditString, largeEntry *widget.Entry)
}
