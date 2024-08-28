package jsondata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testgui/internal/utils"
	jsFile "testgui/pkg/json"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

type ConstantJsonFile struct {
	nameFile string
}

func NewDataBase(name string) jsFile.JsonFile {
	return &ConstantJsonFile{
		nameFile: name,
	}
}

func (j *ConstantJsonFile) Open() (*os.File, error) {
	file, err := os.OpenFile(j.nameFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return file, err
	}
	return file, nil
}

func (j *ConstantJsonFile) Read(state jsFile.JsonInformation) error {
	file, err := os.Open(j.nameFile)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	if err := json.Unmarshal(byteValue, state); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return nil
}

func (j *ConstantJsonFile) Add(path string, nameProject string, commentProject string, window fyne.Window) (error, bool) {
	var state jsFile.JsonInformation

	fileInfo, err := os.Stat("data.json")
	if err != nil {
		if os.IsNotExist(err) {
			state = jsFile.JsonInformation{
				RecentProjects: []jsFile.Project{},
			}
		} else {
			return fmt.Errorf("failed to get file info: %v", err), false
		}
	} else if fileInfo.Size() != 0 {
		err = j.Read(state)
		if err != nil {
			return err, false
		}
	}

	for _, address := range state.RecentProjects {
		if path == address.FileAddress {
			m := fmt.Sprintf("This database has already been added to your projects under the name '%s'", address.Name)
			dialog.ShowInformation("error", m, window)
			return nil, true
		}
	}

	newProject := jsFile.Project{
		Name:        nameProject,
		Comment:     commentProject,
		FileAddress: path,
	}

	state.RecentProjects = append(state.RecentProjects, newProject)

	file, err := os.OpenFile("data.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err), false
	}
	defer file.Close()

	err = utils.WriteJsonFile(file, &state)
	if err != nil {
		return fmt.Errorf("failed to write JSON: %v", err), false
	}

	return nil, false
}

func (j *ConstantJsonFile) Remove(projectName string) error {
	var state jsFile.JsonInformation

	err := j.Read(state)
	if err != nil {
		return err
	}

	for i, project := range state.RecentProjects {
		if project.Name == projectName {
			state.RecentProjects = append(state.RecentProjects[:i], state.RecentProjects[i+1:]...)
			break
		}
	}

	file, err := os.OpenFile("data.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	err = utils.WriteJsonFile(file, &state)
	if err != nil {
		return fmt.Errorf("failed to write JSON: %v", err)
	}

	return nil
}
