package modules

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
)

var (
	urls []string
	mux  sync.Mutex
)

func Scrape_main(domain string, outputFile string, allowedDomains string, template string) {
	aDomain := strings.Split(allowedDomains, ",")
	aDomain = append(aDomain, domain)

	if template == "nrk" {
		aDomain = append(aDomain, "gfx.nrk.no", "tv.nrk.no", "info.nrk.no", "nrk.no", "www.nrk.no")
	}

	fmt.Println("Allowed domains:", aDomain)

	c := colly.NewCollector(
		colly.Async(true),
	)

	c.OnHTML("[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		absoluteLink := e.Request.AbsoluteURL(link)
		//fmt.Println("Absolute url of HREF element:", absoluteLink)
		for _, d := range aDomain {
			if strings.Contains(absoluteLink, d) {
				mux.Lock()
				urls = append(urls, absoluteLink)
				mux.Unlock()
				c.Visit(absoluteLink)
				writeToFile(outputFile, absoluteLink)
			}
		}
	})

	c.OnHTML("[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		absoluteLink := e.Request.AbsoluteURL(link)
		//fmt.Println("Absolute url of SRC element:", absoluteLink)
		for _, d := range aDomain {
			if strings.Contains(absoluteLink, d) {
				mux.Lock()
				urls = append(urls, absoluteLink)
				mux.Unlock()
				c.Visit(absoluteLink)
				writeToFile(outputFile, absoluteLink)
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit(domain)
	c.Wait()
}

func writeToFile(outputFile string, data string) {
	f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if _, err := f.WriteString(data + "\n"); err != nil {
		log.Fatal(err)
	}
}
