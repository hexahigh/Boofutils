package modules

import (
	"bufio"
	"embed"
	"fmt"
	"net/http"
	"context"
	"strings"
	"sync"
	"time"

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
			subdomain := fmt.Sprintf("http://%s.%s", sub, domain)
			client := http.Client{
				Timeout: time.Second * 10, // Timeout after N seconds
			}
			resp, err := client.Get(subdomain)
			if err != nil {
				if err != context.DeadlineExceeded {
					fmt.Println(err)
				}
				return
			}
			if resp.StatusCode == http.StatusOK {
				fmt.Printf("%s is reachable\n", subdomain)
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