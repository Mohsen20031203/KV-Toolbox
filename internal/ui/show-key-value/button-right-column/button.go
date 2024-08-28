package buttonrightcolumn

import (
	"fmt"
	"testgui/internal/utils"

	"fyne.io/fyne/container"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

var CurrentPage int
var ItemsPerPage = 20
var NextButton, PrevButton *widget.Button
var PageLabel *widget.Label

func UpdatePage(rightColumnContent *fyne.Container) {
	if !utils.CheckCondition(rightColumnContent) {
		rightColumnContent.Objects = []fyne.CanvasObject{}
		rightColumnContent.Refresh()
	}

	err, data := currentDBClient.Read()
	if err != nil {
		fmt.Println(err)
	}

	startIndex := currentPage * itemsPerPage
	endIndex := startIndex + itemsPerPage

	if endIndex > len(data) {
		endIndex = len(data)
	}

	// پاک کردن محتوای قبلی
	rightColumnContent.Objects = nil

	for _, item := range data[startIndex:endIndex] {
		truncatedKey := truncateString(item.Key, 20)
		truncatedValue := truncateString(item.Value, 50)

		valueLabel := buidLableKeyAndValue("value", item.Key, item.Value, truncatedValue, folderPath, rightColumnContent)
		keyLabel := buidLableKeyAndValue("key", item.Key, item.Value, truncatedKey, folderPath, rightColumnContent)

		buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
		rightColumnContent.Add(buttonRow)
	}

	// به‌روزرسانی شماره صفحه
	PageLabel.SetText(fmt.Sprintf("Page %d", currentPage+1))

	// غیرفعال کردن دکمه‌ها بر اساس موقعیت فعلی
	prevButton.Disable()
	NextButton.Disable()

	if currentPage > 0 {
		prevButton.Enable()
	}
	if endIndex < len(data) {
		NextButton.Enable()
	}

	rightColumnContent.Refresh()
}
