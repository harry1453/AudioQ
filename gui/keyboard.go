package gui

import (
	"fmt"
	"github.com/harry1453/audioQ/project"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func showKeyboardControlWindow() error {
	_, err := MainWindow{
		Title:      "Keyboard Control",
		Name:       "Keyboard Control",
		Persistent: true,
		Size: Size{
			Width:  256,
			Height: 144,
		},
		OnKeyDown: handleKeystroke,
		Layout:    VBox{Alignment: AlignHCenterVCenter},
		Children: []Widget{
			TextLabel{
				Alignment: AlignHCenterVCenter,
				Text:      "Keyboard Control",
				OnKeyDown: handleKeystroke, // To be sure, in case this somehow gains focus
			},
		},
	}.Run()
	return err
}

func handleKeystroke(key walk.Key) {
	switch key {
	case walk.KeyD:
		project.ForwardsOne()
	case walk.KeyA:
		project.BackwardsOne()
	case walk.KeySpace:
		project.PlayNext()
	case walk.KeyS:
		project.StopPlaying()
	}
	fmt.Println(key)
}
