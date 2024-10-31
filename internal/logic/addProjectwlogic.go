package logic

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	variable "testgui"
	"testgui/internal/utils"

	// "testgui/internal/logic/mainwindowlagic"

	"fyne.io/fyne/v2/widget"
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
