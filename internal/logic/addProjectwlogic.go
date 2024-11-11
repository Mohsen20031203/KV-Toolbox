package logic

import (
	variable "DatabaseDB"
	"DatabaseDB/internal/utils"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	// "DatabaseDB/internal/logic/mainwindowlagic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/gabriel-vasile/mimetype"
)

func HasManifestFile(folderPath string) bool {
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		log.Fatal("Error reading folder:", err)
		return false
	}
	var count uint8
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "MANIFEST-") || filepath.Ext(file.Name()) == ".log" {
			count++
		}

		if count == 2 {
			return true
		}
	}
	return false
}

func CreatFile(value bool, openButton *widget.Button, testConnectionButton *widget.Button) {
	if value {
		openButton.Disable()
		testConnectionButton.Disable()
		variable.CreatDatabase = value
	} else {
		openButton.Enable()
		testConnectionButton.Enable()
		variable.CreatDatabase = value
	}
}

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

	if len(data) == 0 {
		return false, err
	}
	utils.CheckCondition(columnEditKey)
	utils.CheckCondition(rightColumnContent)
	for _, item := range data {

		value, err := variable.CurrentDBClient.Get(item)
		if err != nil {
			return false, err
		}
		truncatedKey := utils.TruncateString(string(item), 20)
		truncatedValue := utils.TruncateString(string(value), 30)

		typeValue := mimetype.Detect([]byte(value))
		if typeValue.Extension() != ".txt" {
			truncatedValue = fmt.Sprintf("* %s . . .", typeValue.Extension())
		}
		valueLabel := BuidLableKeyAndValue("value", item, value, truncatedValue, rightColumnContent, columnEditKey, saveKey, mainWindow)
		keyLabel := BuidLableKeyAndValue("key", item, value, truncatedKey, rightColumnContent, columnEditKey, saveKey, mainWindow)
		rightColumnContent.Refresh()
		buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
		rightColumnContent.Add(buttonRow)
	}

	defer variable.CurrentDBClient.Close()

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
			log.Fatal("this err for func DeletKeyLogic part else delete || err : ", err)
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
			log.Fatal("error : this error in func addkeylogic for add key in database")
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
