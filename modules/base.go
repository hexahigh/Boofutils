package modules

import (
	"encoding/base64"
)

func b64_string(input string, eod string) string {
	if eod == "e" || eod == "" {
		return base64.StdEncoding.EncodeToString([]byte(input))
	}
	if eod == "d" {
		decoded, err := base64.StdEncoding.DecodeString(input)
		if err != nil {
			return input
		}
		return string(decoded)
	}
	return "An error occured while encoding/decoding"
}

func b64_data(input []byte, eod string) string {
	if eod == "e" || eod == "" {
		return base64.StdEncoding.EncodeToString(input)
	}
	if eod == "d" {
		decoded, err := base64.StdEncoding.DecodeString(string(input))
		if err != nil {
			return string(input)
		}
		return string(decoded)
	}
	return "An error occured while encoding/decoding"
}
