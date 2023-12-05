package modules

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
)

// Base_main_string takes an input string, a boolean value saying if the input should be encoded or decoded.
// The boolean should be true if it should be encoded and false if it should be decoded.
// It returns the encoded or decoded string.
func Base_main_string(input string, eod bool, base int) string {
	switch base {
	case 16:
		return B16(input, eod)
	case 32:
		return B32(input, eod)
	case 64:
		return B64(input, eod)
	}
	return "An error occured while encoding/decoding"
}

func B64(input string, eod bool) string {
	if eod == true {
		return base64.StdEncoding.EncodeToString([]byte(input))
	}
	if eod == false {
		decoded, err := base64.StdEncoding.DecodeString(input)
		if err != nil {
			return input
		}
		return string(decoded)
	}
	return "An error occured while encoding/decoding"
}

func B32(input string, eod bool) string {
	if eod == true {
		return base32.StdEncoding.EncodeToString([]byte(input))
	}
	if eod == false {
		decoded, err := base32.StdEncoding.DecodeString(input)
		if err != nil {
			return input
		}
		return string(decoded)
	}
	return "An error occured while encoding/decoding"
}

func B16(input string, eod bool) string {
	if eod == true {
		return hex.EncodeToString([]byte(input))
	}
	if eod == false {
		decoded, err := hex.DecodeString(input)
		if err != nil {
			return input
		}
		return string(decoded)
	}
	return "An error occured while encoding/decoding"
}
