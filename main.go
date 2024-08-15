package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"github.com/syndtr/goleveldb/leveldb"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type JsonInformation struct {
	RecentProjects []Project `json:"recentProjects"`
}

type Project struct {
	Name        string `json:"name"`
	Comment     string `json:"comment"`
	FileAddress string `json:"fileAddress"`
}

func writeJsonFile(file *os.File, state interface{}) error {
	file.Truncate(0)
	file.Seek(0, 0)
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(&state); err != nil {
		return fmt.Errorf("failed to encode JSON: %v", err)
	}
	return nil
}

func readJsonFile(file *os.File, state interface{}) error {
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&state); err != nil {
		return err
	}
	return nil
}

func openFileJson() (*os.File, error) {
	file, err := os.OpenFile("data.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return file, err
	}
	return file, nil
}

func removeProjectFromJsonFile(projectName string) error {

	file, err := openFileJson()
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	var state JsonInformation

	err = readJsonFile(file, state)
	if err != nil {
		return err
	}

	// جستجوی پروژه و حذف آن
	for i, project := range state.RecentProjects {
		if project.Name == projectName {
			state.RecentProjects = append(state.RecentProjects[:i], state.RecentProjects[i+1:]...)
			break
		}
	}

	err = writeJsonFile(file, state)
	if err != nil {
		return fmt.Errorf("failed to decode JSON: %v", err)
	}

	return nil
}

func projectButton(inputText string, lastColumnContent *fyne.Container) *fyne.Container {
	projectButton := widget.NewButton(inputText, func() {
		fmt.Println("Selected project:", inputText)
	})

	buttonContainer := container.NewHBox()

	// دکمه ضربدر برای حذف دکمه
	closeButton := widget.NewButton("✖", func() {
		err := removeProjectFromJsonFile(inputText)
		if err != nil {
			fmt.Println("Failed to remove project from JSON:", err)
		} else {
			// حذف دکمه از UI
			lastColumnContent.Remove(buttonContainer)
			lastColumnContent.Refresh()
		}
	})

	// قرار دادن دکمه اصلی و دکمه ضربدر در یک کانتینر افقی
	buttonContainer = container.NewBorder(nil, nil, nil, closeButton, projectButton)

	return buttonContainer
}

func loadJsonData(fileName string) (JsonInformation, error) {
	var jsonData JsonInformation

	file, err := os.Open(fileName)
	if err != nil {
		return jsonData, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return jsonData, fmt.Errorf("failed to read file: %v", err)
	}

	if err := json.Unmarshal(byteValue, &jsonData); err != nil {
		return jsonData, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return jsonData, nil
}

func addProjectToJsonFile(projectPath *widget.Entry, name *widget.Entry, comment *widget.Entry) (error, bool) {

	file, err := openFileJson()
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	err = handleButtonClick(projectPath.Text)
	if err != nil {
		return err, false
	}

	var state JsonInformation

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err), false
	}

	if fileInfo.Size() == 0 {
		state = JsonInformation{
			RecentProjects: []Project{},
		}
	} else {
		err := readJsonFile(file, state)
		if err != nil {
			return err, false
		}
	}

	for i, addres := range state.RecentProjects {
		if projectPath.Text == addres.FileAddress {
			state.RecentProjects[i].Comment = comment.Text
			file.Truncate(0)
			file.Seek(0, 0)

			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "    ")
			if err := encoder.Encode(&state); err != nil {
				return fmt.Errorf("failed to encode JSON: %v", err), false
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
		return fmt.Errorf("failed to decode JSON: %v", err), true
	}

	return nil, false
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

func openNewWindow(a fyne.App, title string, lastColumnContent *fyne.Container) {
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
	buttonApply := widget.NewButton("Apply", func() {
		fmt.Println("buttonApply")
	})
	buttonOk := widget.NewButton("Ok", func() {
		err, addButton := addProjectToJsonFile(pathEntry2, pathEntry, pathEntryComment)
		if err != nil {
			dialog.ShowInformation("Invalid Folder", "The selected folder does not contain a valid LevelDB manifest file.", newWindow)
			newWindow.Close()
		} else {

			if !addButton {

				buttonContainer := projectButton(pathEntry.Text, lastColumnContent)

				lastColumnContent.Add(buttonContainer)
				lastColumnContent.Refresh()

			}

			newWindow.Close()
		}
	})

	rowBotton := container.NewVBox(
		layout.NewSpacer(),
		container.NewGridWithColumns(3, buttonCancel, buttonApply, buttonOk),
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
	lastColumnContentInWindow := container.NewVBox()

	columns := container.NewHSplit(lastColumnContentInWindow, rightColumnContent)
	columns.SetOffset(0.25)
	newWindow.Resize(fyne.NewSize(900, 500))
	newWindow.CenterOnScreen()
	newWindow.SetContent(columns)
	newWindow.Show()
}

func main() {

	myApp := app.New()
	myWindow := myApp.NewWindow("Non-Scrollable List")

	iconResource := theme.FyneLogo()
	myApp.SetIcon(iconResource)
	myWindow.SetIcon(iconResource)

	lastColumnContent := container.NewVBox() // Move this initialization before `openNewWindow` is called

	plus := widget.NewSelect([]string{"levelDB", "DynamoDB", "SQL", "DocumentDB"}, func(value string) {
		log.Println("Select set to", value)
		openNewWindow(myApp, value, lastColumnContent) // Use the initialized `lastColumnContent`
	})
	plus.PlaceHolder = "Data Source"
	plus.Alignment = fyne.TextAlignCenter

	buttonOne := widget.NewButton("DDL data source", func() {
		fmt.Println("buttonOne")
	})

	buttonTwo := widget.NewButton("Data source from URL", func() {
		fmt.Println("buttonTwo")
	})

	buttonThree := widget.NewButton("Data source from path", func() {
		fmt.Println("buttonThree")
	})

	buttonFor := widget.NewButton("Import from Clipboard", func() {
		fmt.Println("buttonFor")
	})

	spacer := widget.NewLabel("")
	spacer.Resize(fyne.NewSize(0, 30))

	plusContent := container.NewVBox(
		plus,
		buttonOne,
		buttonTwo,
		buttonThree,
		buttonFor,
	)

	pluss := widget.NewButton("+", func() {
		if plusContent.Visible() {
			plusContent.Hide()
		} else {
			plusContent.Show()
		}
	})
	lastColumnContentt := container.NewVBox(
		pluss,
		plusContent,
		spacer,
	)
	plusContent.Hide()

	darkButton := widget.NewButton("Dark", func() {
		myApp.Settings().SetTheme(theme.DarkTheme())
	})
	lightButton := widget.NewButton("Light", func() {
		myApp.Settings().SetTheme(theme.LightTheme())
	})

	darkLight := container.NewVBox(
		layout.NewSpacer(),
		container.NewGridWithColumns(2, lightButton, darkButton),
	)

	jsonData, err := loadJsonData("data.json")
	if err != nil {
		fmt.Println("Error loading JSON data:", err)
	} else {
		for _, project := range jsonData.RecentProjects {

			buttonContainer := projectButton(project.Name, lastColumnContent)

			lastColumnContent.Add(buttonContainer)
		}
	}

	rightColumnContent := container.NewVBox()

	lastColumnScrollable := container.NewScroll(lastColumnContent)
	rightColumnScrollable := container.NewScroll(rightColumnContent)

	mainContent := container.NewBorder(lastColumnContentt, darkLight, nil, nil, lastColumnScrollable)

	columns := container.NewHSplit(mainContent, rightColumnScrollable)
	columns.SetOffset(0.25)

	scroll := container.NewScroll(columns)
	myWindow.CenterOnScreen()
	myWindow.SetContent(scroll)
	myWindow.Resize(fyne.NewSize(1200, 800))
	myWindow.ShowAndRun()
}
