package deletkeyui

import (
	"testgui/internal/logic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func DeleteKeyUi(rightColumnContent *fyne.Container) {
	editWindow := fyne.CurrentApp().NewWindow("Delete in the database")
	editWindow.Resize(fyne.NewSize(600, 300))

	valueEntry := widget.NewMultiLineEntry()
	valueEntry.Resize(fyne.NewSize(500, 500))
	valueEntry.SetPlaceHolder("Key for Delete")

	buttomDelete := widget.NewButton("Delete", func() {

		logic.DeleteKeyLogic(valueEntry, editWindow, rightColumnContent)

	})

	editContent := container.NewVBox(
		widget.NewLabel("Enter the desired key"),
		valueEntry,
		layout.NewSpacer(),
		buttomDelete,
	)
	editWindow.SetContent(editContent)
	editWindow.Show()
}
