package main

import (
	"boofutils/modules"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
)

const AppVersion = "0.2.6 beta"

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

	if !modules.CheckConfigFileExists() {
		fmt.Println("Boofutils has not been configured yet. Would you like to answer some quick questions to get started?")
		fmt.Println("Y/N (Default: Y)")
		if modules.AskInput() == "y" || modules.AskInput() == "Y" {
			modules.AskUserQuestions()
		} else {
			modules.GenerateDefaultConfig()
		}
	}

	username, err := modules.GetOptionFromConfig("name")
	if err != nil {
		log.Fatalf(err.Error())
	}

	if skipTo == "" {
		fmt.Println(modules.Greet(), username+"!", "Welcome to Boofutils.")
		fmt.Println("What would you like to do today?")
		fmt.Println("[1] Calculate hashes of file")
		fmt.Println("[2] Print a file as hexadecimal (Base16)")
		fmt.Println("[3] Subdomain Finder")
		fmt.Println("[9] Reconfigure Boofutils")
		fmt.Println("[0] Exit")
		checkInputAndDoStuff(modules.AskInput())
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

func checkInputAndDoStuff(input string) {
	switch input {
	case "1":
		modules.Hf_main()
	case "2":
		modules.Hex_main()
	case "3":
		modules.SubD_main()
	case "9":
		modules.AskUserQuestions()
	case "0":
		os.Exit(0)
	default:
		fmt.Println("Invalid input")
		os.Exit(0)
	}
}
