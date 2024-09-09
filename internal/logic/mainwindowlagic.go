package logic

import (
	"encoding/json"
	"fmt"
	variable "testgui"

	// "testgui/internal/logic/addProjectwindowlogic"

	jsondata "testgui/internal/logic/json/jsonData"
	"testgui/internal/utils"
	dbpak "testgui/pkg/db"
	leveldbb "testgui/pkg/db/leveldb"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var count int
var lastkey dbpak.Database

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

func UpdatePage(rightColumnContent *fyne.Container) {
	if !utils.CheckCondition(rightColumnContent) {
		rightColumnContent.Objects = []fyne.CanvasObject{}
		rightColumnContent.Refresh()
	}

	err, data := variable.CurrentDBClient.Read()
	if err != nil {
		fmt.Println(err)
	}

	StartIndex := variable.CurrentPage * variable.ItemsPerPage
	EndIndex := StartIndex + variable.ItemsPerPage

	if EndIndex > len(data) {
		EndIndex = len(data)
	}

	rightColumnContent.Objects = nil

	for _, item := range data[StartIndex:EndIndex] {
		truncatedKey := utils.TruncateString(item.Key, 20)
		truncatedValue := utils.TruncateString(item.Value, 50)

		valueLabel := BuidLableKeyAndValue("value", item.Key, item.Value, truncatedValue, variable.FolderPath, rightColumnContent)
		keyLabel := BuidLableKeyAndValue("key", item.Key, item.Value, truncatedKey, variable.FolderPath, rightColumnContent)

		buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
		rightColumnContent.Add(buttonRow)
	}

	variable.PageLabel.SetText(fmt.Sprintf("Page %d", variable.CurrentPage+1))

	variable.PrevButton.Disable()
	variable.NextButton.Disable()

	if variable.CurrentPage > 0 {
		variable.PrevButton.Enable()
	}
	if EndIndex < len(data) {
		variable.NextButton.Enable()
	}

	rightColumnContent.Refresh()
}

func ProjectButton(inputText string, lastColumnContent *fyne.Container, path string, rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button) *fyne.Container {
	projectButton := widget.NewButton(inputText, func() {
		variable.CurrentDBClient = leveldbb.NewDataBase(path)
		variable.PrevButton.Disable()
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

	buttonAdd.Enable()
	if !utils.CheckCondition(rightColumnContent) {
		newObjects := []fyne.CanvasObject{}

		rightColumnContent.Objects = newObjects

		rightColumnContent.Refresh()
	}

	err, data := variable.CurrentDBClient.Read()
	if err != nil {
		fmt.Println("Failed to read database:", err)
		return
	}

	if len(data)/variable.ItemsPerPage > 1 {
		variable.NextButton.Enable()
		variable.CurrentPage = 0
	}
	for _, item := range data {
		if count >= variable.ItemsPerPage {
			count = 0
			break
		}
		count++

		truncatedKey := utils.TruncateString(item.Key, 20)
		truncatedValue := utils.TruncateString(item.Value, 50)

		valueLabel := BuidLableKeyAndValue("value", item.Key, item.Value, truncatedValue, dbPath, rightColumnContent)
		keyLabel := BuidLableKeyAndValue("key", item.Key, item.Value, truncatedKey, dbPath, rightColumnContent)

		buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
		rightColumnContent.Add(buttonRow)
		lastkey = item

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
