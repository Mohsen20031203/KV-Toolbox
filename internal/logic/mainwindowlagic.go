package logic

import (
	variable "DatabaseDB"
	"fmt"
	"strings"

	// "DatabaseDB/internal/logic/addProjectwindowlogic"

	dbpak "DatabaseDB/internal/Databaces"
	"DatabaseDB/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gabriel-vasile/mimetype"
)

func SetupLastColumn(rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button, inputEditString, largeEntry *widget.Entry) *fyne.Container {
	lastColumnContent := container.NewVBox()

	jsonDataa, err := variable.CurrentJson.Load()
	if err != nil {
		println("Error loading JSON data:", err)
	} else {
		for _, project := range jsonDataa.RecentProjects {

			buttonContainer := ProjectButton(project.Name, lastColumnContent, project.FileAddress, rightColumnContentORG, nameButtonProject, buttonAdd, project.Databace, inputEditString, largeEntry)
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
	Orgdata   []dbpak.KVData
	lastPage  int
)

func UpdatePage(rightColumnContent *fyne.Container, inputEditString, largeEntry *widget.Entry) {

	var data = make([]dbpak.KVData, 0)
	var err error
	err = variable.CurrentDBClient.Open()
	if err != nil {
		return
	}
	defer variable.CurrentDBClient.Close()

	if lastEnd == nil && lastStart == nil {
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
		if len(rightColumnContent.Objects) >= variable.ItemsPerPage*3 {
			Orgdata = Orgdata[len(data):]
		}

		Orgdata = append(Orgdata, data...)
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
		valueLabel := BuidLableKeyAndValue("value", item.Key, item.Value, truncatedValue, rightColumnContent, inputEditString, largeEntry)
		keyLabel := BuidLableKeyAndValue("key", item.Key, item.Value, truncatedKey, rightColumnContent, inputEditString, largeEntry)

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

func ProjectButton(inputText string, lastColumnContent *fyne.Container, path string, rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button, nameDatabace string, inputEditString, largeEntry *widget.Entry) *fyne.Container {
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
		UpdatePage(rightColumnContentORG, inputEditString, largeEntry)

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

func BuidLableKeyAndValue(eidtKeyAbdValue string, key []byte, value []byte, nameLable string, rightColumnContent *fyne.Container, inputEditString, largeEntry *widget.Entry) *utils.TappableLabel {
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

			case strings.HasPrefix(typeVlaue.String(), "text/") || strings.HasPrefix(typeVlaue.String(), "application/"):

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
