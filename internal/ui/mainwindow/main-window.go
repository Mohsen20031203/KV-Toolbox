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

var leveldbButton *widget.Button
var BottomDatabase []*widget.Button

func MainWindow(myApp fyne.App) {

	mainWindow := myApp.NewWindow("master")

	iconResource := theme.FyneLogo()
	myApp.SetIcon(iconResource)
	mainWindow.SetIcon(iconResource)

	spacer := widget.NewLabel("")

	// right column
	rightColumnAll := container.NewVBox()

	line := canvas.NewLine(color.Black)
	line.StrokeWidth = 2

	// key top window for colunm keys
	keyRightColunm := widget.NewLabelWithStyle("key", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// value top window for colunm values
	valueRightColunm := widget.NewLabelWithStyle("value", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// column key and value
	keyAndRight := container.NewGridWithColumns(2, keyRightColunm, valueRightColunm)

	// name bottom project in colunm left
	nameButtonProject := widget.NewLabelWithStyle(
		"",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	searchButton := widget.NewButton("Search", func() {

		searchkeyui.SearchKeyUi(rightColumnAll)

	})

	buttonAdd := widget.NewButton("Add", func() {
		addkeyui.OpenWindowAddButton(myApp, rightColumnAll)
	})
	buttonAdd.Disable()

	buttonDelete := widget.NewButton("Delete", func() {
		deletkeyui.DeleteKeyUi(rightColumnAll)
	})

	variable.PageLabel = widget.NewLabel(fmt.Sprintf("Page %d", variable.CurrentPage+1))

	variable.NextButton = widget.NewButton("next", func() {
		variable.CurrentPage++
		variable.PrevButton.Enable()
		logic.UpdatePage(rightColumnAll)
	})
	variable.NextButton.Disable()

	variable.PrevButton = widget.NewButton("prev", func() {
		if variable.CurrentPage > 0 {
			variable.CurrentPage--
			logic.UpdatePage(rightColumnAll)
			variable.NextButton.Enable()
		}
	})

	pageLabelposition := container.NewHBox(
		layout.NewSpacer(),
		variable.PageLabel,
		layout.NewSpacer(),
	)

	rawPrev_Label_Next := container.NewVBox(
		layout.NewSpacer(),
		container.NewGridWithColumns(3, variable.PrevButton, pageLabelposition, variable.NextButton),
	)

	topRightColumn := container.NewVBox(
		nameButtonProject,
		line,
		spacer,
		container.NewGridWithColumns(3, buttonDelete, searchButton, buttonAdd),
		keyAndRight,
	)

	// left column
	leftColumnAll := logic.SetupLastColumn(rightColumnAll, nameButtonProject, buttonAdd)
	spacer.Resize(fyne.NewSize(0, 30))

	for _, m := range variable.NameDatabase {

		leveldbButton = widget.NewButton(m, func() {
			addProjectwindowui.OpenNewWindow(myApp, m, leftColumnAll, rightColumnAll, nameButtonProject, buttonAdd)
		})
		BottomDatabase = append(BottomDatabase, leveldbButton)
	}

	buttonsVisible := false

	toggleButtonsContainer := container.NewVBox()

	pluss := widget.NewButton("+", func() {
		if buttonsVisible {

			toggleButtonsContainer.Objects = nil
		} else {

			for _, m := range BottomDatabase {

				toggleButtonsContainer.Add(m)
			}
		}
		buttonsVisible = !buttonsVisible
		toggleButtonsContainer.Refresh()
	})

	topLeftColumn := container.NewVBox(
		pluss,
		toggleButtonsContainer,
		spacer,
	)

	darkLight := logic.SetupThemeButtons(myApp)

	// all window
	containerAll := ColumnContent(rightColumnAll, leftColumnAll, topLeftColumn, darkLight, topRightColumn, rawPrev_Label_Next)
	mainWindow.CenterOnScreen()
	mainWindow.SetContent(containerAll)
	mainWindow.Resize(fyne.NewSize(1200, 800))
	mainWindow.ShowAndRun()
}

func LeftColumn(leftColumnAll *fyne.Container, topLeftColumn *fyne.Container, darkLight *fyne.Container) *fyne.Container {
	lastColumnScrollable := container.NewScroll(leftColumnAll)

	mainContent := container.NewBorder(topLeftColumn, darkLight, nil, nil, lastColumnScrollable)
	return mainContent
}

func RightColumn(rightColumnAll *fyne.Container, topRightColumn *fyne.Container, rawSearchAndAdd *fyne.Container) fyne.CanvasObject {
	rightColumnScrollable := container.NewVScroll(rightColumnAll)
	mainContent := container.NewBorder(topRightColumn, rawSearchAndAdd, nil, nil, rightColumnScrollable)

	return mainContent
}

func ColumnContent(rightColumnAll *fyne.Container, leftColumnAll *fyne.Container, topLeftColumn *fyne.Container, darkLight *fyne.Container, topRightColumn *fyne.Container, rawSearchAndAdd *fyne.Container) fyne.CanvasObject {

	mainContent := LeftColumn(leftColumnAll, topLeftColumn, darkLight)

	rightColumnScrollable := RightColumn(rightColumnAll, topRightColumn, rawSearchAndAdd)

	columns := container.NewHSplit(mainContent, rightColumnScrollable)
	columns.SetOffset(0.25)

	container.NewScroll(columns)
	return columns
}
