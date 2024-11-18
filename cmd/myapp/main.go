package main

import (
	variable "DatabaseDB"
	jsondata "DatabaseDB/internal/config/jsonconfig"
	"DatabaseDB/internal/ui/mainwindow"

	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()

	variable.CurrentJson = jsondata.NewDataBase()

	mainwindow.MainWindow(myApp)
}
