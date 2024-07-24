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
	"fyne.io/fyne/v2/widget"
)

var data = []string{"Canvas", "Animation", "Animation", "Animation", "Animation", "Animation", "TemeIcon", "TemeIcon", "TemeIcon", "TemeIcon", "TemeIcon", "TemeIcon", "TemeIcon", "TemeIcon", "TemeIcon", "TemeIcon"}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Non-Scrollable List")

	createSeparator := func() *canvas.Line {
		line := canvas.NewLine(color.Black)
		line.StrokeWidth = 1
		return line
	}

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i])
		})

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

	inputAdress := widget.NewEntry()
	inputAdress.SetPlaceHolder("URL : ")

	openButton := widget.NewButton("Open File", func() {
		dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			if file == nil {
				dialog.ShowInformation("No File Selected", "No file was selected.", myWindow)
				return
			}

			// Set file path in the entry widget
			inputAdress.SetText(file.URI().String())

		}, myWindow).Show()
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
		openButton,
		inputAdress,
	)

	// --- Colunm Left -----

	columns := container.NewHSplit(list, rightColumnContent)
	columns.SetOffset(0.25)

	scrol := container.NewScroll(columns)

	myWindow.SetContent(scrol)
	myWindow.Resize(fyne.NewSize(900, 600))
	myWindow.ShowAndRun()
}
