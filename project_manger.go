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

func checkCondition(rightColumnContent *fyne.Container) bool {
	if len(rightColumnContent.Objects) > 2 {
		return false
	}
	return true
}

type datebace struct {
	key   string
	value string
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

// یک تابع برای کوتاه کردن متن و اضافه کردن ... در انتهای آن
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

		// اضافه کردن دکمه‌ها به ستون سمت راست
		buttonRow := container.NewGridWithColumns(2, keyButton, valueButton)
		rightColumnContent.Add(buttonRow)
	}

	rightColumnContent.Refresh()
}

func removeProjectFromJsonFile(projectName string) error {
	file, err := openFileJson()
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	var state *JsonInformation

	err = readJsonFile(file, &state)
	if err != nil {
		return err
	}

	for i, project := range state.RecentProjects {
		if project.Name == projectName {
			state.RecentProjects = append(state.RecentProjects[:i], state.RecentProjects[i+1:]...)
			break
		}
	}

	err = writeJsonFile(file, &state)
	if err != nil {
		return fmt.Errorf("failed to decode JSON: %v", err)
	}

	return nil
}

func projectButton(inputText string, lastColumnContent *fyne.Container, path string, rightColumnContentORG *fyne.Container) *fyne.Container {
	projectButton := widget.NewButton(inputText, func() {
		handleProjectSelection(path, rightColumnContentORG)

	})

	buttonContainer := container.NewHBox()

	closeButton := widget.NewButton("✖", func() {
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

func addProjectToJsonFile(projectPath *widget.Entry, name *widget.Entry, comment *widget.Entry, Window fyne.Window) (error, bool) {
	file, err := openFileJson()
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	err = handleButtonClick(projectPath.Text)
	if err != nil {
		return err, false
	}

	var state *JsonInformation

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err), false
	}

	if fileInfo.Size() == 0 {
		state = &JsonInformation{
			RecentProjects: []Project{},
		}
	} else {
		err := readJsonFile(file, &state)
		if err != nil {
			return err, false
		}
	}

	for _, addres := range state.RecentProjects {
		if projectPath.Text == addres.FileAddress {
			m := fmt.Sprintf("This database has already been added to your projects under the name '%s'", addres.Name)
			dialog.ShowInformation("error", m, Window)

			err = writeJsonFile(file, state)
			if err != nil {
				return fmt.Errorf("failed to decode JSON: %v", err), false
			}
			return nil, true
		}
	}
	newActivity := Project{
		Name:        name.Text,
		Comment:     comment.Text,
		FileAddress: projectPath.Text,
	}

	state.RecentProjects = append(state.RecentProjects, newActivity)

	err = writeJsonFile(file, state)
	if err != nil {
		return fmt.Errorf("failed to decode JSON: %v", err), false
	}
	return nil, false
}

func openNewWindow(a fyne.App, title string, lastColumnContent *fyne.Container, rightColumnContentORG *fyne.Container) {
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

				buttonContainer := projectButton(pathEntry.Text, lastColumnContent, pathEntry2.Text, rightColumnContentORG)
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
