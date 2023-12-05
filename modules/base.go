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

// B64 encodes or decodes a given input string using base64 encoding.
//
// Parameters:
//   - input: The string to be encoded or decoded.
//   - eod: A boolean indicating whether to encode or decode the input string.
//     If true, the input string will be encoded. If false, the input string will be decoded.
//
// Returns:
// - The encoded or decoded string, depending on the value of the 'eod' parameter.
// - If an error occurs during decoding, the original input string is returned.
// - If an error occurs during encoding or decoding, an error message is returned.
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

// B32 encodes or decodes a string using base32 encoding.
//
// The function takes two parameters:
// - input: the string to be encoded or decoded.
// - eod: a boolean flag indicating whether to encode or decode.
//
// The function returns a string.
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

// B16 encodes or decodes a string using the base 16 (hexadecimal) encoding.
//
// It takes two parameters:
// - input: the string to be encoded or decoded.
// - eod: a boolean flag indicating whether to encode or decode the input.
//
// If eod is true, the function will encode the input using hexadecimal encoding and return the encoded string.
// If eod is false, the function will decode the input using hexadecimal encoding and return the decoded string.
// If an error occurs during decoding, the function will return the original input string.
//
// The function returns a string representing the encoded or decoded result.
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
