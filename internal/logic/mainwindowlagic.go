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

func SetupLastColumn(rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button, columnEditKey *fyne.Container, saveKey *widget.Button, mainWindow fyne.Window) *fyne.Container {
	lastColumnContent := container.NewVBox()

	jsonDataa, err := variable.CurrentJson.Load()
	if err != nil {
		println("Error loading JSON data:", err)
	} else {
		for _, project := range jsonDataa.RecentProjects {

			buttonContainer := ProjectButton(project.Name, lastColumnContent, project.FileAddress, rightColumnContentORG, nameButtonProject, buttonAdd, project.Databace, columnEditKey, saveKey, mainWindow)
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
	radioLast *widget.RadioGroup
)

func UpdatePage(rightColumnContent *fyne.Container, columnEditKey *fyne.Container, saveKey *widget.Button, mainWindow fyne.Window) {

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

	var truncatedValue string
	var arrayContainer []fyne.CanvasObject
	for _, item := range data {

		truncatedKey := utils.TruncateString(string(item.Key), 20)

		typeValue := mimetype.Detect(item.Value)
		if typeValue.Extension() != ".txt" {

			truncatedValue = fmt.Sprintf("* %s . . .", typeValue.Extension())
		} else {
			truncatedValue = utils.TruncateString(string(item.Value), 30)

		}
		radio := widget.NewRadioGroup([]string{" "}, nil)
		radio.Disable()

		valueLabel := BuidLableKeyAndValue("value", item.Key, item.Value, truncatedValue, rightColumnContent, columnEditKey, saveKey, mainWindow, radio)
		keyLabel := BuidLableKeyAndValue("key", item.Key, item.Value, truncatedKey, rightColumnContent, columnEditKey, saveKey, mainWindow, radio)
		m := container.NewHBox(
			radio,
			keyLabel,
		)

		buttonRow := container.NewGridWithColumns(2, m, valueLabel)
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

func ProjectButton(inputText string, lastColumnContent *fyne.Container, path string, rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button, nameDatabace string, columnEditKey *fyne.Container, saveKey *widget.Button, mainWindow fyne.Window) *fyne.Container {
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
		utils.CheckCondition(columnEditKey)
		UpdatePage(rightColumnContentORG, columnEditKey, saveKey, mainWindow)

		nameButtonProject.Text = ""
		nameButtonProject.Text = inputText + " - " + nameDatabace

		nameButtonProject.Refresh()

	})

	buttonContainer := container.NewHBox()

	closeButton := widget.NewButton("✖", func() {

		if nameButtonProject.Text == inputText+" - "+nameDatabace {
			utils.CheckCondition(rightColumnContentORG)
			utils.CheckCondition(columnEditKey)

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

func BuidLableKeyAndValue(eidtKeyAbdValue string, key []byte, value []byte, nameLable string, rightColumn *fyne.Container, columnEditKey *fyne.Container, saveKey *widget.Button, mainWindow fyne.Window, radio *widget.RadioGroup) *utils.TappableLabel {
	var lableKeyAndValue *utils.TappableLabel
	var valueEntry *widget.Entry
	var truncatedKey2 string

	lableKeyAndValue = utils.NewTappableLabel(nameLable, func() {
		if radioLast != nil {

			radioLast.SetSelected("")
		}
		radio.SetSelected(" ")
		radioLast = radio
		radioLast.Refresh()
		utils.CheckCondition(columnEditKey)

		typeValue := mimetype.Detect([]byte(value))
		columnEditKey.Add(widget.NewLabel(fmt.Sprintf("Edit %s - %s", eidtKeyAbdValue, nameLable)))

		/*typeKey := widget.NewSelect([]string{"Byte", "other"}, func(s string) {
			if s == "Byte" {
				columnEditKey.Objects = columnEditKey.Objects[:2]

				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))
				columnEditKey.Add(widget.NewLabel("20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 20 49 0a 10 f0 s7 "))

			} else {
				valueEntry := configureEntry(columnEditKey, string(value))
				value = []byte(valueEntry.Text)
			}
		})
		columnEditKey.Add(typeKey)*/

		if eidtKeyAbdValue == "value" {

			switch {
			case strings.HasPrefix(typeValue.String(), "image/"):
				utils.ImageShow([]byte(key), []byte(value), columnEditKey, mainWindow)
				truncatedKey2 = fmt.Sprintf("* %s . . .", typeValue.Extension())

			case strings.HasPrefix(typeValue.String(), "text/") || strings.HasPrefix(typeValue.String(), "application/"):
				valueEntry = configureEntry(columnEditKey, string(value))
				value = []byte(valueEntry.Text)

			}

		} else {

			valueEntry = configureEntry(columnEditKey, string(key))
		}

		saveKey.OnTapped = func() {
			err := variable.CurrentDBClient.Open()
			if err != nil {
				fmt.Println("error Open")
				return
			}
			defer variable.CurrentDBClient.Close()

			saveValue := func() {
				if strings.HasPrefix(typeValue.String(), "text/") {
					value = []byte(valueEntry.Text)
					truncatedKey2 = utils.TruncateString(valueEntry.Text, 30)
				} else if utils.ValueImage != nil {
					value = utils.ValueImage
					utils.ValueImage = nil
				}
				rightColumn.Refresh()
				if err := variable.CurrentDBClient.Add(key, value); err != nil {
					fmt.Println(err)
				}
			}

			updateKey := func() {
				valueBefor, err := variable.CurrentDBClient.Get(key)
				if err != nil {
					return
				}
				if err := variable.CurrentDBClient.Delete(key); err != nil {
					return
				}

				key = []byte(utils.CleanInput(valueEntry.Text))
				if err := variable.CurrentDBClient.Add(key, valueBefor); err != nil {
					fmt.Println(err)
				}
				truncatedKey2 = utils.TruncateString(string(key), 20)
			}

			if eidtKeyAbdValue == "value" {
				saveValue()
			} else {
				updateKey()
			}

			lableKeyAndValue.SetText(truncatedKey2)
			lableKeyAndValue.Refresh()
		}

		columnEditKey.Refresh()
	})
	return lableKeyAndValue
}

func configureEntry(columnEditKey *fyne.Container, content string) *widget.Entry {
	entry := widget.NewMultiLineEntry()
	entry.Resize(fyne.NewSize(500, 500))
	entry.SetText(content)
	scrollableEntry := container.NewScroll(entry)
	scrollableEntry.SetMinSize(fyne.NewSize(200, 300))
	columnEditKey.Add(scrollableEntry)
	return entry
}
