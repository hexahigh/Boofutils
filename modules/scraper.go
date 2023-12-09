package modules

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
)

var (
	urls []string
	mux  sync.Mutex
)

func RandomString(userAgentList []string) string {
	randomIndex := rand.Intn(len(userAgentList))
	return userAgentList[randomIndex]
}

func Scrape_main(domain string, outputFile string, allowedDomains string, template string) {
	aDomain := strings.Split(allowedDomains, ",")
	aDomain = append(aDomain, domain)

	userAgentList := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36 (compatible; Boofutils crawler; +admin@boofdev.eu)",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_4_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Mobile/15E148 Safari/604.1 (compatible; Boofutils crawler; +admin@boofdev.eu)",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36 Edg/87.0.664.75 (compatible; Boofutils crawler; +admin@boofdev.eu)",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36 Edge/18.18363 (compatible; Boofutils crawler; +admin@boofdev.eu)",
	}

	if template == "nrk" {
		aDomain = append(aDomain, "https://gfx.nrk.no", "https://tv.nrk.no", "https://info.nrk.no", "https://nrk.no", "https://www.nrk.no", "https://blog.nrk.no", "https://api.nrk.no", "https://img.nrk.no")
	}

	fmt.Println("Allowed domains:", aDomain)

	c := colly.NewCollector(
		colly.Async(true),
		//colly.AllowURLRevisit(),
	)

	c.OnHTML("[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		absoluteLink := e.Request.AbsoluteURL(link)
		absoluteURL, err := url.Parse(absoluteLink)
		if err != nil {
			log.Fatal(err)
		}
		for _, d := range aDomain {
			domainURL, err := url.Parse(d)
			if err != nil {
				log.Fatal(err)
			}
			if absoluteURL.Host == domainURL.Host {
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
		absoluteURL, err := url.Parse(absoluteLink)
		if err != nil {
			log.Fatal(err)
		}
		for _, d := range aDomain {
			domainURL, err := url.Parse(d)
			if err != nil {
				log.Fatal(err)
			}
			if absoluteURL.Host == domainURL.Host {
				mux.Lock()
				urls = append(urls, absoluteLink)
				mux.Unlock()
				c.Visit(absoluteLink)
				writeToFile(outputFile, absoluteLink)
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		if r.StatusCode == 429 {
			fmt.Println("Got ratelimited.")
			r.Request.Retry()
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
		r.Headers.Set("User-Agent", RandomString(userAgentList))
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

func extractURLs(text string) []string {
	re := regexp.MustCompile(`http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\\(\\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+`)
	return re.FindAllString(text, -1)
}
