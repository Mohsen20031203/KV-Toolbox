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

var firstKey datebace
var ball bool
var currentPage int
var itemsPerPage = 20
var nextButton, prevButton *widget.Button
var pageLabel *widget.Label // برچسب برای نمایش شماره صفحه

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

	pageLabel = widget.NewLabel(fmt.Sprintf("Page %d", currentPage+1)) // ایجاد برچسب برای شماره صفحه

	nextButton = widget.NewButton("next", func() {
		currentPage++
		updatePage(rightColumnContent)
	})

	prevButton = widget.NewButton("prev", func() {
		if currentPage > 0 {
			currentPage--
			updatePage(rightColumnContent)
		}
	})

	centeredContainer := container.NewHBox(
		layout.NewSpacer(),
		nameButtonProject,
		layout.NewSpacer(),
	)
	pageLabelposition := container.NewHBox(
		layout.NewSpacer(),
		pageLabel,
		layout.NewSpacer(),
	)

	rawSearchAndAdd := container.NewVBox(
		layout.NewSpacer(),
		container.NewGridWithColumns(3, prevButton, pageLabelposition, nextButton), // افزودن شماره صفحه بین دکمه‌ها
		container.NewGridWithColumns(2, searchButton, buttonAdd),
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

func updatePage(rightColumnContent *fyne.Container) {
	if !checkCondition(rightColumnContent) {
		rightColumnContent.Objects = []fyne.CanvasObject{}
		rightColumnContent.Refresh()
	}

	err, data := readDatabace(folderPath)
	if err != nil {
		fmt.Println("Failed to read database:", err)
		return
	}

	startIndex := currentPage * itemsPerPage
	endIndex := startIndex + itemsPerPage

	if endIndex > len(data) {
		endIndex = len(data)
	}

	// پاک کردن محتوای قبلی
	rightColumnContent.Objects = nil

	for _, item := range data[startIndex:endIndex] {
		truncatedKey := truncateString(item.key, 20)
		truncatedValue := truncateString(item.value, 50)

		valueLabel := buidLableKeyAndValue("value", item.key, item.value, truncatedValue, folderPath, rightColumnContent)
		keyLabel := buidLableKeyAndValue("key", item.key, item.value, truncatedKey, folderPath, rightColumnContent)

		buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
		rightColumnContent.Add(buttonRow)
	}

	// به‌روزرسانی شماره صفحه
	pageLabel.SetText(fmt.Sprintf("Page %d", currentPage+1))

	// غیرفعال کردن دکمه‌ها بر اساس موقعیت فعلی
	prevButton.Disable()
	nextButton.Disable()

	if currentPage > 0 {
		prevButton.Enable()
	}
	if endIndex < len(data) {
		nextButton.Enable()
	}

	rightColumnContent.Refresh()
}
