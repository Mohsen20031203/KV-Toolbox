package main

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/syndtr/goleveldb/leveldb"
)

func projectButton(inputText string, lastColumnContent *fyne.Container, path string, rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label) *fyne.Container {
	projectButton := widget.NewButton(inputText, func() {
		handleProjectSelection(path, rightColumnContentORG)
		if nameButtonProject.Text == "" {
			nameButtonProject.Text = inputText
		} else {
			nameButtonProject.Text = ""
			nameButtonProject.Text = inputText
		}
		nameButtonProject.Refresh()
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

		if !checkCondition(rightColumnContentORG) && nameButtonProject.Text == inputText {
			newObjects := []fyne.CanvasObject{}

			// حذف تمام ویجت‌ها از کانتینر
			rightColumnContentORG.Objects = newObjects

			// بروزرسانی محتوا
			nameButtonProject.Text = ""
			nameButtonProject.Refresh()
			rightColumnContentORG.Refresh()
		}

		err := removeProjectFromJsonFile(inputText)
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

func openNewWindow(a fyne.App, title string, lastColumnContent *fyne.Container, rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label) {
	newWindow := a.NewWindow(title)

	createSeparator := func() *canvas.Line {
		line := canvas.NewLine(color.Black)
		line.StrokeWidth = 1
		return line
	}
	line1 := createSeparator()

	lableName := widget.NewLabel("Name :")
	pathEntry := widget.NewEntry()
	pathEntry.PlaceHolder = "Name"
	nameContent := container.NewBorder(nil, nil, lableName, nil, pathEntry)

	lableComment := widget.NewLabel("Comment :")
	pathEntryComment := widget.NewEntry()
	pathEntryComment.PlaceHolder = "Comment"
	commentContent := container.NewBorder(nil, nil, lableComment, nil, pathEntryComment)

	pathEntry2 := widget.NewEntry()
	pathEntry2.SetPlaceHolder("No folder selected")

	testConnectionButton := widget.NewButton("Test Connection", func() {

		err := handleButtonClick(pathEntry2.Text)
		if err != nil {
			dialog.ShowError(err, newWindow)
		} else {
			dialog.ShowInformation("Success", "Test connection successful.", newWindow)
		}
	})
	testConnectionButton.Disable()

	pathEntry2.OnChanged = func(text string) {
		if text != "" {
			testConnectionButton.Enable()
		} else {
			testConnectionButton.Disable()
		}
	}

	openButton := widget.NewButton("Open Folder", func() {
		folderDialog := dialog.NewFileOpen(func(dir fyne.URIReadCloser, err error) {
			if err != nil {
				fmt.Println("Error opening folder:", err)
				return
			}
			if dir == nil {
				fmt.Println("No folder selected")
				return
			}
			filePath := dir.URI().Path()

			folderPath := filepath.Dir(filePath)

			if hasManifestFile(folderPath) {
				pathEntry2.SetText(folderPath)
				testConnectionButton.Enable()
			} else {
				dialog.ShowInformation("Invalid Folder", "The selected folder does not contain a valid LevelDB manifest file.", newWindow)
			}

		}, newWindow)
		folderDialog.SetFilter(storage.NewExtensionFileFilter([]string{".log"}))
		folderDialog.Show()
	})

	testOpenButton := container.NewVBox(
		layout.NewSpacer(),
		container.NewGridWithColumns(2, openButton, testConnectionButton),
	)

	buttonCancel := widget.NewButton("Cancel", func() {
		newWindow.Close()
	})

	buttonOk := widget.NewButton("Add", func() {
		err, addButton := addProjectToJsonFile(pathEntry2, pathEntry, pathEntryComment, newWindow)
		if err != nil {
			dialog.ShowInformation("Error ", "There is something wrong with your file and I can't connect to it", newWindow)
		} else {

			if !addButton {

				buttonContainer := projectButton(pathEntry.Text, lastColumnContent, pathEntry2.Text, rightColumnContentORG, nameButtonProject)
				lastColumnContent.Add(buttonContainer)
				lastColumnContent.Refresh()

				handleProjectSelection(pathEntry2.Text, rightColumnContentORG)
				rightColumnContentORG.Refresh()

				newWindow.Close()
			}
		}
	})

	rowBotton := container.NewVBox(
		layout.NewSpacer(),
		container.NewGridWithColumns(2, buttonCancel, buttonOk),
	)

	rightColumnContent := container.NewVBox(
		layout.NewSpacer(),
		nameContent,
		layout.NewSpacer(),
		commentContent,
		layout.NewSpacer(),
		line1,
		layout.NewSpacer(),
		pathEntry2,
		layout.NewSpacer(),
		testOpenButton,
		layout.NewSpacer(),
		rowBotton,
	)

	newWindow.Resize(fyne.NewSize(700, 400))
	newWindow.CenterOnScreen()
	newWindow.SetContent(rightColumnContent)
	newWindow.Show()
}

func handleButtonClick(test string) error {
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

func hasManifestFile(folderPath string) bool {
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
