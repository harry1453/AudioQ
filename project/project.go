package project

import (
	"fmt"
	"github.com/harry1453/audioQ/audio"
)

type Cue struct {
	Name string
	Audio audio.AudioFile
}

type Settings struct {

}

type Project struct {
	Name            string
	Settings        Settings
	Cues            []Cue
	currentCue      uint
	nextCuePlayable *audio.Playable
}

func (project *Project) Init() error {
	project.currentCue = 0
	return project.loadNextCue()
}

func (project *Project) AddCue(name string, fileName string) error {
	cueAudio, err := audio.ParseFile(fileName)
	if err != nil {
		return err
	}
	wasAtEnd := project.isAtEndOfQueue()
	project.Cues = append(project.Cues, Cue{name, cueAudio})
	if wasAtEnd {
		return project.loadNextCue()
	}
	return nil
}

// Begins playing the next song and then attempts to advance the queue
func (project *Project) PlayNext() (chan struct{}, error) {
	if channel, err := project.playNext(); err != nil {
		return nil, err
	} else {
		project.advanceQueue()
		if err = project.loadNextCue(); err != nil {
			return nil, err
		} else {
			return channel, nil
		}
	}
}

// Begins playing the next song in the queue
func (project *Project) playNext() (chan struct{}, error) {
	if project.nextCuePlayable != nil {
		if project.nextCuePlayable.IsPlaying() {
			return nil, fmt.Errorf("next cue already playing")
		} else {
			return project.nextCuePlayable.Play()
		}
	} else {
		return nil, fmt.Errorf("no cue loaded")
	}
}

func (project *Project) JumpTo(cueNumber uint) error {
	if cueNumber >= uint(len(project.Cues)) {
		return fmt.Errorf("cue number outside of range of cues: %d", cueNumber)
	}
	if cueNumber == 0 {
		project.currentCue = 0
	} else {
		project.currentCue = cueNumber - 1
	}
	return project.loadNextCue()
}

func (project *Project) advanceQueue() {
	project.currentCue++
	if project.currentCue == uint(len(project.Cues)) {
		project.currentCue = 0
	}
}

// Returns whether the current cue is the last one or not
// First check is to avoid uint from underflowing
func (project *Project) isAtEndOfQueue() bool {
	if len(project.Cues) == 0 {
		return true
	}
	length := uint(len(project.Cues)-1)
	if project.currentCue == length {
		return true
	}
	return false
}

// Loads the next cue as a playable or sets it to nil
// If the end of the queue has been reached
func (project *Project) loadNextCue() error {
	n := project.currentCue
	if n != 0 {
		n++
	}

	if n == uint(len(project.Cues)) {
		n -= 1
	}

	if playable, err := project.Cues[n].Audio.Decode(); err != nil {
		return err
	} else {
		project.nextCuePlayable = &playable
	}
	return nil
}
