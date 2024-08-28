package windowaddkeyvalue

import (
	"testgui/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/syndtr/goleveldb/leveldb"
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

		truncatedKey := utils.TruncateString(iputKey.Text, 20)
		truncatedValue := utils.TruncateString(iputvalue.Text, 50)

		db, err := leveldb.OpenFile(folderPath, nil)
		if err != nil {
			return
		}
		defer db.Close()
		db.Put([]byte(iputKey.Text), []byte(iputvalue.Text), nil)

		valueLabel := buidLableKeyAndValue("value", iputKey.Text, iputvalue.Text, truncatedValue, folderPath, rightColumnContent)
		keyLabel := buidLableKeyAndValue("key", iputKey.Text, iputvalue.Text, truncatedKey, folderPath, rightColumnContent)

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
