package main

import (
	_ "net/http/pprof"
	variable "testgui"
	jsondata "testgui/internal/json/jsonData"
	"testgui/internal/ui/mainwindow"

	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()

	variable.CurrentJson = jsondata.NewDataBase()

	mainwindow.MainWindow(myApp)
}
