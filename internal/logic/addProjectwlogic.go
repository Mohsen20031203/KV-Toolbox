package logic

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	variable "testgui"
	"testgui/internal/utils"

	// "testgui/internal/logic/mainwindowlagic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

var CreatFileBool bool

type TappableLabel struct {
	widget.Label
	onTapped func()
}

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

func NewTappableLabel(text string, tapped func()) *TappableLabel {
	label := &TappableLabel{
		Label: widget.Label{
			Text: text,
		},
		onTapped: tapped,
	}
	label.ExtendBaseWidget(label)
	return label
}

func (t *TappableLabel) Tapped(_ *fyne.PointEvent) {
	t.onTapped()
}

func CreatFile(value bool, openButton *widget.Button, testConnectionButton *widget.Button) {
	if value {
		openButton.Disable()
		testConnectionButton.Disable()
		CreatFileBool = value
		variable.CreatDatabase = value
	} else {
		openButton.Enable()
		testConnectionButton.Enable()
		CreatFileBool = value
		variable.CreatDatabase = value
	}
}

func HandleButtonClick(test string, nameDatabace string) error {

	err := utils.Checkdatabace(test, nameDatabace)
	if err != nil {
		return err
	}
	err, date := variable.CurrentDBClient.Read(nil, nil, 1)
	nun := HasManifestFile(test)

	if CreatFileBool {
		return nil
	}
	if (len(date) == 0 || err != nil) && !nun {
		return err
	}
	return nil

}
