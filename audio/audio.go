package audio

import (
	"fmt"
	"github.com/faiface/beep/speaker"
	"io"
	"io/ioutil"
	"strings"
)

func StopAll() {
	speaker.Clear()
}

func ParseFile(fileName string, file io.Reader) (File, error) {
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return File{}, err
	}
	var encoding Encoding
	if strings.HasSuffix(fileName, ".mp3") {
		encoding = MP3
	} else if strings.HasSuffix(fileName, ".ogg") {
		encoding = VORBIS
	} else if strings.HasSuffix(fileName, ".wav") {
		encoding = WAV
	} else if strings.HasSuffix(fileName, ".flac") {
		encoding = FLAC
	} else {
		return File{}, fmt.Errorf("invalid file type: %s", fileName)
	}

	return File{
		Encoding: encoding,
		Data:     toBase64(fileContents),
	}, nil
}
