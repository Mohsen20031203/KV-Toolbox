package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func init() {
	app.SetMetadata(fyne.AppMetadata{
		ID: "com.example.LevelManage",
		Name: "LevelManage",
		Version: "0.0.1",
		Build: 1,
		Icon: &fyne.StaticResource{
	StaticName: "icon.png",
	StaticContent: []byte{
		
		Release: false,
		Custom: map[string]string{
			
		},
		
	})
}
