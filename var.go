package variable

import (
	dbpak "testgui/internal/Databaces"
	"testgui/internal/filterdatabase"
	jsFile "testgui/internal/json"

	"fyne.io/fyne/v2/widget"
)

var (
	CurrentDBClient dbpak.DBClient
	CurrentPage     uint8
	ItemsPerPage    = 20
	FolderPath      string
	CurrentJson     jsFile.JsonFile
	NameData        filterdatabase.FilterData
)

var NextButton, PrevButton *widget.Button
var PageLabel *widget.Label

var (
	NameDatabase = []string{
		"levelDB",
		"Pebble",
		"Badger",
	}
)

var CreatDatabase bool

/*var (
	SearchName     = "Search"
	AddName        = "Add"
	TopColumnKey   = "Key"
	TopColumnValue = "value"
	DeleteName     = "Delete"
	NextName       = "Next"
	PrevName       = "Prev"
	nameDatabase   = []string{
		"leveldb",
		"pebble",
	}
)

var (
	NameAddProject     = "Name :"
	CommentAddProject  = "Comment :"
	TestConnection     = "Test Connection"
	OpenFile           = "Open Folder"
	CreatDatabaseCheck = "Create Database"
	CancelAddProject   = "Cancel"
)*/
