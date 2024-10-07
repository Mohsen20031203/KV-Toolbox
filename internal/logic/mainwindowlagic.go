package logic

import (
	"encoding/json"
	"fmt"
	"log"
	variable "testgui"
	"time"

	// "testgui/internal/logic/addProjectwindowlogic"

	dbpak "testgui/internal/Databaces"
	"testgui/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var next_prev bool

func SetupLastColumn(rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button) *fyne.Container {
	lastColumnContent := container.NewVBox()

	jsonDataa, err := variable.CurrentJson.Load()
	if err != nil {
		println("Error loading JSON data:", err)
	} else {
		for _, project := range jsonDataa.RecentProjects {
			buttonContainer := ProjectButton(project.Name, lastColumnContent, project.FileAddress, rightColumnContentORG, nameButtonProject, buttonAdd, project.Databace)
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

	utils.CheckCondition(rightColumnContent)

	var data = make([]dbpak.KVData, 0)
	var err error
	if lastPage <= variable.CurrentPage {
		//next page

		//The reason why "variable.ItemsPerPage" is added by one is that we want to see if the next pages have a value to enable or disable the next or prev key.
		err, data = variable.CurrentDBClient.Read(lastEnd, nil, variable.ItemsPerPage+1)
		if err != nil {
			fmt.Println(err)
		}
		if len(data) > variable.ItemsPerPage {
			variable.NextButton.Enable()
			next_prev = true
		} else {
			variable.NextButton.Disable()
		}
	} else {
		//last page
		if len(currentData) == 0 {
			return
		}

		//The reason why "variable.ItemsPerPage" is added by one is that we want to see if the next pages have a value to enable or disable the next or prev key.
		err, data = variable.CurrentDBClient.Read(nil, lastStart, variable.ItemsPerPage+1)
		if err != nil {
			fmt.Println(err)
		}

		if len(data) > variable.ItemsPerPage {
			variable.PrevButton.Enable()
			next_prev = false
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
	if next_prev {
		lastStart = &data[0].Key
		if len(data) == variable.ItemsPerPage+1 {

			lastEnd = &data[len(data)-2].Key
		} else {
			lastEnd = &data[len(data)-1].Key

		}

	} else {
		lastEnd = &data[len(data)-1].Key
		if len(data) >= variable.ItemsPerPage+1 {
			lastStart = &data[1].Key
			data = data[1:]
		} else {
			lastStart = &data[0].Key
		}
	}

	number := 0

	for _, item := range data {
		if number == variable.ItemsPerPage {
			break
		}
		number++
		truncatedKey := utils.TruncateString(item.Key, 20)
		truncatedValue := utils.TruncateString(item.Value, 50)

		valueLabel := BuidLableKeyAndValue("value", item.Key, item.Value, truncatedValue, rightColumnContent)
		keyLabel := BuidLableKeyAndValue("key", item.Key, item.Value, truncatedKey, rightColumnContent)

		buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
		rightColumnContent.Add(buttonRow)
	}

	variable.PageLabel.SetText(fmt.Sprintf("Page %d", variable.CurrentPage+1))

	rightColumnContent.Refresh()
}

func ProjectButton(inputText string, lastColumnContent *fyne.Container, path string, rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button, nameDatabace string) *fyne.Container {
	projectButton := widget.NewButton(inputText, func() {
		utils.Checkdatabace(path, nameDatabace)
		variable.PrevButton.Disable()
		lastPage = 0
		variable.CurrentPage = 0
		variable.NextButton.Enable()
		lastEnd = nil
		lastStart = nil
		variable.PageLabel.Text = "Page 1"
		variable.FolderPath = path
		HandleProjectSelection(path, rightColumnContentORG, buttonAdd)
		if nameButtonProject.Text == "" {
			nameButtonProject.Text = inputText + " - " + nameDatabace
		} else {
			nameButtonProject.Text = ""
			nameButtonProject.Text = inputText + " - " + nameDatabace
		}
		nameButtonProject.Refresh()
		variable.PageLabel.Refresh()

	})

	if nameButtonProject.Text == "" {
		nameButtonProject.Text = inputText + " - " + nameDatabace
	} else {
		nameButtonProject.Text = ""
		nameButtonProject.Text = inputText + " - " + nameDatabace
	}
	nameButtonProject.Refresh()

	buttonContainer := container.NewHBox()

	closeButton := widget.NewButton("✖", func() {

		if nameButtonProject.Text == inputText+" - "+nameDatabace {
			utils.CheckCondition(rightColumnContentORG)

			buttonAdd.Disable()

			nameButtonProject.Text = ""
			nameButtonProject.Refresh()
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
	utils.CheckCondition(rightColumnContent)

	//The reason why "variable.ItemsPerPage" is added by one is that we want to see if the next pages have a value to enable or disable the next or prev key.

	UpdatePage(rightColumnContent)
}

func BuidLableKeyAndValue(eidtKeyAbdValue string, key string, value string, nameLable string, rightColumnContent *fyne.Container) *TappableLabel {
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
				err := variable.CurrentDBClient.Add(key, valueEntry.Text)
				if err != nil {
					fmt.Println(err)
				}
				truncatedKey2 = utils.TruncateString(valueEntry.Text, 50)

			} else {
				valueBefor, err := variable.CurrentDBClient.Get(key)
				if err != nil {
					return
				}
				err = variable.CurrentDBClient.Delete(key)
				if err != nil {
					return
				}

				key = valueEntry.Text

				err = variable.CurrentDBClient.Add(key, valueBefor)
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

func SearchDatabase(valueEntry *widget.Entry, editWindow fyne.Window, rightColumnContent *fyne.Container) {
	defer variable.CurrentDBClient.Close()

	valueSearch, err := QueryKey(valueEntry)
	if valueSearch == "" && err != nil {
		dialog.ShowError(fmt.Errorf("The key - "+valueEntry.Text+" - does not exist in your database"), editWindow)
		valueEntry.Text = ""
		valueEntry.Refresh()
	} else {
		editWindow.Close()

		utils.CheckCondition(rightColumnContent)

		truncatedKey := utils.TruncateString(valueEntry.Text, 20)
		truncatedValue := utils.TruncateString(valueSearch, 50)

		valueLabel := BuidLableKeyAndValue("value", valueEntry.Text, valueSearch, truncatedValue, rightColumnContent)
		keyLabel := BuidLableKeyAndValue("key", valueEntry.Text, valueSearch, truncatedKey, rightColumnContent)
		buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
		rightColumnContent.Add(buttonRow)
		variable.NextButton.Disable()
		variable.PrevButton.Disable()
	}
}

func DeleteKeyLogic(valueEntry *widget.Entry, editWindow fyne.Window, rightColumnContent *fyne.Container) {
	defer variable.CurrentDBClient.Close()

	valueSearch, err := QueryKey(valueEntry)
	if valueSearch == "" && err != nil {
		dialog.ShowInformation("Error", "This key does not exist in the database", editWindow)
	} else {
		err = variable.CurrentDBClient.Delete(valueEntry.Text)
		if err != nil {
			log.Fatal("this err for func DeletKeyLogic part else delete || err : ", err)
			return
		}
		dialog.ShowInformation("successful", "The operation was successful", editWindow)
		time.Sleep(2 * time.Second)
		editWindow.Close()
	}
}

func AddKeyLogic(iputKey *widget.Entry, iputvalue *widget.Entry, windowAdd fyne.Window) {
	if iputKey.Text == "" && iputvalue.Text == "" {
		dialog.ShowInformation("Error", "Please enter both the key and the value", windowAdd)
	} else if iputvalue.Text != "" && iputKey.Text == "" {
		dialog.ShowInformation("Error", "You cannot leave either the key or both fields empty.", windowAdd)

	}
	defer variable.CurrentDBClient.Close()

	checkNow, err := QueryKey(iputKey)
	if checkNow != "" || err == nil {
		dialog.ShowInformation("Error", "This key has already been added to your database", windowAdd)

	} else {
		err = variable.CurrentDBClient.Add(iputKey.Text, iputvalue.Text)
		if err != nil {
			log.Fatal("error : this error in func addkeylogic for add key in database")
		}
		dialog.ShowInformation("successful", "The operation was successful", windowAdd)
		time.Sleep(2 * time.Second)

		windowAdd.Close()
	}
}

func QueryKey(iputKey *widget.Entry) (string, error) {
	var err error
	err = variable.CurrentDBClient.Open()
	if err != nil {
		return "", err
	}

	checkNow, err := variable.CurrentDBClient.Get(iputKey.Text)
	if err != nil {
		fmt.Println("error : delete func logic for get key in databace")
	}
	return checkNow, err
}
