package logic

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	// "testgui/internal/logic/mainwindowlagic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type TappableLabel struct {
	widget.Label
	onTapped func()
}

func HasManifestFile(folderPath string) bool {
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

func NewTappableLabel(text string, tapped func()) *TappableLabel {
	label := &TappableLabel{
		Label: widget.Label{
			Text: text,
		},
		onTapped: tapped,
	}
	label.ExtendBaseWidget(label)
	return label
}

func (t *TappableLabel) Tapped(_ *fyne.PointEvent) {
	t.onTapped()
}
