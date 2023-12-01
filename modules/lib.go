package modules

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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

func YNtoBool(input string) bool {
	input = strings.ToLower(input)
	return input == "y" || input == "Y" || input == "yes" || input == "true"
}
