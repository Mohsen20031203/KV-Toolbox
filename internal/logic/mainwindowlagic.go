package logic

import (
	"encoding/json"
	"fmt"
	variable "testgui"

	// "testgui/internal/logic/addProjectwindowlogic"

	dbpak "testgui/internal/db"
	leveldbb "testgui/internal/db/leveldb"
	jsondata "testgui/internal/json/jsonData"
	"testgui/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var count int

func SetupLastColumn(rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button) *fyne.Container {
	lastColumnContent := container.NewVBox()

	jsonnew := jsondata.NewDataBase()
	jsonDataa, err := jsonnew.Load()
	if err != nil {
		println("Error loading JSON data:", err)
	} else {
		for _, project := range jsonDataa.RecentProjects {
			buttonContainer := ProjectButton(project.Name, lastColumnContent, project.FileAddress, rightColumnContentORG, nameButtonProject, buttonAdd)
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

var lastStart *string
var lastEnd *string
var lastPage int
var currentData []dbpak.KVData
var lastcurrentData []dbpak.KVData

func UpdatePage(rightColumnContent *fyne.Container) {

	if !utils.CheckCondition(rightColumnContent) {
		rightColumnContent.Objects = []fyne.CanvasObject{}
		rightColumnContent.Refresh()
	}
	if lastStart == nil {
		lastEnd = &currentData[variable.ItemsPerPage-1].Key
	}
	var data = make([]dbpak.KVData, 0)
	var err error
	if lastPage <= variable.CurrentPage {
		//next page
		err, data = variable.CurrentDBClient.Read(lastEnd, nil, variable.ItemsPerPage+1)
		if err != nil {
			fmt.Println(err)
		}
		if len(data) > variable.ItemsPerPage {
			variable.NextButton.Enable()
		} else {
			variable.NextButton.Disable()
		}
	} else {
		//last page
		if len(currentData) == 0 {
			return
		}
		err, data = variable.CurrentDBClient.Read(nil, &currentData[0].Key, variable.ItemsPerPage+1)
		if err != nil {
			fmt.Println(err)
		}
		if len(data) > variable.ItemsPerPage {
			variable.PrevButton.Enable()
		} else {
			variable.PrevButton.Disable()
		}
	}
	lastcurrentData = make([]dbpak.KVData, len(currentData))
	copy(lastcurrentData, currentData)
	currentData = make([]dbpak.KVData, len(data))
	copy(currentData, data)
	if len(data) == 0 {
		return
	}

	lastPage = variable.CurrentPage
	lastStart = &data[0].Key
	lastEnd = &data[len(data)-1].Key

	/*
		StartIndex := variable.CurrentPage * variable.ItemsPerPage
		EndIndex := StartIndex + variable.ItemsPerPage

		if EndIndex > len(data) {
			EndIndex = len(data)
		}

		rightColumnContent.Objects = nil

	*/

	number := 0
	for _, item := range data {
		if number == variable.ItemsPerPage {
			break
		}
		number++
		truncatedKey := utils.TruncateString(item.Key, 20)
		truncatedValue := utils.TruncateString(item.Value, 50)

		valueLabel := BuidLableKeyAndValue("value", item.Key, item.Value, truncatedValue, variable.FolderPath, rightColumnContent)
		keyLabel := BuidLableKeyAndValue("key", item.Key, item.Value, truncatedKey, variable.FolderPath, rightColumnContent)

		buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
		rightColumnContent.Add(buttonRow)
	}

	variable.PageLabel.SetText(fmt.Sprintf("Page %d", variable.CurrentPage+1))

	rightColumnContent.Refresh()
}

func ProjectButton(inputText string, lastColumnContent *fyne.Container, path string, rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button) *fyne.Container {
	projectButton := widget.NewButton(inputText, func() {
		variable.CurrentDBClient = leveldbb.NewDataBase(path)
		variable.PrevButton.Disable()
		lastStart = nil
		lastPage = 0
		variable.CurrentPage = 0
		variable.NextButton.Enable()
		variable.PageLabel.Text = "Page 1"
		variable.FolderPath = path
		HandleProjectSelection(path, rightColumnContentORG, buttonAdd)
		if nameButtonProject.Text == "" {
			nameButtonProject.Text = inputText
		} else {
			nameButtonProject.Text = ""
			nameButtonProject.Text = inputText
		}

		variable.CurrentDBClient = leveldbb.NewDataBase(path)

		nameButtonProject.Refresh()
		variable.PageLabel.Refresh()

	})

	if nameButtonProject.Text == "" {
		nameButtonProject.Text = inputText
	} else {
		nameButtonProject.Text = ""
		nameButtonProject.Text = inputText
	}
	nameButtonProject.Refresh()

	buttonContainer := container.NewHBox()

	closeButton := widget.NewButton("✖", func() {

		if !utils.CheckCondition(rightColumnContentORG) && nameButtonProject.Text == inputText {
			newObjects := []fyne.CanvasObject{}

			rightColumnContentORG.Objects = newObjects
			buttonAdd.Disable()

			nameButtonProject.Text = ""
			nameButtonProject.Refresh()
			rightColumnContentORG.Refresh()
		}

		err := variable.CurrentJson.Remove(inputText)
		if err != nil {
			fmt.Println("Failed to remove project from JSON:", err)
		} else {

			lastColumnContent.Remove(buttonContainer)
			lastColumnContent.Refresh()
		}
	})

	buttonContainer = container.NewBorder(nil, nil, nil, closeButton, projectButton)
	return buttonContainer
}

func HandleProjectSelection(dbPath string, rightColumnContent *fyne.Container, buttonAdd *widget.Button) {

	opts := &opt.Options{
		ReadOnly: true,
	}
	DBDB, err := leveldb.OpenFile(dbPath, opts)
	if err != nil {
		return
	}
	defer DBDB.Close()

	buttonAdd.Enable()
	if !utils.CheckCondition(rightColumnContent) {
		newObjects := []fyne.CanvasObject{}

		rightColumnContent.Objects = newObjects

		rightColumnContent.Refresh()
	}

	err, data := variable.CurrentDBClient.Read(nil, nil, variable.ItemsPerPage+1)
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		fmt.Println("Failed to read database:", err)
		return
	}

	if len(data)/variable.ItemsPerPage > 1 {
		variable.NextButton.Enable()
		variable.CurrentPage = 0
	}
	currentData = currentData[:0]
	count = 0
	for _, item := range data {
		if count >= variable.ItemsPerPage {
			count = 0
			break
		}
		currentData = append(currentData, item)
		count++

		truncatedKey := utils.TruncateString(item.Key, 20)
		truncatedValue := utils.TruncateString(item.Value, 50)

		valueLabel := BuidLableKeyAndValue("value", item.Key, item.Value, truncatedValue, dbPath, rightColumnContent)
		keyLabel := BuidLableKeyAndValue("key", item.Key, item.Value, truncatedKey, dbPath, rightColumnContent)

		buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
		rightColumnContent.Add(buttonRow)

	}

	rightColumnContent.Refresh()
}

func BuidLableKeyAndValue(eidtKeyAbdValue string, key string, value string, nameLable string, Addres string, rightColumnContent *fyne.Container) *TappableLabel {
	var lableKeyAndValue *TappableLabel

	lableKeyAndValue = NewTappableLabel(nameLable, func() {
		editWindow := fyne.CurrentApp().NewWindow("Edit" + eidtKeyAbdValue)
		editWindow.Resize(fyne.NewSize(600, 600))

		valueEntry := widget.NewMultiLineEntry()
		valueEntry.Resize(fyne.NewSize(500, 500))
		if eidtKeyAbdValue == "value" {
			if utils.IsValidJSON(value) {
				var formattedJSON map[string]interface{}
				json.Unmarshal([]byte(value), &formattedJSON)
				jsonString, _ := json.MarshalIndent(formattedJSON, "", "  ")
				valueEntry.SetText(string(jsonString))
			} else {
				valueEntry.SetText(value)
			}
		} else {
			valueEntry.SetText(key)
		}
		scrollableEntry := container.NewScroll(valueEntry)
		mainContainer := container.NewBorder(nil, nil, nil, nil, scrollableEntry)

		scrollableEntry.SetMinSize(fyne.NewSize(600, 500))
		saveButton := widget.NewButton("Save", func() {
			var truncatedKey2 string

			err := variable.CurrentDBClient.Open()
			if err != nil {
				fmt.Println("error Open")
			}
			defer variable.CurrentDBClient.Close()

			if eidtKeyAbdValue == "value" {
				err := variable.CurrentDBClient.Add(key, value)
				if err != nil {
					fmt.Println(err)
				}
				truncatedKey2 = utils.TruncateString(valueEntry.Text, 50)

			} else {
				valueBefor := variable.CurrentDBClient.Get(key)

				err = variable.CurrentDBClient.Delet(key)
				if err != nil {
					return
				}

				key = valueEntry.Text

				err := variable.CurrentDBClient.Add(key, valueBefor)
				if err != nil {
					fmt.Println(err)
				}
				truncatedKey2 = utils.TruncateString(key, 20)
			}

			lableKeyAndValue.SetText(truncatedKey2)
			lableKeyAndValue.Refresh()

			editWindow.Close()
			rightColumnContent.Refresh()
		})

		cancelButton := widget.NewButton("Cancel", func() {
			editWindow.Close()
		})

		m := container.NewGridWithColumns(2, cancelButton, saveButton)
		b := container.NewBorder(nil, m, nil, nil)

		editContent := container.NewVBox(
			widget.NewLabel("Edit "+eidtKeyAbdValue+" :"),
			mainContainer,
			layout.NewSpacer(),
			b,
		)

		editWindow.SetContent(editContent)
		editWindow.Show()
	})
	return lableKeyAndValue
}