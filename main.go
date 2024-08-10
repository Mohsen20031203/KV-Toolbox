package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"

	/*
		"fyne.io/fyne/v2/canvas"
		"fyne.io/fyne/v2/dialog"*/
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func hasManifestFile(folderPath string) bool {
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		fmt.Println("Error reading folder:", err)
		return false
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "MANIFEST-") {
			return true
		}
	}
	return false
}

func openNewWindow(a fyne.App, title string) {
	newWindow := a.NewWindow(title)

	createSeparator := func() *canvas.Line {
		line := canvas.NewLine(color.Black)
		line.StrokeWidth = 1
		return line
	}
	line1 := createSeparator()

	lableName := widget.NewLabel("Name :")
	pathEntry := widget.NewEntry()
	pathEntry.PlaceHolder = "Name"
	nameContent := container.NewBorder(nil, nil, lableName, nil, pathEntry)

	lableComment := widget.NewLabel("Commert :")
	pathEntryComment := widget.NewEntry()
	pathEntryComment.PlaceHolder = "Commint"
	commentContent := container.NewBorder(nil, nil, lableComment, nil, pathEntryComment)

	pathEntry2 := widget.NewEntry()
	pathEntry2.SetPlaceHolder("No folder selected")

	openButton := widget.NewButton("Open Folder", func() {

		folderDialog := dialog.NewFileOpen(func(dir fyne.URIReadCloser, err error) {
			if err != nil {
				fmt.Println("Error opening folder:", err)
				return
			}
			if dir == nil {
				fmt.Println("No folder selected")
				return
			}
			filePath := dir.URI().Path()

			// استخراج پوشه از مسیر فایل
			folderPath := filepath.Dir(filePath)

			if hasManifestFile(folderPath) {
				pathEntry2.SetText(folderPath)
			} else {
				dialog.ShowInformation("Invalid Folder", "The selected folder does not contain a valid LevelDB manifest file.", newWindow)
			}

		}, newWindow)
		folderDialog.Show()
		folderDialog.SetFilter(storage.NewExtensionFileFilter([]string{".log"}))
	})

	buttonCancel := widget.NewButton("Cancel", func() {
		fmt.Println("buttonCancel")
	})
	buttonApply := widget.NewButton("Apply", func() {
		fmt.Println("buttonApply")
	})
	buttonOk := widget.NewButton("Ok", func() {
		fmt.Println("buttonOk")
	})

	rowBotton := container.NewVBox(
		layout.NewSpacer(),
		container.NewGridWithColumns(3, buttonCancel, buttonApply, buttonOk),
	)

	rightColumnContent := container.NewVBox(
		layout.NewSpacer(), // Add space here
		nameContent,
		layout.NewSpacer(), // Add space here
		commentContent,
		layout.NewSpacer(), // Add space here
		line1,
		layout.NewSpacer(), // Add space here
		line1,
		line1,
		layout.NewSpacer(), // Add space here
		layout.NewSpacer(), // Add space here
		pathEntry2,
		layout.NewSpacer(), // Add space here
		openButton,
		layout.NewSpacer(), // Add space here
		layout.NewSpacer(), // Add space here
		layout.NewSpacer(), // Add space here
		rowBotton,
	)
	lastColumnContent := container.NewVBox()

	columns := container.NewHSplit(lastColumnContent, rightColumnContent)
	columns.SetOffset(0.25)
	newWindow.Resize(fyne.NewSize(1000, 600))
	newWindow.CenterOnScreen()
	newWindow.SetContent(columns)
	newWindow.Show()
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Non-Scrollable List")

	iconResource := theme.FyneLogo()

	myApp.SetIcon(iconResource)

	myWindow.SetIcon(iconResource)

	//---------------------------------------------- right -----------------------------------------------
	/*
		createSeparator := func() *canvas.Line {
			line := canvas.NewLine(color.Black)
			line.StrokeWidth = 1
			return line
		}

		// --- Items -----
		item1 := widget.NewLabel("Item 1")
		separator1 := createSeparator()

		item2 := widget.NewLabel("Item 2")

		// --- Select -----
		combo := widget.NewSelect([]string{"Option 1", "Option 2"}, func(value string) {
			log.Println("Select set to", value)
		})

		// --- Select And Type -----
		options := []string{"Option 1", "Option 2", "Option 3"}
		selectEntry := widget.NewSelectEntry(options)

		selectEntry.SetText("Option 2")

		selectEntry.OnChanged = func(s string) {
			fmt.Println("Changed to:", s)
		}

		// --- Check Box Normal -----
		check := widget.NewCheck("Check", func(value bool) {
			fmt.Println("Check : ", value)
		})

		// --- Check Box Disable -----
		disableCheck := widget.NewCheck("disableCheck", func(value bool) {
			fmt.Println("disableCheckprint :", value)
		})
		disableCheck.Disable()

		// --- Check Box Group -----
		option := []string{"Mashhad", "Tehran", "Esfahan"}
		checkBoxes := []*widget.Check{}
		for _, m := range option {

			groupCheck := widget.NewCheck(m, func(value bool) {
				fmt.Println("Group Check Box : ", value)
			})
			checkBoxes = append(checkBoxes, groupCheck)
		}

		horizontalContainer := container.NewHBox()
		for _, checkBox := range checkBoxes {
			horizontalContainer.Add(checkBox)
		}

		// --- Redio Group -----
		RadionOption := []string{"Football", "Basketball", "Baceball"}
		var radioButtons []fyne.CanvasObject
		for _, option := range RadionOption {
			radio := widget.NewRadioGroup([]string{option}, func(value string) {
				fmt.Println("Radio : ", value)
			})
			radioButtons = append(radioButtons, radio)
		}
		radioContainer := container.NewGridWithColumns(len(radioButtons), radioButtons...)

		// --- Radio Disable -----
		disableRadioItem := []string{"wrestling", "swimming"}
		var radioDisableButtons []fyne.CanvasObject

		for _, h := range disableRadioItem {

			disableRadio := widget.NewRadioGroup([]string{h}, func(sport string) {
				fmt.Println("disable : ", sport)
			})

			disableRadio.Disable()
			radioDisableButtons = append(radioDisableButtons, disableRadio)

		}
		radioContainerDisable := container.NewGridWithColumns(len(radioDisableButtons)*2, radioDisableButtons...)

		// --- Navbar Slider -----

		Slider := widget.NewSlider(0, 100)
		Slider.SetValue(20)

		// --- Select File -----

		pathEntry := widget.NewEntry()
		pathEntry.SetPlaceHolder("No folder selected")

		openButton := widget.NewButton("Open Folder", func() {
			folderDialog := dialog.NewFolderOpen(func(dir fyne.ListableURI, err error) {
				if err != nil {
					fmt.Println("Error opening folder:", err)
					return
				}
				if dir == nil {
					fmt.Println("No folder selected")
					return
				}
				pathEntry.SetText(dir.Path())
			}, myWindow)
			folderDialog.Show()
		})


	*/
	//---------------------------------------------- left -----------------------------------------------
	// --- buttonPlus -----

	plus := widget.NewSelect([]string{"levelDB", "DynamoDB", "SQL", "DucomentDB", "PostgreSQL", "IBM Information Management System", "Integrated Data Store", "ObjectDB", "Apache Cassandra"}, func(value string) {
		log.Println("Select set to", value)
		openNewWindow(myApp, value)
	})
	plus.PlaceHolder = "Data Source"      // تغییر متن Placeholder
	plus.Alignment = fyne.TextAlignCenter // مرکز کردن متن

	buttonOne := widget.NewButton("DDL data source", func() {
		fmt.Println("buttonOne")
	})
	buttonTwo := widget.NewButton("Data source from URL", func() {
		fmt.Println("buttonTwo")
	})
	buttonThree := widget.NewButton("Data source from path", func() {
		fmt.Println("buttonThree")
	})
	buttonFor := widget.NewButton("Import from Clipboard", func() {
		fmt.Println("buttonFor")
	})

	// --- Bottuns Light And Dark -----

	darkBottum := widget.NewButton("Dark", func() {
		myApp.Settings().SetTheme(theme.DarkTheme())
	})

	lightBottum := widget.NewButton("Light", func() {
		myApp.Settings().SetTheme(theme.LightTheme())
	})

	darkLight := container.NewVBox(
		layout.NewSpacer(),
		container.NewGridWithColumns(2, lightBottum, darkBottum),
	)

	// --- Colunm Right -----
	rightColumnContent := container.NewVBox(
	/*item1,
	separator1,
	item2,
	combo,
	selectEntry,
	check,
	disableCheck,
	horizontalContainer,
	radioContainer,
	radioContainerDisable,
	Slider,
	pathEntry,
	openButton,*/
	)

	// --- Colunm Left -----

	lastColumnContent := container.NewVBox(
		plus,
		buttonOne,
		buttonTwo,
		buttonThree,
		buttonFor,
		layout.NewSpacer(), // Add space here

		darkLight,
	)

	columns := container.NewHSplit(lastColumnContent, rightColumnContent)
	columns.SetOffset(0.25)

	scrol := container.NewScroll(columns)
	myWindow.CenterOnScreen()
	myWindow.SetContent(scrol)
	myWindow.Resize(fyne.NewSize(1200, 800))
	myWindow.ShowAndRun()
}
