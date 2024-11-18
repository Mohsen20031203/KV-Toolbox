package variable

import (
	dbpak "DatabaseDB/internal/Databaces"
	jsFile "DatabaseDB/internal/config"
	"DatabaseDB/internal/filterdatabase"
)

var (
	CurrentDBClient dbpak.DBClient
	CurrentPage     int
	ItemsPerPage    = 16
	FolderPath      string
	CurrentJson     jsFile.ConfigHandler
	NameData        filterdatabase.FilterData
	ItemsAdded      bool
	PreviousOffsetY float32
	ResultSearch    bool
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
