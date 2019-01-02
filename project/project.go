package project

import (
	"fmt"
	"github.com/harry1453/audioQ/audio"
	"io"
)

type Cue struct {
	Name  string
	Audio audio.AudioFile
}

type CueInfo struct {
	Name string
}

type Settings struct {
	BufferSize uint
}

type Project struct {
	Name               string
	Settings           Settings
	Cues               []Cue
	isClosed           bool
	currentCue         uint
	nextCuePlayable    *audio.Playable
	cueFinishedChannel chan bool
}

type ProjectInfo struct {
	Name     string
	Settings Settings
	Cues     []CueInfo
}

func (project *Project) Init() error {
	project.Name = "Untitled"
	project.currentCue = 0
	project.isClosed = false
	if err := project.loadNextCue(); err != nil {
		return err
	}
	project.cueFinishedChannel = make(chan bool)
	go project.monitorCueFinishedChannel()
	return nil
}

func (project *Project) Close() {
	project.isClosed = true
	audio.StopAll()
	if project.nextCuePlayable != nil {
		project.nextCuePlayable.Close()
		project.nextCuePlayable = nil
	}
	close(project.cueFinishedChannel)
}

func (cue *Cue) getInfo() CueInfo {
	return CueInfo{
		Name: cue.Name,
	}
}

func (project *Project) GetInfo() ProjectInfo {
	cues := make([]CueInfo, len(project.Cues))
	for i := 0; i < len(project.Cues); i++ {
		cues[i] = project.Cues[i].getInfo()
	}
	return ProjectInfo{
		Name:     project.Name,
		Settings: project.Settings,
		Cues:     cues,
	}
}

func (project *Project) StopPlaying() {
	audio.StopAll()
	project.cueFinishedChannel <- true
}

func (project *Project) monitorCueFinishedChannel() {
	for !project.isClosed {
		if project.nextCuePlayable != nil {
			project.nextCuePlayable.Initialize(project.Settings.BufferSize)
		}
		<-project.cueFinishedChannel
	}
}

func (project *Project) AddCue(name string, fileName string, file io.Reader) error {
	cueAudio, err := audio.ParseFile(fileName, file)
	if err != nil {
		return err
	}
	wasAtEnd := project.isAtEndOfQueue()
	project.Cues = append(project.Cues, Cue{name, cueAudio})
	if wasAtEnd {
		if err := project.loadNextCue(); err != nil {
			return err
		} else {
			project.cueFinishedChannel <- true
			return nil
		}
	}
	return nil
}

func (project *Project) RemoveCue(cueNumber int) error {
	if cueNumber < 0 || cueNumber >= len(project.Cues) {
		return fmt.Errorf("cue number out of range: %d", cueNumber)
	}
	if project.isAtEndOfQueue() {
		project.currentCue--
	}
	copy(project.Cues[cueNumber:], project.Cues[cueNumber+1:])
	project.Cues[len(project.Cues)-1] = Cue{} // or the zero value of T
	project.Cues = project.Cues[:len(project.Cues)-1]
	return project.loadNextCue()
}

// Begins playing the next song and then attempts to advance the queue
func (project *Project) PlayNext() error {
	if err := project.playNext(); err != nil {
		return err
	} else {
		project.advanceQueue()
		if err = project.loadNextCue(); err != nil {
			return err
		} else {
			return nil
		}
	}
}

// Begins playing the next song in the queue
func (project *Project) playNext() error {
	if project.nextCuePlayable != nil {
		if project.nextCuePlayable.IsPlaying() {
			return fmt.Errorf("next cue already playing")
		} else {
			return project.nextCuePlayable.Play(project.cueFinishedChannel, project.Settings.BufferSize)
		}
	} else {
		return fmt.Errorf("no cue loaded")
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
	length := uint(len(project.Cues) - 1)
	if project.currentCue == length {
		return true
	}
	return false
}

// Loads the next cue as a playable or sets it to nil
// If the end of the queue has been reached
func (project *Project) loadNextCue() error {
	if len(project.Cues) == 0 {
		return nil
	}

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
