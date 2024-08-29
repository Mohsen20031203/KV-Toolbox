package mainwindow

import (
	"fmt"
	"image/color"

	"testgui/internal/logic/mainwindowlagic"
	"testgui/internal/ui/addProjectwindowui"
	buttomaddkvui "testgui/internal/ui/buttomaddKVui"
	buttonrightcolumn "testgui/internal/ui/showkeyvalue/buttonrightcolumn"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var NextButton, PrevButton *widget.Button
var PageLabel *widget.Label

func MainWindow(myApp fyne.App) {

	myWindow := myApp.NewWindow("Non-Scrollable List")

	iconResource := theme.FyneLogo()
	myApp.SetIcon(iconResource)
	myWindow.SetIcon(iconResource)

	spacer := widget.NewLabel("")

	// right column
	rightColumnContent := container.NewVBox()

	line := canvas.NewLine(color.Black)
	line.StrokeWidth = 2

	keyRightColunm := widget.NewButton("key", func() {})
	valueRightColunm := widget.NewButton("value", func() {})
	nameButtonProject := widget.NewLabelWithStyle(
		"",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	searchButton := widget.NewButton("Search", func() {})
	buttonAdd := widget.NewButton("Add", func() {
		buttomaddkvui.OpenWindowAddButton(myApp, rightColumnContent, myWindow)
	})
	buttonAdd.Disable()

	keyAndRight := container.NewGridWithColumns(2, keyRightColunm, valueRightColunm)

	PageLabel = widget.NewLabel(fmt.Sprintf("Page %d", buttonrightcolumn.CurrentPage+1))

	NextButton = widget.NewButton("next", func() {
		buttonrightcolumn.CurrentPage++
		buttonrightcolumn.UpdatePage(rightColumnContent)
	})

	PrevButton = widget.NewButton("prev", func() {
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
		PageLabel,
		layout.NewSpacer(),
	)

	rawSearchAndAdd := container.NewVBox(
		layout.NewSpacer(),
		container.NewGridWithColumns(3, PrevButton, pageLabelposition, NextButton),
		container.NewGridWithColumns(2, searchButton, buttonAdd),
	)

	rightColumnContenttt := container.NewVBox(
		spacer,
		line,
		centeredContainer,
		keyAndRight,
	)

	// left column
	lastColumnContent := mainwindowlagic.SetupLastColumn(rightColumnContent, nameButtonProject, buttonAdd)
	spacer.Resize(fyne.NewSize(0, 30))

	pluss := widget.NewButton("+", func() {
		addProjectwindowui.OpenNewWindow(myApp, "levelDB", lastColumnContent, rightColumnContent, nameButtonProject, buttonAdd)
	})
	lastColumnContentt := container.NewVBox(
		pluss,
		spacer,
	)

	darkLight := mainwindowlagic.SetupThemeButtons(myApp)

	// all window
	containerAll := ColumnContent(rightColumnContent, lastColumnContent, lastColumnContentt, darkLight, rightColumnContenttt, rawSearchAndAdd)
	myWindow.CenterOnScreen()
	myWindow.SetContent(containerAll)
	myWindow.Resize(fyne.NewSize(1200, 800))
	myWindow.ShowAndRun()
}

func LeftColumn(lastColumnContent *fyne.Container, lastColumnContentt *fyne.Container, darkLight *fyne.Container) *fyne.Container {
	lastColumnScrollable := container.NewScroll(lastColumnContent)

	mainContent := container.NewBorder(lastColumnContentt, darkLight, nil, nil, lastColumnScrollable)
	return mainContent
}

func RightColumn(rightColumnContent *fyne.Container, rightColumnContenttt *fyne.Container, rawSearchAndAdd *fyne.Container) fyne.CanvasObject {
	rightColumnScrollable := container.NewVScroll(rightColumnContent)
	mainContent := container.NewBorder(rightColumnContenttt, rawSearchAndAdd, nil, nil, rightColumnScrollable)

	return mainContent
}

func ColumnContent(rightColumnContent *fyne.Container, lastColumnContent *fyne.Container, lastColumnContentt *fyne.Container, darkLight *fyne.Container, rightColumnContenttt *fyne.Container, rawSearchAndAdd *fyne.Container) fyne.CanvasObject {
	mainContent := LeftColumn(lastColumnContent, lastColumnContentt, darkLight)
	rightColumnScrollable := RightColumn(rightColumnContent, rightColumnContenttt, rawSearchAndAdd)
	columns := container.NewHSplit(mainContent, rightColumnScrollable)
	columns.SetOffset(0.25)

	container.NewScroll(columns)
	return columns
}
