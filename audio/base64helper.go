package audio

import (
	"encoding/base64"
)

func toBase64(input []byte) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func fromBase64(input string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(input)
}
