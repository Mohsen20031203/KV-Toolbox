package jsondata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	variable "testgui"
	jsFile "testgui/internal/json"
	"testgui/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

type ConstantJsonFile struct {
	nameFile string
}

func NewDataBase() jsFile.JsonFile {
	return &ConstantJsonFile{
		nameFile: "data.json",
	}
}

func (j *ConstantJsonFile) Open() (*os.File, error) {
	return os.OpenFile(j.nameFile, os.O_RDWR|os.O_CREATE, 0644)
}

func (j *ConstantJsonFile) Read(state *jsFile.JsonInformation) error {
	file, err := j.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	return json.Unmarshal(byteValue, &state)
}

func (j *ConstantJsonFile) Write(state interface{}) error {
	file, err := j.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	file.Truncate(0)
	file.Seek(0, 0)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	return encoder.Encode(&state)
}

func (j *ConstantJsonFile) Add(path string, nameProject string, commentProject string, window fyne.Window, nameDatabace string) (error, bool) {
	var state jsFile.JsonInformation

	err := HandleButtonClick(path, nameDatabace)
	if err != nil {
		return err, false
	}

	err = j.Read(&state)
	if err != nil && err.Error() != "unexpected end of JSON input" {
		return err, false
	}

	for _, addres := range state.RecentProjects {
		if path == addres.FileAddress {
			m := fmt.Sprintf("This database has already been added to your projects under the name '%s'", addres.Name)
			dialog.ShowInformation("error", m, window)
			return nil, true
		}
	}

	newActivity := jsFile.Project{
		Name:        nameProject,
		Comment:     commentProject,
		FileAddress: path,
		Databace:    nameDatabace,
	}

	state.RecentProjects = append(state.RecentProjects, newActivity)

	return j.Write(state), false
}

func (j *ConstantJsonFile) Remove(projectName string) error {
	var state jsFile.JsonInformation

	err := j.Read(&state)
	if err != nil {
		return err
	}

	for i, project := range state.RecentProjects {
		if project.Name == projectName {
			state.RecentProjects = append(state.RecentProjects[:i], state.RecentProjects[i+1:]...)
			break
		}
	}

	return j.Write(state)
}

func (j *ConstantJsonFile) Load() (jsFile.JsonInformation, error) {
	var jsonData jsFile.JsonInformation

	err := j.Read(&jsonData)
	if err != nil {
		return jsonData, err
	}

	return jsonData, nil
}

func HandleButtonClick(test string, nameDatabace string) error {

	utils.Checkdatabace(test, nameDatabace)
	err := variable.CurrentDBClient.Open()
	if err != nil {
		return fmt.Errorf("no entries found in the database")
	}
	defer variable.CurrentDBClient.Close()
	return nil

}
