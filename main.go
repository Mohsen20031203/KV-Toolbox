package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var lastKey datebace
var ball bool

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Non-Scrollable List")

	iconResource := theme.FyneLogo()
	myApp.SetIcon(iconResource)
	myWindow.SetIcon(iconResource)

	line := canvas.NewLine(color.Black)
	line.StrokeWidth = 2

	rightColumnContent := container.NewVBox()

	keyRightColunm := widget.NewButton("key", func() {})
	valueRightColunm := widget.NewButton("value", func() {})
	nameButtonProject := widget.NewLabelWithStyle(
		"",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	searchButton := widget.NewButton("Search", func() {})
	buttonAdd := widget.NewButton("Add", func() {
		openWindowAddButton(myApp, rightColumnContent, myWindow)
	})
	buttonAdd.Disable()

	t := container.NewGridWithColumns(2, keyRightColunm, valueRightColunm)

	lastColumnContent := setupLastColumn(rightColumnContent, nameButtonProject, buttonAdd)
	spacer := widget.NewLabel("")
	spacer.Resize(fyne.NewSize(0, 30))

	pluss := widget.NewButton("+", func() {
		openNewWindow(myApp, "levelDB", lastColumnContent, rightColumnContent, nameButtonProject, buttonAdd)
	})
	lastColumnContentt := container.NewVBox(
		pluss,
		spacer,
	)

	nextButton := widget.NewButton("next", func() {
		if !checkCondition(rightColumnContent) {
			newObjects := []fyne.CanvasObject{}

			rightColumnContent.Objects = newObjects

			rightColumnContent.Refresh()
		}

		err, data := readDatabace(folderPath)
		if err != nil {
			fmt.Println("Failed to read database:", err)
			return
		}

		for _, item := range data {

			if lastkey != item && !ball {
				continue
			}
			ball = true

			if count >= 5 {
				count = 0
				ball = false
				break
			}
			lastkey = item
			count++

			if ball {

				truncatedKey := truncateString(item.key, 20)
				truncatedValue := truncateString(item.value, 50)

				valueLabel := buidLableKeyAndValue("value", item.key, item.value, truncatedValue, folderPath, rightColumnContent)
				keyLabel := buidLableKeyAndValue("key", item.key, item.value, truncatedKey, folderPath, rightColumnContent)

				buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
				rightColumnContent.Add(buttonRow)
			}

		}

		rightColumnContent.Refresh()

	})
	preButton := widget.NewButton("prev", func() {
		// چک کردن وضعیت و پاکسازی محتوای ستون راست
		if !checkCondition(rightColumnContent) {
			newObjects := []fyne.CanvasObject{}
			rightColumnContent.Objects = newObjects
			rightColumnContent.Refresh()
		}

		// خواندن داده‌ها از دیتابیس
		err, data := readDatabace(folderPath)
		if err != nil {
			fmt.Println("Failed to read database:", err)
			return
		}

		// پیمایش لیست از انتها به ابتدا
		for i := len(data) - 1; i >= 0; i-- {
			item := data[i]

			// ادامه تا رسیدن به lastkey
			if lastkey != item && !ball {
				continue
			}
			ball = true

			// نمایش حداکثر 5 آیتم و سپس توقف
			if count >= 5 {
				count = 0
				ball = false
				break
			}
			lastkey = item
			count++

			// ساخت و اضافه کردن برچسب‌ها به ستون راست
			truncatedKey := truncateString(item.key, 20)
			truncatedValue := truncateString(item.value, 50)

			valueLabel := buidLableKeyAndValue("value", item.key, item.value, truncatedValue, folderPath, rightColumnContent)
			keyLabel := buidLableKeyAndValue("key", item.key, item.value, truncatedKey, folderPath, rightColumnContent)

			buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
			rightColumnContent.Add(buttonRow)
		}

		// تازه‌سازی محتوای ستون راست
		rightColumnContent.Refresh()
	})

	centeredContainer := container.NewHBox(
		layout.NewSpacer(),
		nameButtonProject,
		layout.NewSpacer(),
	)

	rawSearchAndAdd := container.NewVBox(
		layout.NewSpacer(),
		container.NewGridWithColumns(2, searchButton, buttonAdd),
		container.NewGridWithColumns(2, preButton, nextButton),
	)

	rightColumnContenttt := container.NewVBox(
		spacer,
		line,
		centeredContainer,
		t,
	)

	darkLight := setupThemeButtons(myApp)

	containerAll := columnContent(rightColumnContent, lastColumnContent, lastColumnContentt, darkLight, rightColumnContenttt, rawSearchAndAdd)
	myWindow.CenterOnScreen()
	myWindow.SetContent(containerAll)
	myWindow.Resize(fyne.NewSize(1200, 800))
	myWindow.ShowAndRun()
}
