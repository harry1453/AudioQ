package audio

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"time"
)

type Playable struct {
	stream beep.StreamSeekCloser
	format beep.Format
	isPlaying bool
	isClosed bool
}

func (playable *Playable) Play() (chan struct{}, error) {
	if playable.isClosed && !playable.isPlaying {
		return nil, fmt.Errorf("playable has already been played")
	}

	playing := make(chan struct{})

	fmt.Println("init")
	speaker.Init(playable.format.SampleRate, playable.format.SampleRate.N(time.Second/10))

	fmt.Println("play")
	speaker.Play(beep.Seq(playable.stream, beep.Callback(func() {
		playable.Close()
		close(playing)
	})))

	playable.isPlaying = true

	return playing, nil
}

func (playable *Playable) IsPlaying() bool {
	return playable.isPlaying
}

func (playable *Playable) Close() error {
	playable.isPlaying = false
	playable.isClosed = true
	return playable.stream.Close()
}
