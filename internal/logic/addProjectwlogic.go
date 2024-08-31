package logic

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	// "testgui/internal/logic/mainwindowlagic"

	dbpak "testgui/pkg/db"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/syndtr/goleveldb/leveldb"
)

type TappableLabel struct {
	widget.Label
	db       dbpak.DBClient
	onTapped func()
}

func (t *TappableLabel) HasManifestFile(folderPath string) bool {
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

func (t *TappableLabel) HandleButtonClick(test string) error {
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
