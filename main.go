package main

import (
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/syndtr/goleveldb/leveldb"
)

var listAll []string

func main() {

	db, err := leveldb.OpenFile("mydb", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	myApp := app.New()
	mywindow := myApp.NewWindow("leveldb with fyne")

	keyEntry := widget.NewEntry()
	keyEntry.SetPlaceHolder("key")

	valueEntry := widget.NewEntry()
	valueEntry.SetPlaceHolder("value")

	saveButton := widget.NewButton("save", func() {
		key := keyEntry.Text
		value := valueEntry.Text

		if key != "" || value != "" {
			err = db.Put([]byte(key), []byte(value), nil)
			if err != nil {
				fmt.Println("save information")
			} else {
				keyEntry.SetText("")
				valueEntry.SetText("")
				fmt.Println(key, value)
				listAll = append(listAll, key)
			}
		}

	})

	searchButton := widget.NewButton("search", func() {
		key := keyEntry.Text

		if key != "" {
			data, err := db.Get([]byte(key), nil)
			if err != nil {
				fmt.Println("search information")
			} else {
				valueEntry.SetText(string(data))
			}
		}

	})

	listEntry := widget.NewEntry()
	listEntry.SetPlaceHolder("list")

	list := widget.NewButton("list All ", func() {

		listEntry.SetText(strings.Join(listAll, ","))

	})

	content := container.NewVBox(
		keyEntry,
		valueEntry,
		saveButton,
		searchButton,
		list,
		listEntry,
	)

	contenttt := container.NewHBox(content)

	mywindow.SetContent(contenttt)

	mywindow.Resize(fyne.NewSize(500, 500))
	mywindow.ShowAndRun()
}
