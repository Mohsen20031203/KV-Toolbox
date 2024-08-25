package main

import (
	"encoding/json"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/syndtr/goleveldb/leveldb"
)

var count int
var lastkey datebace

type datebace struct {
	key   string
	value string
}

func isValidJSON(data string) bool {
	var js json.RawMessage
	m := json.Unmarshal([]byte(data), &js) == nil
	return m
}

func checkCondition(rightColumnContent *fyne.Container) bool {
	if len(rightColumnContent.Objects) > 2 {
		return false
	}
	return true
}

func readDatabace(Addres string) (error, []datebace) {
	var Item []datebace
	db, err := leveldb.OpenFile(Addres, nil)
	if err != nil {
		return err, Item
	}
	defer db.Close()

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		key := string(iter.Key())
		value := string(iter.Value())
		Item = append(Item, datebace{key: key, value: value})
	}
	iter.Release()

	return nil, Item
}

func truncateString(str string, maxLength int) string {
	if len(str) > maxLength {
		return str[:maxLength] + "..."
	}
	return str
}

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

func handleProjectSelection(dbPath string, rightColumnContent *fyne.Container, buttonAdd *widget.Button) {

	buttonAdd.Enable()
	if !checkCondition(rightColumnContent) {
		newObjects := []fyne.CanvasObject{}

		rightColumnContent.Objects = newObjects

		rightColumnContent.Refresh()
	}

	err, data := readDatabace(dbPath)
	if err != nil {
		fmt.Println("Failed to read database:", err)
		return
	}

	for _, item := range data {
		if count >= itemsPerPage {
			count = 0
			break
		}
		count++

		truncatedKey := truncateString(item.key, 20)
		truncatedValue := truncateString(item.value, 50)

		valueLabel := buidLableKeyAndValue("value", item.key, item.value, truncatedValue, dbPath, rightColumnContent)
		keyLabel := buidLableKeyAndValue("key", item.key, item.value, truncatedKey, dbPath, rightColumnContent)

		buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
		rightColumnContent.Add(buttonRow)
		lastkey = item

	}

	rightColumnContent.Refresh()
}

func buidLableKeyAndValue(eidtKeyAbdValue string, key string, value string, nameLable string, Addres string, rightColumnContent *fyne.Container) *TappableLabel {
	var lableKeyAndValue *TappableLabel

	lableKeyAndValue = NewTappableLabel(nameLable, func() {
		editWindow := fyne.CurrentApp().NewWindow("Edit" + eidtKeyAbdValue)
		editWindow.Resize(fyne.NewSize(600, 600))

		valueEntry := widget.NewMultiLineEntry()
		valueEntry.Resize(fyne.NewSize(500, 500))
		if eidtKeyAbdValue == "value" {
			if isValidJSON(value) {
				var formattedJSON map[string]interface{}
				json.Unmarshal([]byte(value), &formattedJSON)
				jsonString, _ := json.MarshalIndent(formattedJSON, "", "  ")
				valueEntry.SetText(string(jsonString))
			} else {
				valueEntry.SetText(value)
			}
		} else {
			valueEntry.SetText(key)
		}
		scrollableEntry := container.NewScroll(valueEntry)
		mainContainer := container.NewBorder(nil, nil, nil, nil, scrollableEntry)

		scrollableEntry.SetMinSize(fyne.NewSize(600, 500))
		saveButton := widget.NewButton("Save", func() {
			var truncatedKey2 string
			db, err := leveldb.OpenFile(Addres, nil)
			if err != nil {
				return
			}
			defer db.Close()

			if eidtKeyAbdValue == "value" {
				db.Put([]byte(key), []byte(valueEntry.Text), nil)
				truncatedKey2 = truncateString(valueEntry.Text, 50)

			} else {
				valueBefor, err := db.Get([]byte(key), nil)
				if err != nil {
					return
				}

				err = db.Delete([]byte(key), nil)
				if err != nil {
					return
				}

				key = valueEntry.Text
				db.Put([]byte(key), []byte(valueBefor), nil)
				truncatedKey2 = truncateString(key, 20)
			}

			lableKeyAndValue.SetText(truncatedKey2)
			lableKeyAndValue.Refresh()

			editWindow.Close()
			rightColumnContent.Refresh()
		})

		cancelButton := widget.NewButton("Cancel", func() {
			editWindow.Close()
		})

		m := container.NewGridWithColumns(2, cancelButton, saveButton)
		b := container.NewBorder(nil, m, nil, nil)

		editContent := container.NewVBox(
			widget.NewLabel("Edit "+eidtKeyAbdValue+" :"),
			mainContainer,
			layout.NewSpacer(),
			b,
		)

		editWindow.SetContent(editContent)
		editWindow.Show()
	})
	return lableKeyAndValue
}
