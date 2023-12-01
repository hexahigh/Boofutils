package modules

import (
	"bufio"
	"embed"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
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
	var subdomains []string
	for scanner.Scan() {
		subdomains = append(subdomains, scanner.Text())
	}

	p := mpb.New(mpb.WithWaitGroup(&sync.WaitGroup{}))

	bar := p.AddBar(int64(len(subdomains)),
		mpb.PrependDecorators(
			decor.CountersNoUnit("%d / %d", decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			decor.Percentage(decor.WCSyncSpace),
		),
	)

	var wg sync.WaitGroup
	for _, sub := range subdomains {
		wg.Add(1)
		go func(sub string) {
			defer wg.Done()
			subdomain := fmt.Sprintf("%s.%s", sub, domain)
			_, err := net.LookupHost(subdomain)
			if err == nil {
				fmt.Printf("%s has A or AAAA record\n", subdomain)
			}
			cname, err := net.LookupCNAME(subdomain)
			if err == nil {
				fmt.Printf("%s has CNAME record: %s\n", subdomain, cname)
			}
			bar.Increment()
		}(sub)
	}

	wg.Wait()
	p.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}
