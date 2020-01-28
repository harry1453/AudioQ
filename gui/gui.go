package gui

import (
	"fmt"
	"github.com/harry1453/audioQ/constants"
	"github.com/harry1453/audioQ/project"
	"github.com/harry1453/go-common-file-dialog/cfd"
	"github.com/harry1453/go-common-file-dialog/cfdutil"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func Initialize() {
	nameUpdateChannel := make(chan string)
	project.AddNameListener(nameUpdateChannel)
	settingsUpdateChannel := make(chan project.Settings)
	project.AddSettingsListener(settingsUpdateChannel)
	settingsStringUpdateChannel := make(chan string)
	go func() {
		for {
			settingsStringUpdateChannel <- strconv.Itoa(int((<-settingsUpdateChannel).BufferSize))
		}
	}()

	var cueTable *walk.TableView
	cueTableModel := NewCueModel()

	project.AddCuesUpdateListener(func() {
		cueTableModel.ResetRows()
	})

	var window *walk.MainWindow

	exit, err := MainWindow{
		AssignTo: &window,
		Title:    "AudioQ " + constants.VERSION,
		Name:     "AudioQ " + constants.VERSION,
		Layout:   VBox{Alignment: AlignHCenterVCenter, Spacing: 5},
		Size: Size{
			Width:  1280,
			Height: 720,
		},
		MenuItems: []MenuItem{
			Menu{
				Text: "File",
				Items: []MenuItem{
					Action{
						Text: "Open",
						OnTriggered: func() {
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
								handleError(window, err)
								return
							}
							if err := project.LoadProjectFile(fileName); err != nil {
								handleError(window, err)
							}
						},
					},
					Action{
						Text: "Save As",
						OnTriggered: func() {
							fileName, err := cfdutil.ShowSaveFileDialog(cfd.DialogConfig{
								Title: "Save AudioQ File",
								FileFilters: []cfd.FileFilter{
									{
										DisplayName: "AudioQ File (*.audioq)",
										Pattern:     "*.audioq",
									},
								},
							})
							if !strings.HasSuffix(fileName, ".audioq") {
								fileName = fileName + ".audioq"
							}
							if err != nil {
								handleError(window, err)
								return
							}
							if err := project.SaveProjectFile(fileName); err != nil {
								handleError(window, err)
							}
						},
					},
				},
			},
		},
		Children: []Widget{
			GroupBox{
				Title:  "Project Settings",
				Layout: VBox{Alignment: AlignHCenterVCenter, Spacing: 5},
				Children: []Widget{
					setting("Project name", func(newName string) {
						project.SetName(newName)
					}, nameUpdateChannel),
					setting("Buffer Size", func(newBufferSize string) {
						n, err := strconv.Atoi(newBufferSize)
						if err != nil {
							handleError(window, fmt.Errorf("failed to parse buffer size: %s, %s", newBufferSize, err))
							return
						}
						project.SetSettings(project.Settings{BufferSize: uint(n)})
					}, settingsStringUpdateChannel),
				},
			},
			GroupBox{
				Title:  "Playback Control",
				Layout: VBox{},
				Children: []Widget{
					TableView{
						AssignTo:         &cueTable,
						AlternatingRowBG: true,
						CheckBoxes:       true,
						ColumnsOrderable: true,
						Columns: []TableViewColumn{
							{Title: "#"},
							{Title: "Sel"},
							{Title: "Name", Width: 400},
						},
						Model: cueTableModel,
					},
					Composite{
						Alignment: AlignHCenterVCenter,
						Layout:    HBox{Alignment: AlignHCenterVCenter, Spacing: 5},
						Children: []Widget{
							PushButton{
								Text: "Jump",
								OnClicked: func() {
									currentIndex, valid := getCurrentCueIndex(window, cueTable)
									if valid {
										if err := project.JumpTo(currentIndex); err != nil {
											handleError(window, err)
										}
									}
								},
							},
							PushButton{
								Text: "Play",
								OnClicked: func() {
									if err := project.PlayNext(); err != nil {
										handleError(window, err)
									}
								},
							},
							PushButton{
								Text: "Stop",
								OnClicked: func() {
									project.StopPlaying()
								},
							},
							PushButton{
								Text: "Keyboard Control",
								OnClicked: func() {
									if err := showKeyboardControlWindow(); err != nil {
										handleError(window, err)
									}
								},
							},
						},
					},
				},
			},
			GroupBox{
				Title:  "Project Manipulation",
				Layout: HBox{Alignment: AlignHCenterVCenter, Spacing: 5},
				Children: []Widget{
					PushButton{
						Text: "Rename",
						OnClicked: func() {
							newName, err := prompt(window, "New name?")
							if err != nil {
								if err != PromptCancelled {
									handleError(window, err)
								}
								return
							}
							currentIndex, valid := getCurrentCueIndex(window, cueTable)
							if valid {
								if err := project.RenameCue(currentIndex, newName); err != nil {
									handleError(window, err)
									return
								}
							}
						},
					},
					PushButton{
						Text: "Move",
						OnClicked: func() {
							from, valid := getCurrentCueIndex(window, cueTable)
							if !valid {
								return
							}
							toString, err := prompt(window, "Index To?")
							if err != nil {
								if err != PromptCancelled {
									handleError(window, err)
								}
								return
							}
							to, err := strconv.Atoi(toString)
							if err != nil {
								handleError(window, err)
								return
							}
							if err := project.MoveCue(from, to); err != nil {
								handleError(window, err)
								return
							}
						},
					},
					PushButton{
						Text: "Delete",
						OnClicked: func() {
							currentIndex, valid := getCurrentCueIndex(window, cueTable)
							if valid {
								if err := project.RemoveCue(currentIndex); err != nil {
									handleError(window, err)
									return
								}
							}
						},
					},
					PushButton{
						Text: "Add Cue",
						OnClicked: func() {
							cueName, err := prompt(window, "Cue Name?")
							if err != nil {
								if err != PromptCancelled {
									handleError(window, err)
								}
								return
							}
							fileName, err := cfdutil.ShowOpenFileDialog(cfd.DialogConfig{
								Title: "Open Cue",
								FileFilters: []cfd.FileFilter{
									{
										DisplayName: "Audio Files (*.wav, *.flac, *.mp3, *.ogg",
										Pattern:     "*.wav;*.flac;*.mp3;*.ogg",
									},
								},
							})
							if err != nil {
								handleError(window, err)
								return
							}
							file, err := os.Open(fileName)
							if err != nil {
								handleError(window, err)
								return
							}
							if err := project.AddCue(cueName, fileName, file); err != nil {
								handleError(window, err)
								return
							}
						},
					},
				},
			},
		},
		OnDropFiles: func(fileNames []string) {
			sort.Strings(fileNames)
			for index, fileName := range fileNames {
				indexString := strconv.Itoa(index + 1)
				file, err := os.Open(fileName)
				if err != nil {
					handleError(window, err)
					continue
				}
				if err := project.AddCue("Cue "+indexString, fileName, file); err != nil {
					handleError(window, err)
				}
			}
		},
	}.Run()
	if err != nil {
		log.Println("GUI Error:", exit, err)
	}
	os.Exit(0)
}
