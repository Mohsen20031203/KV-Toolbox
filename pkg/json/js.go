package jsFile

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
}

type JsonFile interface {
	Open() (*os.File, error)
	Read(state *JsonInformation) error
	Add(path string, name string, comment string, window fyne.Window) (error, bool)
	Remove(projectName string) error
	Load() (JsonInformation, error)
	Write(state interface{}) error
}
