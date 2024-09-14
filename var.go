package variable

import (
	dbpak "testgui/internal/db"
	jsFile "testgui/internal/json"

	"fyne.io/fyne/v2/widget"
)

var (
	CurrentDBClient dbpak.DBClient
	CurrentPage     int
	ItemsPerPage    = 10
	FolderPath      string
	CurrentJson     jsFile.JsonFile
)

var NextButton, PrevButton *widget.Button
var PageLabel *widget.Label
