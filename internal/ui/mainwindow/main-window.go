package mainwindow

import (
	variable "DatabaseDB"
	"image/color"

	Filterbadger "DatabaseDB/internal/filterdatabase/badger"
	FilterLeveldb "DatabaseDB/internal/filterdatabase/leveldb"
	Filterpebbledb "DatabaseDB/internal/filterdatabase/pebble"
	"DatabaseDB/internal/logic"
	addkeyui "DatabaseDB/internal/ui/addKeyui"
	deletkeyui "DatabaseDB/internal/ui/deletKeyUi"
	searchkeyui "DatabaseDB/internal/ui/searchKeyui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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

	// right column show key
	rightColumnAll := container.NewVBox()

	// right column Edit
	rightColumEdit := container.NewVBox()

	line := canvas.NewLine(color.Black)
	line.StrokeWidth = 2

	// key top window for colunm keys
	keyRightColunm := widget.NewLabelWithStyle("key", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// value top window for colunm values
	valueRightColunm := widget.NewLabelWithStyle("value", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// value top window for colunm values
	editRightColunm := widget.NewLabelWithStyle("Edit", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// column key and value
	keyAndRight := container.NewGridWithColumns(3, keyRightColunm, valueRightColunm, editRightColunm)

	// name bottom project in colunm right
	nameButtonProject := widget.NewLabelWithStyle(
		"",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	largeEntry := widget.NewMultiLineEntry()

	largeEntry.Wrapping = fyne.TextWrapWord

	scroll := container.NewVScroll(largeEntry)
	scroll.SetMinSize(fyne.NewSize(400, 400))

	saveEditKey := widget.NewButton("Save", nil)
	cancelEditKey := widget.NewButton("Cancle", nil)

	saveAndCancle := container.NewGridWithColumns(2, saveEditKey, cancelEditKey)

	rightColumEdit = container.NewVBox(
		widget.NewLabel("key"),
		widget.NewLabel("String"),
		widget.NewEntry(),
		widget.NewLabel("Byte"),
		scroll,
	)

	columnEdit := container.NewBorder(nil, saveAndCancle, nil, nil, rightColumEdit)

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

			switch m {
			case "levelDB":
				variable.NameData = FilterLeveldb.NewFileterLeveldb()
			case "Pebble":
				variable.NameData = Filterpebbledb.NewFileterLeveldb()
			case "Badger":
				variable.NameData = Filterbadger.NewFileterBadger()
				//case "Redis":
				//	variable.NameData = Filterredis.NewFileterRedis()

			}
			variable.NameData.FormCreate(myApp, m, leftColumnAll, rightColumnAll, nameButtonProject, buttonAdd)
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
	containerAll := ColumnContent(rightColumnAll, columnEdit, leftColumnAll, topLeftColumn, darkLight, topRightColumn)
	mainWindow.CenterOnScreen()
	mainWindow.SetContent(containerAll)
	mainWindow.Resize(fyne.NewSize(1300, 800))
	mainWindow.ShowAndRun()
}

func LeftColumn(leftColumnAll *fyne.Container, topLeftColumn *fyne.Container, darkLight *fyne.Container) *fyne.Container {
	lastColumnScrollable := container.NewScroll(leftColumnAll)

	mainContent := container.NewBorder(topLeftColumn, darkLight, nil, nil, lastColumnScrollable)
	return mainContent
}

func RightColumn(rightColumnAll *fyne.Container, topRightColumn *fyne.Container, rightColumEdit *fyne.Container) fyne.CanvasObject {
	rightColumnScrollable := container.NewVScroll(rightColumnAll)

	up := false

	rightColumnScrollable.OnScrolled = func(p fyne.Position) {
		maxScroll := rightColumnAll.MinSize().Height - rightColumnScrollable.Size().Height

		if up && p.Y == 0 {
			variable.CurrentPage--
			if variable.CurrentPage < 3 {
				up = false
				variable.CurrentPage = 3
				return
			}
			numberLast := len(rightColumnAll.Objects)
			logic.UpdatePage(rightColumnAll)

			rightColumnAll.Objects = rightColumnAll.Objects[:numberLast]

			rightColumnScrollable.Offset.Y = maxScroll / 2
			rightColumnScrollable.Refresh()

		} else if p.Y == maxScroll && !variable.ItemsAdded {
			return
		} else if p.Y == maxScroll && variable.ItemsAdded {

			variable.CurrentPage++
			numberLast := len(rightColumnAll.Objects)
			logic.UpdatePage(rightColumnAll)
			rightColumnScrollable.Offset.Y = maxScroll / 2

			if len(rightColumnAll.Objects) > (variable.ItemsPerPage)*3 {
				rightColumnAll.Objects = rightColumnAll.Objects[len(rightColumnAll.Objects)-numberLast:]
				up = true
			}

		}

	}

	columns := container.NewHSplit(rightColumnScrollable, rightColumEdit)
	columns.SetOffset(0.65)
	mainContent := container.NewBorder(topRightColumn, nil, nil, nil, columns)

	return mainContent
}

func ColumnContent(rightColumnAll *fyne.Container, rightColumEdit *fyne.Container, leftColumnAll *fyne.Container, topLeftColumn *fyne.Container, darkLight *fyne.Container, topRightColumn *fyne.Container) fyne.CanvasObject {

	mainContent := LeftColumn(leftColumnAll, topLeftColumn, darkLight)

	rightColumnScrollable := RightColumn(rightColumnAll, topRightColumn, rightColumEdit)

	columns := container.NewHSplit(mainContent, rightColumnScrollable)
	columns.SetOffset(0.20)

	container.NewScroll(columns)
	return columns
}
