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
	report "github.com/hexahigh/boofutils/modules/report"
)

//go:embed LICENSE
var LICENSE embed.FS

const AppVersion = "1.6.3"

var subD_threads int
var skipTo, subD_domain, FIA_in, FIA_out, bua_in, bua_out, ansiimg_filename, ansiimg_output string
var version, showLicense *bool
var FIA_decode, FIA_compress, update_binary, bua_encode, bua_b2, update_allow_win, bua_mute, ansivid_gifMode, ansivid_gifAsciiMode, ansivid_blockMode, ansivid_py_install, donutRainbow, donutColor, chachacha_decrypt, chachacha_mute bool

var ansivid_musicFile, ansivid_gifFile, ansivid_gifSeq, ansivid_py_strat, ansivid_py_in, scraper_allowedDomains, scraper_template, chachacha_in, chachacha_out, chachacha_password string
var ansivid_duration, ansivid_gifWidth, ansivid_gifHeight, ansivid_loopNum int
var ansivid_gifContrast, ansivid_gifSigma, donutSpeed float64
var ansiimg_width, ansiimg_height uint

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
		fmt.Println("bua")
		fmt.Println("ansivid")
		fmt.Println("ansiimg")
		fmt.Println("chachacha")
		fmt.Println("scraper")
		fmt.Println("donut")
		fmt.Println("report")
	}

	// Subcommands
	subdomainCommand := flag.NewFlagSet("subdomain", flag.ExitOnError)
	subdomainCommand.IntVar(&subD_threads, "t", 10, "Number of threads to use")
	subdomainCommand.StringVar(&subD_domain, "d", "undef", "Domain to scan")

	updateCommand := flag.NewFlagSet("update", flag.ExitOnError)

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
	ansiimgCommand := flag.NewFlagSet("ansiimg", flag.ExitOnError)
	donutCommand := flag.NewFlagSet("donut", flag.ExitOnError)
	ansividCommand := flag.NewFlagSet("ansivid", flag.ExitOnError)
	ansivid_pyCommand := flag.NewFlagSet("ansivid-py", flag.ExitOnError)
	scraperCommand := flag.NewFlagSet("scraper", flag.ExitOnError)
	chachachaCommand := flag.NewFlagSet("chachacha", flag.ExitOnError)
	urlCommand := flag.NewFlagSet("url", flag.ExitOnError)
	reportCommand := flag.NewFlagSet("report", flag.ExitOnError)

	flag.Parse()

	subdomainCommand.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", "subdomain")
		subdomainCommand.PrintDefaults()
	}

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
		fmt.Println(m.ColorCyanBold24bit, "Boofutils", m.ColorReset)
		fmt.Println("Version:", AppVersion)
		fmt.Println("Boofdev")
		fmt.Println("github.com/hexahigh/boofutils")
		fmt.Println(m.ColorPurpleBold24bit, "\nRandom cat fact:", m.ColorOrangeBold24bit, m.RandomCatFact())
		os.Exit(0)
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "report":
			out := reportCommand.String("o", "report.json", "Output file")
			stdout := reportCommand.Bool("s", false, "Print to stdout")
			pl := reportCommand.Int("pl", 0, "Print level. -1: no output, 0: only errors, 1: verbose, 2: very verbose")
			reportCommand.Parse(os.Args[2:])
			report.Report(*out, *stdout, *pl)
			os.Exit(0)
		case "subdomain":
			subdomainCommand.Parse(os.Args[2:])
			m.SubD_main(subD_threads, subD_domain)
			os.Exit(0)
		case "url":
			var urlCommandURL string
			var urlCommandBrute bool
			var urlCommandThreads int
			urlCommand.StringVar(&urlCommandURL, "u", "", "URL to scan")
			urlCommand.IntVar(&urlCommandThreads, "t", 10, "Number of threads to use")
			urlCommand.BoolVar(&urlCommandBrute, "b", false, "Bruteforce")
			urlCommand.Parse(os.Args[2:])
			m.Url_main(urlCommandThreads, urlCommandURL, urlCommandBrute)
		case "update":
			var update_ignore_req, update_allow_win, update_binary bool
			updateCommand.BoolVar(&update_binary, "b", false, "Update using a pre-compiled binary")
			updateCommand.BoolVar(&update_allow_win, "w", false, "Allow Windows")
			updateCommand.BoolVar(&update_ignore_req, "ignore-req", false, "Ignore requirements")
			updateCommand.Parse(os.Args[2:])
			m.Upd_main(update_binary, update_allow_win, update_ignore_req)
		case "fileinaudio":
			fileinaudioCommand.Parse(os.Args[2:])
			m.Fileinaudio_main(FIA_in, FIA_out, FIA_decode, FIA_compress)
			os.Exit(0)
		case "fileinimage":
			fileinimageCommand.Parse(os.Args[2:])
			m.Fileinimage_main(FIA_in, FIA_out, FIA_decode, FIA_compress)
			os.Exit(0)
		case "bua":
			buaCommand.StringVar(&bua_in, "i", "", "Comma separated list of input files/folders")
			buaCommand.StringVar(&bua_out, "o", "", "Output file/folder")
			buaCommand.BoolVar(&bua_encode, "e", false, "Create archive")
			buaCommand.BoolVar(&bua_b2, "b2", false, "Use bzip2 compression")
			buaCommand.BoolVar(&bua_mute, "m", false, "Mute audio")
			buaCommand.Parse(os.Args[2:])
			m.Bua_main(bua_in, bua_out, bua_encode, bua_b2, bua_mute)
			os.Exit(0)
		case "ansiimg":
			ansiimgCommand.StringVar(&ansiimg_filename, "i", "", "Input file")
			ansiimgCommand.StringVar(&ansiimg_output, "o", "", "Output file")
			ansiimgCommand.UintVar(&ansiimg_width, "w", 100, "Width")
			ansiimgCommand.UintVar(&ansiimg_height, "h", 0, "Height")
			ansiimgCommand.Parse(os.Args[2:])
			m.Ansiimg_main(ansiimg_filename, ansiimg_output, ansiimg_width, ansiimg_height)
			os.Exit(0)
		case "ansivid":
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
			ansividCommand.IntVar(&ansivid_loopNum, "l", 1, "GIF loop number")
			ansividCommand.Parse(os.Args[2:])
			m_ansivid.Ansivid_main(ansivid_musicFile, ansivid_gifWidth, ansivid_gifHeight, ansivid_duration, ansivid_gifFile, ansivid_gifSeq, ansivid_loopNum, ansivid_gifMode, ansivid_gifContrast, ansivid_gifAsciiMode, ansivid_gifSigma, ansivid_blockMode)
			os.Exit(0)
		case "donut":
			donutCommand.BoolVar(&donutColor, "tc", false, "Use true color")
			donutCommand.Float64Var(&donutSpeed, "s", 1, "Speed")
			donutCommand.BoolVar(&donutRainbow, "r", false, "Rainbow")
			donutCommand.Parse(os.Args[2:])
			m.Donut_main(donutSpeed, donutRainbow, donutColor)
			os.Exit(0)
		case "ansivid-py":
			ansivid_pyCommand.StringVar(&ansivid_py_in, "i", "", "Input file")
			ansivid_pyCommand.StringVar(&ansivid_py_strat, "s", "1", "Strategy")
			ansivid_pyCommand.BoolVar(&ansivid_py_install, "ins", false, "Install")
			ansivid_pyCommand.Parse(os.Args[2:])
			m.Ansivid_py_main(ansivid_py_in, ansivid_py_strat, ansivid_py_install)
			os.Exit(0)
		case "scraper":
			domainPtr := scraperCommand.String("domain", "https://example.com", "The domain to scan")
			outputFilePtr := scraperCommand.String("output", "urls.txt", "The output file")
			scraperCommand.StringVar(&scraper_allowedDomains, "allowed", "", "CS Allowed domains")
			scraperCommand.StringVar(&scraper_template, "t", "", "Template")
			scraperCommand.Parse(os.Args[2:])
			m.Scrape_main(*domainPtr, *outputFilePtr, scraper_allowedDomains, scraper_template)
			os.Exit(0)
		case "chachacha":
			var chachacha_in, chachacha_out, chachacha_password, chachacha_keyfile string
			var chachacha_verbose bool
			chachachaCommand.StringVar(&chachacha_in, "i", "", "Input file")
			chachachaCommand.StringVar(&chachacha_out, "o", "", "Output file")
			chachachaCommand.BoolVar(&chachacha_decrypt, "d", false, "Decrypt")
			chachachaCommand.StringVar(&chachacha_password, "p", "", "Password")
			chachachaCommand.BoolVar(&chachacha_mute, "m", false, "Mute audio")
			chachachaCommand.StringVar(&chachacha_keyfile, "k", "", "Use a keyfile as password")
			chachachaCommand.BoolVar(&chachacha_verbose, "v", false, "Verbose")
			chachachaCommand.Parse(os.Args[2:])
			m.Chacha_main(chachacha_password, chachacha_decrypt, chachacha_in, chachacha_out, chachacha_mute, chachacha_keyfile, chachacha_verbose)
			os.Exit(0)
		default:
		}
	}
}

func main() {
	if !m.CheckConfigFileExists() {
		fmt.Println("Boofutils has not been configured yet. Would you like to answer some quick questions to get started?")
		fmt.Println("Y/N (Default: Y)")
		if m.YNtoBool(m.AskInput()) == true {
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
