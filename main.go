package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
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
		windowAdd := myApp.NewWindow("add Key and Value")
		iputKey := widget.NewEntry()
		iputKey.SetPlaceHolder("Key")
		iputvalue := widget.NewEntry()
		iputvalue.SetPlaceHolder("Value")
		ButtonAddAdd := widget.NewButton("Add", func() {})

		cont := container.NewVBox(
			iputKey,
			iputvalue,
			ButtonAddAdd,
		)
		windowAdd.SetContent(cont)
		windowAdd.Resize(fyne.NewSize(900, 500))
		windowAdd.Show()
	})
	buttonAdd.Disable()
	m := container.NewGridWithColumns(2, buttonAdd, searchButton)

	t := container.NewGridWithColumns(2, keyRightColunm, valueRightColunm)

	lastColumnContent := setupLastColumn(myApp, rightColumnContent, nameButtonProject, buttonAdd)
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

	rightColumnContenttt := container.NewVBox(
		centeredContainer,
		m,
		spacer,
		t,
	)

	darkLight := setupThemeButtons(myApp)

	containerAll := columnContent(rightColumnContent, lastColumnContent, lastColumnContentt, darkLight, rightColumnContenttt)
	myWindow.CenterOnScreen()
	myWindow.SetContent(containerAll)
	myWindow.Resize(fyne.NewSize(1200, 800))
	myWindow.ShowAndRun()
}
