// internal/utils/helpers.go
package utils

import (
	variable "DatabaseDB"
	"DatabaseDB/internal/Databaces/PebbleDB"
	badgerDB "DatabaseDB/internal/Databaces/badger"
	leveldbb "DatabaseDB/internal/Databaces/leveldb"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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
	nameData := input
	if len(nameData) > length {
		nameData = nameData[:length] + ". . ."
	}
	parts := strings.Split(nameData, "\n")
	if len(parts) > 1 {

		nameData = parts[0] + " . . ."
	}

	return nameData
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

func ImageShow(key []byte, value []byte, mainContainer *fyne.Container, editWindow fyne.Window) {
	var lableAddpicture *widget.Button
	var image *canvas.Image

	image = canvas.NewImageFromResource(fyne.NewStaticResource("placeholder.png", value))
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(300, 300))
	mainContainer.Add(image)

	lableAddpicture = widget.NewButton("+", func() {
		folderPath := dialog.NewFileOpen(func(dir fyne.URIReadCloser, err error) {
			if err != nil || dir == nil {
				log.Fatal("Error opening folder or no folder selected")
				return
			}

			go func() {

				valueFinish, err := ioutil.ReadAll(dir)
				if err != nil {
					log.Fatal("Error reading file:", err)
					return
				}

				image.Resource = fyne.NewStaticResource("image.png", valueFinish)
				image.Refresh()
				ValueImage = valueFinish
			}()

		}, editWindow)

		folderPath.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg", ".gif"}))
		folderPath.Show()
	})
	mainContainer.Add(lableAddpicture)

}
