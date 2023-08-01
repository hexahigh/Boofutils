package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
)

func main() {
	fmt.Println("Hello", getName()+"!", "Welcome to Boofutils.")
	fmt.Println("What would you like to do today?")
	fmt.Println("[1] Calculate hashes of file")
	fmt.Println("[0] Exit")
	checkInputAndDoStuff(askInput())
}

func getName() string {
	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}
	username := user.Username
	return username
}

func askInputOLD() string {
	var input string
	fmt.Scanln(&input)
	return input
}

func askInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
	}
	// Trim the newline character at the end
	input = strings.TrimSpace(input)
	return input
}

func checkInputAndDoStuff(input string) {
	if input == "1" {
		hf_main()
	}
	if input == "0" {
		os.Exit(0)
	}
	if input != "1" && input != "0" {
		fmt.Println("Invalid input")
		checkInputAndDoStuff(askInput())
	}
}
