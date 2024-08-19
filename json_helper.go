package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type JsonInformation struct {
	RecentProjects []Project `json:"recentProjects"`
}

type Project struct {
	Name        string `json:"name"`
	Comment     string `json:"comment"`
	FileAddress string `json:"fileAddress"`
}

func writeJsonFile(file *os.File, state interface{}) error {
	file.Truncate(0)
	file.Seek(0, 0)
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(&state); err != nil {
		return fmt.Errorf("failed to encode JSON: %v", err)
	}
	return nil
}

func readJsonFile(file *os.File, state interface{}) error {
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&state); err != nil {
		return err
	}
	return nil
}

func openFileJson() (*os.File, error) {
	file, err := os.OpenFile("data.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return file, err
	}
	return file, nil
}

func loadJsonData(fileName string) (JsonInformation, error) {
	var jsonData JsonInformation

	file, err := os.Open(fileName)
	if err != nil {
		return jsonData, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return jsonData, fmt.Errorf("failed to read file: %v", err)
	}

	if err := json.Unmarshal(byteValue, &jsonData); err != nil {
		return jsonData, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return jsonData, nil
}

func addProjectToJsonFile(projectPath *widget.Entry, name *widget.Entry, comment *widget.Entry, Window fyne.Window) (error, bool) {
	file, err := openFileJson()
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	err = handleButtonClick(projectPath.Text)
	if err != nil {
		return err, false
	}

	var state *JsonInformation

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err), false
	}

	if fileInfo.Size() == 0 {
		state = &JsonInformation{
			RecentProjects: []Project{},
		}
	} else {
		err := readJsonFile(file, &state)
		if err != nil {
			return err, false
		}
	}

	for _, addres := range state.RecentProjects {
		if projectPath.Text == addres.FileAddress {
			m := fmt.Sprintf("This database has already been added to your projects under the name '%s'", addres.Name)
			dialog.ShowInformation("error", m, Window)

			err = writeJsonFile(file, state)
			if err != nil {
				return fmt.Errorf("failed to decode JSON: %v", err), false
			}
			return nil, true
		}
	}
	newActivity := Project{
		Name:        name.Text,
		Comment:     comment.Text,
		FileAddress: projectPath.Text,
	}

	state.RecentProjects = append(state.RecentProjects, newActivity)

	err = writeJsonFile(file, state)
	if err != nil {
		return fmt.Errorf("failed to decode JSON: %v", err), false
	}
	return nil, false
}

func removeProjectFromJsonFile(projectName string) error {
	file, err := openFileJson()
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	var state *JsonInformation

	err = readJsonFile(file, &state)
	if err != nil {
		return err
	}

	for i, project := range state.RecentProjects {
		if project.Name == projectName {
			state.RecentProjects = append(state.RecentProjects[:i], state.RecentProjects[i+1:]...)
			break
		}
	}

	err = writeJsonFile(file, &state)
	if err != nil {
		return fmt.Errorf("failed to decode JSON: %v", err)
	}

	return nil
}
