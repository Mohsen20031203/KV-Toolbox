package buttomaddkvui

import (
	addprojectWL "testgui/internal/logic/addProject-window-logic"
	"testgui/internal/logic/mainwindowlagic"
	"testgui/internal/utils"
	leveldbb "testgui/pkg/db/leveldb"

	"fyne.io/fyne/container"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func OpenWindowAddButton(myApp fyne.App, rightColumnContent *fyne.Container, myWindow fyne.Window) {
	var db *leveldbb.ConstantDatabase
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

		truncatedKey := utils.TruncateString(iputKey.Text, 20)
		truncatedValue := utils.TruncateString(iputvalue.Text, 50)

		err := mainwindowlagic.CurrentDBClient.Open()
		if err != nil {
			return
		}
		defer db.DB.Close()
		db.DB.Put([]byte(iputKey.Text), []byte(iputvalue.Text), nil)

		valueLabel := addprojectWL.BuidLableKeyAndValue("value", iputKey.Text, iputvalue.Text, truncatedValue, mainwl.FolderPath, rightColumnContent)
		keyLabel := addprojectWL.BuidLableKeyAndValue("key", iputKey.Text, iputvalue.Text, truncatedKey, mainwl.FolderPath, rightColumnContent)

		buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
		rightColumnContent.Add(buttonRow)
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
