package searchkeyui

import (
	"DatabaseDB/internal/logic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func SearchKeyUi(rightColumnContent *fyne.Container, inputEditString, largeEntry *widget.Entry) {
	editWindow := fyne.CurrentApp().NewWindow("Search in the database")
	editWindow.Resize(fyne.NewSize(600, 300))

	valueEntry := widget.NewMultiLineEntry()
	valueEntry.Resize(fyne.NewSize(500, 500))
	valueEntry.SetPlaceHolder("Key for Search")

	buttomSearch := widget.NewButton("Search", func() {

		result, _ := logic.SearchDatabase(valueEntry, editWindow, rightColumnContent, inputEditString, largeEntry)
		if !result {
			dialog.ShowInformation("Error", "Such a key is not available in the database", editWindow)
		}

	})

	editContent := container.NewVBox(
		widget.NewLabel("Enter the desired key"),
		valueEntry,
		layout.NewSpacer(),
		buttomSearch,
	)
	editWindow.SetContent(editContent)
	editWindow.Show()
}
