package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/syndtr/goleveldb/leveldb"
)

type datebace struct {
	key   string
	value string
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
		newObjects := []fyne.CanvasObject{
			rightColumnContent.Objects[0], // ویجت اول
			rightColumnContent.Objects[1], // ویجت دوم
		}

		// حذف تمام ویجت‌ها از کانتینر
		rightColumnContent.Objects = newObjects

		// بروزرسانی محتوا
		rightColumnContent.Refresh()
	}
	// خواندن داده‌ها از دیتابیس
	err, data := readDatabace(dbPath)
	if err != nil {
		fmt.Println("Failed to read database:", err)
		return
	}

	for _, item := range data {
		// کوتاه کردن key و value در صورت نیاز
		truncatedKey := truncateString(item.key, 20)
		truncatedValue := truncateString(item.value, 50)

		valueLabel := buidLableKeyAndValue("value", item.key, item.value, truncatedValue, dbPath, rightColumnContent)
		keyLabel := buidLableKeyAndValue("key", item.key, item.value, truncatedKey, dbPath, rightColumnContent)

		// افزودن برچسب‌ها به شبکه
		buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
		rightColumnContent.Add(buttonRow)
	}

	rightColumnContent.Refresh()
}

// valueOrKey ---> if valueOrKey = true ----> value -&&- if valueOrKey = false ----> key
func buidLableKeyAndValue(eidtKeyAbdValue string, key string, value string, nameLable string, Addres string, rightColumnContent *fyne.Container) *TappableLabel {
	var lableKeyAndValue *TappableLabel

	// ایجاد برچسب قابل کلیک برای کلید
	lableKeyAndValue = NewTappableLabel(nameLable, func() {
		editWindow := fyne.CurrentApp().NewWindow("Edit" + eidtKeyAbdValue)
		editWindow.Resize(fyne.NewSize(600, 600))

		valueEntry := widget.NewMultiLineEntry()
		valueEntry.Resize(fyne.NewSize(500, 500))
		if eidtKeyAbdValue == "value" {

			valueEntry.SetText(value)
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

			// بروز‌رسانی متن برچسب
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
