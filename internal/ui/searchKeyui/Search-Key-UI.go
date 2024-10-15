package searchkeyui

import (
	"testgui/internal/logic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func SearchKeyUi(rightColumnContent *fyne.Container) {
	editWindow := fyne.CurrentApp().NewWindow("Search in the database")
	editWindow.Resize(fyne.NewSize(600, 300))

	valueEntry := widget.NewMultiLineEntry()
	valueEntry.Resize(fyne.NewSize(500, 500))
	valueEntry.SetPlaceHolder("Key for Search")

	buttomSearch := widget.NewButton("Search", func() {

		result, _ := logic.SearchDatabase(valueEntry, editWindow, rightColumnContent)
		if result {
			dialog.ShowInformation("Error", "Such a key is not available in the database", editWindow)
		}

	})

	/*valueEntry.OnChanged = func(s string) {
		if len(s) > 0 && s[len(s)-1] == '\n' {

			valueEntry.SetText(s[:len(s)-1])

			buttomSearch.OnTapped()
		}
	}*/

	editContent := container.NewVBox(
		widget.NewLabel("Enter the desired key"),
		valueEntry,
		layout.NewSpacer(),
		buttomSearch,
	)
	editWindow.SetContent(editContent)
	editWindow.Show()
}
