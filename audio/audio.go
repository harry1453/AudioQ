package audio

import (
	"fmt"
	"github.com/faiface/beep/speaker"
	"io/ioutil"
	"strings"
)

func StopAll() {
	speaker.Clear()
}

func ParseFile(fileName string) (AudioFile, error) {
	fileContents, err := ioutil.ReadFile(fileName)
	if err != nil {
		return AudioFile{}, err
	}
	var encoding AudioEncoding
	if strings.HasSuffix(fileName, ".mp3") {
		encoding = MP3
	} else if strings.HasSuffix(fileName, ".ogg") {
		encoding = VORBIS
	} else if strings.HasSuffix(fileName, ".wav") {
		encoding = WAV
	} else if strings.HasSuffix(fileName, ".flac") {
		encoding = FLAC
	} else {
		return AudioFile{}, fmt.Errorf("invalid file type: %s", fileName)
	}

	return AudioFile{
		Encoding: encoding,
		Data: toBase64(fileContents),
	}, nil
}
