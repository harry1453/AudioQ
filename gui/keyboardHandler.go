package gui

import (
	"github.com/harry1453/audioQ/project"
	"github.com/lxn/walk"
)

func handleKeystroke(key walk.Key) {
	switch key {
	case walk.KeyD:
		project.ForwardsOne()
	case walk.KeyA:
		project.BackwardsOne()
	case walk.KeySpace:
		project.PlayNext()
	case walk.KeyEscape:
		project.StopPlaying()
	}
}
