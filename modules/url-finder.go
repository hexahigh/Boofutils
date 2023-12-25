package modules

import (
	"bufio"
	"embed"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/klauspost/compress/zstd"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

//go:embed embed/urls.txt.zst
var urlList embed.FS

func Url_main(threads int, domain string, bruteForce bool) {
	var foundURLs []string

	if domain == "undef" {
		fmt.Println("What domain would you like to scan?")
		domain = AskInput()
	}
	fmt.Println("Disable the progress bar? Y/N (Default: N)")
	quiet := YNtoBool(AskInput())

	// remove http and https from url
	strings.Replace(domain, "http://", "", -1)
	strings.Replace(domain, "https://", "", -1)

	var urls []string
	if bruteForce {
		charset := "abcdefghijklmnopqrstuvwxyz0123456789.,-_ABCDEFGHIJKLMNOPQRSTUVWXYZØÆÅæåø"
		for i := 1; i <= 10; i++ {
			urls = append(urls, generateCombinations(charset, i)...)
		}
	} else {
		f, _ := urlList.Open("embed/urls.txt.zst")
		defer f.Close()

		d, err := zstd.NewReader(f)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer d.Close()

		scanner := bufio.NewScanner(d)
		for scanner.Scan() {
			urls = append(urls, scanner.Text())
		}
	}

	var p *mpb.Progress
	var bar *mpb.Bar
	if !quiet {
		p = mpb.New()
		bar = p.AddBar(int64(len(urls)),
			mpb.PrependDecorators(
				decor.CountersNoUnit("%d / %d", decor.WCSyncSpace),
			),
			mpb.AppendDecorators(
				decor.Percentage(decor.WCSyncSpace),
			),
		)
	}

	fmt.Println("Starting scan with " + fmt.Sprint(threads) + " threads...")
	var wg sync.WaitGroup
	sem := make(chan struct{}, threads)

	for _, u := range urls {
		sem <- struct{}{} // Acquire a semaphore slot
		wg.Add(1)
		go func(u string) {
			defer func() {
				<-sem // Release the semaphore slot
				wg.Done()
			}()
			fullURL := fmt.Sprintf("%s/%s", domain, u)
			_, err := url.ParseRequestURI(fullURL)
			if err == nil {
				resp, err := http.Get(fullURL)
				if err == nil && resp.StatusCode == http.StatusOK {
					fmt.Printf("%s is a valid URL\n", fullURL)
					foundURLs = append(foundURLs, fullURL)
				}
			}
			if !quiet {
				bar.Increment()
			}
		}(u)
	}

	wg.Wait()
	if !quiet {
		p.Wait()
	}
	fmt.Println("Done! Found URLs:")
	fmt.Println(foundURLs)
}

func generateCombinations(alphabet string, length int) []string {
	if length == 0 {
		return []string{""}
	}

	var result []string
	for _, letter := range alphabet {
		for _, combination := range generateCombinations(alphabet, length-1) {
			result = append(result, string(letter)+combination)
		}
	}
	return result
}
