package main

import (
	variable "testgui"
	"testgui/internal/ui/mainwindow"
	jsondata "testgui/pkg/json/jsonData"

	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()

	variable.CurrentJson = jsondata.NewDataBase()

	mainwindow.MainWindow(myApp)
}
