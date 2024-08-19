package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Non-Scrollable List")

	iconResource := theme.FyneLogo()
	myApp.SetIcon(iconResource)
	myWindow.SetIcon(iconResource)

	rightColumnContent := container.NewVBox()

	keyRightColunm := widget.NewButton("key", func() {})
	valueRightColunm := widget.NewButton("value", func() {})
	nameButtonProject := widget.NewButton("", func() {})

	t := container.NewGridWithColumns(2, keyRightColunm, valueRightColunm)

	lastColumnContent := setupLastColumn(myApp, rightColumnContent, nameButtonProject)
	spacer := widget.NewLabel("")
	spacer.Resize(fyne.NewSize(0, 30))

	pluss := widget.NewButton("+", func() {
		openNewWindow(myApp, "levelDB", lastColumnContent, rightColumnContent, nameButtonProject)
	})
	lastColumnContentt := container.NewVBox(
		pluss,
		spacer,
	)

	rightColumnContenttt := container.NewVBox(
		t,
		spacer,
		nameButtonProject,
		spacer,
	)

	darkLight := setupThemeButtons(myApp)

	containerAll := columnContent(rightColumnContent, lastColumnContent, lastColumnContentt, darkLight, rightColumnContenttt)
	myWindow.CenterOnScreen()
	myWindow.SetContent(containerAll)
	myWindow.Resize(fyne.NewSize(1200, 800))
	myWindow.ShowAndRun()
}
