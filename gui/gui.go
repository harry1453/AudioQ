package gui

import (
	"github.com/harry1453/audioQ/constants"
	"github.com/harry1453/audioQ/project"
	"github.com/harry1453/go-common-file-dialog/cfd"
	"github.com/harry1453/go-common-file-dialog/cfdutil"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"strconv"
)

func Initialize() {
	nameUpdateChannel := make(chan string)
	project.Instance.AddNameListener(nameUpdateChannel)
	settingsUpdateChannel := make(chan project.Settings)
	project.Instance.AddSettingsListener(settingsUpdateChannel)
	settingsStringUpdateChannel := make(chan string)
	go func() {
		for {
			settingsStringUpdateChannel <- strconv.Itoa(int((<-settingsUpdateChannel).BufferSize))
		}
	}()
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
							Composite{
								Layout: HBox{Spacing: 1},
								Children: []Widget{
									PushButton{
										Text: "Open",
										OnClicked: func() {
											fileName, err := cfdutil.ShowOpenFileDialog(cfd.DialogConfig{
												Title: "Open AudioQ File",
												FileFilters: []cfd.FileFilter{
													{
														DisplayName: "AudioQ File (*.audioq)",
														Pattern:     "*.audioq",
													},
												},
											})
											if err != nil {
												log.Println("Error showing open file dialog:", err)
												return
											}
											project.LoadProject(fileName)
										},
									},
									PushButton{
										Text: "Save",
										OnClicked: func() {
											fileName, err := cfdutil.ShowSaveFileDialog(cfd.DialogConfig{
												Title: "Save AudioQ File",
												FileFilters: []cfd.FileFilter{
													{
														DisplayName: "AudioQ File (*.audioq)",
														Pattern:     "*.audioq",
													},
												},
											})
											if err != nil {
												log.Println("Error showing open file dialog:", err)
												return
											}
											project.SaveProject(fileName)
										},
									},
								},
							},
							Setting("Project name", func(newName string) {
								project.Instance.SetName(newName)
							}, nameUpdateChannel),
							Setting("Buffer Size", func(newBufferSize string) {
								n, err := strconv.Atoi(newBufferSize)
								if err != nil {
									log.Println("Failed to parse buffer size:", newBufferSize, err)
									return
								}
								project.Instance.SetSettings(project.Settings{BufferSize: uint(n)})
							}, settingsStringUpdateChannel),
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
