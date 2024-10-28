// internal/utils/helpers.go
package utils

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	variable "testgui"
	"testgui/internal/Databaces/PebbleDB"
	Redisdb "testgui/internal/Databaces/Redis"
	badgerDB "testgui/internal/Databaces/badger"
	leveldbb "testgui/internal/Databaces/leveldb"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

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
	parts := strings.Split(test, "|-|")

	switch nameDatabace {
	case "levelDB":
		variable.CurrentDBClient = leveldbb.NewDataBaseLeveldb(parts[0])
	case "Pebble":
		variable.CurrentDBClient = PebbleDB.NewDataBasePebble(parts[0])
	case "Badger":
		variable.CurrentDBClient = badgerDB.NewDataBaseBadger(parts[0])
	case "Redis":

		variable.CurrentDBClient = Redisdb.NewDataBaseRedis(parts[0], parts[1], parts[2])
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

func ImageShow(key []byte, value []byte, nameLable string, mainContainer *fyne.Container, editWindow fyne.Window) {

	imgReader := bytes.NewReader([]byte(value))
	image := canvas.NewImageFromReader(imgReader, "image.png")

	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(400, 400))

	mainContainer.Add(image)
}
