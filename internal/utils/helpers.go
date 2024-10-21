// internal/utils/helpers.go
package utils

import (
	"encoding/json"
	"os"
	"strings"
	variable "testgui"
	"testgui/internal/Databaces/PebbleDB"
	Redisdb "testgui/internal/Databaces/Redis"
	badgerDB "testgui/internal/Databaces/badger"
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
