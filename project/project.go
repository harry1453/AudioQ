package project

// TODO Mutexes to prevent data races on a project, especially for writing

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
	name                    string
	nameUpdateListeners     []chan<- string
	settings                Settings
	settingsUpdateListeners []chan<- Settings
	Cues                    []Cue
	isClosed                bool
	currentCue              uint
	nextCuePlayable         *audio.Playable
	cueFinishedChannel      chan bool
}

type ProjectInfo struct {
	Name       string
	Settings   Settings
	Cues       []CueInfo
	CurrentCue uint
}

func (project *Project) Init() error {
	if project.name == "" {
		project.name = "Untitled"
	}
	if project.settings.BufferSize == 0 {
		project.settings.BufferSize = 100
	}
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
		Name:       project.name,
		Settings:   project.settings,
		Cues:       cues,
		CurrentCue: project.currentCue,
	}
}

func (project *Project) AddNameListener(listener chan<- string) {
	project.nameUpdateListeners = append(project.nameUpdateListeners, listener)
}

func (project *Project) GetName() string {
	return project.name
}

func (project *Project) SetName(name string) {
	project.name = name
	for _, listener := range project.nameUpdateListeners {
		listener <- name
	}
}

func (project *Project) AddSettingsListener(listener chan<- Settings) {
	project.settingsUpdateListeners = append(project.settingsUpdateListeners, listener)
}

func (project *Project) GetSettings() Settings {
	return project.settings
}

func (project *Project) SetSettings(settings Settings) {
	project.settings = settings
	for _, listener := range project.settingsUpdateListeners {
		listener <- settings
	}
}

func (project *Project) IsCueNumberInRange(cueNumber int) bool {
	return cueNumber >= 0 && cueNumber < len(project.Cues)
}

func (project *Project) StopPlaying() {
	audio.StopAll()
	project.cueFinishedChannel <- true
}

func (project *Project) monitorCueFinishedChannel() {
	for !project.isClosed {
		if project.nextCuePlayable != nil {
			project.nextCuePlayable.Initialize(project.settings.BufferSize)
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
	if !project.IsCueNumberInRange(cueNumber) {
		return fmt.Errorf("cue number out of range: %d", cueNumber)
	}
	if project.isAtEndOfQueue() {
		project.currentCue--
	}
	copy(project.Cues[cueNumber:], project.Cues[cueNumber+1:])
	project.Cues[len(project.Cues)-1] = Cue{}
	project.Cues = project.Cues[:len(project.Cues)-1]
	return project.loadNextCue()
}

func (project *Project) RenameCue(cueNumber int, name string) error {
	if !project.IsCueNumberInRange(cueNumber) {
		return fmt.Errorf("cue number out of range: %d", cueNumber)
	}
	project.Cues[cueNumber].Name = name
	return nil
}

func (project *Project) MoveCue(from, to int) error {
	if !project.IsCueNumberInRange(from) {
		return fmt.Errorf("cue number out of range: %d", from)
	}
	if !project.IsCueNumberInRange(to) {
		return fmt.Errorf("cue number out of range: %d", from)
	}
	// Get Cue to move
	cue := project.Cues[from]

	// Remove cue
	copy(project.Cues[from:], project.Cues[from+1:])
	project.Cues[len(project.Cues)-1] = Cue{}
	project.Cues = project.Cues[:len(project.Cues)-1]

	// Insert the Cue again
	project.Cues = append(project.Cues, Cue{} /* use the zero value of the element type */)
	copy(project.Cues[to+1:], project.Cues[to:])
	project.Cues[to] = cue
	return nil
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
			return project.nextCuePlayable.Play(project.cueFinishedChannel, project.settings.BufferSize)
		}
	} else {
		return fmt.Errorf("no cue loaded")
	}
}

func (project *Project) JumpTo(cueNumber int) error {
	if !project.IsCueNumberInRange(cueNumber) {
		return fmt.Errorf("cue number outside of range of cues: %d", cueNumber)
	}
	project.currentCue = uint(cueNumber)
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
	if n >= uint(len(project.Cues)) {
		n = 0
	}
	if playable, err := project.Cues[n].Audio.Decode(); err != nil {
		return err
	} else {
		project.nextCuePlayable = &playable
	}
	return nil
}
