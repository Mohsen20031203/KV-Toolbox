package Filterredis

import (
	"fmt"
	"image/color"
	variable "testgui"
	"testgui/internal/filterdatabase"
	"testgui/internal/logic"
	"testgui/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type NameDatabaseredis struct{}

func NewFileterRedis() filterdatabase.FilterData {
	return &NameDatabaseredis{}
}

func (l *NameDatabaseredis) FilterFile(path string) bool {
	return true
}

func (l *NameDatabaseredis) FilterFormat(folderDialog *dialog.FileDialog) {}

func (l *NameDatabaseredis) FormCreate(a fyne.App, title string, lastColumnContent *fyne.Container, rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button) {
	newWindow := a.NewWindow(title)

	createSeparator := func() *canvas.Line {
		line := canvas.NewLine(color.Black)
		line.StrokeWidth = 1
		return line
	}
	line1 := createSeparator()

	lableHost := widget.NewLabel("Host :")
	pathEntry := widget.NewEntry()
	pathEntry.PlaceHolder = "Host"
	nameContent := container.NewBorder(nil, nil, lableHost, nil, pathEntry)

	lablePort := widget.NewLabel("Port :")
	pathEntryport := widget.NewEntry()
	pathEntryport.PlaceHolder = "Port"
	pornContent := container.NewBorder(nil, nil, lablePort, nil, pathEntryport)

	lableUsername := widget.NewLabel("Username :")
	pathEntryUsername := widget.NewEntry()
	pathEntryUsername.PlaceHolder = "Username"
	usernameContent := container.NewBorder(nil, nil, lableUsername, nil, pathEntryUsername)

	lablePassword := widget.NewLabel("Password :")
	pathEntryPassword := widget.NewEntry()
	pathEntryPassword.PlaceHolder = "Password"
	passwordContent := container.NewBorder(nil, nil, lablePassword, nil, pathEntryPassword)

	testConnectionButton := widget.NewButton("Test Connection", func() {})
	testConnectionButton.Disable()

	buttonCancel := widget.NewButton("Cancel", func() {
		newWindow.Close()
	})

	buttonOk := widget.NewButton("Add", func() {
		datajson, err := variable.CurrentJson.Load()
		if err != nil {
			fmt.Println("error ", err)
		}
		for _, m := range datajson.RecentProjects {
			if pathEntry.Text == m.Name {
				dialog.ShowInformation("Error ", "Your database name is duplicate", newWindow)
				return
			}
		}

		var addButton bool
		path := fmt.Sprintf("%s|-|%s|-|%s|-|%s", pathEntry.Text, pathEntryport.Text, pathEntryUsername.Text, pathEntryPassword.Text)
		err = logic.HandleButtonClick(path, title)
		if err == nil {

			err, addButton = variable.CurrentJson.Add("", pathEntry.Text, pathEntryport.Text, newWindow, title)
		}

		if err != nil {
			dialog.ShowInformation("Error ", string(err.Error()), newWindow)
		} else {
			if !addButton {

				utils.CheckCondition(rightColumnContentORG)

				path := fmt.Sprintf("%s|-|%s|-|%s|-|%s", pathEntry.Text, pathEntryport.Text, pathEntryUsername.Text, pathEntryPassword.Text)
				buttonContainer := logic.ProjectButton(path, lastColumnContent, "", rightColumnContentORG, nameButtonProject, buttonAdd, title)
				lastColumnContent.Add(buttonContainer)
				lastColumnContent.Refresh()

				variable.CreatDatabase = false
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
		pornContent,
		layout.NewSpacer(),
		usernameContent,
		layout.NewSpacer(),
		passwordContent,
		layout.NewSpacer(),
		line1,
		layout.NewSpacer(),
		rowBotton,
	)

	newWindow.Resize(fyne.NewSize(700, 400))
	newWindow.CenterOnScreen()
	newWindow.SetContent(rightColumnContent)
	newWindow.Show()
}
