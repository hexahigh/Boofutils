package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"

	m "github.com/hexahigh/boofutils/modules"
)

//go:embed LICENSE
var LICENSE embed.FS

const AppVersion = "0.4.1 beta"

var subD_threads int
var skipTo, subD_domain, FIA_in, FIA_out string
var version, showLicense *bool
var FIA_decode, FIA_compress, update_binary, update_allow_win bool

func init() {
	version = flag.Bool("v", false, "Prints the current version")
	flag.StringVar(&skipTo, "s", "", "Skip the main menu and go to the selected task. Example Usage: -s 1")
	showLicense = flag.Bool("l", false, "Print the license")

	// usage
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", "Boofutils")
		flag.PrintDefaults()
		fmt.Println("Subcommands:")
		fmt.Println("subdomain -t <threads> -d <domain>")
		fmt.Println("update")
		fmt.Println("fileinaudio")
		fmt.Println("fileinimage")
	}

	// Subcommands
	subdomainCommand := flag.NewFlagSet("subdomain", flag.ExitOnError)
	subdomainCommand.IntVar(&subD_threads, "t", 10, "Number of threads to use")
	subdomainCommand.StringVar(&subD_domain, "d", "undef", "Domain to scan")

	updateCommand := flag.NewFlagSet("update", flag.ExitOnError)
	updateCommand.BoolVar(&update_binary, "b", false, "Update using a pre-compiled binary")
	updateCommand.BoolVar(&update_allow_win, "w", false, "Allow Windows")

	fileinaudioCommand := flag.NewFlagSet("fileinaudio", flag.ExitOnError)
	fileinaudioCommand.StringVar(&FIA_in, "i", "", "Input file")
	fileinaudioCommand.StringVar(&FIA_out, "o", "", "Output file")
	fileinaudioCommand.BoolVar(&FIA_decode, "d", false, "Decode")
	fileinaudioCommand.BoolVar(&FIA_compress, "nc", false, "Disable compression when encoding/decoding")

	fileinimageCommand := flag.NewFlagSet("fileinimage", flag.ExitOnError)
	fileinimageCommand.StringVar(&FIA_in, "i", "", "Input file")
	fileinimageCommand.StringVar(&FIA_out, "o", "", "Output file")
	fileinimageCommand.BoolVar(&FIA_decode, "d", false, "Decode")
	fileinimageCommand.BoolVar(&FIA_compress, "nc", false, "Disable compression when encoding/decoding")

	flag.Parse()

	subdomainCommand.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", "subdomain")
		subdomainCommand.PrintDefaults()
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "subdomain":
			subdomainCommand.Parse(os.Args[2:])
			m.SubD_main(subD_threads, subD_domain)
		case "update":
			updateCommand.Parse(os.Args[2:])
			m.Upd_main(update_binary, update_allow_win)
		case "fileinaudio":
			fileinaudioCommand.Parse(os.Args[2:])
			m.Fileinaudio_main(FIA_in, FIA_out, FIA_decode, FIA_compress)
		case "fileinimage":
			fileinimageCommand.Parse(os.Args[2:])
			m.Fileinimage_main(FIA_in, FIA_out, FIA_decode, FIA_compress)
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
		fmt.Println("\033[36mBoofutils\033[0m")
		fmt.Println("Version:", AppVersion)
		fmt.Println("2023 - Boofdev")
		fmt.Println("github.com/hexahigh/boofutils")
		fmt.Println(m.ColorPurple, "\nRandom cat fact:", m.ColorNone, m.RandomCatFact())
		os.Exit(0)
	}

	if !m.CheckConfigFileExists() {
		fmt.Println("Boofutils has not been configured yet. Would you like to answer some quick questions to get started?")
		fmt.Println("Y/N (Default: Y)")
		if m.AskInput() == "y" || m.AskInput() == "Y" {
			m.AskUserQuestions()
		} else {
			m.GenerateDefaultConfig()
		}
	}

	username, err := m.GetOptionFromConfig("name")
	if err != nil {
		log.Fatalf(err.Error())
	}

	if skipTo == "" {
		fmt.Println(m.Greet(), username+"!", "Welcome to Boofutils.")
		fmt.Println("What would you like to do today?")
		fmt.Println("[\033[36m1\033[0m] Print subcommands")
		fmt.Println("[\033[36m9\033[0m] Reconfigure Boofutils")
		fmt.Println("[\033[36m0\033[0m] Exit")
		checkInputAndDoStuff(m.AskInput())
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
		flag.Usage()
		os.Exit(0)
	case "9":
		m.AskUserQuestions()
	case "0":
		os.Exit(0)
	default:
		fmt.Println("Invalid input")
		os.Exit(0)
	}
}
