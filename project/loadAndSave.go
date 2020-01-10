package project

import (
	"encoding/json"
	"io"
	"os"
)

type jsonProject struct {
	Name     string
	Settings Settings
	Cues     []Cue
}

func LoadProjectFile(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	return LoadProject(file)
}

func SaveProjectFile(fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	return SaveProject(file)
}

func LoadProject(projectReader io.Reader) error {
	Close()
	var jsonProject jsonProject
	if err := json.NewDecoder(projectReader).Decode(&jsonProject); err != nil {
		return err
	}
	name = jsonProject.Name
	Cues = jsonProject.Cues
	settings = jsonProject.Settings
	if err := Init(); err != nil {
		return err
	}
	return nil
}

func SaveProject(projectWriter io.Writer) error {
	jsonProject := jsonProject{
		Name:     name,
		Settings: settings,
		Cues:     Cues,
	}
	if err := json.NewEncoder(projectWriter).Encode(jsonProject); err != nil {
		return err
	}
	return nil
}
