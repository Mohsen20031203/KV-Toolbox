package deletkeyui

import (
	"DatabaseDB/internal/logic"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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

		message := fmt.Sprintf("Are you sure you want to delete the key: _ %s _?", valueEntry.Text)

		dialog.ShowConfirm("Confirm Delete", message,
			func(response bool) {
				if response {
					logic.DeleteKeyLogic(valueEntry, editWindow, rightColumnContent)
				} else {

				}
			}, editWindow)

	})
	buttomDelete.Importance = widget.HighImportance
	editContent := container.NewVBox(
		widget.NewLabel("Enter the desired key"),
		valueEntry,
		layout.NewSpacer(),
		buttomDelete,
	)
	editWindow.SetContent(editContent)
	editWindow.Show()
}
