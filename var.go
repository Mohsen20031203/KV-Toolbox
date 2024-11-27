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

/*
1. GOOS='darwin'
2. GOARCH='arm64'
3. CC='clang'


-------------------
1. export GOOS=windows
2. export GOARCH=amd64 || export GOARCH=386
3. export CC=x86_64-w64-mingw32-gcc || export CC=i686-w64-mingw32-gcc


*/
