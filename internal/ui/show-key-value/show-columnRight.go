package main

import (
	"encoding/json"
	"fmt"
	dbpak "testgui/pkg"

	"testgui/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var currentDBClient dbpak.DBClient
var count int
var lastkey dbpak.Database

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
	if !utils.CheckCondition(rightColumnContent) {
		newObjects := []fyne.CanvasObject{}

		rightColumnContent.Objects = newObjects

		rightColumnContent.Refresh()
	}

	//currentDBClient = newFunc
	currentPage = 0
	prevButton.Disable()

	err, data := currentDBClient.Read()
	if err != nil {
		fmt.Println("Failed to read database:", err)
		return
	}

	for _, item := range data {
		if currentPage < len(data)/itemsPerPage {
			nextButton.Enable()
		}
		if count >= itemsPerPage {
			count = 0
			break
		}
		count++

		truncatedKey := utils.TruncateString(item.Key, 20)
		truncatedValue := utils.TruncateString(item.Value, 50)

		valueLabel := buidLableKeyAndValue("value", item.Key, item.Value, truncatedValue, dbPath, rightColumnContent)
		keyLabel := buidLableKeyAndValue("key", item.Key, item.Value, truncatedKey, dbPath, rightColumnContent)

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

			err := currentDBClient.Open()
			defer currentDBClient.Close()

			if eidtKeyAbdValue == "value" {
				err := currentDBClient.Add(key, value)
				if err != nil {
					fmt.Println(err)
				}
				truncatedKey2 = utils.TruncateString(valueEntry.Text, 50)

			} else {
				valueBefor := currentDBClient.Get(key)

				err = currentDBClient.Delet(key)
				if err != nil {
					return
				}

				key = valueEntry.Text

				err := currentDBClient.Add(key, valueBefor)
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
