// internal/utils/helpers.go
package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
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

func Checkdatabace(test string, nameDatabace string) error {
	if nameDatabace == "levelDB" {

		variable.CurrentDBClient = leveldbb.NewDataBaseLeveldb(test)

	} else if nameDatabace == "Pebble" {

		variable.CurrentDBClient = PebbleDB.NewDataBasePebble(test)

	}
	if _, err := os.Stat(test); os.IsNotExist(err) && !variable.CreatDatabase {

		return fmt.Errorf("dont found file")
	} else {

		err := variable.CurrentDBClient.Open()
		if err != nil {
			return fmt.Errorf("i cant connection in database")
		}
		variable.CurrentDBClient.Close()
		return nil
	}

}

func CleanInput(input string) string {
	cleaned := strings.TrimSpace(input)
	return cleaned
}
