package jsondata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	jsFile "testgui/pkg/json"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/syndtr/goleveldb/leveldb"
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
	file, err := os.OpenFile(j.nameFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return file, err
	}
	return file, nil
}

func (j *ConstantJsonFile) Read(file *os.File, state interface{}) error {
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&state); err != nil {
		return err
	}
	return nil
}

func (j *ConstantJsonFile) Add(path string, nameProject string, commentProject string, window fyne.Window) (error, bool) {
	var state jsFile.JsonInformation

	file, err := j.Open()
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	err = handleButtonClick(path)
	if err != nil {
		return err, false
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err), false
	}

	if fileInfo.Size() == 0 {
		state = jsFile.JsonInformation{
			RecentProjects: []jsFile.Project{},
		}
	} else {
		err := j.Read(file, &state)
		if err != nil {
			return err, false
		}
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
	}

	state.RecentProjects = append(state.RecentProjects, newActivity)

	err = j.Read(file, state)
	if err != nil {
		return fmt.Errorf("failed to decode JSON: %v", err), false
	}
	return nil, false
}

func (j *ConstantJsonFile) Remove(projectName string) error {
	var state jsFile.JsonInformation

	file, err := j.Open()
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	err = j.Read(file, &state)
	if err != nil {
		return err
	}

	for i, project := range state.RecentProjects {
		if project.Name == projectName {
			state.RecentProjects = append(state.RecentProjects[:i], state.RecentProjects[i+1:]...)
			break
		}
	}

	return nil
}

func (j *ConstantJsonFile) Load() (jsFile.JsonInformation, error) {
	var jsonData jsFile.JsonInformation

	file, err := os.Open(j.nameFile)
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

func handleButtonClick(test string) error {
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
