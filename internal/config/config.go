package config

import (
	"os"

	"fyne.io/fyne/v2"
)

type JsonInformation struct {
	RecentProjects []Project `json:"recentProjects"`
}

type Project struct {
	Name        string `json:"name"`
	Comment     string `json:"comment"`
	FileAddress string `json:"fileAddress"`
	Databace    string `json:"databace"`
}

type ConfigHandler interface {
	Open() (*os.File, error)
	Read(state *JsonInformation) error
	Add(dataINF map[string]string, window fyne.Window, nameDatabace string) (error, bool)
	Remove(projectName string) error
	Load() (JsonInformation, error)
	Write(state interface{}) error
}
