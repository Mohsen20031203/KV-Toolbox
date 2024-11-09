package variable

import (
	dbpak "DatabaseDB/internal/Databaces"
	"DatabaseDB/internal/filterdatabase"
	jsFile "DatabaseDB/internal/json"
)

var (
	CurrentDBClient dbpak.DBClient
	CurrentPage     int
	ItemsPerPage    = 7
	FolderPath      string
	CurrentJson     jsFile.JsonFile
	NameData        filterdatabase.FilterData
	ItemsAdded      bool
	PreviousOffsetY float32
)

var (
	NameDatabase = []string{
		"levelDB",
		"Pebble",
		"Badger",
		//"Redis",
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
