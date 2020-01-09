package console

import (
	"bufio"
	"fmt"
	"github.com/harry1453/audioQ/project"
	"os"
	"strconv"
)

var input = bufio.NewReader(os.Stdin)

func Initialize() {
	for {
		switch getCommand() {
		case "play":
			play()
			break
		case "stop":
			stop()
			break
		case "add":
			addSong()
			break
		case "jump":
			jump()
			break
		case "save":
			save()
			break
		case "load":
			load()
			break
		}
	}
}

func save() {
	fmt.Println("Where to?")
	project.SaveProject(getCommand())
}

func load() {
	fmt.Println("Where from?")
	project.LoadProject(getCommand())
}

func play() {
	fmt.Println("next!")
	err := project.Instance.PlayNext()
	fmt.Println("Error: ", err)
}

func stop() {
	fmt.Println("stop!")
	project.Instance.StopPlaying()
}

func jump() {
	fmt.Println("jump to where?")
	n, err := strconv.Atoi(getCommand())
	if err != nil {
		fmt.Println("could not parse")
		return
	}
	project.Instance.JumpTo(n)
}

func addSong() {
	fmt.Println("OK! File name?")
	fileName := getCommand()
	//fmt.Println("OK! Cue name?")
	//cueName := getCommand()
	cueName := "cue"
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Error: ", project.Instance.AddCue(cueName, fileName, file))
}

func getCommand() string {
	s, err := input.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return s[:len(s)-1]
}
