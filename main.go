package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"

	m "github.com/hexahigh/boofutils/modules"
	m_ansivid "github.com/hexahigh/boofutils/modules/ansivid"
)

//go:embed LICENSE
var LICENSE embed.FS

const AppVersion = "1.0.0"

var subD_threads int
var skipTo, subD_domain, FIA_in, FIA_out, bua_in, bua_out, ansiimg_filename, ansiimg_output string
var version, showLicense *bool
var FIA_decode, FIA_compress, update_binary, bua_encode, bua_b2, update_allow_win bool

var ansivid_musicFile, ansivid_gifFile, ansivid_gifSeq string
var ansivid_duration, ansivid_gifWidth, ansivid_gifHeight, ansivid_loopNum int
var ansivid_gifContrast, ansivid_gifSigma float64
var ansivid_gifMode, ansivid_gifAsciiMode, ansivid_blockMode bool

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

	buaCommand := flag.NewFlagSet("bua", flag.ExitOnError)
	buaCommand.StringVar(&bua_in, "i", "", "Comma separated list of input files/folders")
	buaCommand.StringVar(&bua_out, "o", "", "Output file/folder")
	buaCommand.BoolVar(&bua_encode, "e", false, "Create archive")
	buaCommand.BoolVar(&bua_b2, "b2", false, "Use bzip2 compression")

	ansiimgCommand := flag.NewFlagSet("ansiimg", flag.ExitOnError)
	ansiimgCommand.StringVar(&ansiimg_filename, "i", "", "Input file")
	ansiimgCommand.StringVar(&ansiimg_output, "o", "", "Output file")

	ansividCommand := flag.NewFlagSet("ansivid", flag.ExitOnError)
	// func Ansivid_main(musicFile string, gifWidth int, gifHeight int, duration int, gifFile string, gifSeq string, loopNum int, gifMode bool, gifContrast float64, gifAsciiMode bool, gifSigma float64, blockMode bool) {
	ansividCommand.StringVar(&ansivid_musicFile, "a", "", "AUdio file")
	ansividCommand.StringVar(&ansivid_gifFile, "g", "", "GIF file")
	ansividCommand.StringVar(&ansivid_gifSeq, "s", "0", "GIF sequence")
	ansividCommand.IntVar(&ansivid_duration, "d", 10, "GIF duration")
	ansividCommand.IntVar(&ansivid_gifWidth, "w", 100, "GIF width")
	ansividCommand.IntVar(&ansivid_gifHeight, "h", 100, "GIF height")
	ansividCommand.Float64Var(&ansivid_gifContrast, "c", 0, "GIF contrast")
	ansividCommand.Float64Var(&ansivid_gifSigma, "sigma", 0, "GIF sigma")
	ansividCommand.BoolVar(&ansivid_gifMode, "m", false, "GIF mode")
	ansividCommand.BoolVar(&ansivid_gifAsciiMode, "ascii", false, "GIF ascii mode")
	ansividCommand.BoolVar(&ansivid_blockMode, "block", false, "GIF block mode")
	ansiimgCommand.IntVar(&ansivid_loopNum, "l", 1, "GIF loop number")

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
		case "bua":
			buaCommand.Parse(os.Args[2:])
			m.Bua_main(bua_in, bua_out, bua_encode, bua_b2)
			os.Exit(0)
		case "ansiimg":
			ansiimgCommand.Parse(os.Args[2:])
			m.Ansiimg_main(ansiimg_filename, ansiimg_output)
			os.Exit(0)
		case "ansivid":
			ansividCommand.Parse(os.Args[2:])
			m_ansivid.Ansivid_main(ansivid_musicFile, ansivid_gifWidth, ansivid_gifHeight, ansivid_duration, ansivid_gifFile, ansivid_gifSeq, ansivid_loopNum, ansivid_gifMode, ansivid_gifContrast, ansivid_gifAsciiMode, ansivid_gifSigma, ansivid_blockMode)
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
