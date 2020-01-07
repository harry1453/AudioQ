package project

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadProject(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}
	defer file.Close()
	Instance.Close()
	if err := json.NewDecoder(file).Decode(Instance); err != nil {
		fmt.Println("ERROR: ", err)
		return
	}
	if err := Instance.Init(); err != nil {
		fmt.Println("ERROR: ", err)
		return
	}
	fmt.Println("OK!")
}

func SaveProject(fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}
	defer file.Close()
	if err := json.NewEncoder(file).Encode(Instance); err != nil {
		fmt.Println("ERROR: ", err)
		return
	}
	fmt.Println("OK!")
}
