package modules

import (
	"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/klauspost/compress/zstd"
)

//go:embed embed/cat_facts.json.zst
var cat embed.FS

// AskInput prompts the user for input from the command line.
//
// It reads input from the standard input using bufio.NewReader and os.Stdin.
// If there is an error reading the input, it prints an error message.
//
// It trims the newline character at the end of the input string using
// strings.TrimSpace.
//
// It returns the input as a string.
func AskInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
	}
	// Trim the newline character at the end
	input = strings.TrimSpace(input)
	return input
}

// YNtoBool converts a string to a boolean value.
//
// It takes a string as input and returns a boolean value. The function
// converts the input string to lowercase and checks if it is equal to
// "y", "yes", or "true". If any of these conditions are true, the
// function returns true; otherwise, it returns false.
func YNtoBool(input string) bool {
	input = strings.ToLower(input)
	return input == "y" || input == "yes" || input == "true"
}

func RandomCatFact() string {
	// This method uses less memory but is slower.
	// I really need to decrease this programs memory usage.
	f, _ := cat.Open("embed/cat_facts.json.zst")
	defer f.Close()

	d, err := zstd.NewReader(f)
	if err != nil {
		fmt.Println(err)
	}
	defer d.Close()
	bytes, err := io.ReadAll(d)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var catJson map[string]interface{}
	json.Unmarshal(bytes, &catJson)

	facts, ok := catJson["facts"].([]interface{})
	if !ok {
		fmt.Println("Key 'facts' not found in JSON")
		return ""
	}

	randFact := facts[rand.Intn(len(facts))]

	return randFact.(string)
}

func CheckIfError(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func FileSize(filePath string) int64 {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return 0
	}
	defer file.Close()

	// Get the file size
	fi, err := file.Stat()
	if err != nil {
		return 0
	}

	return fi.Size()
}

func waitForKeypress() {
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	_, _, err = keyboard.GetKey()
	if err != nil {
		panic(err)
	}
	os.Exit(0)
}

func verbosePrintln(s string, v bool) {
	if v {
		fmt.Println(s)
	}
}
