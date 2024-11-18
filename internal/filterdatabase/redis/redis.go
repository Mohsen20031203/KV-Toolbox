package Filterredis

import (
	variable "DatabaseDB"
	"DatabaseDB/internal/filterdatabase"
	"DatabaseDB/internal/logic"
	"DatabaseDB/internal/utils"
	"fmt"
	"image/color"

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

func (l *NameDatabaseredis) FormCreate(a fyne.App, title string, lastColumnContent *fyne.Container, rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button, columnEditKey *fyne.Container, saveKey *widget.Button, mainWindow fyne.Window) {
	newWindow := a.NewWindow(title)

	createSeparator := func() *canvas.Line {
		line := canvas.NewLine(color.Black)
		line.StrokeWidth = 1
		return line
	}
	line1 := createSeparator()

	lableName := widget.NewLabel("Name :")
	pathEntryName := widget.NewEntry()
	pathEntryName.PlaceHolder = "Name"
	nameContent := container.NewBorder(nil, nil, lableName, nil, pathEntryName)

	lableComment := widget.NewLabel("Comment :")
	pathEntryComment := widget.NewEntry()
	pathEntryComment.PlaceHolder = "Comment"
	commentContent := container.NewBorder(nil, nil, lableComment, nil, pathEntryComment)

	lableHost := widget.NewLabel("Addres :")
	pathEntry := widget.NewEntry()
	pathEntry.PlaceHolder = "Addres"
	HostContent := container.NewBorder(nil, nil, lableHost, nil, pathEntry)

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
	testConnectionButton.Importance = widget.HighImportance

	buttonCancel := widget.NewButton("Cancel", func() {
		newWindow.Close()
	})

	buttonOk := widget.NewButton("Add", func() {
		data := map[string]string{
			"Name":     pathEntryName.Text,
			"Comment":  pathEntryComment.Text,
			"Addres":   pathEntry.Text,
			"Database": title,
			"Username": pathEntryUsername.Text,
			"Password": pathEntryPassword.Text,
		}
		datajson, err := variable.CurrentJson.Load()
		if err != nil {
			fmt.Println("error ", err)
		}
		for _, m := range datajson.RecentProjects {
			if pathEntryName.Text == m.Name {
				dialog.ShowInformation("Error ", "Your database name is duplicate", newWindow)
				return
			}
		}

		path := fmt.Sprintf("%s|-|%s|-|%s", data["Addres"], data["Username"], data["Password"])

		var addButton bool
		err = logic.HandleButtonClick(path, title)
		if err == nil {

			err, addButton = variable.CurrentJson.Add(data, newWindow, title)
		}

		if err != nil {
			dialog.ShowInformation("Error ", string(err.Error()), newWindow)
		} else {
			if !addButton {

				utils.CheckCondition(rightColumnContentORG)

				//buttonContainer := logic.ProjectButton(pathEntryName.Text, lastColumnContent, path, rightColumnContentORG, nameButtonProject, buttonAdd, title)
				//lastColumnContent.Add(buttonContainer)
				lastColumnContent.Refresh()

				variable.CreatDatabase = false
				newWindow.Close()
			}
		}
	})
	buttonOk.Importance = widget.HighImportance

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
		HostContent,
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
