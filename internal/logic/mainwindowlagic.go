package logic

import (
	"fmt"
	"log"
	"strings"
	variable "testgui"
	"time"

	// "testgui/internal/logic/addProjectwindowlogic"

	dbpak "testgui/internal/Databaces"
	"testgui/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gabriel-vasile/mimetype"
)

func SetupLastColumn(rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button) *fyne.Container {
	lastColumnContent := container.NewVBox()

	jsonDataa, err := variable.CurrentJson.Load()
	if err != nil {
		println("Error loading JSON data:", err)
	} else {
		for _, project := range jsonDataa.RecentProjects {

			buttonContainer := ProjectButton(project.Name, lastColumnContent, project.FileAddress, rightColumnContentORG, nameButtonProject, buttonAdd, project.Databace)
			lastColumnContent.Add(buttonContainer)
		}
	}

	return lastColumnContent
}

func SetupThemeButtons(app fyne.App) *fyne.Container {
	darkButton := widget.NewButton("Dark", func() {
		app.Settings().SetTheme(theme.DarkTheme())
	})
	lightButton := widget.NewButton("Light", func() {
		app.Settings().SetTheme(theme.LightTheme())
	})

	darkLight := container.NewVBox(
		layout.NewSpacer(),
		container.NewGridWithColumns(2, lightButton, darkButton),
	)
	return darkLight
}

var (
	lastStart *[]byte
	lastEnd   *[]byte
	count     int
	Orgdata   []dbpak.KVData
	lastPage  int
)

func UpdatePage(rightColumnContent *fyne.Container) {

	var data = make([]dbpak.KVData, 0)
	var err error
	err = variable.CurrentDBClient.Open()
	if err != nil {
		return
	}
	defer variable.CurrentDBClient.Close()

	if lastEnd == nil && lastStart == nil {
		count = 0
		Orgdata = Orgdata[:0]
	}
	if lastPage < variable.CurrentPage {
		//next page

		//The reason why "variable.ItemsPerPage" is added by one is that we want to see if the next pages have a value to enable or disable the next or prev key.
		err, data = variable.CurrentDBClient.Read(lastEnd, nil, variable.ItemsPerPage+1)
		if err != nil {
			fmt.Println(err)
		}

		if len(data) == variable.ItemsPerPage+1 {
			data = data[:variable.ItemsPerPage]
			variable.ItemsAdded = true

		} else {
			variable.ItemsAdded = false

		}
		if len(data) == 0 {
			return
		}
		if count > 2 {
			Orgdata = Orgdata[len(data):]
		}

		Orgdata = append(Orgdata, data...)
		count++
	} else {

		//The reason why "variable.ItemsPerPage" is added by one is that we want to see if the next pages have a value to enable or disable the next or prev key.
		err, data = variable.CurrentDBClient.Read(nil, lastStart, variable.ItemsPerPage+1)
		if err != nil {
			fmt.Println(err)
		}

		if len(data) == variable.ItemsPerPage+1 {
			data = data[1:]
			variable.ItemsAdded = true
		}
		if len(data) == 0 {
			return
		}
		Orgdata = Orgdata[:len(Orgdata)-len(data)]
		Orgdata = append(data, Orgdata...)

	}

	lastStart = &Orgdata[0].Key
	lastEnd = &Orgdata[len(Orgdata)-1].Key

	var arrayContainer []fyne.CanvasObject
	for _, item := range data {

		truncatedKey := utils.TruncateString(string(item.Key), 20)
		truncatedValue := utils.TruncateString(string(item.Value), 30)

		typeValue := mimetype.Detect(item.Value)
		if typeValue.Extension() != ".txt" {

			truncatedValue = fmt.Sprintf("* %s . . .", typeValue.Extension())
		}
		valueLabel := BuidLableKeyAndValue("value", item.Key, item.Value, truncatedValue, rightColumnContent)
		keyLabel := BuidLableKeyAndValue("key", item.Key, item.Value, truncatedKey, rightColumnContent)

		buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
		arrayContainer = append(arrayContainer, buttonRow)
	}
	if lastPage > variable.CurrentPage {

		rightColumnContent.Objects = append(arrayContainer, rightColumnContent.Objects...)
	} else {

		rightColumnContent.Objects = append(rightColumnContent.Objects, arrayContainer...)

	}

	data = data[:0]
	rightColumnContent.Refresh()
	lastPage = variable.CurrentPage
}

func ProjectButton(inputText string, lastColumnContent *fyne.Container, path string, rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button, nameDatabace string) *fyne.Container {
	projectButton := widget.NewButton(inputText+" - "+nameDatabace, func() {
		variable.ItemsAdded = true
		utils.Checkdatabace(path, nameDatabace)
		buttonAdd.Enable()
		variable.FolderPath = path
		lastEnd = nil
		variable.CurrentPage = 1
		lastPage = 0
		variable.PreviousOffsetY = 0
		lastStart = nil
		utils.CheckCondition(rightColumnContentORG)
		UpdatePage(rightColumnContentORG)

		nameButtonProject.Text = ""
		nameButtonProject.Text = inputText + " - " + nameDatabace

		nameButtonProject.Refresh()

	})

	buttonContainer := container.NewHBox()

	closeButton := widget.NewButton("✖", func() {

		if nameButtonProject.Text == inputText+" - "+nameDatabace {
			utils.CheckCondition(rightColumnContentORG)

			buttonAdd.Disable()

			nameButtonProject.Text = ""
			nameButtonProject.Refresh()
		}

		err := variable.CurrentJson.Remove(inputText)
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

func BuidLableKeyAndValue(eidtKeyAbdValue string, key []byte, value []byte, nameLable string, rightColumnContent *fyne.Container) *utils.TappableLabel {
	var lableKeyAndValue *utils.TappableLabel
	var contentType *fyne.Container
	var valueEntry *widget.Entry
	var truncatedKey2 string

	lableKeyAndValue = utils.NewTappableLabel(nameLable, func() {
		editWindow := fyne.CurrentApp().NewWindow("Edit" + eidtKeyAbdValue)
		editWindow.Resize(fyne.NewSize(600, 600))
		mainContainer := container.NewVBox()

		typeVlaue := mimetype.Detect([]byte(value))
		if eidtKeyAbdValue == "value" {

			switch {
			case strings.HasPrefix(typeVlaue.String(), "image/"):
				contentType = utils.ImageShow([]byte(key), []byte(value), nameLable, mainContainer, editWindow)

				typeValue := mimetype.Detect([]byte(value))
				truncatedKey2 = fmt.Sprintf("* %s . . .", typeValue.Extension())

			case strings.HasPrefix(typeVlaue.String(), "video/"):
				fmt.Println("video")

			case strings.HasPrefix(typeVlaue.String(), "audio/"):
				fmt.Println("audio")

			case strings.HasPrefix(typeVlaue.String(), "application/"):
				valueEntry2 := widget.NewMultiLineEntry()
				valueEntry2.Resize(fyne.NewSize(500, 500))
				valueEntry2.SetText(string(value))
				mainContainer.Add(valueEntry2)

				contentType = container.NewVBox(widget.NewLabel(""))

				value = []byte(valueEntry2.Text)

			case strings.HasPrefix(typeVlaue.String(), "text/"):

				valueEntry = widget.NewMultiLineEntry()
				valueEntry.Resize(fyne.NewSize(500, 500))
				valueEntry.SetText(string(value))
				scrollableEntry := container.NewScroll(valueEntry)
				mainContainer = container.NewBorder(nil, nil, nil, nil, scrollableEntry)
				scrollableEntry.SetMinSize(fyne.NewSize(600, 500))
				mainContainer.Add(scrollableEntry)

				contentType = container.NewVBox(widget.NewLabel(""))

				value = []byte(valueEntry.Text)
			case strings.HasPrefix(typeVlaue.String(), "font/"):
				fmt.Println("font")
			}

		} else {
			valueEntry = widget.NewMultiLineEntry()
			valueEntry.Resize(fyne.NewSize(500, 500))
			valueEntry.SetText(string(key))
			scrollableEntry := container.NewScroll(valueEntry)
			mainContainer = container.NewBorder(nil, nil, nil, nil, scrollableEntry)
			scrollableEntry.SetMinSize(fyne.NewSize(600, 500))
			mainContainer.Add(scrollableEntry)
			contentType = container.NewVBox(widget.NewLabel(""))

		}

		saveButton := widget.NewButton("Save", func() {

			err := variable.CurrentDBClient.Open()
			if err != nil {
				fmt.Println("error Open")
			}
			defer variable.CurrentDBClient.Close()

			if eidtKeyAbdValue == "value" {

				if strings.HasPrefix(typeVlaue.String(), "text/") {
					value = []byte(valueEntry.Text)
					truncatedKey2 = utils.TruncateString(valueEntry.Text, 30)
				} else if utils.ValueImage != nil {
					value = utils.ValueImage

				}

				err := variable.CurrentDBClient.Add(key, value)
				if err != nil {
					fmt.Println(err)
				}

			} else {

				valueBefor, err := variable.CurrentDBClient.Get(key)
				if err != nil {
					return
				}
				err = variable.CurrentDBClient.Delete(key)
				if err != nil {
					return
				}

				key = []byte(utils.CleanInput(valueEntry.Text))

				err = variable.CurrentDBClient.Add(key, valueBefor)
				if err != nil {
					fmt.Println(err)
				}
				truncatedKey2 = utils.TruncateString(string(key), 20)
			}

			lableKeyAndValue.SetText(truncatedKey2)
			lableKeyAndValue.Refresh()

			editWindow.Close()
			rightColumnContent.Refresh()
		})

		cancelButton := widget.NewButton("Cancel", func() {
			editWindow.Close()
		})

		rowBottom := container.NewVBox(
			contentType,
			container.NewBorder(nil, container.NewGridWithColumns(2, cancelButton, saveButton), nil, nil),
		)
		editContentScr := container.NewScroll(mainContainer)
		coulumnORG := container.NewBorder(
			widget.NewLabel("Edit "+eidtKeyAbdValue+" :"), rowBottom, nil, nil, editContentScr,
		)

		editWindow.SetContent(coulumnORG)
		editWindow.Show()
	})
	return lableKeyAndValue
}

func SearchDatabase(valueEntry *widget.Entry, editWindow fyne.Window, rightColumnContent *fyne.Container) (bool, error) {

	err := variable.CurrentDBClient.Open()
	if err != nil {
		return false, err
	}

	key := utils.CleanInput(valueEntry.Text)
	err, data := variable.CurrentDBClient.Search([]byte(key))
	if err != nil {
		return false, err
	}

	defer variable.CurrentDBClient.Close()

	if len(data) == 0 {
		return false, err
	}
	utils.CheckCondition(rightColumnContent)
	for _, item := range data {

		value, err := variable.CurrentDBClient.Get(item)
		if err != nil {
			return false, err
		}
		truncatedKey := utils.TruncateString(string(item), 20)
		truncatedValue := utils.TruncateString(string(value), 30)

		typeValue := mimetype.Detect([]byte(value))
		if typeValue.Extension() != ".txt" {
			truncatedValue = fmt.Sprintf("* %s . . .", typeValue.Extension())
		}
		valueLabel := BuidLableKeyAndValue("value", item, value, truncatedValue, rightColumnContent)
		keyLabel := BuidLableKeyAndValue("key", item, value, truncatedKey, rightColumnContent)
		buttonRow := container.NewGridWithColumns(2, keyLabel, valueLabel)
		rightColumnContent.Add(buttonRow)
	}

	editWindow.Close()
	return true, nil
}

func DeleteKeyLogic(valueEntry *widget.Entry, editWindow fyne.Window, rightColumnContent *fyne.Container) {
	defer variable.CurrentDBClient.Close()

	key := utils.CleanInput(valueEntry.Text)

	valueSearch, err := QueryKey(valueEntry.Text)
	if valueSearch == "" && err != nil {
		dialog.ShowInformation("Error", "This key does not exist in the database", editWindow)
	} else {
		err = variable.CurrentDBClient.Delete([]byte(key))
		if err != nil {
			log.Fatal("this err for func DeletKeyLogic part else delete || err : ", err)
			return
		}
		dialog.ShowInformation("successful", "The operation was successful", editWindow)
		time.Sleep(2 * time.Second)
		editWindow.Close()
	}
}

func AddKeyLogic(iputKey string, valueFinish []byte, windowAdd fyne.Window) {

	key := utils.CleanInput(iputKey)

	defer variable.CurrentDBClient.Close()

	checkNow, err := QueryKey(iputKey)
	if checkNow != "" || err == nil {
		dialog.ShowInformation("Error", "This key has already been added to your database", windowAdd)

	} else {
		err = variable.CurrentDBClient.Add([]byte(key), valueFinish)
		if err != nil {
			log.Fatal("error : this error in func addkeylogic for add key in database")
		}

		windowAdd.Close()
	}
}

func QueryKey(iputKey string) (string, error) {
	var err error

	key := utils.CleanInput(iputKey)

	err = variable.CurrentDBClient.Open()
	if err != nil {
		return "", err
	}
	checkNow, err := variable.CurrentDBClient.Get([]byte(key))
	if err != nil {
		fmt.Println("error : delete func logic for get key in databace")
	}
	return string(checkNow), err
}
