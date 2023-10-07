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

const AppVersion = "0.2.2 beta"

var skipTo string
var version *bool

func init() {
	version = flag.Bool("v", false, "Prints the current version")
	flag.StringVar(&skipTo, "skip", "", "Skip the main menu and go to the selected task. Example Usage: -skip 1")
	flag.Parse()
}

func main() {

	if *version {
		fmt.Println(AppVersion)
		os.Exit(0)
	}
	if skipTo == "" {
		fmt.Println(Greet(), getName()+"!", "Welcome to Boofutils.")
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
	switch input {
	case "1":
		hf_main()
	case "2":
		hex_main()
	case "0":
		os.Exit(0)
	default:
		fmt.Println("Invalid input")
		os.Exit(0)
	}
}
