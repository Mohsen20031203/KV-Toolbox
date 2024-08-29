package main

import (
	mainW "testgui/internal/ui/main-window"

	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()

	mainW.MainWindow(myApp)
}
