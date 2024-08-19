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

func handleProjectSelection(dbPath string, rightColumnContent *fyne.Container) {

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
	const maxValueLength = 30

	// ایجاد دکمه‌ها برای هر رکورد و اضافه کردن آنها به ستون سمت راست
	for _, item := range data {
		// کوتاه کردن key و value در صورت نیاز
		truncatedKey := truncateString(item.key, maxKeyLength)
		truncatedValue := truncateString(item.value, maxValueLength)

		keyButton := widget.NewButton(truncatedKey, func() {
			editWindow := fyne.CurrentApp().NewWindow("Edit Value")
			editWindow.Resize(fyne.NewSize(600, 600))

			valueEntry := widget.NewMultiLineEntry()
			valueEntry.Resize(fyne.NewSize(500, 500))
			scrollableEntry := container.NewScroll(valueEntry)
			mainContainer := container.NewBorder(nil, nil, nil, nil, scrollableEntry)

			scrollableEntry.SetMinSize(fyne.NewSize(600, 500))
			valueEntry.SetText(item.key)

			saveButton := widget.NewButton("Save", func() {
				// ذخیره مقدار جدید
				item.key = valueEntry.Text
				editWindow.Close()
				rightColumnContent.Refresh()
			})

			cancelButton := widget.NewButton("Cancel", func() {
				editWindow.Close()
			})

			m := container.NewGridWithColumns(2, cancelButton, saveButton)
			b := container.NewBorder(nil, m, nil, nil)

			editContent := container.NewVBox(
				widget.NewLabel("Edit Value:"),
				mainContainer,
				layout.NewSpacer(),
				b,
			)

			editWindow.SetContent(editContent)
			editWindow.Show()
		})

		valueButton := widget.NewButton(truncatedValue, func() {
			editWindow := fyne.CurrentApp().NewWindow("Edit Value")
			editWindow.Resize(fyne.NewSize(600, 600))

			valueEntry := widget.NewMultiLineEntry()
			valueEntry.Resize(fyne.NewSize(500, 500))
			scrollableEntry := container.NewScroll(valueEntry)
			mainContainer := container.NewBorder(nil, nil, nil, nil, scrollableEntry)

			scrollableEntry.SetMinSize(fyne.NewSize(600, 500))
			valueEntry.SetText(item.value)

			saveButton := widget.NewButton("Save", func() {
				// ذخیره مقدار جدید
				item.value = valueEntry.Text
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

		buttonRow := container.NewGridWithColumns(2, keyButton, valueButton)
		rightColumnContent.Add(buttonRow)
	}

	rightColumnContent.Refresh()
}
