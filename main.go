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

func addProjectToJsonFile(file *os.File, projectPath *widget.Entry, name *widget.Entry, comment *widget.Entry) error {

	err := handleButtonClick(projectPath.Text)
	if err != nil {
		return err
	}
	var state JsonInformation
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&state); err != nil {
		return fmt.Errorf("failed to decode JSON: %v", err)
	}

	newactivity := Project{
		Name:        name.Text,
		Comment:     comment.Text,
		FileAddress: projectPath.Text,
	}

	state.RecentProjects = append(state.RecentProjects, newactivity)

	file.Truncate(0) // پاک کردن محتوای فعلی فایل
	file.Seek(0, 0)  // بازنشانی موقعیت فایل به ابتدای آن

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ") // تعیین فرمت خوانا برای JSON
	if err := encoder.Encode(&state); err != nil {
		return fmt.Errorf("failed to encode JSON: %v", err)
	}

	return nil
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

func openNewWindow(a fyne.App, title string, fileJson *os.File) {
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
				testConnectionButton.Enable() // Enable after selecting a valid folder
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
		err := addProjectToJsonFile(fileJson, pathEntry2, pathEntry, pathEntryComment)
		if err != nil {
			dialog.ShowInformation("Invalid Folder", "The selected folder does not contain a valid LevelDB manifest file.", newWindow)
			newWindow.Close()
		}
		newWindow.Close()
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
	lastColumnContent := container.NewVBox()

	columns := container.NewHSplit(lastColumnContent, rightColumnContent)
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

	file, err := os.OpenFile("data.json", os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// --- Dropdown and Buttons ---
	plus := widget.NewSelect([]string{"levelDB", "DynamoDB", "SQL", "DocumentDB"}, func(value string) {
		log.Println("Select set to", value)
		openNewWindow(myApp, value, file)
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

	// --- Theme Buttons ---
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

	// --- Layout ---
	lastColumnContent := container.NewVBox(
		plus,
		buttonOne,
		buttonTwo,
		buttonThree,
		buttonFor,
		layout.NewSpacer(),
		darkLight,
	)

	rightColumnContent := container.NewVBox()

	columns := container.NewHSplit(lastColumnContent, rightColumnContent)
	columns.SetOffset(0.25)

	scroll := container.NewScroll(columns)
	myWindow.CenterOnScreen()
	myWindow.SetContent(scroll)
	myWindow.Resize(fyne.NewSize(1200, 800))
	myWindow.ShowAndRun()
}
