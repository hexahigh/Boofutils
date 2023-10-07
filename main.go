package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
)

var skipTo string

func init() {
	flag.StringVar(&skipTo, "skip", "", "Skip the main menu and go to the selected task. Usage: -skip 1")
	flag.Parse()
}

func main() {
	fmt.Println(skipTo)
	if skipTo == "" {
		fmt.Println("Hello", getName()+"!", "Welcome to Boofutils.")
		fmt.Println("What would you like to do today?")
		fmt.Println("[1] Calculate hashes of file")
		fmt.Println("[2] Print a file as hexadecimal (Base16)")
		fmt.Println("[0] Exit")
		checkInputAndDoStuff(askInput())
	} else {
		checkInputAndDoStuff(skipTo)
	}
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
	if input == "2" {
		hex_main()
	}
	if input == "0" {
		os.Exit(0)
	}
	// There is definetely a better way to do this but it works
	if input != "1" && input != "2" && input != "0" {
		fmt.Println("Invalid input")
		os.Exit(0)
	}
}
