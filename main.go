package main

import (
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Non-Scrollable List")

	createSeparator := func() *canvas.Line {
		line := canvas.NewLine(color.Black)
		line.StrokeWidth = 1
		return line
	}

	item1 := widget.NewLabel("Item 1")
	separator1 := createSeparator()

	item2 := widget.NewLabel("Item 2")

	listContainer := container.NewVBox(
		item1,
		separator1,
		item2,
	)

	combo := widget.NewSelect([]string{"Option 1", "Option 2"}, func(value string) {
		log.Println("Select set to", value)
	})

	options := []string{"Option 1", "Option 2", "Option 3"}
	selectEntry := widget.NewSelectEntry(options)

	selectEntry.SetText("Option 2")

	selectEntry.OnChanged = func(s string) {
		fmt.Println("Changed to:", s)
	}

	check := widget.NewCheck("Check", func(value bool) {
		fmt.Println("Check : ", value)
	})

	disableCheck := widget.NewCheck("disableCheck", func(value bool) {
		fmt.Println("disableCheckprint :", value)
	})
	disableCheck.Disable()

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

	rightColumnContent := container.NewVBox(
		listContainer,
		combo,
		selectEntry,
		check,
		disableCheck,
		horizontalContainer,
	)

	leftColumnContent := container.NewVBox(
		widget.NewLabel("Left Column: 1/4 Width"),
		widget.NewButton("Button 1", func() { println("Button 1 clicked") }),
		widget.NewEntry(),
	)

	columns := container.NewHSplit(leftColumnContent, rightColumnContent)

	columns.SetOffset(0.25)

	myWindow.SetContent(columns)
	myWindow.Resize(fyne.NewSize(800, 400))
	myWindow.ShowAndRun()
}
