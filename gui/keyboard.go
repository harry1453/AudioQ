package gui

import (
	"fmt"
	"github.com/harry1453/audioQ/project"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func showKeyboardControlWindow(owner walk.Form) error {
	_, err := Dialog{
		Title:     "Keyboard Control",
		Name:      "Keyboard Control",
		OnKeyDown: handleKeystroke,
		Layout:    Grid{},
		Children: []Widget{
			TextLabel{
				Text:      "Keyboard Control",
				OnKeyDown: handleKeystroke, // To be sure, in case this somehow gains focus
			},
		},
	}.Run(owner)
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
