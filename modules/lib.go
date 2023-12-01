package modules

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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
