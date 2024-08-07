package main

import (
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Non-Scrollable List")

	iconResource := theme.FyneLogo()

	myApp.SetIcon(iconResource)

	myWindow.SetIcon(iconResource)

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

	// --- Bottuns Light And Dark -----

	darkBottum := widget.NewButton("Dark", func() {
		myApp.Settings().SetTheme(theme.DarkTheme())
	})

	lightBottum := widget.NewButton("Light", func() {
		myApp.Settings().SetTheme(theme.LightTheme())
	})

	darkLight := container.NewGridWithColumns(2, lightBottum, darkBottum)

	// --- buttonPlus -----

	plusButton := widget.NewButton("+", func() {
		// وقتی دکمه کلیک شد، منو باز شود
		item1 := fyne.NewMenuItem("Data Source", nil)
		item2 := fyne.NewMenuItem("Details", nil)
		item3 := fyne.NewMenuItem("Home", nil)
		item4 := fyne.NewMenuItem("Run", nil)

		// منوی فرزند برای آیتم اول
		item1.ChildMenu = fyne.NewMenu(
			"", // برچسب خالی
			fyne.NewMenuItem("LevelBD", func() {
				filter := storage.NewExtensionFileFilter([]string{".ldb"})

				fileDialog := dialog.NewFileOpen(func(uri fyne.URIReadCloser, err error) {
					if err != nil || uri == nil {
						return
					}
					defer uri.Close()
					// انجام عملیات با فایل انتخاب شده
					filePath := uri.URI().Path()
					dialog.ShowInformation("Selected File", filePath, myWindow)
				}, myWindow)

				fileDialog.SetFilter(filter)
				fileDialog.Show()
			}), // آیتم‌های منوی فرزند
			fyne.NewMenuItem("MongoDB", nil),
			fyne.NewMenuItem("MYSQL", nil),
		)

		// منوی اصلی
		mainMenu := fyne.NewMainMenu(
			fyne.NewMenu("File", item1, item2, item3, item4),
			fyne.NewMenu("Help",
				fyne.NewMenuItem("About", func() {
					fyne.CurrentApp().SendNotification(&fyne.Notification{
						Title:   "About",
						Content: "This is a simple Fyne application",
					})
				}),
			),
		)

		// تنظیم منوی اصلی برای پنجره
		myWindow.SetMainMenu(mainMenu)
	})

	// --- Colunm Right -----
	rightColumnContent := container.NewVBox(
		item1,
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
		openButton,
		darkLight,
	)

	// --- Colunm Left -----

	lastColumnContent := container.NewVBox(
		plusButton,
	)

	columns := container.NewHSplit(lastColumnContent, rightColumnContent)
	columns.SetOffset(0.25)

	scrol := container.NewScroll(columns)

	myWindow.SetContent(scrol)
	myWindow.Resize(fyne.NewSize(900, 600))
	myWindow.ShowAndRun()
}
