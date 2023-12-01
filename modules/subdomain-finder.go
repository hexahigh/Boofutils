package modules

import (
	"bufio"
	"embed"
	"fmt"
	"net/http"
	"strings"
	"time"
)

//go:embed embed/subdomains.txt
var content embed.FS

func SubD_main() {
	fmt.Println("What domain would you like to scan?")
	domain := AskInput()

	// remove http and https from url
	strings.Replace(domain, "http://", "", -1)
	strings.Replace(domain, "https://", "", -1)

	data, err := content.ReadFile("embed/subdomains.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		sub := scanner.Text()
		subdomain := fmt.Sprintf("http://%s.%s", sub, domain)
		client := http.Client{
			Timeout: time.Second * 2, // Timeout after 2 seconds
		}
		resp, err := client.Get(subdomain)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if resp.StatusCode == http.StatusOK {
			fmt.Printf("%s is reachable\n", subdomain)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}