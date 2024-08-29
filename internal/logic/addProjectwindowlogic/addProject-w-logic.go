package addProjectwindowlogic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testgui/internal/logic/mainwindowlagic"
	"testgui/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/syndtr/goleveldb/leveldb"
)

type TappableLabel struct {
	widget.Label
	onTapped func()
}

func HasManifestFile(folderPath string) bool {
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		fmt.Println("Error reading folder:", err)
		return false
	}
	var count int64
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

func HandleButtonClick(test string) error {
	db, err := leveldb.OpenFile(test, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	iter := db.NewIterator(nil, nil)
	defer iter.Release()

	if iter.First() {
		key := iter.Key()
		value, err := db.Get(key, nil)
		if err != nil {
			return fmt.Errorf("failed to get value for key %s: %v", key, err)
		}

		fmt.Printf("First key: %s, value: %s\n", key, value)
		return nil
	}
	return fmt.Errorf("no entries found in the database")
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

func BuidLableKeyAndValue(eidtKeyAbdValue string, key string, value string, nameLable string, Addres string, rightColumnContent *fyne.Container) *TappableLabel {
	var lableKeyAndValue *TappableLabel

	lableKeyAndValue = NewTappableLabel(nameLable, func() {
		editWindow := fyne.CurrentApp().NewWindow("Edit" + eidtKeyAbdValue)
		editWindow.Resize(fyne.NewSize(600, 600))

		valueEntry := widget.NewMultiLineEntry()
		valueEntry.Resize(fyne.NewSize(500, 500))
		if eidtKeyAbdValue == "value" {
			if utils.IsValidJSON(value) {
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

			err := mainwindowlagic.CurrentDBClient.Open()
			defer mainwindowlagic.CurrentDBClient.Close()

			if eidtKeyAbdValue == "value" {
				err := mainwindowlagic.CurrentDBClient.Add(key, value)
				if err != nil {
					fmt.Println(err)
				}
				truncatedKey2 = utils.TruncateString(valueEntry.Text, 50)

			} else {
				valueBefor := mainwindowlagic.CurrentDBClient.Get(key)

				err = mainwindowlagic.CurrentDBClient.Delet(key)
				if err != nil {
					return
				}

				key = valueEntry.Text

				err := mainwindowlagic.CurrentDBClient.Add(key, valueBefor)
				if err != nil {
					fmt.Println(err)
				}
				truncatedKey2 = utils.TruncateString(key, 20)
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
