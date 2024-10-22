package variable

import (
	dbpak "testgui/internal/Databaces"
	"testgui/internal/filterdatabase"
	jsFile "testgui/internal/json"
)

var (
	CurrentDBClient dbpak.DBClient
	CurrentPage     int
	ItemsPerPage    = 20
	FolderPath      string
	CurrentJson     jsFile.JsonFile
	NameData        filterdatabase.FilterData
	ItemsAdded      bool
)

var (
	NameDatabase = []string{
		"levelDB",
		"Pebble",
		"Badger",
		"Redis",
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
