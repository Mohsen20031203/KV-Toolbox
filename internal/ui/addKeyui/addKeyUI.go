package addkeyui

import (
	"log"
	variable "testgui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func OpenWindowAddButton(myApp fyne.App, rightColumnContent *fyne.Container, myWindow fyne.Window) {
	windowAdd := myApp.NewWindow("add Key and Value")
	iputKey := widget.NewEntry()
	iputKey.SetPlaceHolder("Key")
	iputvalue := widget.NewMultiLineEntry()
	iputvalue.SetPlaceHolder("value")
	iputvalue.Resize(fyne.NewSize(500, 500))

	scrollableEntry := container.NewScroll(iputvalue)

	ButtonAddAdd := widget.NewButton("Add", func() {

		if iputKey.Text == "" && iputvalue.Text == "" {
			dialog.ShowInformation("Error", "Please enter both the key and the value", myWindow)
		} else if iputvalue.Text != "" && iputKey.Text == "" {
			dialog.ShowInformation("Error", "You cannot leave either the key or both fields empty.", myWindow)

		}

		err := variable.CurrentDBClient.Open()
		if err != nil {
			return
		}
		defer variable.CurrentDBClient.Close()

		checkNow := variable.CurrentDBClient.Get(iputKey.Text)
		if checkNow != "" {
			dialog.ShowInformation("Error", "This key has already been added to your database", windowAdd)
			return
		}

		err = variable.CurrentDBClient.Add(iputKey.Text, iputvalue.Text)
		if err != nil {
			log.Fatal("error in main window line 172")
		}

		windowAdd.Close()
	})
	cont := container.NewVBox(
		iputKey,
	)
	m := container.NewBorder(cont, ButtonAddAdd, nil, nil, scrollableEntry)

	windowAdd.SetContent(m)
	windowAdd.Resize(fyne.NewSize(900, 500))
	windowAdd.Show()
}
