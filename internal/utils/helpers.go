// internal/utils/helpers.go
package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	variable "testgui"
	"testgui/internal/Databaces/PebbleDB"
	badgerDB "testgui/internal/Databaces/badger"
	leveldbb "testgui/internal/Databaces/leveldb"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

var ValueImage []byte

type TappableLabel struct {
	widget.Label
	onTapped func()
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

func TruncateString(input string, length int) string {
	if len(input) > length {
		return input[:length] + "..."
	}
	return input
}

func IsValidJSON(data string) bool {
	var js json.RawMessage
	m := json.Unmarshal([]byte(data), &js) == nil
	return m
}

func CheckCondition(rightColumnContent *fyne.Container) {
	if len(rightColumnContent.Objects) > 0 {
		rightColumnContent.Objects = []fyne.CanvasObject{}
		rightColumnContent.Refresh()
	}
}

func Checkdatabace(test string, nameDatabace string) error {
	//parts := strings.Split(test, "|-|")

	switch nameDatabace {
	case "levelDB":
		variable.CurrentDBClient = leveldbb.NewDataBaseLeveldb(test)
	case "Pebble":
		variable.CurrentDBClient = PebbleDB.NewDataBasePebble(test)
	case "Badger":
		variable.CurrentDBClient = badgerDB.NewDataBaseBadger(test)
	case "Redis":

		//variable.CurrentDBClient = Redisdb.NewDataBaseRedis(parts[0], parts[1], parts[2])
	}

	if nameDatabace != "Redis" {

		if _, err := os.Stat(test); os.IsNotExist(err) && !variable.CreatDatabase {

			return err
		}
	}

	return nil
}

func CleanInput(input string) string {
	cleaned := strings.TrimSpace(input)
	return cleaned
}

func ImageShow(key []byte, value []byte, nameLable string, mainContainer *fyne.Container, editWindow fyne.Window) *fyne.Container {
	var contentt *fyne.Container
	var lableAddpicture *widget.Button

	imgReader := bytes.NewReader([]byte(value))
	image := canvas.NewImageFromReader(imgReader, "image.png")

	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(400, 400))

	mainContainer.Add(image)

	lableAddpicture = widget.NewButton("+", func() {
		folderPath := dialog.NewFileOpen(func(dir fyne.URIReadCloser, err error) {
			if err != nil {
				fmt.Println("Error opening folder:", err)
				return
			}
			if dir == nil {
				fmt.Println("No folder selected")
				return
			}

			valueFinish, err := ioutil.ReadAll(dir)
			if err != nil {
				fmt.Println("Error reading file:", err)
				return
			}
			ValueImage = valueFinish

			imgReader := bytes.NewReader(valueFinish)
			image := canvas.NewImageFromReader(imgReader, "image.png")

			image.FillMode = canvas.ImageFillContain
			image.SetMinSize(fyne.NewSize(400, 400))

			if len(mainContainer.Objects) >= 1 {
				mainContainer.Objects = mainContainer.Objects[:0]
			}

			mainContainer.Add(image)
			mainContainer.Refresh()

		}, editWindow)
		folderPath.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg", ".gif"}))

		folderPath.Show()
	})

	contentt = container.NewVBox(
		lableAddpicture,
	)
	return contentt
}
