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

	// محدودیت طول برای کلید و مقدار
	const maxKeyLength = 20
	const maxValueLength = 50

	// ایجاد برچسب‌ها برای هر رکورد و اضافه کردن آنها به ستون سمت راست
	for _, item := range data {
		// کوتاه کردن key و value در صورت نیاز
		truncatedKey := truncateString(item.key, maxKeyLength)
		truncatedValue := truncateString(item.value, maxValueLength)

		// تعریف متغیرهای برچسب‌ها به صورت قابل تغییر
		var keyLabel, valueLabel *TappableLabel

		// ایجاد برچسب قابل کلیک برای کلید
		keyLabel = NewTappableLabel(truncatedKey, func() {
			editWindow := fyne.CurrentApp().NewWindow("Edit Key")
			editWindow.Resize(fyne.NewSize(600, 600))

			valueEntry := widget.NewMultiLineEntry()
			valueEntry.Resize(fyne.NewSize(500, 500))
			scrollableEntry := container.NewScroll(valueEntry)
			mainContainer := container.NewBorder(nil, nil, nil, nil, scrollableEntry)

			scrollableEntry.SetMinSize(fyne.NewSize(600, 500))
			valueEntry.SetText(item.key)
			saveButton := widget.NewButton("Save", func() {
				db, err := leveldb.OpenFile(dbPath, nil)
				if err != nil {
					return
				}
				defer db.Close()

				valueBefor, err := db.Get([]byte(item.key), nil)
				if err != nil {
					return
				}

				err = db.Delete([]byte(item.key), nil)
				if err != nil {
					return
				}

				item.key = valueEntry.Text
				db.Put([]byte(item.key), []byte(valueBefor), nil)

				truncatedKey2 := truncateString(item.key, maxKeyLength)
				// بروز‌رسانی متن برچسب
				keyLabel.SetText(truncatedKey2)
				keyLabel.Refresh()

				editWindow.Close()
				rightColumnContent.Refresh()
			})

			cancelButton := widget.NewButton("Cancel", func() {
				editWindow.Close()
			})

			m := container.NewGridWithColumns(2, cancelButton, saveButton)
			b := container.NewBorder(nil, m, nil, nil)

			editContent := container.NewVBox(
				widget.NewLabel("Edit Key:"),
				mainContainer,
				layout.NewSpacer(),
				b,
			)

			editWindow.SetContent(editContent)
			editWindow.Show()
		})

		// ایجاد برچسب قابل کلیک برای مقدار
		valueLabel = NewTappableLabel(truncatedValue, func() {
			editWindow := fyne.CurrentApp().NewWindow("Edit Value")
			editWindow.Resize(fyne.NewSize(600, 600))

			valueEntry := widget.NewMultiLineEntry()
			valueEntry.Resize(fyne.NewSize(500, 500))
			scrollableEntry := container.NewScroll(valueEntry)
			mainContainer := container.NewBorder(nil, nil, nil, nil, scrollableEntry)

			scrollableEntry.SetMinSize(fyne.NewSize(600, 500))
			valueEntry.SetText(item.value)

			saveButton := widget.NewButton("Save", func() {
				db, err := leveldb.OpenFile(dbPath, nil)
				if err != nil {
					return
				}
				defer db.Close()
				db.Put([]byte(item.key), []byte(valueEntry.Text), nil)

				// ذخیره مقدار جدید
				item.value = valueEntry.Text

				truncatedValue2 := truncateString(item.value, maxValueLength)
				// بروز‌رسانی متن برچسب
				valueLabel.SetText(truncatedValue2)
				valueLabel.Refresh()

				editWindow.Close()
				rightColumnContent.Refresh()
			})

			cancelButton := widget.NewButton("Cancel", func() {
				editWindow.Close()
			})

			bottomButtons := container.NewGridWithColumns(2, cancelButton, saveButton)
			positionBottomButtons := container.NewBorder(nil, bottomButtons, nil, nil)

			editContent := container.NewVBox(
				widget.NewLabel("Edit Value:"),
				mainContainer,
				layout.NewSpacer(),
				positionBottomButtons,
			)

			editWindow.SetContent(editContent)
			editWindow.Show()
		})

		// افزودن برچسب‌ها به شبکه
		buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
		rightColumnContent.Add(buttonRow)
	}

	rightColumnContent.Refresh()
}
