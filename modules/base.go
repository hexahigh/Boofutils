package modules

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
)

func Base_main(base int, eod bool, stringORbin bool, input string) {
	if stringORbin == true {
		switch base {
		case 16:
			B16_string(input, eod)
		case 32:
			B32_string(input, eod)
		case 64:
			B64_string(input, eod)
		}
	}
}

func B64_string(input string, eod bool) string {
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

func B64_data(input []byte, eod bool) string {
	if eod == true {
		return base64.StdEncoding.EncodeToString(input)
	}
	if eod == false {
		decoded, err := base64.StdEncoding.DecodeString(string(input))
		if err != nil {
			return string(input)
		}
		return string(decoded)
	}
	return "An error occured while encoding/decoding"
}

func B32_string(input string, eod bool) string {
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

func B32_data(input []byte, eod bool) string {
	if eod == true {
		return base32.StdEncoding.EncodeToString(input)
	}
	if eod == false {
		decoded, err := base32.StdEncoding.DecodeString(string(input))
		if err != nil {
			return string(input)
		}
		return string(decoded)
	}
	return "An error occured while encoding/decoding"
}

func B16_string(input string, eod bool) string {
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

func B16_data(input []byte, eod bool) string {
	if eod == true {
		return hex.EncodeToString(input)
	}
	if eod == false {
		decoded, err := hex.DecodeString(string(input))
		if err != nil {
			return string(input)
		}
		return string(decoded)
	}
	return "An error occured while encoding/decoding"
}
