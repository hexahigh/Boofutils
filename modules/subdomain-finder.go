package modules

import (
	"bufio"
	"embed"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/klauspost/compress/zstd"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

//go:embed embed/subdomains.txt.zst
var content embed.FS

func SubD_main() {
	var foundDomains []string

	fmt.Println("What domain would you like to scan?")
	domain := AskInput()

	fmt.Println("Disable the progress bar? Y/N (Default: N)")
	quiet := YNtoBool(AskInput())

	// remove http and https from url
	strings.Replace(domain, "http://", "", -1)
	strings.Replace(domain, "https://", "", -1)

	f, _ := content.Open("embed/subdomains.txt.zst")
	defer f.Close()

	d, err := zstd.NewReader(f)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer d.Close()

	scanner := bufio.NewScanner(d)
	var subdomains []string
	for scanner.Scan() {
		subdomains = append(subdomains, scanner.Text())
	}

	var p *mpb.Progress
	var bar *mpb.Bar
	if !quiet {
		p = mpb.New()
		bar = p.AddBar(int64(len(subdomains)),
			mpb.PrependDecorators(
				decor.CountersNoUnit("%d / %d", decor.WCSyncSpace),
			),
			mpb.AppendDecorators(
				decor.Percentage(decor.WCSyncSpace),
			),
		)
	}

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
			foundDomains = append(foundDomains, subdomain)
			cname, err := net.LookupCNAME(subdomain)
			if err == nil {
				fmt.Printf("%s has CNAME record: %s\n", subdomain, cname)
				foundDomains = append(foundDomains, subdomain)
			}
			if !quiet {
				bar.Increment()
			}
		}(sub)
	}

	wg.Wait()
	if !quiet {
		p.Wait()
	}

	wg.Wait()
	p.Wait()
	fmt.Println("Done! Found subdomains:")
	fmt.Println(foundDomains)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}
