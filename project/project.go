package project

// TODO Mutexes to prevent data races on a project, especially for writing

import (
	"fmt"
	"github.com/harry1453/audioQ/audio"
	"io"
)

type Cue struct {
	Name  string
	Audio audio.File
}

type CueInfo struct {
	Name string
}

type Settings struct {
	BufferSize uint
}

var (
	nameUpdateListeners     []chan<- string
	settingsUpdateListeners []chan<- Settings
	cuesUpdateListeners     []func()
	name                    string
	settings                Settings
	Cues                    []Cue
	isClosed                bool
	CurrentCue              uint
	nextCuePlayable         *audio.Playable
	cueFinishedChannel      chan bool
)

type ProjectInfo struct {
	Name       string
	Settings   Settings
	Cues       []CueInfo
	CurrentCue uint
}

func Init() error {
	if name == "" {
		name = "Untitled"
	}
	if settings.BufferSize == 0 {
		settings.BufferSize = 100
	}
	CurrentCue = 0
	isClosed = false
	if err := loadNextCue(); err != nil {
		return err
	}
	cueFinishedChannel = make(chan bool)
	go monitorCueFinishedChannel()
	for _, listener := range nameUpdateListeners {
		listener <- name
	}
	for _, listener := range settingsUpdateListeners {
		listener <- settings
	}
	updateCuesListeners()
	return nil
}

func Close() {
	isClosed = true
	audio.StopAll()
	if nextCuePlayable != nil {
		_ = nextCuePlayable.Close()
		nextCuePlayable = nil
	}
	close(cueFinishedChannel)
}

func (cue *Cue) getInfo() CueInfo {
	return CueInfo{
		Name: cue.Name,
	}
}

func GetInfo() ProjectInfo {
	cues := make([]CueInfo, len(Cues))
	for i := 0; i < len(Cues); i++ {
		cues[i] = Cues[i].getInfo()
	}
	return ProjectInfo{
		Name:       name,
		Settings:   settings,
		Cues:       cues,
		CurrentCue: CurrentCue,
	}
}

func AddNameListener(listener chan<- string) {
	nameUpdateListeners = append(nameUpdateListeners, listener)
}

func GetName() string {
	return name
}

func SetName(name string) {
	name = name
	for _, listener := range nameUpdateListeners {
		listener <- name
	}
}

func AddSettingsListener(listener chan<- Settings) {
	settingsUpdateListeners = append(settingsUpdateListeners, listener)
}

func GetSettings() Settings {
	return settings
}

func SetSettings(newSettings Settings) {
	settings = newSettings
	for _, listener := range settingsUpdateListeners {
		listener <- settings
	}
}

func AddCuesUpdateListener(listener func()) {
	cuesUpdateListeners = append(cuesUpdateListeners, listener)
}

func IsCueNumberInRange(cueNumber int) bool {
	return cueNumber >= 0 && cueNumber < len(Cues)
}

func StopPlaying() {
	audio.StopAll()
	cueFinishedChannel <- true
}

func monitorCueFinishedChannel() {
	for !isClosed {
		if nextCuePlayable != nil {
			if err := nextCuePlayable.Initialize(settings.BufferSize); err != nil {
				fmt.Println("Error initializing playable", err)
			}
		}
		<-cueFinishedChannel
	}
}

func updateCuesListeners() {
	for _, listener := range cuesUpdateListeners {
		listener()
	}
}

func AddCue(name string, fileName string, file io.Reader) error {
	cueAudio, err := audio.ParseFile(fileName, file)
	if err != nil {
		return err
	}
	wasAtEnd := isAtEndOfQueue()
	Cues = append(Cues, Cue{name, cueAudio})
	defer updateCuesListeners()
	if wasAtEnd {
		if err := loadNextCue(); err != nil {
			return err
		} else {
			cueFinishedChannel <- true
			return nil
		}
	}
	return nil
}

func RemoveCue(cueNumber int) error {
	if !IsCueNumberInRange(cueNumber) {
		return fmt.Errorf("cue number out of range: %d", cueNumber)
	}
	if isAtEndOfQueue() {
		CurrentCue--
	}
	copy(Cues[cueNumber:], Cues[cueNumber+1:])
	Cues[len(Cues)-1] = Cue{}
	Cues = Cues[:len(Cues)-1]
	defer updateCuesListeners()
	return loadNextCue()
}

func RenameCue(cueNumber int, name string) error {
	if !IsCueNumberInRange(cueNumber) {
		return fmt.Errorf("cue number out of range: %d", cueNumber)
	}
	Cues[cueNumber].Name = name
	defer updateCuesListeners()
	return nil
}

func MoveCue(from, to int) error {
	if !IsCueNumberInRange(from) {
		return fmt.Errorf("cue number out of range: %d", from)
	}
	if !IsCueNumberInRange(to) {
		return fmt.Errorf("cue number out of range: %d", from)
	}
	// Get Cue to move
	cue := Cues[from]

	// Remove cue
	copy(Cues[from:], Cues[from+1:])
	Cues[len(Cues)-1] = Cue{}
	Cues = Cues[:len(Cues)-1]

	// Insert the Cue again
	Cues = append(Cues, Cue{} /* use the zero value of the element type */)
	copy(Cues[to+1:], Cues[to:])
	Cues[to] = cue
	defer updateCuesListeners()
	return nil
}

// Begins playing the next song and then attempts to advance the queue
func PlayNext() error {
	if err := playNext(); err != nil {
		return err
	} else {
		advanceQueue()
		if err = loadNextCue(); err != nil {
			return err
		} else {
			return nil
		}
	}
}

// Begins playing the next song in the queue
func playNext() error {
	if nextCuePlayable != nil {
		if nextCuePlayable.IsPlaying() {
			return fmt.Errorf("next cue already playing")
		} else {
			return nextCuePlayable.Play(cueFinishedChannel, settings.BufferSize)
		}
	} else {
		return fmt.Errorf("no cue loaded")
	}
}

func JumpTo(cueNumber int) error {
	if !IsCueNumberInRange(cueNumber) {
		return fmt.Errorf("cue number outside of range of cues: %d", cueNumber)
	}
	CurrentCue = uint(cueNumber)
	defer updateCuesListeners()
	return loadNextCue()
}

func advanceQueue() {
	CurrentCue++
	if CurrentCue == uint(len(Cues)) {
		CurrentCue = 0
	}
	defer updateCuesListeners()
}

// Returns whether the current cue is the last one or not
// First check is to avoid uint from underflowing
func isAtEndOfQueue() bool {
	if len(Cues) == 0 {
		return true
	}
	length := uint(len(Cues) - 1)
	if CurrentCue == length {
		return true
	}
	return false
}

// Loads the next cue as a playable or sets it to nil
// If the end of the queue has been reached
func loadNextCue() error {
	if len(Cues) == 0 {
		return nil
	}

	n := CurrentCue
	if n >= uint(len(Cues)) {
		n = 0
	}
	if playable, err := Cues[n].Audio.Decode(); err != nil {
		return err
	} else {
		nextCuePlayable = &playable
	}
	return nil
}
