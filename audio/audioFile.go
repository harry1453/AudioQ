package audio

import (
	"bytes"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
	"io"
	"io/ioutil"
)

type Encoding uint

const (
	WAV Encoding = iota
	MP3
	FLAC
	VORBIS
)

type File struct {
	Encoding Encoding
	Data     string
}

func (file *File) Decode() (Playable, error) {
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
		decode = vorbis.Decode
		break
	default:
		return Playable{}, fmt.Errorf("invalid encoding: %d", file.Encoding)
	}

	stream, format, err := decode(ioutil.NopCloser(bytes.NewReader(data)))
	if err != nil {
		return Playable{}, err
	}
	return Playable{stream, format, false, false, false}, err
}
