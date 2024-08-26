package main

import (
	"encoding/json"
	"fmt"
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

// FileHandler تعریف می‌کند که باید چه عملیات‌هایی را روی فایل انجام دهد
type FileHandler interface {
	Open(name string) (*os.File, error)
	Close(file *os.File) error
	Stat(file *os.File) (os.FileInfo, error)
}

// JsonHandler تعریف می‌کند که باید چه عملیات‌هایی را روی JSON انجام دهد
type JsonHandler interface {
	Encode(file *os.File, state interface{}) error
	Decode(file *os.File, state interface{}) error
}

// DefaultFileHandler پیاده‌سازی پیش‌فرض برای عملیات روی فایل‌ها
type DefaultFileHandler struct{}

func (d *DefaultFileHandler) Open(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
}

func (d *DefaultFileHandler) Close(file *os.File) error {
	return file.Close()
}

func (d *DefaultFileHandler) Stat(file *os.File) (os.FileInfo, error) {
	return file.Stat()
}

// DefaultJsonHandler پیاده‌سازی پیش‌فرض برای عملیات روی JSON
type DefaultJsonHandler struct{}

func (d *DefaultJsonHandler) Encode(file *os.File, state interface{}) error {
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	return encoder.Encode(&state)
}

func (d *DefaultJsonHandler) Decode(file *os.File, state interface{}) error {
	decoder := json.NewDecoder(file)
	return decoder.Decode(&state)
}

func writeJsonFile(fileHandler FileHandler, jsonHandler JsonHandler, state interface{}) error {
	f, err := fileHandler.Open("data.json")
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer fileHandler.Close(f)

	return jsonHandler.Encode(f, state)
}

func readJsonFile(fileHandler FileHandler, jsonHandler JsonHandler, state interface{}) error {
	f, err := fileHandler.Open("data.json")
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer fileHandler.Close(f)

	return jsonHandler.Decode(f, state)
}

func addProjectToJsonFile(projectPath *widget.Entry, name *widget.Entry, comment *widget.Entry, Window fyne.Window, fileHandler FileHandler, jsonHandler JsonHandler) (error, bool) {
	f, err := fileHandler.Open("data.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer fileHandler.Close(f)

	err = handleButtonClick(projectPath.Text)
	if err != nil {
		return err, false
	}

	var state JsonInformation

	fileInfo, err := fileHandler.Stat(f)
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err), false
	}

	if fileInfo.Size() == 0 {
		state = JsonInformation{
			RecentProjects: []Project{},
		}
	} else {
		err := readJsonFile(fileHandler, jsonHandler, &state)
		if err != nil {
			return err, false
		}
	}

	for _, addres := range state.RecentProjects {
		if projectPath.Text == addres.FileAddress {
			m := fmt.Sprintf("This database has already been added to your projects under the name '%s'", addres.Name)
			dialog.ShowInformation("error", m, Window)
			return nil, true
		}
	}
	newActivity := Project{
		Name:        name.Text,
		Comment:     comment.Text,
		FileAddress: projectPath.Text,
	}

	state.RecentProjects = append(state.RecentProjects, newActivity)

	err = writeJsonFile(fileHandler, jsonHandler, state)
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %v", err), false
	}
	return nil, false
}

func removeProjectFromJsonFile(projectName string, fileHandler FileHandler, jsonHandler JsonHandler) error {
	f, err := fileHandler.Open("data.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer fileHandler.Close(f)

	var state JsonInformation

	err = readJsonFile(fileHandler, jsonHandler, &state)
	if err != nil {
		return err
	}

	for i, project := range state.RecentProjects {
		if project.Name == projectName {
			state.RecentProjects = append(state.RecentProjects[:i], state.RecentProjects[i+1:]...)
			break
		}
	}

	err = writeJsonFile(fileHandler, jsonHandler, state)
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %v", err)
	}

	return nil
}

func loadJsonData(fileName string, fileHandler FileHandler, jsonHandler JsonHandler) (JsonInformation, error) {
	var jsonData JsonInformation

	// باز کردن فایل
	file, err := fileHandler.Open(fileName)
	if err != nil {
		return jsonData, fmt.Errorf("failed to open file: %v", err)
	}
	defer fileHandler.Close(file)

	// بارگذاری داده‌های JSON
	if err := jsonHandler.Decode(file, &jsonData); err != nil {
		return jsonData, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return jsonData, nil
}
