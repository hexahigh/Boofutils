package main

import (
	"boofutils/modules"
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
)

//go:embed LICENSE
var LICENSE embed.FS

const AppVersion = "0.3.3 beta"

var subD_threads int
var skipTo, subD_domain string
var version, showLicense *bool

func init() {
	version = flag.Bool("v", false, "Prints the current version")
	flag.StringVar(&skipTo, "s", "", "Skip the main menu and go to the selected task. Example Usage: -s 1")
	showLicense = flag.Bool("l", false, "Print the license")

	// Subcommands
	subdomainCommand := flag.NewFlagSet("subdomain", flag.ExitOnError)
	subdomainCommand.IntVar(&subD_threads, "t", 10, "Number of threads to use")
	subdomainCommand.StringVar(&subD_domain, "d", "undef", "Domain to scan")

	updateCommand := flag.NewFlagSet("update", flag.ExitOnError)

	flag.Parse()

	subdomainCommand.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", "subdomain")
		subdomainCommand.PrintDefaults()
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "subdomain":
			subdomainCommand.Parse(os.Args[2:])
			modules.SubD_main(subD_threads, subD_domain)
		case "update":
			updateCommand.Parse(os.Args[2:])
			modules.Upd_main()
		default:
		}
	}
}

func main() {

	if *showLicense {
		data, err := LICENSE.ReadFile("LICENSE")
		if err != nil {
			fmt.Println("Error reading file:", err)
			os.Exit(1)
		}
		fmt.Println(string(data))
		os.Exit(0)
	}

	if *version {
		fmt.Println("\033[36mBoofutils\033[0m\nVersion:", AppVersion, "\nBoof Development")
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
		fmt.Println("[\033[36m1\033[0m] Calculate hashes of file")
		fmt.Println("[\033[36m2\033[0m] Print a file as hexadecimal (Base16)")
		fmt.Println("[\033[36m3\033[0m] Subdomain Finder")
		fmt.Println("[\033[36m9\033[0m] Reconfigure Boofutils")
		fmt.Println("[\033[36m0\033[0m] Exit")
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
		modules.SubD_main(subD_threads, subD_domain)
	case "9":
		modules.AskUserQuestions()
	case "0":
		os.Exit(0)
	default:
		fmt.Println("Invalid input")
		os.Exit(0)
	}
}
