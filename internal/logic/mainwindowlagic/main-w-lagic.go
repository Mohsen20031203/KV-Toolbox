package mainwindowlagic

import (
	"fmt"
	"testgui/internal/logic/addProjectwindowlogic"
	"testgui/internal/ui/mainwindow"
	"testgui/internal/utils"
	dbpak "testgui/pkg/db"
	jsFile "testgui/pkg/json"
	jsondata "testgui/pkg/json/jsonData"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var CurrentDBClient dbpak.DBClient
var CurrentPage int
var ItemsPerPage = 20
var FolderPath string
var CurrentJson jsFile.JsonFile

func SetupLastColumn(rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button) *fyne.Container {
	lastColumnContent := container.NewVBox()

	jsonnew := jsondata.NewDataBase()
	jsonDataa, err := jsonnew.Load()
	if err != nil {
		println("Error loading JSON data:", err)
	} else {
		for _, project := range jsonDataa.RecentProjects {
			buttonContainer := utils.ProjectButton(project.Name, lastColumnContent, project.FileAddress, rightColumnContentORG, nameButtonProject, buttonAdd)
			lastColumnContent.Add(buttonContainer)
		}
	}

	return lastColumnContent
}

func SetupThemeButtons(app fyne.App) *fyne.Container {
	darkButton := widget.NewButton("Dark", func() {
		app.Settings().SetTheme(theme.DarkTheme())
	})
	lightButton := widget.NewButton("Light", func() {
		app.Settings().SetTheme(theme.LightTheme())
	})

	darkLight := container.NewVBox(
		layout.NewSpacer(),
		container.NewGridWithColumns(2, lightButton, darkButton),
	)
	return darkLight
}

func UpdatePage(rightColumnContent *fyne.Container) {
	if !utils.CheckCondition(rightColumnContent) {
		rightColumnContent.Objects = []fyne.CanvasObject{}
		rightColumnContent.Refresh()
	}

	err, data := CurrentDBClient.Read()
	if err != nil {
		fmt.Println(err)
	}

	StartIndex := CurrentPage * ItemsPerPage
	EndIndex := StartIndex + ItemsPerPage

	if EndIndex > len(data) {
		EndIndex = len(data)
	}

	rightColumnContent.Objects = nil

	for _, item := range data[StartIndex:EndIndex] {
		truncatedKey := utils.TruncateString(item.Key, 20)
		truncatedValue := utils.TruncateString(item.Value, 50)

		valueLabel := addProjectwindowlogic.BuidLableKeyAndValue("value", item.Key, item.Value, truncatedValue, FolderPath, rightColumnContent)
		keyLabel := addProjectwindowlogic.BuidLableKeyAndValue("key", item.Key, item.Value, truncatedKey, FolderPath, rightColumnContent)

		buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
		rightColumnContent.Add(buttonRow)
	}

	// به‌روزرسانی شماره صفحه
	mainwindow.PageLabel.SetText(fmt.Sprintf("Page %d", CurrentPage+1))

	// غیرفعال کردن دکمه‌ها بر اساس موقعیت فعلی
	mainwindow.PrevButton.Disable()
	mainwindow.NextButton.Disable()

	if CurrentPage > 0 {
		mainwindow.PrevButton.Enable()
	}
	if EndIndex < len(data) {
		mainwindow.NextButton.Enable()
	}

	rightColumnContent.Refresh()
}

/*
func handleButtonClick(test string) error {
	db, err := leveldb.OpenFile(test, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	iter := db.NewIterator(nil, nil)
	defer iter.Release()

	if iter.First() {
		key := iter.Key()
		value, err := db.Get(key, nil)
		if err != nil {
			return fmt.Errorf("failed to get value for key %s: %v", key, err)
		}

		fmt.Printf("First key: %s, value: %s\n", key, value)
		return nil
	}
	return fmt.Errorf("no entries found in the database")
}
*/
