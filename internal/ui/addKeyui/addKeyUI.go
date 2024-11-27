package addkeyui

import (
	"DatabaseDB/internal/logic"
	"io/ioutil"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func OpenWindowAddButton(myApp fyne.App, rightColumnContent *fyne.Container) {
	windowAdd := myApp.NewWindow("add Key and Value")
	iputKey := widget.NewEntry()
	iputKey.SetPlaceHolder("Key")

	iputvalue := widget.NewMultiLineEntry()
	iputvalue.SetPlaceHolder("value")

	nameFile := widget.NewButton("Name File", nil)

	var valueFinish []byte
	uploadFile := widget.NewButton("UploadFile", func() {
		folderPath := dialog.NewFileOpen(func(dir fyne.URIReadCloser, err error) {
			if err != nil {
				log.Println("Error opening folder:", err)
				return
			}
			if dir == nil {
				log.Println("No folder selected")
				return
			}

			filename := dir.URI().Name()

			valueFinish, err = ioutil.ReadAll(dir)
			if err != nil {
				log.Println(err.Error())
				return
			}

			nameFile.SetText(filename)
			nameFile.Refresh()
		}, windowAdd)
		folderPath.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".gif", ".txt", ".json", ".go", ".html", ".css", ".js"}))
		folderPath.Show()
	})

	uploadFile.Disable()
	iputvalue.Disable()
	nameFile.Disable()

	typeValue := widget.NewLabel("Select the type of file you want")
	redioType := widget.NewRadioGroup([]string{"Text", "File"}, func(typeRedio string) {
		switch typeRedio {
		case "Text":
			iputvalue.Enable()
			uploadFile.Disable()
			nameFile.Disable()
		case "File":
			uploadFile.Enable()
			nameFile.Enable()
			iputvalue.Disable()

		}
	})

	redioType.Horizontal = true
	rowRedio := container.NewHBox(typeValue, redioType)

	columns := container.NewHSplit(uploadFile, nameFile)
	columns.SetOffset(0.80)

	ButtonAddAdd := widget.NewButton("Add", func() {
		if uploadFile.Disabled() {
			valueFinish = []byte(iputvalue.Text)
		}
		logic.AddKeyLogic(iputKey.Text, valueFinish, windowAdd)
	})
	ButtonAddAdd.Importance = widget.HighImportance
	cont := container.NewVBox(
		iputKey,
		rowRedio,
		iputvalue,
		columns,
	)
	m := container.NewBorder(cont, ButtonAddAdd, nil, nil, nil)

	windowAdd.SetContent(m)
	windowAdd.Resize(fyne.NewSize(600, 400))
	windowAdd.Show()
}
