package mainwindow

import (
	"fmt"
	"image/color"
	"log"
	variable "testgui"

	"testgui/internal/logic"
	"testgui/internal/ui/addProjectwindowui"
	deletkeyui "testgui/internal/ui/deletKeyUi"
	searchkeyui "testgui/internal/ui/searchKeyui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

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

	searchButton := widget.NewButton("Search", func() {

		searchkeyui.SearchKeyUi(rightColumnContent)

	})

	buttonAdd := widget.NewButton("Add", func() {
		OpenWindowAddButton(myApp, rightColumnContent, myWindow)
	})
	buttonAdd.Disable()

	buttonDelet := widget.NewButton("Delet", func() {
		deletkeyui.DeleteKeyUi(rightColumnContent)
	})

	keyAndRight := container.NewGridWithColumns(2, keyRightColunm, valueRightColunm)

	variable.PageLabel = widget.NewLabel(fmt.Sprintf("Page %d", variable.CurrentPage+1))

	variable.NextButton = widget.NewButton("next", func() {
		variable.CurrentPage++
		variable.PrevButton.Enable()
		logic.UpdatePage(rightColumnContent)
	})

	variable.PrevButton = widget.NewButton("prev", func() {
		if variable.CurrentPage > 0 {
			variable.CurrentPage--
			logic.UpdatePage(rightColumnContent)
			variable.NextButton.Enable()
		}
	})

	centeredContainer := container.NewHBox(
		layout.NewSpacer(),
		nameButtonProject,
		layout.NewSpacer(),
	)
	pageLabelposition := container.NewHBox(
		layout.NewSpacer(),
		variable.PageLabel,
		layout.NewSpacer(),
	)

	rawSearchAndAdd := container.NewVBox(
		layout.NewSpacer(),
		container.NewGridWithColumns(3, variable.PrevButton, pageLabelposition, variable.NextButton),
		container.NewGridWithColumns(3, buttonDelet, searchButton, buttonAdd),
	)

	rightColumnContenttt := container.NewVBox(
		spacer,
		line,
		centeredContainer,
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

func OpenWindowAddButton(myApp fyne.App, rightColumnContent *fyne.Container, myWindow fyne.Window) {
	windowAdd := myApp.NewWindow("add Key and Value")
	iputKey := widget.NewEntry()
	iputKey.SetPlaceHolder("Key")
	iputvalue := widget.NewMultiLineEntry()
	iputvalue.SetPlaceHolder("value")
	iputvalue.Resize(fyne.NewSize(500, 500))

	scrollableEntry := container.NewScroll(iputvalue)

	ButtonAddAdd := widget.NewButton("Add", func() {

		if iputKey.Text == "" && iputvalue.Text == "" {
			dialog.ShowInformation("Error", "Please enter both the key and the value", myWindow)
		} else if iputvalue.Text != "" && iputKey.Text == "" {
			dialog.ShowInformation("Error", "You cannot leave either the key or both fields empty.", myWindow)

		}

		err := variable.CurrentDBClient.Open()
		if err != nil {
			return
		}
		defer variable.CurrentDBClient.Close()

		err = variable.CurrentDBClient.Add(iputKey.Text, iputvalue.Text)
		if err != nil {
			log.Fatal("error in main window line 172")
		}

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
