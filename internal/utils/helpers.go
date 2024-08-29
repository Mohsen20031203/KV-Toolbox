// internal/utils/helpers.go
package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	mainwl "testgui/internal/logic/main-window-lagic"
	mainW "testgui/internal/ui/main-window"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// TruncateString کوتاه کردن رشته به طول مشخص
func TruncateString(input string, length int) string {
	if len(input) > length {
		return input[:length] + "..."
	}
	return input
}

// SanitizeString پاک‌سازی رشته
func SanitizeString(input string) string {
	return strings.TrimSpace(input)
}

// SanitizeString پاک‌سازی رشته
func WriteJsonFile(file *os.File, state interface{}) error {
	file.Truncate(0)
	file.Seek(0, 0)
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(&state); err != nil {
		return fmt.Errorf("failed to encode JSON: %v", err)
	}
	return nil
}

func IsValidJSON(data string) bool {
	var js json.RawMessage
	m := json.Unmarshal([]byte(data), &js) == nil
	return m
}

func CheckCondition(rightColumnContent *fyne.Container) bool {
	if len(rightColumnContent.Objects) > 2 {
		return false
	}
	return true
}

func ProjectButton(inputText string, lastColumnContent *fyne.Container, path string, rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button) *fyne.Container {
	projectButton := widget.NewButton(inputText, func() {
		mainW.PageLabel.Text = "Page 1"
		mainwl.FolderPath = path
		HandleProjectSelection(path, rightColumnContentORG, buttonAdd)
		if nameButtonProject.Text == "" {
			nameButtonProject.Text = inputText
		} else {
			nameButtonProject.Text = ""
			nameButtonProject.Text = inputText
		}

		nameButtonProject.Refresh()
		mainW.PageLabel.Refresh()

	})

	if nameButtonProject.Text == "" {
		nameButtonProject.Text = inputText
	} else {
		nameButtonProject.Text = ""
		nameButtonProject.Text = inputText
	}
	nameButtonProject.Refresh()

	buttonContainer := container.NewHBox()

	closeButton := widget.NewButton("✖", func() {

		if !CheckCondition(rightColumnContentORG) && nameButtonProject.Text == inputText {
			newObjects := []fyne.CanvasObject{}

			rightColumnContentORG.Objects = newObjects
			buttonAdd.Disable()

			nameButtonProject.Text = ""
			nameButtonProject.Refresh()
			rightColumnContentORG.Refresh()
		}

		err := mainwl.CurrentJson
		if err != nil {
			fmt.Println("Failed to remove project from JSON:", err)
		} else {

			lastColumnContent.Remove(buttonContainer)
			lastColumnContent.Refresh()
		}
	})

	buttonContainer = container.NewBorder(nil, nil, nil, closeButton, projectButton)
	return buttonContainer
}

func ProjectButton(inputText string, lastColumnContent *fyne.Container, path string, rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button) *fyne.Container {
	projectButton := widget.NewButton(inputText, func() {
		mainW.PageLabel.Text = "Page 1"
		FolderPath = path
		HandleProjectSelection(path, rightColumnContentORG, buttonAdd)
		if nameButtonProject.Text == "" {
			nameButtonProject.Text = inputText
		} else {
			nameButtonProject.Text = ""
			nameButtonProject.Text = inputText
		}

		CurrentDBClient = leveldbb.NewDataBase(path)

		nameButtonProject.Refresh()
		mainW.PageLabel.Refresh()

	})

	if nameButtonProject.Text == "" {
		nameButtonProject.Text = inputText
	} else {
		nameButtonProject.Text = ""
		nameButtonProject.Text = inputText
	}
	nameButtonProject.Refresh()

	buttonContainer := container.NewHBox()

	closeButton := widget.NewButton("✖", func() {

		if !utils.CheckCondition(rightColumnContentORG) && nameButtonProject.Text == inputText {
			newObjects := []fyne.CanvasObject{}

			rightColumnContentORG.Objects = newObjects
			buttonAdd.Disable()

			nameButtonProject.Text = ""
			nameButtonProject.Refresh()
			rightColumnContentORG.Refresh()
		}

		err := CurrentJson.Remove(inputText)
		if err != nil {
			fmt.Println("Failed to remove project from JSON:", err)
		} else {

			lastColumnContent.Remove(buttonContainer)
			lastColumnContent.Refresh()
		}
	})

	buttonContainer = container.NewBorder(nil, nil, nil, closeButton, projectButton)
	return buttonContainer
}
