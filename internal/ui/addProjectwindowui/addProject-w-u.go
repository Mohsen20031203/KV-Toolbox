package addProjectwindowui

import (
	"fmt"
	"image/color"
	"path/filepath"
	variable "testgui"

	//"testgui/internal/logic/logic"

	"testgui/internal/logic"
	"testgui/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func OpenNewWindow(a fyne.App, title string, lastColumnContent *fyne.Container, rightColumnContentORG *fyne.Container, nameButtonProject *widget.Label, buttonAdd *widget.Button) {

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

		err := logic.HandleButtonClick(pathEntry2.Text)
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

			variable.FolderPath = filepath.Dir(filePath)

			if logic.HasManifestFile(variable.FolderPath) {
				pathEntry2.SetText(variable.FolderPath)
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
		err, addButton := variable.CurrentJson.Add(pathEntry2.Text, pathEntry.Text, pathEntryComment.Text, newWindow)
		if err != nil {
			dialog.ShowInformation("Error ", "There is something wrong with your file and I can't connect to it", newWindow)
		} else {
			if !addButton {

				if !utils.CheckCondition(rightColumnContentORG) {
					newObjects := []fyne.CanvasObject{}

					rightColumnContentORG.Objects = newObjects

					rightColumnContentORG.Refresh()
				}

				buttonContainer := logic.ProjectButton(pathEntry.Text, lastColumnContent, pathEntry2.Text, rightColumnContentORG, nameButtonProject, buttonAdd)
				lastColumnContent.Add(buttonContainer)
				lastColumnContent.Refresh()

				/*handleProjectSelection(pathEntry2.Text, rightColumnContentORG, buttonAdd)
				rightColumnContentORG.Refresh()*/

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
