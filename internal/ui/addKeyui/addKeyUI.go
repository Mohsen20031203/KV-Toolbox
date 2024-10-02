package addkeyui

import (
	"testgui/internal/logic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func OpenWindowAddButton(myApp fyne.App, rightColumnContent *fyne.Container) {
	windowAdd := myApp.NewWindow("add Key and Value")
	iputKey := widget.NewEntry()
	iputKey.SetPlaceHolder("Key")
	iputvalue := widget.NewMultiLineEntry()
	iputvalue.SetPlaceHolder("value")
	iputvalue.Resize(fyne.NewSize(500, 500))

	scrollableEntry := container.NewScroll(iputvalue)

	ButtonAddAdd := widget.NewButton("Add", func() {

		logic.AddKeyLogic(iputKey, iputvalue, windowAdd)
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
