package modules

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
)

func B64_string(input string, eod string) string {
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

func B64_data(input []byte, eod string) string {
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

func B32_string(input string, eod string) string {
	if eod == "e" || eod == "" {
		return base32.StdEncoding.EncodeToString([]byte(input))
	}
	if eod == "d" {
		decoded, err := base32.StdEncoding.DecodeString(input)
		if err != nil {
			return input
		}
		return string(decoded)
	}
	return "An error occured while encoding/decoding"
}

func B32_data(input []byte, eod string) string {
	if eod == "e" || eod == "" {
		return base32.StdEncoding.EncodeToString(input)
	}
	if eod == "d" {
		decoded, err := base32.StdEncoding.DecodeString(string(input))
		if err != nil {
			return string(input)
		}
		return string(decoded)
	}
	return "An error occured while encoding/decoding"
}

func B16_string(input string, eod string) string {
	if eod == "e" || eod == "" {
		return hex.EncodeToString([]byte(input))
	}
	if eod == "d" {
		decoded, err := hex.DecodeString(input)
		if err != nil {
			return input
		}
		return string(decoded)
	}
	return "An error occured while encoding/decoding"
}

func B16_data(input []byte, eod string) string {
	if eod == "e" || eod == "" {
		return hex.EncodeToString(input)
	}
	if eod == "d" {
		decoded, err := hex.DecodeString(string(input))
		if err != nil {
			return string(input)
		}
		return string(decoded)
	}
	return "An error occured while encoding/decoding"
}
