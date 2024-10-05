package mainwindow

import (
	"fmt"
	"image/color"
	variable "testgui"

	"testgui/internal/logic"
	addkeyui "testgui/internal/ui/addKeyui"
	"testgui/internal/ui/addProjectwindowui"
	deletkeyui "testgui/internal/ui/deletKeyUi"
	searchkeyui "testgui/internal/ui/searchKeyui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func MainWindow(myApp fyne.App) {

	myWindow := myApp.NewWindow("master")

	iconResource := theme.FyneLogo()
	myApp.SetIcon(iconResource)
	myWindow.SetIcon(iconResource)

	spacer := widget.NewLabel("")

	// right column
	rightColumnContent := container.NewVBox()

	line := canvas.NewLine(color.Black)
	line.StrokeWidth = 2

	keyRightColunm := widget.NewLabelWithStyle("key", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	valueRightColunm := widget.NewLabelWithStyle("value", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	nameButtonProject := widget.NewLabelWithStyle(
		"",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	searchButton := widget.NewButton("Search", func() {

		searchkeyui.SearchKeyUi(rightColumnContent)

	})

	buttonAdd := widget.NewButton("Add", func() {
		addkeyui.OpenWindowAddButton(myApp, rightColumnContent)
	})
	buttonAdd.Disable()

	buttonDelete := widget.NewButton("Delete", func() {
		deletkeyui.DeleteKeyUi(rightColumnContent)
	})

	keyAndRight := container.NewGridWithColumns(2, keyRightColunm, valueRightColunm)

	variable.PageLabel = widget.NewLabel(fmt.Sprintf("Page %d", variable.CurrentPage+1))

	variable.NextButton = widget.NewButton("next", func() {
		variable.CurrentPage++
		variable.PrevButton.Enable()
		logic.UpdatePage(rightColumnContent)
	})
	variable.NextButton.Disable()

	variable.PrevButton = widget.NewButton("prev", func() {
		if variable.CurrentPage > 0 {
			variable.CurrentPage--
			logic.UpdatePage(rightColumnContent)
			variable.NextButton.Enable()
		}
	})

	pageLabelposition := container.NewHBox(
		layout.NewSpacer(),
		variable.PageLabel,
		layout.NewSpacer(),
	)

	rawSearchAndAdd := container.NewVBox(
		layout.NewSpacer(),
		container.NewGridWithColumns(3, variable.PrevButton, pageLabelposition, variable.NextButton),
	)

	rightColumnContenttt := container.NewVBox(
		nameButtonProject,
		line,
		spacer,
		container.NewGridWithColumns(3, buttonDelete, searchButton, buttonAdd),
		keyAndRight,
	)

	// left column
	lastColumnContent := logic.SetupLastColumn(rightColumnContent, nameButtonProject, buttonAdd)
	spacer.Resize(fyne.NewSize(0, 30))

	leveldbButton := widget.NewButton("levelDB", func() {
		addProjectwindowui.OpenNewWindow(myApp, "levelDB", lastColumnContent, rightColumnContent, nameButtonProject, buttonAdd)
	})

	radisButton := widget.NewButton("Pebble", func() {
		addProjectwindowui.OpenNewWindow(myApp, "Pebble", lastColumnContent, rightColumnContent, nameButtonProject, buttonAdd)
	})
	buttonsVisible := false

	toggleButtonsContainer := container.NewVBox()

	pluss := widget.NewButton("+", func() {
		if buttonsVisible {

			toggleButtonsContainer.Objects = nil
		} else {

			toggleButtonsContainer.Add(radisButton)
			toggleButtonsContainer.Add(leveldbButton)
		}
		buttonsVisible = !buttonsVisible
		toggleButtonsContainer.Refresh()
	})

	lastColumnContentt := container.NewVBox(
		pluss,
		toggleButtonsContainer,
		spacer,
	)

	darkLight := logic.SetupThemeButtons(myApp)

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
