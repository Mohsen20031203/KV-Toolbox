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
	CreatDatabase   bool
)

var (
	NameDatabase = []string{
		"levelDB",
		"Pebble",
		"Badger",
		//"Redis",
	}
)
