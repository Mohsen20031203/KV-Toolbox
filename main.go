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

	combo := widget.NewSelect([]string{"Option 1", "Option 2"}, func(value string) {
		log.Println("Select set to", value)
	})

	options := []string{"Option 1", "Option 2", "Option 3"}
	selectEntry := widget.NewSelectEntry(options)

	selectEntry.SetText("Option 2")

	selectEntry.OnChanged = func(s string) {
		fmt.Println("Changed to:", s)
	}

	listContainer := container.NewVBox(
		item1,
		separator1,
		item2,
		combo,
	)

	rightColumnContent := container.NewVBox(
		listContainer,
		selectEntry,
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
