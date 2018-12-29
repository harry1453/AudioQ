package audio

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"time"
)

type Playable struct {
	stream        beep.StreamSeekCloser
	format        beep.Format
	isPlaying     bool
	isClosed      bool
	isInitialized bool
}

func (playable *Playable) Initialize() error {
	fmt.Println("init")
	if !playable.isInitialized {
		err := speaker.Init(playable.format.SampleRate, playable.format.SampleRate.N(time.Second/10))
		playable.isInitialized = err == nil
		return err
	}
	return nil
}

func (playable *Playable) Play(cueFinishedChannel chan bool) error {
	if playable.isClosed && !playable.isPlaying {
		return fmt.Errorf("playable has already been played")
	}

	if !playable.isInitialized {
		if err := playable.Initialize(); err != nil {
			return err
		}
	}

	fmt.Println("play")
	speaker.Play(beep.Seq(playable.stream, beep.Callback(func() {
		playable.Close()
		cueFinishedChannel <- true
	})))

	playable.isPlaying = true

	return nil
}

func (playable *Playable) IsPlaying() bool {
	return playable.isPlaying
}

func (playable *Playable) Close() error {
	playable.isPlaying = false
	playable.isClosed = true
	return playable.stream.Close()
}
