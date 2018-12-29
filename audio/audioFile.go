package audio

import (
	"bytes"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/wav"
	"io"
	"io/ioutil"
)

type AudioEncoding uint

const (
	WAV = iota
	MP3
	FLAC
	VORBIS
)

type AudioFile struct {
	Encoding AudioEncoding
	Data string
}

func (file *AudioFile) Decode() (Playable, error) {
	data, err := fromBase64(file.Data)
	if err != nil {
		return Playable{}, err
	}

	var decode func(io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error)
	switch file.Encoding {
	case WAV:
		decode = wav.Decode
		break
	case MP3:
		decode = mp3.Decode
		break
	case FLAC:
		decode = flac.Decode
		break
	case VORBIS:
		decode = flac.Decode
		break
	default:
		return Playable{}, fmt.Errorf("invalid encoding: %d", file.Encoding)
	}

	stream, format, err := decode(ioutil.NopCloser(bytes.NewReader(data)))
	if err != nil {
		return Playable{}, err
	}
	return Playable{stream, format, false, false}, err
}
