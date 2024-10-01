package searchkeyui

import (
	"testgui/internal/logic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func SearchKeyUi(rightColumnContent *fyne.Container) {
	editWindow := fyne.CurrentApp().NewWindow("Enter the desired key")
	editWindow.Resize(fyne.NewSize(600, 300))

	valueEntry := widget.NewMultiLineEntry()
	valueEntry.Resize(fyne.NewSize(500, 500))

	buttomSearch := widget.NewButton("Search", func() {

		logic.SearchDatabase(valueEntry, editWindow, rightColumnContent)
		editWindow.Close()

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
