package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func setupLastColumn(app fyne.App) *fyne.Container {
	lastColumnContent := container.NewVBox()

	go func() {
		jsonData, err := loadJsonData("data.json")
		if err != nil {
			println("Error loading JSON data:", err)
		} else {
			for _, project := range jsonData.RecentProjects {
				buttonContainer := projectButton(project.Name, lastColumnContent)
				lastColumnContent.Add(buttonContainer)
			}
		}
	}()

	return lastColumnContent
}

func setupThemeButtons(app fyne.App) *fyne.Container {
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

func leftColumn(lastColumnContent *fyne.Container, lastColumnContentt *fyne.Container, darkLight *fyne.Container) *fyne.Container {
	lastColumnScrollable := container.NewScroll(lastColumnContent)

	mainContent := container.NewBorder(lastColumnContentt, darkLight, nil, nil, lastColumnScrollable)
	return mainContent
}

func rightColumn(rightColumnContent *fyne.Container) fyne.CanvasObject {
	rightColumnScrollable := container.NewScroll(rightColumnContent)
	return rightColumnScrollable
}

func columnContent(rightColumnContent *fyne.Container, lastColumnContent *fyne.Container, lastColumnContentt *fyne.Container, darkLight *fyne.Container) fyne.CanvasObject {
	mainContent := leftColumn(lastColumnContent, lastColumnContentt, darkLight)
	rightColumnScrollable := rightColumn(rightColumnContent)
	columns := container.NewHSplit(mainContent, rightColumnScrollable)
	columns.SetOffset(0.25)

	container.NewScroll(columns)
	return columns
}