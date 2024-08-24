package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

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

	nextButton := widget.NewButton("next", func() {
		handleProjectSelection(folderPath, rightColumnContent, buttonAdd)
	})
	preButton := widget.NewButton("prev", func() {
		handleProjectSelection(folderPath, rightColumnContent, buttonAdd)
	})

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

	centeredContainer := container.NewHBox(
		layout.NewSpacer(), // Spacer برای قرار دادن فاصله در بالا
		nameButtonProject,  // لیبل
		layout.NewSpacer(), // Spacer برای قرار دادن فاصله در پایین
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
