// internal/utils/helpers.go
package utils

import (
	"encoding/json"
	variable "testgui"
	"testgui/internal/Databaces/PebbleDB"
	leveldbb "testgui/internal/Databaces/leveldb"

	"fyne.io/fyne/v2"
)

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
		newObjects := []fyne.CanvasObject{}
		rightColumnContent.Objects = newObjects
		rightColumnContent.Refresh()
	}
}

func Checkdatabace(test string, nameDatabace string) {
	if nameDatabace == "levelDB" {

		variable.CurrentDBClient = leveldbb.NewDataBaseLeveldb(test)
	} else if nameDatabace == "Pebble" {

		variable.CurrentDBClient = PebbleDB.NewDataBasePebble(test)
	}
}
