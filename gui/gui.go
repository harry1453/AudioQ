package gui

import (
	"github.com/harry1453/audioQ/constants"
	"github.com/harry1453/audioQ/project"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
)

func Initialize() {
	nameUpdateChannel := make(chan string)
	project.Instance.AddNameListener(nameUpdateChannel)
	settingsUpdateChannel := make(chan project.Settings)
	project.Instance.AddSettingsListener(settingsUpdateChannel)
	mainWindow := MainWindow{
		Title:  "AudioQ " + constants.VERSION,
		Name:   "AudioQ " + constants.VERSION,
		Layout: HBox{},
		Size: Size{
			Width:  100,
			Height: 100,
		},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					Composite{
						Name:     "Control View",
						Layout:   VBox{},
						Children: []Widget{},
					},
					Composite{
						Name:   "Project View",
						Layout: VBox{},
						Children: []Widget{
							Setting("Project name", func(newName string) {
								project.Instance.SetName(newName)
							}, nameUpdateChannel),
						},
					},
				},
			},
		},
	}
	exit, err := mainWindow.Run()
	if err != nil {
		log.Println("GUI Error:", exit, err)
	}
}

func Setting(name string, onUpdate func(string), newValueChannel <-chan string) Widget {
	var textEdit *walk.TextEdit
	go func() {
		for {
			textEdit.SetText(<-newValueChannel)
		}
	}()
	return Composite{
		Layout: HBox{Spacing: 1},
		Children: []Widget{
			TextLabel{
				Text: name + ": ",
			},
			TextEdit{
				AssignTo:      &textEdit,
				CompactHeight: true,
				OnTextChanged: func() {
					onUpdate(textEdit.Text())
				},
			},
		},
	}
}
