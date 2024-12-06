package logic

import (
	variable "DatabaseDB"
	"DatabaseDB/internal/utils"
	"fmt"

	// "DatabaseDB/internal/logic/mainwindowlagic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/gabriel-vasile/mimetype"
)

func HandleButtonClick(test string, nameDatabace string) error {
	err := utils.Checkdatabace(test, nameDatabace)
	if err != nil {
		return err
	}

	if !variable.CreatDatabase {

		nun := variable.NameData.FilterFile(test)
		if !nun {
			return fmt.Errorf("error for no found files database")
		}
	}
	err = variable.CurrentDBClient.Open()
	if err != nil {
		return err
	}
	defer variable.CurrentDBClient.Close()

	return nil
}

func SearchDatabase(valueEntry *widget.Entry, editWindow fyne.Window, rightColumnContent *fyne.Container, columnEditKey *fyne.Container, saveKey *widget.Button, mainWindow fyne.Window) (bool, error) {

	err := variable.CurrentDBClient.Open()
	if err != nil {
		return false, err
	}

	key := utils.CleanInput(valueEntry.Text)
	err, data := variable.CurrentDBClient.Search([]byte(key))
	if err != nil {
		return false, err
	}
	defer variable.CurrentDBClient.Close()

	if len(data) == 0 {
		return false, err
	}
	utils.CheckCondition(columnEditKey)
	utils.CheckCondition(rightColumnContent)
	var truncatedValue string
	var count int
	for _, item := range data {

		if count > 40 {
			dialog.ShowInformation("Error", "The result of your keys is more than 60 and I will only show the first 60.If your key is not among these, please search more precisely.", mainWindow)
			count = 0
			break
		}
		count++

		value, err := variable.CurrentDBClient.Get(item)
		if err != nil {
			return false, err
		}
		truncatedKey := utils.TruncateString(string(item), 20)

		typeValue := mimetype.Detect([]byte(value))
		if typeValue.Extension() != ".txt" {
			truncatedValue = fmt.Sprintf("* %s . . .", typeValue.Extension())
		} else {
			truncatedValue = utils.TruncateString(string(value), 20)

		}

		valueLabel := BuidLableKeyAndValue("value", item, value, truncatedValue, rightColumnContent, columnEditKey, saveKey, mainWindow)
		keyLabel := BuidLableKeyAndValue("key", item, value, truncatedKey, rightColumnContent, columnEditKey, saveKey, mainWindow)

		rightColumnContent.Refresh()
		buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
		rightColumnContent.Add(buttonRow)
	}
	editWindow.Close()
	return true, nil
}

func DeleteKeyLogic(valueEntry *widget.Entry, editWindow fyne.Window, rightColumnContent *fyne.Container) {
	defer variable.CurrentDBClient.Close()

	key := utils.CleanInput(valueEntry.Text)

	valueSearch, err := QueryKey(valueEntry.Text)
	if valueSearch == "" && err != nil {
		dialog.ShowInformation("Error", "This key does not exist in the database", editWindow)
	} else {
		err = variable.CurrentDBClient.Delete([]byte(key))
		if err != nil {
			fmt.Print(err.Error())
			return
		}
		editWindow.Close()
	}
}

func AddKeyLogic(iputKey string, valueFinish []byte, windowAdd fyne.Window) {

	key := utils.CleanInput(iputKey)

	defer variable.CurrentDBClient.Close()

	checkNow, err := QueryKey(iputKey)
	if checkNow != "" || err == nil {
		dialog.ShowInformation("Error", "This key has already been added to your database", windowAdd)

	} else {
		err = variable.CurrentDBClient.Add([]byte(key), valueFinish)
		if err != nil {
			fmt.Print(err.Error())
		}

		windowAdd.Close()
	}
}

func QueryKey(iputKey string) (string, error) {
	var err error

	key := utils.CleanInput(iputKey)

	err = variable.CurrentDBClient.Open()
	if err != nil {
		return "", err
	}
	checkNow, err := variable.CurrentDBClient.Get([]byte(key))
	if err != nil {
		fmt.Println("error : delete func logic for get key in databace")
	}
	return string(checkNow), err
}
