package main

import (
	"fmt"
	"image/color"
	uimain "testgui/internal/ui"
	buttonrightcolumn "testgui/internal/ui/show-key-value/button-right-column"
	windowaddkeyvalue "testgui/internal/ui/window-add-key-value"
	windowaddproject "testgui/internal/ui/window-add-project"

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
		windowaddkeyvalue.OpenWindowAddButton(myApp, rightColumnContent, myWindow)
	})
	buttonAdd.Disable()

	t := container.NewGridWithColumns(2, keyRightColunm, valueRightColunm)

	lastColumnContent := uimain.SetupLastColumn(rightColumnContent, nameButtonProject, buttonAdd)
	spacer := widget.NewLabel("")
	spacer.Resize(fyne.NewSize(0, 30))

	pluss := widget.NewButton("+", func() {
		windowaddproject.OpenNewWindow(myApp, "levelDB", lastColumnContent, rightColumnContent, nameButtonProject, buttonAdd)
	})
	lastColumnContentt := container.NewVBox(
		pluss,
		spacer,
	)

	buttonrightcolumn.PageLabel = widget.NewLabel(fmt.Sprintf("Page %d", buttonrightcolumn.CurrentPage+1))

	buttonrightcolumn.NextButton = widget.NewButton("next", func() {
		buttonrightcolumn.CurrentPage++
		buttonrightcolumn.UpdatePage(rightColumnContent)
	})

	buttonrightcolumn.PrevButton = widget.NewButton("prev", func() {
		if buttonrightcolumn.CurrentPage > 0 {
			buttonrightcolumn.CurrentPage--
			buttonrightcolumn.UpdatePage(rightColumnContent)
		}
	})

	centeredContainer := container.NewHBox(
		layout.NewSpacer(),
		nameButtonProject,
		layout.NewSpacer(),
	)
	pageLabelposition := container.NewHBox(
		layout.NewSpacer(),
		buttonrightcolumn.PageLabel,
		layout.NewSpacer(),
	)

	rawSearchAndAdd := container.NewVBox(
		layout.NewSpacer(),
		container.NewGridWithColumns(3, buttonrightcolumn.PrevButton, pageLabelposition, buttonrightcolumn.NextButton), // افزودن شماره صفحه بین دکمه‌ها
		container.NewGridWithColumns(2, searchButton, buttonAdd),
	)

	rightColumnContenttt := container.NewVBox(
		spacer,
		line,
		centeredContainer,
		t,
	)

	darkLight := uimain.SetupThemeButtons(myApp)

	containerAll := uimain.ColumnContent(rightColumnContent, lastColumnContent, lastColumnContentt, darkLight, rightColumnContenttt, rawSearchAndAdd)
	myWindow.CenterOnScreen()
	myWindow.SetContent(containerAll)
	myWindow.Resize(fyne.NewSize(1200, 800))
	myWindow.ShowAndRun()
}

/*
package main

import (
	"testgui/internal/ui"

	"fyne.io/fyne/v2/widget"
)

var currentPage int
var itemsPerPage = 20
var nextButton, prevButton *widget.Button
var pageLabel *widget.Label // برچسب برای نمایش شماره صفحه

func main() {
	ui.NewUi()
}

*/
